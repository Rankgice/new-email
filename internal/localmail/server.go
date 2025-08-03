package localmail

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// LocalMailServer 本地邮件服务器
type LocalMailServer struct {
	smtpPort   int
	imapPort   int
	smtpServer *SimpleSMTPServer
	imapServer *SimpleIMAPServer
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewLocalMailServer 创建本地邮件服务器
func NewLocalMailServer(smtpPort, imapPort int) *LocalMailServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &LocalMailServer{
		smtpPort:   smtpPort,
		imapPort:   imapPort,
		smtpServer: NewSimpleSMTPServer(smtpPort),
		imapServer: NewSimpleIMAPServer(imapPort),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动本地邮件服务器
func (s *LocalMailServer) Start() error {
	log.Printf("🚀 启动本地邮件服务器...")

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

	log.Printf("✅ 本地邮件服务器启动成功")
	log.Printf("📧 SMTP服务器: localhost:%d", s.smtpPort)
	log.Printf("📬 IMAP服务器: localhost:%d", s.imapPort)

	return nil
}

// Stop 停止本地邮件服务器
func (s *LocalMailServer) Stop() error {
	log.Printf("🛑 停止本地邮件服务器...")

	s.cancel()
	s.wg.Wait()

	log.Printf("✅ 本地邮件服务器已停止")
	return nil
}

// testPorts 测试端口是否可用
func (s *LocalMailServer) testPorts() error {
	// 测试SMTP端口
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.smtpPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("SMTP端口 %d 不可用: %v", s.smtpPort, err)
	}
	conn.Close()

	// 测试IMAP端口
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.imapPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("IMAP端口 %d 不可用: %v", s.imapPort, err)
	}
	conn.Close()

	return nil
}
