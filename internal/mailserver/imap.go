package mailserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/emersion/go-imap/server"
)

// IMAPServer 基于go-imap库的IMAP服务器实现
type IMAPServer struct {
	port        int
	domain      string
	storage     *MailStorage
	server      *server.Server
	tlsCertPath string
	tlsKeyPath  string
}

// NewIMAPServer 创建IMAP服务器
func NewIMAPServer(config Config, storage *MailStorage) *IMAPServer {
	backend := NewCustomBackend(storage)
	s := server.New(backend)
	s.Addr = fmt.Sprintf(":%d", config.IMAPPort)

	// 配置TLS
	if config.TLSCertPath != "" && config.TLSKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(config.TLSCertPath, config.TLSKeyPath)
		if err != nil {
			log.Fatalf("无法加载TLS证书: %v", err)
		}
		s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		s.AllowInsecureAuth = false
	} else {
		s.AllowInsecureAuth = true
	}

	return &IMAPServer{
		port:        config.IMAPPort,
		domain:      config.Domain,
		storage:     storage,
		server:      s,
		tlsCertPath: config.TLSCertPath,
		tlsKeyPath:  config.TLSKeyPath,
	}
}

// Start 启动IMAP服务器
func (s *IMAPServer) Start(ctx context.Context) error {
	useTLS := s.tlsCertPath != "" && s.tlsKeyPath != ""

	if useTLS {
		log.Printf("✅ IMAP服务器 (TLS) 启动成功，监听端口: %d", s.port)
	} else {
		log.Printf("⚠️ IMAP服务器 (非TLS) 启动成功，监听端口: %d", s.port)
	}

	// 在goroutine中启动服务器
	go func() {
		var err error
		if useTLS {
			err = s.server.ListenAndServeTLS()
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil {
			log.Printf("IMAP服务器运行错误: %v", err)
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	log.Printf("IMAP服务器收到停止信号")

	// 关闭服务器
	if err := s.server.Close(); err != nil {
		log.Printf("关闭IMAP服务器失败: %v", err)
		return err
	}

	log.Printf("✅ IMAP服务器已停止")
	return nil
}

// Stop 停止IMAP服务器
func (s *IMAPServer) Stop() error {
	return s.server.Close()
}
