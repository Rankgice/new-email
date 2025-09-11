package mailserver

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/emersion/go-message"
	"github.com/emersion/go-sasl"
	gosmtp "github.com/emersion/go-smtp"
)

// SMTPSession 实现 smtp.Session 和 smtp.AuthSession 接口
type SMTPSession struct {
	backend       *SMTPBackend
	conn          *gosmtp.Conn
	from          string
	to            []string
	serverType    SMTPServerType // 服务器类型
	authenticated bool           // 认证状态
	requireAuth   bool           // 是否要求认证
	authUser      string         // 已认证的用户
}

// AuthMechanisms 返回支持的认证机制
func (s *SMTPSession) AuthMechanisms() []string {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}

	// 目前只支持PLAIN认证机制
	mechanisms := []string{"PLAIN", "LOGIN"}

	log.Printf("🔐 AuthMechanisms被调用 [%s]: 返回支持的认证机制 %v", serverTypeStr, mechanisms)
	return mechanisms
}

// Auth 处理指定的认证机制
func (s *SMTPSession) Auth(mech string) (sasl.Server, error) {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}

	log.Printf("🔐 Auth方法被调用 [%s]: 请求认证机制 %s", serverTypeStr, mech)

	switch strings.ToUpper(mech) {
	case "PLAIN":
		// 创建PLAIN认证服务器
		return sasl.NewPlainServer(func(identity, username, password string) error {
			log.Printf("🔐 PLAIN认证请求 [%s]: identity=%s, username=%s", serverTypeStr, identity, username)

			// 验证邮箱格式
			if !strings.Contains(username, "@") {
				log.Printf("❌ 认证失败: 无效的邮箱格式 %s [%s]", username, serverTypeStr)
				return fmt.Errorf("invalid email format")
			}

			// 使用存储层验证凭据
			if !s.backend.storage.ValidatePassword(username, password) {
				log.Printf("❌ 认证失败: 用户名或密码错误 %s [%s]", username, serverTypeStr)
				return fmt.Errorf("invalid credentials")
			}

			// 认证成功
			s.authenticated = true
			s.authUser = username
			log.Printf("✅ PLAIN认证成功: %s [%s]", username, serverTypeStr)
			return nil
		}), nil

	case "LOGIN":
		// LOGIN认证机制 - 暂时不支持，因为go-sasl没有直接的NewLoginServer
		log.Printf("⚠️  LOGIN认证机制暂不支持，请使用PLAIN认证 [%s]", serverTypeStr)
		return nil, fmt.Errorf("LOGIN authentication not supported, please use PLAIN")

	default:
		log.Printf("❌ 不支持的认证机制: %s [%s]", mech, serverTypeStr)
		return nil, fmt.Errorf("unsupported authentication mechanism: %s", mech)
	}
}

// Mail 处理MAIL FROM命令
func (s *SMTPSession) Mail(from string, opts *gosmtp.MailOptions) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📤 MAIL FROM: %s [%s]", from, serverTypeStr)

	// MSA服务器必须要求认证
	if s.requireAuth && !s.authenticated {
		log.Printf("❌ MSA服务器要求认证，但未认证 [%s]", serverTypeStr)
		return fmt.Errorf("authentication required")
	}

	// 验证发件人地址格式
	if _, err := mail.ParseAddress(from); err != nil {
		log.Printf("❌ 无效的发件人地址: %s, 错误: %v [%s]", from, err, serverTypeStr)
		return fmt.Errorf("invalid sender address: %v", err)
	}

	// MSA服务器需要验证发件人权限
	if s.serverType == SMTPServerTypeSubmit && s.authenticated {
		// TODO: 验证认证用户是否有权限使用此发件人地址
		// 这里应该检查认证用户是否匹配发件人地址或有权限代发
		log.Printf("🔍 验证发件人权限: %s (认证用户: %s)", from, s.authUser)
	}

	s.from = from
	s.to = []string{} // 重置收件人列表

	log.Printf("✅ 发件人设置成功: %s [%s]", from, serverTypeStr)
	return nil
}

