package mailserver

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/mail"
	"strings"
	"time"

	"github.com/emersion/go-message"
	"github.com/emersion/go-smtp"
)

// SMTPSession 实现 smtp.Session 接口
type SMTPSession struct {
	backend       *SMTPBackend
	conn          *smtp.Conn
	from          string
	to            []string
	serverType    SMTPServerType // 服务器类型
	authenticated bool           // 认证状态
	requireAuth   bool           // 是否要求认证
	authUser      string         // 已认证的用户
}

// AuthPlain 处理PLAIN认证
func (s *SMTPSession) AuthPlain(username, password string) error {
	// 记录认证请求
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("🔐 SMTP认证请求 [%s]: %s", serverTypeStr, username)

	// MTA服务器(25端口)通常不需要认证，但如果有认证请求也要处理
	// MSA服务器(587端口)必须要求认证
	if s.serverType == SMTPServerTypeReceive {
		// MTA: 可选认证，主要用于中继控制
		log.Printf("⚠️  MTA服务器收到认证请求，将验证但不强制要求")
	}

	// 验证邮箱格式
	if !strings.Contains(username, "@") {
		log.Printf("❌ 认证失败: 无效的邮箱格式 %s [%s]", username, serverTypeStr)
		return fmt.Errorf("invalid email format")
	}

	// 使用存储层验证凭据
	if !s.backend.storage.ValidateCredentials(username, password) {
		log.Printf("❌ 认证失败: 用户名或密码错误 %s [%s]", username, serverTypeStr)
		return fmt.Errorf("invalid credentials")
	}

	// 认证成功
	s.authenticated = true
	s.authUser = username
	log.Printf("✅ 认证成功: %s [%s]", username, serverTypeStr)
	return nil
}

// Mail 处理MAIL FROM命令
func (s *SMTPSession) Mail(from string, opts *smtp.MailOptions) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📤 MAIL FROM: %s [%s]", from, serverTypeStr)

	// MSA服务器必须要求认证
	if s.requireAuth && !s.authenticated {
		log.Printf("❌ MSA服务器要求认证，但未认证 [%s]", serverTypeStr)
		return fmt.Errorf("authentication required")
	}

	// 验证发件人地址格式
	if _, err := mail.ParseAddress(from); err != nil {
		log.Printf("❌ 无效的发件人地址: %s, 错误: %v [%s]", from, err, serverTypeStr)
		return fmt.Errorf("invalid sender address: %v", err)
	}

	// MSA服务器需要验证发件人权限
	if s.serverType == SMTPServerTypeSubmit && s.authenticated {
		// TODO: 验证认证用户是否有权限使用此发件人地址
		// 这里应该检查认证用户是否匹配发件人地址或有权限代发
		log.Printf("🔍 验证发件人权限: %s (认证用户: %s)", from, s.authUser)
	}

	s.from = from
	s.to = []string{} // 重置收件人列表

	log.Printf("✅ 发件人设置成功: %s [%s]", from, serverTypeStr)
	return nil
}

// Rcpt 处理RCPT TO命令
func (s *SMTPSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📥 RCPT TO: %s [%s]", to, serverTypeStr)

	// 验证收件人地址格式
	if _, err := mail.ParseAddress(to); err != nil {
		log.Printf("❌ 无效的收件人地址: %s, 错误: %v [%s]", to, err, serverTypeStr)
		return fmt.Errorf("invalid recipient address: %v", err)
	}

	// MTA服务器需要检查是否为本地域名
	if s.serverType == SMTPServerTypeReceive {
		// TODO: 检查收件人是否为本地域名的邮箱
		// 如果不是本地域名，应该拒绝接收（防止成为开放中继）
		log.Printf("🔍 MTA检查收件人域名: %s", to)
	}

	// 检查是否超过最大收件人数量
	maxRecipients := 50
	if s.serverType == SMTPServerTypeReceive {
		maxRecipients = 100 // MTA可以接受更多收件人
	}

	if len(s.to) >= maxRecipients {
		log.Printf("❌ 收件人数量超过限制: %d [%s]", len(s.to), serverTypeStr)
		return fmt.Errorf("too many recipients")
	}

	s.to = append(s.to, to)
	log.Printf("✅ 收件人添加成功: %s (总数: %d) [%s]", to, len(s.to), serverTypeStr)
	return nil
}

// Data 处理DATA命令
func (s *SMTPSession) Data(r io.Reader) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📨 开始接收邮件数据... [%s]", serverTypeStr)

	if s.from == "" {
		return fmt.Errorf("no sender specified")
	}

	if len(s.to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// 解析邮件
	msg, err := message.Read(r)
	if err != nil {
		log.Printf("❌ 解析邮件失败: %v [%s]", err, serverTypeStr)
		return fmt.Errorf("failed to parse message: %v", err)
	}

	// 读取邮件正文
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("❌ 读取邮件正文失败: %v [%s]", err, serverTypeStr)
		return fmt.Errorf("failed to read message body: %v", err)
	}
	log.Printf("📊 邮件数据大小: %d 字节 [%s]", len(body), serverTypeStr)

	// MIME 头解码器， 解码标题
	decoder := new(mime.WordDecoder)
	subject, err := decoder.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		subject = msg.Header.Get("Subject") // 解码失败就用原文
	}

	// 根据服务器类型进行不同处理
	if s.serverType == SMTPServerTypeSubmit {
		// MSA: 用户提交的邮件，需要添加发送者信息和DKIM签名
		log.Printf("📤 处理用户提交邮件: %s", subject)
		// TODO: 添加DKIM签名、设置发送时间等
	} else {
		// MTA: 接收的外部邮件，需要进行垃圾邮件检查
		log.Printf("📥 处理接收邮件: %s", subject)
		// TODO: 垃圾邮件检查、病毒扫描等
	}

	// 创建存储邮件对象
	storedMail := &StoredMail{
		MessageID:   generateMessageID(s.backend.domain),
		From:        s.from,
		To:          s.to,
		Subject:     subject,
		Body:        string(body),
		ContentType: msg.Header.Get("Content-Type"),
		Size:        len(body),
		Received:    time.Now(),
		IsRead:      false,
		Folder:      "INBOX",
	}

	// 存储邮件
	if err := s.backend.storage.StoreMail(storedMail); err != nil {
		log.Printf("❌ 存储邮件失败: %v [%s]", err, serverTypeStr)
		return fmt.Errorf("failed to store message: %v", err)
	}

	log.Printf("✅ 邮件存储成功: %s [%s]", storedMail.MessageID, serverTypeStr)
	log.Printf("📧 发件人: %s", s.from)
	log.Printf("📧 收件人: %v", s.to)
	log.Printf("📧 主题: %s", storedMail.Subject)

	return nil
}

// Reset 重置会话状态
func (s *SMTPSession) Reset() {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("🔄 重置SMTP会话状态 [%s]", serverTypeStr)
	s.from = ""
	s.to = []string{}
}

// Logout 处理会话注销
func (s *SMTPSession) Logout() error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("👋 SMTP会话注销 [%s]", serverTypeStr)
	return nil
}

// generateMessageID 生成消息ID
func generateMessageID(domain string) string {
	return fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), domain)
}

// formatHeaders 格式化邮件头
func formatHeaders(header mail.Header) string {
	var headers []string
	for key, values := range header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(headers, "\n")
}
