package mailserver

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// IMAPServer IMAP服务器
type IMAPServer struct {
	port    int
	domain  string
	storage *MailStorage
}

// IMAPSession IMAP会话
type IMAPSession struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	domain        string
	storage       *MailStorage
	authenticated bool
	username      string
	selectedBox   string
	state         string
}

// NewIMAPServer 创建IMAP服务器
func NewIMAPServer(port int, domain string, storage *MailStorage) *IMAPServer {
	return &IMAPServer{
		port:    port,
		domain:  domain,
		storage: storage,
	}
}

// Start 启动IMAP服务器
func (s *IMAPServer) Start(ctx context.Context) error {
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
func (s *IMAPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	session := &IMAPSession{
		conn:          conn,
		reader:        bufio.NewReader(conn),
		writer:        bufio.NewWriter(conn),
		domain:        s.domain,
		storage:       s.storage,
		authenticated: false,
		state:         "NOT_AUTHENTICATED",
	}

	// 发送欢迎消息
	session.writeLine(fmt.Sprintf("* OK [CAPABILITY IMAP4rev1] %s IMAP Server Ready", s.domain))

	for {
		line, err := session.reader.ReadString('\n')
		if err != nil {
			log.Printf("IMAP读取失败: %v", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log.Printf("IMAP收到: %s", line)

		if err := session.handleCommand(line); err != nil {
			log.Printf("IMAP命令处理失败: %v", err)
			break
		}

		if session.state == "LOGOUT" {
			break
		}
	}
}

// handleCommand 处理IMAP命令
func (s *IMAPSession) handleCommand(line string) error {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		s.writeLine("* BAD Invalid command")
		return nil
	}

	tag := parts[0]
	command := strings.ToUpper(parts[1])
	args := ""
	if len(parts) > 2 {
		args = parts[2]
	}

	switch command {
	case "CAPABILITY":
		return s.handleCapability(tag)
	case "LOGIN":
		return s.handleLogin(tag, args)
	case "LIST":
		return s.handleList(tag, args)
	case "SELECT":
		return s.handleSelect(tag, args)
	case "FETCH":
		return s.handleFetch(tag, args)
	case "SEARCH":
		return s.handleSearch(tag, args)
	case "STORE":
		return s.handleStore(tag, args)
	case "LOGOUT":
		return s.handleLogout(tag)
	case "NOOP":
		return s.handleNoop(tag)
	default:
		s.writeLine(fmt.Sprintf("%s BAD Command not recognized", tag))
		return nil
	}
}

// handleCapability 处理CAPABILITY命令
func (s *IMAPSession) handleCapability(tag string) error {
	s.writeLine("* CAPABILITY IMAP4rev1")
	s.writeLine(fmt.Sprintf("%s OK CAPABILITY completed", tag))
	return nil
}

// handleLogin 处理LOGIN命令
func (s *IMAPSession) handleLogin(tag, args string) error {
	if s.authenticated {
		s.writeLine(fmt.Sprintf("%s BAD Already authenticated", tag))
		return nil
	}

	// 解析用户名和密码
	parts := strings.Fields(args)
	if len(parts) != 2 {
		s.writeLine(fmt.Sprintf("%s BAD Invalid arguments", tag))
		return nil
	}

	username := strings.Trim(parts[0], "\"")
	password := strings.Trim(parts[1], "\"")

	// 验证凭据
	if s.storage.ValidateCredentials(username, password) {
		s.authenticated = true
		s.username = username
		s.state = "AUTHENTICATED"
		s.writeLine(fmt.Sprintf("%s OK LOGIN completed", tag))
	} else {
		s.writeLine(fmt.Sprintf("%s NO LOGIN failed", tag))
	}

	return nil
}

// handleList 处理LIST命令
func (s *IMAPSession) handleList(tag, args string) error {
	if !s.authenticated {
		s.writeLine(fmt.Sprintf("%s NO Not authenticated", tag))
		return nil
	}

	// 简单实现，只返回INBOX
	s.writeLine("* LIST () \"/\" \"INBOX\"")
	s.writeLine(fmt.Sprintf("%s OK LIST completed", tag))
	return nil
}

// handleSelect 处理SELECT命令
func (s *IMAPSession) handleSelect(tag, args string) error {
	if !s.authenticated {
		s.writeLine(fmt.Sprintf("%s NO Not authenticated", tag))
		return nil
	}

	mailbox := strings.Trim(args, "\"")
	if mailbox != "INBOX" {
		s.writeLine(fmt.Sprintf("%s NO Mailbox does not exist", tag))
		return nil
	}

	// 获取邮件数量
	mails, err := s.storage.GetMails(s.username, mailbox, 1000)
	if err != nil {
		s.writeLine(fmt.Sprintf("%s NO Select failed", tag))
		return nil
	}

	s.selectedBox = mailbox
	s.state = "SELECTED"

	s.writeLine(fmt.Sprintf("* %d EXISTS", len(mails)))
	s.writeLine("* 0 RECENT")
	s.writeLine("* OK [UIDVALIDITY 1] UIDs valid")
	s.writeLine("* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)")
	s.writeLine("* OK [PERMANENTFLAGS (\\Deleted \\Seen \\*)] Limited")
	s.writeLine(fmt.Sprintf("%s OK [READ-WRITE] SELECT completed", tag))

	return nil
}

// handleFetch 处理FETCH命令
func (s *IMAPSession) handleFetch(tag, args string) error {
	if s.state != "SELECTED" {
		s.writeLine(fmt.Sprintf("%s NO No mailbox selected", tag))
		return nil
	}

	// 简单解析FETCH参数
	parts := strings.SplitN(args, " ", 2)
	if len(parts) != 2 {
		s.writeLine(fmt.Sprintf("%s BAD Invalid arguments", tag))
		return nil
	}

	seqRange := parts[0]
	items := parts[1]

	// 获取邮件列表
	mails, err := s.storage.GetMails(s.username, s.selectedBox, 100)
	if err != nil {
		s.writeLine(fmt.Sprintf("%s NO Fetch failed", tag))
		return nil
	}

	// 解析序列号范围
	start, end := 1, len(mails)
	if seqRange != "*" {
		if strings.Contains(seqRange, ":") {
			rangeParts := strings.Split(seqRange, ":")
			if len(rangeParts) == 2 {
				if s, err := strconv.Atoi(rangeParts[0]); err == nil {
					start = s
				}
				if e, err := strconv.Atoi(rangeParts[1]); err == nil {
					end = e
				}
			}
		} else {
			if s, err := strconv.Atoi(seqRange); err == nil {
				start = s
				end = s
			}
		}
	}

	// 返回邮件信息
	for i := start; i <= end && i <= len(mails); i++ {
		if i < 1 {
			continue
		}
		mail := mails[i-1]

		if strings.Contains(strings.ToUpper(items), "ENVELOPE") {
			s.writeLine(fmt.Sprintf("* %d FETCH (ENVELOPE (\"%s\" \"%s\" ((\"%s\" NIL \"%s\" \"%s\")) ((\"%s\" NIL \"%s\" \"%s\")) NIL NIL NIL \"%s\"))",
				i, mail.Received.Format("02-Jan-2006 15:04:05 -0700"), mail.Subject,
				mail.From, mail.From, mail.From,
				mail.To[0], mail.To[0], mail.To[0],
				mail.MessageID))
		}

		if strings.Contains(strings.ToUpper(items), "FLAGS") {
			flags := ""
			if mail.IsRead {
				flags = "\\Seen"
			}
			s.writeLine(fmt.Sprintf("* %d FETCH (FLAGS (%s))", i, flags))
		}

		if strings.Contains(strings.ToUpper(items), "BODY") {
			s.writeLine(fmt.Sprintf("* %d FETCH (BODY[TEXT] {%d}", i, len(mail.Body)))
			s.writeLine(mail.Body)
			s.writeLine(")")
		}
	}

	s.writeLine(fmt.Sprintf("%s OK FETCH completed", tag))
	return nil
}

// handleSearch 处理SEARCH命令
func (s *IMAPSession) handleSearch(tag, args string) error {
	if s.state != "SELECTED" {
		s.writeLine(fmt.Sprintf("%s NO No mailbox selected", tag))
		return nil
	}

	// 简单实现，返回所有邮件
	mails, err := s.storage.GetMails(s.username, s.selectedBox, 100)
	if err != nil {
		s.writeLine(fmt.Sprintf("%s NO Search failed", tag))
		return nil
	}

	result := "* SEARCH"
	for i := range mails {
		result += fmt.Sprintf(" %d", i+1)
	}
	s.writeLine(result)
	s.writeLine(fmt.Sprintf("%s OK SEARCH completed", tag))

	return nil
}

// handleStore 处理STORE命令
func (s *IMAPSession) handleStore(tag, args string) error {
	if s.state != "SELECTED" {
		s.writeLine(fmt.Sprintf("%s NO No mailbox selected", tag))
		return nil
	}

	s.writeLine(fmt.Sprintf("%s OK STORE completed", tag))
	return nil
}

// handleLogout 处理LOGOUT命令
func (s *IMAPSession) handleLogout(tag string) error {
	s.writeLine("* BYE IMAP4rev1 Server logging out")
	s.writeLine(fmt.Sprintf("%s OK LOGOUT completed", tag))
	s.state = "LOGOUT"
	return nil
}

// handleNoop 处理NOOP命令
func (s *IMAPSession) handleNoop(tag string) error {
	s.writeLine(fmt.Sprintf("%s OK NOOP completed", tag))
	return nil
}

// writeLine 写入一行数据
func (s *IMAPSession) writeLine(line string) {
	s.writer.WriteString(line + "\r\n")
	s.writer.Flush()
	log.Printf("IMAP发送: %s", line)
}