// Rcpt 处理RCPT TO命令
func (s *SMTPSession) Rcpt(to string, opts *gosmtp.RcptOptions) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📥 RCPT TO: %s [%s]", to, serverTypeStr)

	// MSA服务器必须要求认证
	if s.requireAuth && !s.authenticated {
		log.Printf("❌ MSA服务器要求认证，但未认证 [%s]", serverTypeStr)
		return fmt.Errorf("authentication required")
	}

	// 验证收件人地址格式
	if _, err := mail.ParseAddress(to); err != nil {
		log.Printf("❌ 无效的收件人地址: %s, 错误: %v [%s]", to, err, serverTypeStr)
		return fmt.Errorf("invalid recipient address: %v", err)
	}

	// MTA服务器需要检查是否为本地域名
	if s.serverType == SMTPServerTypeReceive {
		// 检查收件人是否为本地域名的邮箱
		if !s.isLocalDomain(to) {
			log.Printf("❌ 收件人不属于本地域名，拒绝接收: %s", to)
			return fmt.Errorf("relay not permitted")
		}
		log.Printf("✅ MTA确认本地域名邮箱: %s", to)
	}

	// 检查是否超过最大收件人数量
	maxRecipients := 50
	if s.serverType == SMTPServerTypeReceive {
		maxRecipients = 100 // MTA可以接受更多收件人
	}

	if len(s.to) >= maxRecipients {
		log.Printf("❌ 收件人数量超过限制: %d [%s]", len(s.to), serverTypeStr)
		return fmt.Errorf("too many recipients")
	}

	s.to = append(s.to, to)
	log.Printf("✅ 收件人添加成功: %s (总数: %d) [%s]", to, len(s.to), serverTypeStr)
	return nil
}

// Data 处理DATA命令
func (s *SMTPSession) Data(r io.Reader) error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("📨 开始接收邮件数据... [%s]", serverTypeStr)

	if s.from == "" {
		return fmt.Errorf("no sender specified")
	}

	if len(s.to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// 解析邮件
	msg, err := message.Read(r)
	if err != nil {
		log.Printf("❌ 解析邮件失败: %v [%s]", err, serverTypeStr)
		return fmt.Errorf("failed to parse message: %v", err)
	}

	// 读取邮件正文
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("❌ 读取邮件正文失败: %v [%s]", err, serverTypeStr)
		return fmt.Errorf("failed to read message body: %v", err)
	}
	log.Printf("📊 邮件数据大小: %d 字节 [%s]", len(body), serverTypeStr)

	// MIME 头解码器， 解码标题
	decoder := new(mime.WordDecoder)
	subject, err := decoder.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		subject = msg.Header.Get("Subject") // 解码失败就用原文
	}

	// 根据服务器类型进行不同处理
	if s.serverType == SMTPServerTypeSubmit {
		// MSA: 用户提交的邮件，需要处理转发逻辑
		log.Printf("📤 处理用户提交邮件: %s", subject)

		// TODO: 实现邮件转发逻辑
		// 分离本地和外部收件人
		localRecipients := []string{}
		externalRecipients := []string{}

		for _, recipient := range s.to {
			//暂时统一当做外部邮件处理
			//if s.isLocalDomain(recipient) {
			//	localRecipients = append(localRecipients, recipient)
			//} else {
			externalRecipients = append(externalRecipients, recipient)
			//}
		}

		log.Printf("📬 本地收件人: %v", localRecipients)
		log.Printf("🌐 外部收件人: %v", externalRecipients)

		// 处理外部收件人 - 转发到外部邮件服务器
		if len(externalRecipients) > 0 {
			log.Printf("🚀 开始转发邮件到外部服务器，收件人: %v", externalRecipients)
			if err := s.relayToExternal(s.from, externalRecipients, msg, body); err != nil {
				log.Printf("❌ 外部邮件转发失败: %v [%s]", err, serverTypeStr)
				// 根据策略决定是否返回错误
				// 选项1: 返回错误，整个邮件发送失败
				// 选项2: 只记录日志，本地邮件仍然成功
				// 这里我们选择返回错误，确保用户知道外部邮件发送失败
				return fmt.Errorf("failed to relay external message: %v", err)
			}
			log.Printf("✅ 外部邮件转发成功，收件人: %v", externalRecipients)
		}

		// 处理本地收件人 - 存储到本地邮箱
		if len(localRecipients) > 0 {
			localMail := &StoredMail{
				MessageID:   generateMessageID(s.backend.domain),
				From:        s.from,
				To:          localRecipients, // 只存储本地收件人
				Subject:     subject,
				Body:        string(body),
				ContentType: msg.Header.Get("Content-Type"),
				Size:        len(body),
				Received:    time.Now(),
				IsRead:      false,
				Folder:      "Sent", // 用户提交的邮件应存储在Sent文件夹
			}

			if err := s.backend.storage.StoreMail(localMail); err != nil {
				log.Printf("❌ 存储本地邮件失败: %v [%s]", err, serverTypeStr)
				return fmt.Errorf("failed to store local message: %v", err)
			}
			log.Printf("✅ 本地邮件存储成功: %s", localMail.MessageID)
		}

		log.Printf("✅ 邮件处理完成 [%s] - 本地:%d, 外部:%d", serverTypeStr, len(localRecipients), len(externalRecipients))

		return nil

	} else {
		// MTA: 接收的外部邮件，需要进行垃圾邮件检查
		log.Printf("📥 处理接收邮件: %s", subject)
		// TODO: 垃圾邮件检查、病毒扫描等

		// 创建存储邮件对象
		storedMail := &StoredMail{
			MessageID:   generateMessageID(s.backend.domain),
			From:        s.from,
			To:          s.to,
			Subject:     subject,
			Body:        string(body),
			ContentType: msg.Header.Get("Content-Type"),
			Size:        len(body),
			Received:    time.Now(),
			IsRead:      false,
			Folder:      "INBOX", // 接收的外部邮件应存储在INBOX文件夹
		}

		// 存储邮件
		if err := s.backend.storage.StoreMail(storedMail); err != nil {
			log.Printf("❌ 存储邮件失败: %v [%s]", err, serverTypeStr)
			return fmt.Errorf("failed to store message: %v", err)
		}

		log.Printf("✅ 邮件存储成功: %s [%s]", storedMail.MessageID, serverTypeStr)
		log.Printf("📧 发件人: %s", s.from)
		log.Printf("📧 收件人: %v", s.to)
		log.Printf("📧 主题: %s", storedMail.Subject)

		return nil
	}
}

