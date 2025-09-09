package service

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// IMAPConfig IMAP配置
type IMAPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseTLS   bool   `json:"useTLS"`
}

// IMAPService IMAP服务
type IMAPService struct {
	config        IMAPConfig
	client        *client.Client
	messageParser *MessageParser
}

// NewIMAPService 创建IMAP服务
func NewIMAPService(config IMAPConfig) *IMAPService {
	return &IMAPService{
		config:        config,
		messageParser: NewMessageParser(),
	}
}

// IMAPEmail IMAP邮件
type IMAPEmail struct {
	UID         uint32           `json:"uid"`
	MessageID   string           `json:"messageId"`
	Subject     string           `json:"subject"`
	From        string           `json:"from"`
	To          []string         `json:"to"`
	Cc          []string         `json:"cc"`
	Bcc         []string         `json:"bcc"`
	Date        time.Time        `json:"date"`
	Body        string           `json:"body"`
	ContentType string           `json:"contentType"`
	IsRead      bool             `json:"isRead"`
	Size        uint32           `json:"size"`
	Attachments []IMAPAttachment `json:"attachments"`
}

// IMAPAttachment IMAP附件
type IMAPAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
	Data        []byte `json:"data"`
}

// Connect 连接到IMAP服务器
func (s *IMAPService) Connect() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var c *client.Client
	var err error

	if s.config.UseTLS {
		// 使用TLS连接
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // 允许跳过证书验证，仅用于开发测试
			ServerName:         s.config.Host,
		}
		c, err = client.DialTLS(addr, tlsConfig)
	} else {
		// 使用普通连接
		c, err = client.Dial(addr)
	}

	if err != nil {
		return fmt.Errorf("连接IMAP服务器失败: %v", err)
	}

	// 登录
	if err := c.Login(s.config.Username, s.config.Password); err != nil {
		c.Logout()
		return fmt.Errorf("IMAP登录失败: %v", err)
	}

	s.client = c
	return nil
}

// Disconnect 断开连接
func (s *IMAPService) Disconnect() error {
	if s.client != nil {
		err := s.client.Logout()
		s.client = nil
		return err
	}
	return nil
}

// TestConnection 测试IMAP连接
func (s *IMAPService) TestConnection() error {
	// 如果没有配置主机，跳过测试
	if s.config.Host == "" {
		return nil
	}

	if err := s.Connect(); err != nil {
		return err
	}
	return s.Disconnect()
}

// ListMailboxes 列出邮箱文件夹
func (s *IMAPService) ListMailboxes() ([]string, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接到IMAP服务器")
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- s.client.List("", "*", mailboxes)
	}()

	var names []string
	for m := range mailboxes {
		names = append(names, m.Name)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return names, nil
}

// SelectMailbox 选择邮箱文件夹
func (s *IMAPService) SelectMailbox(mailbox string) (*imap.MailboxStatus, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接到IMAP服务器")
	}

	mbox, err := s.client.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("选择邮箱失败: %v", err)
	}

	return mbox, nil
}

// FetchEmails 获取邮件列表
func (s *IMAPService) FetchEmails(mailbox string, limit uint32) ([]*IMAPEmail, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接到IMAP服务器")
	}

	// 选择邮箱
	mbox, err := s.SelectMailbox(mailbox)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []*IMAPEmail{}, nil
	}

	// 计算获取范围
	from := uint32(1)
	to := mbox.Messages
	if limit > 0 && limit < mbox.Messages {
		from = mbox.Messages - limit + 1
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	// 获取邮件
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- s.client.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchFlags,
			imap.FetchRFC822Size,
			imap.FetchUid,
		}, messages)
	}()

	var emails []*IMAPEmail
	for msg := range messages {
		email := &IMAPEmail{
			UID:     msg.Uid,
			Subject: msg.Envelope.Subject,
			From:    s.formatAddresses(msg.Envelope.From),
			To:      s.formatAddressList(msg.Envelope.To),
			Cc:      s.formatAddressList(msg.Envelope.Cc),
			Bcc:     s.formatAddressList(msg.Envelope.Bcc),
			Date:    msg.Envelope.Date,
			Size:    msg.Size,
			IsRead:  !s.hasFlag(msg.Flags, "\\Seen"),
		}

		if msg.Envelope.MessageId != "" {
			email.MessageID = msg.Envelope.MessageId
		}

		emails = append(emails, email)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return emails, nil
}

// FetchEmailBody 获取邮件正文
func (s *IMAPService) FetchEmailBody(uid uint32) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("未连接到IMAP服务器")
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- s.client.UidFetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
	}()

	msg := <-messages
	if msg == nil {
		return "", fmt.Errorf("邮件不存在")
	}

	if err := <-done; err != nil {
		return "", err
	}

	// 解析邮件内容
	r := msg.GetBody(&imap.BodySectionName{})
	if r == nil {
		return "", fmt.Errorf("无法获取邮件内容")
	}

	// 使用新的邮件解析器
	parsed, err := s.messageParser.ParseMessage(r)
	if err != nil {
		return "", fmt.Errorf("解析邮件失败: %v", err)
	}

	// 返回文本内容
	return s.messageParser.ExtractTextContent(parsed), nil
}

// GetParsedEmail 获取解析后的完整邮件
func (s *IMAPService) GetParsedEmail(uid uint32) (*ParsedMessage, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接到IMAP服务器")
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- s.client.UidFetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
	}()

	msg := <-messages
	if msg == nil {
		return nil, fmt.Errorf("邮件不存在")
	}

	if err := <-done; err != nil {
		return nil, err
	}

	// 解析邮件内容
	r := msg.GetBody(&imap.BodySectionName{})
	if r == nil {
		return nil, fmt.Errorf("无法获取邮件内容")
	}

	// 使用新的邮件解析器
	return s.messageParser.ParseMessage(r)
}

// MarkAsRead 标记邮件为已读
func (s *IMAPService) MarkAsRead(uid uint32) error {
	if s.client == nil {
		return fmt.Errorf("未连接到IMAP服务器")
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{"\\Seen"}

	return s.client.UidStore(seqset, item, flags, nil)
}

// MarkAsUnread 标记邮件为未读
func (s *IMAPService) MarkAsUnread(uid uint32) error {
	if s.client == nil {
		return fmt.Errorf("未连接到IMAP服务器")
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	item := imap.FormatFlagsOp(imap.RemoveFlags, true)
	flags := []interface{}{"\\Seen"}

	return s.client.UidStore(seqset, item, flags, nil)
}

// formatAddresses 格式化地址
func (s *IMAPService) formatAddresses(addresses []*imap.Address) string {
	if len(addresses) == 0 {
		return ""
	}
	return fmt.Sprintf("%s@%s", addresses[0].MailboxName, addresses[0].HostName)
}

// formatAddressList 格式化地址列表
func (s *IMAPService) formatAddressList(addresses []*imap.Address) []string {
	var result []string
	for _, addr := range addresses {
		result = append(result, fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName))
	}
	return result
}

// hasFlag 检查是否有指定标志
func (s *IMAPService) hasFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}
