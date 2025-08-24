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
	backend *SMTPBackend
	conn    *smtp.Conn
	from    string
	to      []string
}

// AuthPlain 处理PLAIN认证
func (s *SMTPSession) AuthPlain(username, password string) error {
	log.Printf("🔐 SMTP认证请求: %s", username)

	// 验证邮箱格式
	if !strings.Contains(username, "@") {
		log.Printf("❌ 认证失败: 无效的邮箱格式 %s", username)
		return fmt.Errorf("invalid email format")
	}

	// 使用存储层验证凭据
	if !s.backend.storage.ValidateCredentials(username, password) {
		log.Printf("❌ 认证失败: 用户名或密码错误 %s", username)
		return fmt.Errorf("invalid credentials")
	}

	log.Printf("✅ 认证成功: %s", username)
	return nil
}

// Mail 处理MAIL FROM命令
func (s *SMTPSession) Mail(from string, opts *smtp.MailOptions) error {
	log.Printf("📤 MAIL FROM: %s", from)

	// 验证发件人地址格式
	if _, err := mail.ParseAddress(from); err != nil {
		log.Printf("❌ 无效的发件人地址: %s, 错误: %v", from, err)
		return fmt.Errorf("invalid sender address: %v", err)
	}

	s.from = from
	s.to = []string{} // 重置收件人列表

	log.Printf("✅ 发件人设置成功: %s", from)
	return nil
}

// Rcpt 处理RCPT TO命令
func (s *SMTPSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	log.Printf("📥 RCPT TO: %s", to)

	// 验证收件人地址格式
	if _, err := mail.ParseAddress(to); err != nil {
		log.Printf("❌ 无效的收件人地址: %s, 错误: %v", to, err)
		return fmt.Errorf("invalid recipient address: %v", err)
	}

	// 检查是否超过最大收件人数量
	if len(s.to) >= 50 {
		log.Printf("❌ 收件人数量超过限制: %d", len(s.to))
		return fmt.Errorf("too many recipients")
	}

	s.to = append(s.to, to)
	log.Printf("✅ 收件人添加成功: %s (总数: %d)", to, len(s.to))
	return nil
}

// Data 处理DATA命令
func (s *SMTPSession) Data(r io.Reader) error {
	log.Printf("📨 开始接收邮件数据...")

	if s.from == "" {
		return fmt.Errorf("no sender specified")
	}

	if len(s.to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// 解析邮件
	msg, err := message.Read(r)
	if err != nil {
		log.Printf("❌ 解析邮件失败: %v", err)
		return fmt.Errorf("failed to parse message: %v", err)
	}

	// 读取邮件正文
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("❌ 读取邮件正文失败: %v", err)
		return fmt.Errorf("failed to read message body: %v", err)
	}
	log.Printf("📊 邮件数据大小: %d 字节", len(body))

	// MIME 头解码器， 解码标题
	decoder := new(mime.WordDecoder)
	subject, err := decoder.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		subject = msg.Header.Get("Subject") // 解码失败就用原文
	}
	// TODO 如果收件人不是自己，并且发件人是用户，应当直接转发

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
		log.Printf("❌ 存储邮件失败: %v", err)
		return fmt.Errorf("failed to store message: %v", err)
	}

	log.Printf("✅ 邮件存储成功: %s", storedMail.MessageID)
	log.Printf("📧 发件人: %s", s.from)
	log.Printf("📧 收件人: %v", s.to)
	log.Printf("📧 主题: %s", storedMail.Subject)

	return nil
}

// Reset 重置会话状态
func (s *SMTPSession) Reset() {
	log.Printf("🔄 重置SMTP会话状态")
	s.from = ""
	s.to = []string{}
}

// Logout 处理会话结束
func (s *SMTPSession) Logout() error {
	log.Printf("👋 SMTP会话结束: %s", s.conn.Conn().RemoteAddr())
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
