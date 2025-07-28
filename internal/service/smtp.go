package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
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
	// 构建邮件地址
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// 设置认证
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// 构建收件人列表
	recipients := make([]string, 0)
	recipients = append(recipients, message.To...)
	recipients = append(recipients, message.Cc...)
	recipients = append(recipients, message.Bcc...)

	// 构建邮件内容
	emailBody := s.buildEmailBody(message)

	// 发送邮件
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, message.From, recipients, []byte(emailBody))
	} else {
		return smtp.SendMail(addr, auth, message.From, recipients, []byte(emailBody))
	}
}

// sendWithTLS 使用TLS发送邮件
func (s *SMTPService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建TLS连接
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return err
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(msg)
	return err
}

// buildEmailBody 构建邮件内容
func (s *SMTPService) buildEmailBody(message EmailMessage) string {
	var body strings.Builder

	// 邮件头
	body.WriteString(fmt.Sprintf("From: %s\r\n", message.From))
	body.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))

	if len(message.Cc) > 0 {
		body.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.Cc, ", ")))
	}

	body.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	body.WriteString("MIME-Version: 1.0\r\n")

	// 内容类型
	if message.ContentType == "" {
		message.ContentType = "text/plain"
	}
	body.WriteString(fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", message.ContentType))
	body.WriteString("\r\n")

	// 邮件正文
	body.WriteString(message.Body)

	return body.String()
}

// TestConnection 测试SMTP连接
func (s *SMTPService) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	if s.config.UseTLS {
		// 测试TLS连接
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.config.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.config.Host)
		if err != nil {
			return err
		}
		defer client.Quit()

		// 测试认证
		if auth != nil {
			return client.Auth(auth)
		}
		return nil
	} else {
		// 测试普通连接
		client, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer client.Quit()

		// 测试认证
		if auth != nil {
			return client.Auth(auth)
		}
		return nil
	}
}

// GetSMTPConfig 获取SMTP配置
func (s *SMTPService) GetSMTPConfig() SMTPConfig {
	// 返回配置副本，隐藏密码
	config := s.config
	config.Password = "***"
	return config
}
