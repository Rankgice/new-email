package localmail

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

// SimpleIMAPServer 简单IMAP服务器
type SimpleIMAPServer struct {
	port int
}

// NewSimpleIMAPServer 创建简单IMAP服务器
func NewSimpleIMAPServer(port int) *SimpleIMAPServer {
	return &SimpleIMAPServer{
		port: port,
	}
}

// Start 启动IMAP服务器
func (s *SimpleIMAPServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("IMAP服务器监听失败: %v", err)
	}
	defer listener.Close()

	log.Printf("✅ IMAP服务器启动成功，监听端口: %d", s.port)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("IMAP连接接受失败: %v", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理IMAP连接
func (s *SimpleIMAPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 发送欢迎消息
	writer.WriteString("* OK [CAPABILITY IMAP4rev1] localhost IMAP Server Ready\r\n")
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

		log.Printf("IMAP收到: %s", line)

		// 解析命令
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 2 {
			writer.WriteString("* BAD Invalid command\r\n")
			writer.Flush()
			continue
		}

		tag := parts[0]
		command := strings.ToUpper(parts[1])

		switch command {
		case "CAPABILITY":
			writer.WriteString("* CAPABILITY IMAP4rev1\r\n")
			writer.WriteString(fmt.Sprintf("%s OK CAPABILITY completed\r\n", tag))
		case "LOGIN":
			// 简单接受任何登录
			writer.WriteString(fmt.Sprintf("%s OK LOGIN completed\r\n", tag))
		case "LIST":
			writer.WriteString("* LIST () \".\" INBOX\r\n")
			writer.WriteString(fmt.Sprintf("%s OK LIST completed\r\n", tag))
		case "SELECT":
			writer.WriteString("* 0 EXISTS\r\n")
			writer.WriteString("* 0 RECENT\r\n")
			writer.WriteString("* OK [UIDVALIDITY 1] UIDs valid\r\n")
			writer.WriteString("* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")
			writer.WriteString("* OK [PERMANENTFLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)] Limited\r\n")
			writer.WriteString(fmt.Sprintf("%s OK [READ-WRITE] SELECT completed\r\n", tag))
		case "FETCH":
			writer.WriteString(fmt.Sprintf("%s OK FETCH completed\r\n", tag))
		case "SEARCH":
			writer.WriteString("* SEARCH\r\n")
			writer.WriteString(fmt.Sprintf("%s OK SEARCH completed\r\n", tag))
		case "LOGOUT":
			writer.WriteString("* BYE IMAP4rev1 Server logging out\r\n")
			writer.WriteString(fmt.Sprintf("%s OK LOGOUT completed\r\n", tag))
			writer.Flush()
			return
		case "NOOP":
			writer.WriteString(fmt.Sprintf("%s OK NOOP completed\r\n", tag))
		default:
			writer.WriteString(fmt.Sprintf("%s BAD Command not recognized\r\n", tag))
		}

		writer.Flush()
		log.Printf("IMAP发送: 响应已发送")
	}
}
