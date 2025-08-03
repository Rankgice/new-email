package localmail

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

// SimpleSMTPServer 简单SMTP服务器
type SimpleSMTPServer struct {
	port int
}

// NewSimpleSMTPServer 创建简单SMTP服务器
func NewSimpleSMTPServer(port int) *SimpleSMTPServer {
	return &SimpleSMTPServer{
		port: port,
	}
}

// Start 启动SMTP服务器
func (s *SimpleSMTPServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("SMTP服务器监听失败: %v", err)
	}
	defer listener.Close()

	log.Printf("✅ SMTP服务器启动成功，监听端口: %d", s.port)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("SMTP连接接受失败: %v", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理SMTP连接
func (s *SimpleSMTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 发送欢迎消息
	writer.WriteString("220 localhost SMTP Server Ready\r\n")
	writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log.Printf("SMTP收到: %s", line)

		// 简单的命令处理
		parts := strings.SplitN(line, " ", 2)
		command := strings.ToUpper(parts[0])

		switch command {
		case "HELO", "EHLO":
			writer.WriteString("250 localhost Hello\r\n")
		case "MAIL":
			writer.WriteString("250 OK\r\n")
		case "RCPT":
			writer.WriteString("250 OK\r\n")
		case "DATA":
			writer.WriteString("354 Start mail input; end with <CRLF>.<CRLF>\r\n")
			// 读取邮件数据直到遇到单独的点
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				dataLine = strings.TrimSpace(dataLine)
				if dataLine == "." {
					break
				}
			}
			writer.WriteString("250 OK Message accepted\r\n")
		case "QUIT":
			writer.WriteString("221 localhost closing connection\r\n")
			writer.Flush()
			return
		case "RSET":
			writer.WriteString("250 OK\r\n")
		case "NOOP":
			writer.WriteString("250 OK\r\n")
		default:
			writer.WriteString("500 Command not recognized\r\n")
		}

		writer.Flush()
		log.Printf("SMTP发送: 响应已发送")
	}
}
