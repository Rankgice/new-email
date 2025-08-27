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
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// SMTPSession 实现 smtp.Session 和 smtp.AuthSession 接口
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

// AuthMechanisms 返回支持的认证机制
func (s *SMTPSession) AuthMechanisms() []string {
	// 目前只支持PLAIN认证机制
	mechanisms := []string{"PLAIN"}

	log.Printf("🔐 支持的认证机制: %v", mechanisms)
	return mechanisms
}

// Auth 处理指定的认证机制
func (s *SMTPSession) Auth(mech string) (sasl.Server, error) {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}

	log.Printf("🔐 请求认证机制 [%s]: %s", serverTypeStr, mech)

	switch strings.ToUpper(mech) {
	case "PLAIN":
		// 创建PLAIN认证服务器
		return sasl.NewPlainServer(func(identity, username, password string) error {
			log.Printf("🔐 PLAIN认证请求 [%s]: identity=%s, username=%s", serverTypeStr, identity, username)

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
			log.Printf("✅ PLAIN认证成功: %s [%s]", username, serverTypeStr)
			return nil
		}), nil

	case "LOGIN":
		// LOGIN认证机制 - 暂时不支持，因为go-sasl没有直接的NewLoginServer
		log.Printf("⚠️  LOGIN认证机制暂不支持，请使用PLAIN认证 [%s]", serverTypeStr)
		return nil, fmt.Errorf("LOGIN authentication not supported, please use PLAIN")

	default:
		log.Printf("❌ 不支持的认证机制: %s [%s]", mech, serverTypeStr)
		return nil, fmt.Errorf("unsupported authentication mechanism: %s", mech)
	}
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

	// MSA服务器必须要求认证
	if s.requireAuth && !s.authenticated {
		log.Printf("❌ MSA服务器要求认证，但未认证 [%s]", serverTypeStr)
		return fmt.Errorf("authentication required")
	}

	// 验证收件人地址格式
	if _, err := mail.ParseAddress(to); err != nil {
		log.Printf("❌ 无效的收件人地址: %s, 错误: %v [%s]", to, err, serverTypeStr)
		return fmt.Errorf("invalid recipient address: %v", err)
	}

	// MTA服务器需要检查是否为本地域名
	if s.serverType == SMTPServerTypeReceive {
		// 检查收件人是否为本地域名的邮箱
		if !s.isLocalDomain(to) {
			log.Printf("❌ 收件人不属于本地域名，拒绝接收: %s", to)
			return fmt.Errorf("relay not permitted")
		}
		log.Printf("✅ MTA确认本地域名邮箱: %s", to)
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

// isLocalDomain 检查是否为本地域名
func (s *SMTPSession) isLocalDomain(email string) bool {
	// 提取邮箱的域名部分
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	// 检查是否为服务器域名
	serverDomain := strings.ToLower(s.backend.domain)
	if domain == serverDomain {
		log.Printf("✅ 匹配服务器域名: %s", domain)

		// 进一步检查邮箱是否存在于数据库中
		if s.backend.storage.isMailboxExists(email) {
			log.Printf("✅ 邮箱存在于数据库: %s", email)
			return true
		} else {
			log.Printf("⚠️  域名匹配但邮箱不存在: %s", email)
			// 对于自建邮箱，即使邮箱不存在也应该接收（可以后续创建）
			return true
		}
	}

	// TODO: 这里应该查询数据库中配置的其他本地域名列表
	// 暂时只检查服务器配置的域名
	log.Printf("❌ 域名不匹配: %s vs %s", domain, serverDomain)
	return false
}

// generateMessageID 生成邮件ID
func generateMessageID(domain string) string {
	return fmt.Sprintf("<%d.%d@%s>", time.Now().Unix(), time.Now().Nanosecond(), domain)
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
