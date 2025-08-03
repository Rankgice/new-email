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

// MailServer 邮件服务器
type MailServer struct {
	config     Config
	smtpServer *SMTPServer
	imapServer *IMAPServer
	storage    *MailStorage
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// Config 邮件服务器配置
type Config struct {
	SMTPPort     int    `yaml:"smtp_port"`
	IMAPPort     int    `yaml:"imap_port"`
	Domain       string `yaml:"domain"`
	DatabasePath string `yaml:"database_path"`
}

// NewMailServer 创建邮件服务器
func NewMailServer(config Config, db *gorm.DB) *MailServer {
	ctx, cancel := context.WithCancel(context.Background())

	storage := NewMailStorage(db)

	return &MailServer{
		config:     config,
		storage:    storage,
		ctx:        ctx,
		cancel:     cancel,
		smtpServer: NewSMTPServer(config.SMTPPort, config.Domain, storage),
		imapServer: NewIMAPServer(config.IMAPPort, config.Domain, storage),
	}
}

// Start 启动邮件服务器
func (s *MailServer) Start() error {
	log.Printf("🚀 启动邮件服务器...")
	log.Printf("📧 SMTP服务器: localhost:%d", s.config.SMTPPort)
	log.Printf("📬 IMAP服务器: localhost:%d", s.config.IMAPPort)
	log.Printf("🌐 域名: %s", s.config.Domain)

	// 启动SMTP服务器
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.smtpServer.Start(s.ctx); err != nil {
			log.Printf("❌ SMTP服务器启动失败: %v", err)
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
	time.Sleep(100 * time.Millisecond)

	// 测试端口是否可用
	if err := s.testPorts(); err != nil {
		return err
	}

	log.Printf("✅ 邮件服务器启动成功")
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
	// 测试SMTP端口
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.config.SMTPPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("SMTP端口 %d 不可用: %v", s.config.SMTPPort, err)
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
