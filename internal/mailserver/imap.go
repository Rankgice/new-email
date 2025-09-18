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
	port        int
	domain      string
	storage     *MailStorage
	server      *imapserver.Server
	listener    net.Listener
	tlsCertPath string
	tlsKeyPath  string
}

// NewIMAPServer 创建IMAP v2服务器
func NewIMAPServer(config Config, storage *MailStorage) *IMAPServer {
	options := &imapserver.Options{
		NewSession: func() imapserver.Session {
			return NewIMAPSession(storage)
		},
	}

	// 配置TLS
	if config.TLSCertPath != "" && config.TLSKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(config.TLSCertPath, config.TLSKeyPath)
		if err != nil {
			log.Fatalf("无法加载TLS证书: %v", err)
		}
		options.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		options.InsecureAuth = false
	} else {
		options.InsecureAuth = true
	}

	server := imapserver.New(options)

	return &IMAPServer{
		port:        config.IMAPPort,
		domain:      config.Domain,
		storage:     storage,
		server:      server,
		tlsCertPath: config.TLSCertPath,
		tlsKeyPath:  config.TLSKeyPath,
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

	useTLS := s.tlsCertPath != "" && s.tlsKeyPath != ""

	if useTLS {
		log.Printf("✅ IMAP服务器 (TLS) 启动成功，监听端口: %d", s.port)
	} else {
		log.Printf("⚠️ IMAP服务器 (非TLS) 启动成功，监听端口: %d", s.port)
	}

	// 在goroutine中启动服务器
	go func() {
		var serveErr error
		if useTLS {
			// 对于TLS，我们需要使用TLS监听器
			cert, err := tls.LoadX509KeyPair(s.tlsCertPath, s.tlsKeyPath)
			if err != nil {
				log.Printf("IMAP服务器TLS证书加载失败: %v", err)
				return
			}
			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
			tlsListener := tls.NewListener(s.listener, tlsConfig)
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
