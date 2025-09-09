package mailserver

import (
	"context"
	"fmt"
	"log"

	"github.com/emersion/go-imap/server"
)

// IMAPServer 基于go-imap库的IMAP服务器实现
type IMAPServer struct {
	port    int
	domain  string
	storage *MailStorage
	server  *server.Server
}

// NewIMAPServer 创建IMAP服务器
func NewIMAPServer(port int, domain string, storage *MailStorage) *IMAPServer {
	// 创建自定义后端
	backend := NewCustomBackend(storage)

	// 创建IMAP服务器
	s := server.New(backend)
	s.Addr = fmt.Sprintf(":%d", port)
	s.AllowInsecureAuth = true // 允许不安全的认证，仅用于测试

	return &IMAPServer{
		port:    port,
		domain:  domain,
		storage: storage,
		server:  s,
	}
}

// Start 启动IMAP服务器
func (s *IMAPServer) Start(ctx context.Context) error {
	log.Printf("✅ IMAP服务器启动成功，监听端口: %d (使用go-imap库)", s.port)

	// 在goroutine中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
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