// Reset 重置会话状态
func (s *SMTPSession) Reset() {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("🔄 重置SMTP会话状态 [%s]", serverTypeStr)
	s.from = ""
	s.to = []string{}
}

// Logout 处理会话注销
func (s *SMTPSession) Logout() error {
	serverTypeStr := "MTA(接收)"
	if s.serverType == SMTPServerTypeSubmit {
		serverTypeStr = "MSA(提交)"
	}
	log.Printf("👋 SMTP会话注销 [%s]", serverTypeStr)
	return nil
}

// isLocalDomain 检查是否为本地域名
func (s *SMTPSession) isLocalDomain(email string) bool {
	// 提取邮箱的域名部分
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	// 检查是否为服务器域名，查询数据库配置的域名列表
	_, err := s.backend.storage.domainModel.GetByName(domain)
	if err != nil {
		log.Printf("❌ 域名不匹配: %s ,err: %v", domain, err)
		return false
	}
	// 进一步检查邮箱是否存在于数据库中
	if s.backend.storage.isMailboxExists(email) {
		log.Printf("✅ 邮箱存在于数据库: %s", email)
		return true
	} else {
		log.Printf("⚠️  域名匹配但邮箱不存在: %s", email)
		// 对于自建邮箱，即使邮箱不存在也应该接收（可以后续创建）
		return true
	}

}

// generateMessageID 生成邮件ID
func generateMessageID(domain string) string {
	return fmt.Sprintf("<%d.%d@%s>", time.Now().Unix(), time.Now().Nanosecond(), domain)
}

// formatHeaders 格式化邮件头
func formatHeaders(header mail.Header) string {
	var headers []string
	for key, values := range header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(headers, "\n")
}

// relayToExternal 转发邮件到外部邮件服务器
func (s *SMTPSession) relayToExternal(from string, recipients []string, originalMsg *message.Entity, body []byte) error {
	log.Printf("🚀 开始转发邮件到外部服务器...")
	log.Printf("   发件人: %s", from)
	log.Printf("   收件人: %v", recipients)

	// 按域名分组收件人，为每个域名单独转发
	domainGroups := make(map[string][]string)
	for _, recipient := range recipients {
		parts := strings.Split(recipient, "@")
		if len(parts) != 2 {
			log.Printf("⚠️  跳过无效收件人地址: %s", recipient)
			continue
		}
		domain := strings.ToLower(parts[1])
		domainGroups[domain] = append(domainGroups[domain], recipient)
	}

	// 为每个域名组转发邮件
	for domain, domainRecipients := range domainGroups {
		if err := s.relayToDomain(domain, from, domainRecipients, originalMsg, body); err != nil {
			log.Printf("❌ 转发到域名 %s 失败: %v", domain, err)
			return fmt.Errorf("failed to relay to domain %s: %v", domain, err)
		}
		log.Printf("✅ 成功转发到域名: %s", domain)
	}

	return nil
}

