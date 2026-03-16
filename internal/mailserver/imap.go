package mailserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/emersion/go-imap/v2/imapserver"
)

// IMAPServer 基于go-imap/v2库的IMAP服务器实现
type IMAPServer struct {
	port      int
	domain    string
	storage   *MailStorage
	server    *imapserver.Server
	listener  net.Listener
	useTLS    bool
	tlsConfig *tls.Config
}

// NewIMAPServer 创建IMAP服务器
func NewIMAPServer(config Config, storage *MailStorage) *IMAPServer {
	options := &imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			session := NewIMAPSession(storage)
			greeting := &imapserver.GreetingData{
				PreAuth: false, // 需要认证
			}
			return session, greeting, nil
		},
	}

	tlsConfig, useTLS := loadOptionalTLSConfig("IMAP服务器", config.IMAPUseTLS, config.IMAPTLSCertPath, config.IMAPTLSKeyPath)
	if tlsConfig != nil {
		options.TLSConfig = tlsConfig
		options.InsecureAuth = false
	} else {
		options.InsecureAuth = true
	}

	server := imapserver.New(options)

	return &IMAPServer{
		port:      config.IMAPPort,
		domain:    config.Domain,
		storage:   storage,
		server:    server,
		useTLS:    useTLS,
		tlsConfig: tlsConfig,
	}
}

// Start 启动IMAP服务器
func (s *IMAPServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)

	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("无法监听端口 %d: %v", s.port, err)
	}

	if s.useTLS {
		log.Printf("✅ IMAP服务器 (TLS) 启动成功，监听端口: %d", s.port)
	} else {
		log.Printf("⚠️ IMAP服务器 (非TLS) 启动成功，监听端口: %d", s.port)
	}

	// 在goroutine中启动服务器
	go func() {
		var serveErr error
		if s.useTLS {
			tlsListener := tls.NewListener(s.listener, s.tlsConfig)
			serveErr = s.server.Serve(tlsListener)
		} else {
			serveErr = s.server.Serve(s.listener)
		}

		if serveErr != nil && serveErr != net.ErrClosed {
			log.Printf("IMAP服务器运行错误: %v", serveErr)
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	log.Printf("IMAP服务器收到停止信号")

	// 关闭服务器
	if err := s.Stop(); err != nil {
		log.Printf("关闭IMAP服务器失败: %v", err)
		return err
	}

	log.Printf("✅ IMAP服务器已停止")
	return nil
}

// Stop 停止IMAP服务器
func (s *IMAPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
