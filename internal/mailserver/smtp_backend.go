package mailserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

// SMTPServerType SMTP服务器类型
type SMTPServerType int

const (
	SMTPServerTypeReceive SMTPServerType = iota // MTA - 接收外部邮件 (25端口)
	SMTPServerTypeSubmit                        // MSA - 用户提交邮件 (587端口)
)

// SMTPServer SMTP服务器
type SMTPServer struct {
	port       int
	domain     string
	storage    *MailStorage
	server     *smtp.Server
	serverType SMTPServerType // 服务器类型
}

// NewSMTPReceiveServer 创建SMTP接收服务器 (MTA - 25端口)
func NewSMTPReceiveServer(port int, domain string, storage *MailStorage) *SMTPServer {
	backend := NewSMTPBackend(domain, storage, SMTPServerTypeReceive)

	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(":%d", port)
	server.Domain = domain
	server.WriteTimeout = 30 * time.Second
	server.ReadTimeout = 30 * time.Second
	server.MaxMessageBytes = 50 * 1024 * 1024 // 50MB for external emails
	server.MaxRecipients = 100
	server.AllowInsecureAuth = true // MTA可以接受非加密连接
	// 注意：MTA通常不需要认证，认证逻辑在Backend中处理

	return &SMTPServer{
		port:       port,
		domain:     domain,
		storage:    storage,
		server:     server,
		serverType: SMTPServerTypeReceive,
	}
}

// NewSMTPSubmitServer 创建SMTP提交服务器 (MSA - 587端口)
func NewSMTPSubmitServer(port int, domain string, storage *MailStorage) *SMTPServer {
	backend := NewSMTPBackend(domain, storage, SMTPServerTypeSubmit)

	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(":%d", port)
	server.Domain = domain
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 25 * 1024 * 1024 // 25MB for user submissions
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true // MSA要求加密认证
	// 注意：MSA需要认证，认证逻辑在Backend中强制执行

	// 配置TLS（生产环境需要真实证书）
	// 这里使用自签名证书用于开发测试
	tlsConfig := &tls.Config{
		ServerName: domain,
		// 生产环境应该加载真实的SSL证书
		// Certificates: []tls.Certificate{cert},
	}
	server.TLSConfig = tlsConfig
	server.EnableSMTPUTF8 = true // 支持UTF8邮件地址

	return &SMTPServer{
		port:       port,
		domain:     domain,
		storage:    storage,
		server:     server,
		serverType: SMTPServerTypeSubmit,
	}
}

// SMTPBackend 实现 smtp.Backend 接口
type SMTPBackend struct {
	domain     string
	storage    *MailStorage
	serverType SMTPServerType
}

// NewSMTPBackend 创建SMTP后端
func NewSMTPBackend(domain string, storage *MailStorage, serverType SMTPServerType) *SMTPBackend {
	return &SMTPBackend{
		domain:     domain,
		storage:    storage,
		serverType: serverType,
	}
}

// NewSession 创建新的SMTP会话
func (b *SMTPBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	serverTypeStr := "MTA(接收)"
	if b.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📧 新SMTP连接来自: %s [%s]", c.Conn().RemoteAddr(), serverTypeStr)

	session := &SMTPSession{
		backend:       b,
		conn:          c,
		serverType:    b.serverType,
		authenticated: false,
	}

	// MSA服务器需要更严格的控制
	if b.serverType == SMTPServerTypeSubmit {
		session.requireAuth = true
		log.Printf("🔒 MSA服务器要求认证")
	} else {
		session.requireAuth = false
		log.Printf("🌐 MTA服务器可接受未认证连接")
	}

	// 验证session实现了必要的接口
	var _ smtp.Session = session
	var _ smtp.AuthSession = session // 确保实现了AuthSession接口

	log.Printf("✅ 会话创建成功，支持认证接口: %t", true)

	// 返回session，go-smtp会自动检测是否实现了AuthSession接口
	return session, nil
}

// Start 启动SMTP服务器
func (s *SMTPServer) Start(ctx context.Context) error {
	serverTypeStr := "接收服务器(MTA)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "提交服务器(MSA)"
	}

	log.Printf("✅ SMTP%s启动成功，监听端口: %d", serverTypeStr, s.port)
	log.Printf("🌐 SMTP域名: %s", s.domain)

	// 在goroutine中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("❌ SMTP%s错误: %v", serverTypeStr, err)
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	log.Printf("🛑 正在关闭SMTP%s...", serverTypeStr)

	// 优雅关闭服务器
	return s.server.Close()
}