// relayToDomain 转发邮件到指定域名的邮件服务器
func (s *SMTPSession) relayToDomain(domain string, from string, recipients []string, originalMsg *message.Entity, body []byte) error {
	// 查找域名的MX记录
	mxHost, err := s.lookupMX(domain)
	if err != nil {
		return fmt.Errorf("MX lookup failed for %s: %v", domain, err)
	}

	log.Printf("🌐 连接到 %s 的邮件服务器: %s", domain, mxHost)

	// 使用go-smtp客户端连接到外部SMTP服务器
	addr := net.JoinHostPort(mxHost, "25")

	// 首先尝试普通连接
	client, err := gosmtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", mxHost, err)
	}
	defer client.Close()

	// 发送EHLO
	if err := client.Hello(s.backend.domain); err != nil {
		return fmt.Errorf("EHLO failed: %v", err)
	}

	// 检查是否支持STARTTLS，如果支持则重新连接使用TLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		log.Printf("✅ 服务器支持STARTTLS，重新连接使用TLS")
		client.Close() // 关闭当前连接

		// 使用STARTTLS重新连接
		tlsConfig := &tls.Config{
			ServerName: mxHost,
		}
		client, err = gosmtp.DialStartTLS(addr, tlsConfig)
		if err != nil {
			log.Printf("⚠️  STARTTLS连接失败，尝试普通连接: %v", err)
			// 如果STARTTLS失败，回退到普通连接
			client, err = gosmtp.Dial(addr)
			if err != nil {
				return fmt.Errorf("failed to connect to %s: %v", mxHost, err)
			}
		} else {
			log.Printf("✅ STARTTLS连接成功，使用加密连接")
		}

		// 重新发送EHLO
		if err := client.Hello(s.backend.domain); err != nil {
			return fmt.Errorf("EHLO failed after STARTTLS: %v", err)
		}
	}

	// 设置发件人
	if err := client.Mail(from, nil); err != nil {
		return fmt.Errorf("MAIL FROM failed: %v", err)
	}

	// 设置收件人
	successfulRecipients := []string{}
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient, nil); err != nil {
			log.Printf("⚠️  收件人 %s 被拒绝: %v", recipient, err)
			// 继续处理其他收件人，不立即返回错误
		} else {
			successfulRecipients = append(successfulRecipients, recipient)
		}
	}

	if len(successfulRecipients) == 0 {
		return fmt.Errorf("所有收件人都被拒绝")
	}

	log.Printf("📧 成功接受的收件人: %v", successfulRecipients)

	// 发送邮件内容
	dataWriter, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %v", err)
	}
	defer dataWriter.Close()

	// 重建完整的邮件内容（包括头部和正文）
	if err := s.writeCompleteMessage(dataWriter, originalMsg, body, from, successfulRecipients); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	// 完成发送
	if err := dataWriter.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	// 发送QUIT
	if err := client.Quit(); err != nil {
		log.Printf("⚠️  QUIT命令失败: %v", err)
	}

	log.Printf("✅ 邮件成功转发到 %s，成功收件人: %v", mxHost, successfulRecipients)
	return nil
}

// lookupMX 查找域名的MX记录
func (s *SMTPSession) lookupMX(domain string) (string, error) {
	// 这里应该实现真正的MX记录查询
	// 为了简化演示，我们使用一些常见的邮件服务器

	switch strings.ToLower(domain) {
	case "gmail.com", "googlemail.com":
		return "gmail-smtp-in.l.google.com", nil
	case "qq.com":
		return "mx1.qq.com", nil
	case "163.com":
		return "mx.163.com", nil
	case "126.com":
		return "mx.126.com", nil
	case "sina.com":
		return "mx.sina.com", nil
	case "hotmail.com", "outlook.com", "live.com":
		return "mx1.hotmail.com", nil
	case "yahoo.com":
		return "mta5.am0.yahoodns.net", nil
	case "example.com":
		return "mx.example.com", nil
	default:
		// 对于其他域名，尝试使用通用的MX记录格式
		// 在生产环境中，应该使用DNS查询
		possibleMX := []string{
			"mx." + domain,
			"mx1." + domain,
			"mail." + domain,
			domain, // 有些域名直接使用主域名作为MX
		}

		// 返回第一个可能的MX记录
		// 在真实环境中，这里应该做DNS查询验证
		log.Printf("💡 使用默认MX记录格式: mx.%s", domain)
		return possibleMX[0], nil
	}
}

// writeCompleteMessage 写入完整的邮件消息
func (s *SMTPSession) writeCompleteMessage(writer io.Writer, originalMsg *message.Entity, body []byte, from string, recipients []string) error {
	// 添加必要的邮件头
	fmt.Fprintf(writer, "From: %s\r\n", from)
	fmt.Fprintf(writer, "To: %s\r\n", strings.Join(recipients, ", "))

	// 复制原始邮件头（除了From和To）
	fields := originalMsg.Header.Fields()
	for fields.Next() {
		key := fields.Key()
		value := fields.Value()
		if strings.ToLower(key) == "from" || strings.ToLower(key) == "to" {
			continue // 跳过，我们已经设置了
		}
		fmt.Fprintf(writer, "%s: %s\r\n", key, value)
	}

	// 添加转发信息
	fmt.Fprintf(writer, "X-Relayed-By: %s\r\n", s.backend.domain)
	fmt.Fprintf(writer, "X-Relayed-At: %s\r\n", time.Now().Format(time.RFC1123Z))

	// 空行分隔头部和正文
	fmt.Fprintf(writer, "\r\n")

	// 写入邮件正文
	_, err := writer.Write(body)
	return err
}
