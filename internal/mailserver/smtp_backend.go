package mailserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

// SMTPBackend 实现 smtp.Backend 接口
type SMTPBackend struct {
	domain  string
	storage *MailStorage
}

// NewSMTPBackend 创建SMTP后端
func NewSMTPBackend(domain string, storage *MailStorage) *SMTPBackend {
	return &SMTPBackend{
		domain:  domain,
		storage: storage,
	}
}

// NewSession 创建新的SMTP会话
func (b *SMTPBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	log.Printf("📧 新SMTP连接来自: %s", c.Conn().RemoteAddr())
	return &SMTPSession{
		backend: b,
		conn:    c,
	}, nil
}

// SMTPServer SMTP服务器
type SMTPServer struct {
	port    int
	domain  string
	storage *MailStorage
	server  *smtp.Server
}

// NewSMTPServer 创建SMTP服务器
func NewSMTPServer(port int, domain string, storage *MailStorage) *SMTPServer {
	backend := NewSMTPBackend(domain, storage)

	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(":%d", port)
	server.Domain = domain
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 10 * 1024 * 1024 // 10MB
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true // 允许非TLS认证（开发环境）

	return &SMTPServer{
		port:    port,
		domain:  domain,
		storage: storage,
		server:  server,
	}
}

// Start 启动SMTP服务器
func (s *SMTPServer) Start(ctx context.Context) error {
	log.Printf("✅ SMTP服务器启动成功，监听端口: %d", s.port)
	log.Printf("🌐 SMTP域名: %s", s.domain)

	// 在goroutine中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("❌ SMTP服务器错误: %v", err)
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	log.Printf("🛑 正在关闭SMTP服务器...")

	// 优雅关闭服务器
	return s.server.Close()
}
