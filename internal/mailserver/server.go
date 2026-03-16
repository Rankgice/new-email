package mailserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Config 邮件服务器配置
type Config struct {
	SMTPReceivePort int    `yaml:"smtp_receive_port"` // 25端口 - 接收外部邮件 (MTA)
	SMTPSubmitPort  int    `yaml:"smtp_submit_port"`  // 587端口 - 用户提交邮件 (MSA)
	IMAPPort        int    `yaml:"imap_port"`         // 993端口 - IMAP访问
	Domain          string `yaml:"domain"`
	DatabasePath    string `yaml:"database_path"`
	SMTPUseTLS      bool   `yaml:"smtp_use_tls"`
	SMTPTLSCertPath string `yaml:"smtp_tls_cert_path"` // SMTP TLS证书路径
	SMTPTLSKeyPath  string `yaml:"smtp_tls_key_path"`  // SMTP TLS密钥路径
	IMAPUseTLS      bool   `yaml:"imap_use_tls"`
	IMAPTLSCertPath string `yaml:"imap_tls_cert_path"` // IMAP TLS证书路径
	IMAPTLSKeyPath  string `yaml:"imap_tls_key_path"`  // IMAP TLS密钥路径
}

// MailServer 邮件服务器
type MailServer struct {
	config            Config
	smtpReceiveServer *SMTPServer // 25端口 - 接收外部邮件
	smtpSubmitServer  *SMTPServer // 587端口 - 用户提交邮件
	imapServer        *IMAPServer
	storage           *MailStorage
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

// NewMailServer 创建邮件服务器
func NewMailServer(config Config, db *gorm.DB) *MailServer {
	ctx, cancel := context.WithCancel(context.Background())

	storage := NewMailStorage(db, config.Domain)

	return &MailServer{
		config:  config,
		storage: storage,
		ctx:     ctx,
		cancel:  cancel,
		// 创建接收服务器 (25端口 - MTA功能)
		smtpReceiveServer: NewSMTPReceiveServer(config.SMTPReceivePort, config.Domain, storage, config.SMTPUseTLS, config.SMTPTLSCertPath, config.SMTPTLSKeyPath),
		// 创建提交服务器 (587端口 - MSA功能)
		smtpSubmitServer: NewSMTPSubmitServer(config.SMTPSubmitPort, config.Domain, storage, config.SMTPUseTLS, config.SMTPTLSCertPath, config.SMTPTLSKeyPath),
		// IMAP服务器
		imapServer: NewIMAPServer(config, storage),
	}
}

// Start 启动邮件服务器
func (s *MailServer) Start() error {
	log.Printf("🚀 启动邮件服务器...")
	log.Printf("📧 SMTP接收服务器 (MTA): localhost:%d - 用于接收外部邮件", s.config.SMTPReceivePort)
	log.Printf("📤 SMTP提交服务器 (MSA): localhost:%d - 用于用户认证提交", s.config.SMTPSubmitPort)
	log.Printf("📬 IMAP服务器: localhost:%d", s.config.IMAPPort)
	log.Printf("🌐 域名: %s", s.config.Domain)
	log.Printf("⚠️  外部邮件应连接到端口%d，用户提交应连接到端口%d", s.config.SMTPReceivePort, s.config.SMTPSubmitPort)

	// 启动SMTP接收服务器 (25端口)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.smtpReceiveServer.Start(s.ctx); err != nil {
			log.Printf("❌ SMTP接收服务器启动失败: %v", err)
		}
	}()

	// 启动SMTP提交服务器 (587端口)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.smtpSubmitServer.Start(s.ctx); err != nil {
			log.Printf("❌ SMTP提交服务器启动失败: %v", err)
		}
	}()

	// 启动IMAP服务器
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.imapServer.Start(s.ctx); err != nil {
			log.Printf("❌ IMAP服务器启动失败: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(200 * time.Millisecond)

	// 测试端口是否可用（可选，可能导致连接冲突）
	// if err := s.testPorts(); err != nil {
	//     return err
	// }

	log.Printf("✅ 邮件服务器启动完成")
	log.Printf("💡 使用说明:")
	log.Printf("   - 外部邮件服务器发送到: localhost:%d (无需认证)", s.config.SMTPReceivePort)
	log.Printf("   - 用户邮件客户端连接: localhost:%d (需要认证)", s.config.SMTPSubmitPort)
	log.Printf("   - IMAP邮件访问: localhost:%d", s.config.IMAPPort)
	return nil
}

// Stop 停止邮件服务器
func (s *MailServer) Stop() error {
	log.Printf("🛑 停止邮件服务器...")

	s.cancel()
	s.wg.Wait()

	log.Printf("✅ 邮件服务器已停止")
	return nil
}

// testPorts 测试端口是否可用
func (s *MailServer) testPorts() error {
	// 测试SMTP接收端口 (25)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.config.SMTPReceivePort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("SMTP接收端口 %d 不可用: %v", s.config.SMTPReceivePort, err)
	}
	conn.Close()

	// 测试SMTP提交端口 (587)
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.config.SMTPSubmitPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("SMTP提交端口 %d 不可用: %v", s.config.SMTPSubmitPort, err)
	}
	conn.Close()

	// 测试IMAP端口
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.config.IMAPPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("IMAP端口 %d 不可用: %v", s.config.IMAPPort, err)
	}
	conn.Close()

	return nil
}
