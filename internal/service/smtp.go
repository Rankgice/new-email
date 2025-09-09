package service

import (
	"crypto/tls"
	"io"

	"github.com/go-mail/mail/v2"
)

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseTLS   bool   `json:"useTLS"`
}

// SMTPService SMTP服务
type SMTPService struct {
	config SMTPConfig
}

// NewSMTPService 创建SMTP服务
func NewSMTPService(config SMTPConfig) *SMTPService {
	return &SMTPService{
		config: config,
	}
}

// EmailMessage 邮件消息
type EmailMessage struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Cc          []string          `json:"cc"`
	Bcc         []string          `json:"bcc"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	ContentType string            `json:"contentType"` // text/plain 或 text/html
	Attachments []EmailAttachment `json:"attachments"`
}

// EmailAttachment 邮件附件
type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

// SendEmail 发送邮件
func (s *SMTPService) SendEmail(message EmailMessage) error {
	// 创建邮件消息
	m := mail.NewMessage()

	// 设置发件人
	m.SetHeader("From", message.From)

	// 设置收件人
	if len(message.To) > 0 {
		m.SetHeader("To", message.To...)
	}

	// 设置抄送
	if len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}

	// 设置密送
	if len(message.Bcc) > 0 {
		m.SetHeader("Bcc", message.Bcc...)
	}

	// 设置主题
	m.SetHeader("Subject", message.Subject)

	// 设置邮件内容
	contentType := message.ContentType
	if contentType == "" {
		contentType = "text/plain"
	}

	if contentType == "text/html" {
		m.SetBody("text/html", message.Body)
	} else {
		m.SetBody("text/plain", message.Body)
	}

	// 添加附件
	for _, attachment := range message.Attachments {
		m.Attach(attachment.Filename, mail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Data)
			return err
		}))
	}

	// 创建SMTP拨号器
	d := mail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	// 配置TLS
	if s.config.UseTLS {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.config.Host,
		}
	} else {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		d.StartTLSPolicy = mail.NoStartTLS
	}

	// 发送邮件
	return d.DialAndSend(m)
}

// TestConnection 测试SMTP连接
func (s *SMTPService) TestConnection() error {
	// 如果没有配置主机，跳过测试
	if s.config.Host == "" {
		return nil
	}

	// 创建SMTP拨号器
	d := mail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	// 配置TLS
	if s.config.UseTLS {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // 允许跳过证书验证，仅用于开发测试
			ServerName:         s.config.Host,
		}
		d.StartTLSPolicy = mail.MandatoryStartTLS // 明确要求STARTTLS
	} else {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		d.StartTLSPolicy = mail.NoStartTLS
	}

	// 测试连接
	conn, err := d.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// GetSMTPConfig 获取SMTP配置
func (s *SMTPService) GetSMTPConfig() SMTPConfig {
	// 返回配置副本，隐藏密码
	config := s.config
	config.Password = "***" // 隐藏密码
	return config
}
