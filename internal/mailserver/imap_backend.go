package mailserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-message"
	"github.com/rankgice/new-email/internal/model"
)

// CustomBackend 自定义IMAP后端
type CustomBackend struct {
	storage *MailStorage
}

// NewCustomBackend 创建自定义后端
func NewCustomBackend(storage *MailStorage) *CustomBackend {
	return &CustomBackend{
		storage: storage,
	}
}

// Login 用户登录
func (bkd *CustomBackend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	log.Printf("IMAP登录尝试: %s", username)

	// 验证用户凭据
	if !bkd.storage.ValidateCredentials(username, password) {
		log.Printf("IMAP登录失败: %s", username)
		return nil, errors.New("认证失败")
	}

	// 获取邮箱信息
	mailbox, err := bkd.storage.findMailboxByEmail(username)
	if err != nil {
		log.Printf("获取邮箱信息失败: %v", err)
		return nil, err
	}
	if mailbox == nil {
		log.Printf("邮箱不存在: %s", username)
		return nil, errors.New("邮箱不存在")
	}

	log.Printf("IMAP登录成功: %s", username)
	return &CustomUser{
		username: username,
		mailbox:  mailbox,
		storage:  bkd.storage,
	}, nil
}

// CustomUser 自定义用户
type CustomUser struct {
	username string
	mailbox  *model.Mailbox
	storage  *MailStorage
}

// Username 返回用户名
func (u *CustomUser) Username() string {
	return u.username
}

// ListMailboxes 列出邮箱文件夹
func (u *CustomUser) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	folders, err := u.storage.folderModel.GetByMailboxId(u.mailbox.Id)
	if err != nil {
		log.Printf("获取邮箱 %s 的文件夹失败: %v", u.username, err)
		return nil, err
	}

	var result []backend.Mailbox
	for _, folder := range folders {
		result = append(result, &CustomMailbox{
			name:    folder.Name,
			user:    u,
			storage: u.storage,
			folder:  folder, // 传递完整的文件夹对象
		})
	}

	return result, nil
}

// GetMailbox 获取邮箱
func (u *CustomUser) GetMailbox(name string) (backend.Mailbox, error) {
	folder, err := u.storage.folderModel.GetByMailboxIdAndName(u.mailbox.Id, name, nil) // 假设顶级文件夹
	if err != nil {
		log.Printf("获取邮箱 %s 的文件夹 %s 失败: %v", u.username, name, err)
		return nil, err
	}
	if folder == nil {
		return nil, errors.New("邮箱不存在")
	}

	return &CustomMailbox{
		name:    folder.Name,
		user:    u,
		storage: u.storage,
		folder:  folder,
	}, nil
}

// CreateMailbox 创建邮箱
func (u *CustomUser) CreateMailbox(name string) error {
	// 检查文件夹是否已存在
	existingFolder, err := u.storage.folderModel.GetByMailboxIdAndName(u.mailbox.Id, name, nil)
	if err != nil {
		return err
	}
	if existingFolder != nil {
		return errors.New("文件夹已存在")
	}

	// 创建新文件夹
	_, err = u.storage.getOrCreateFolder(u.mailbox.Id, name, nil, false)
	if err != nil {
		log.Printf("为邮箱 %s 创建文件夹 %s 失败: %v", u.username, name, err)
		return err
	}
	log.Printf("为邮箱 %s 成功创建文件夹: %s", u.username, name)
	return nil
}

// DeleteMailbox 删除邮箱
func (u *CustomUser) DeleteMailbox(name string) error {
	folder, err := u.storage.folderModel.GetByMailboxIdAndName(u.mailbox.Id, name, nil)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("文件夹不存在")
	}
	if folder.IsSystem {
		return errors.New("不能删除系统文件夹")
	}

	// TODO: 检查文件夹是否为空，或者是否需要移动邮件到Trash
	// 目前简化处理，直接删除
	if err := u.storage.folderModel.Delete(folder.Id); err != nil {
		log.Printf("删除邮箱 %s 的文件夹 %s 失败: %v", u.username, name, err)
		return err
	}
	log.Printf("成功删除邮箱 %s 的文件夹: %s", u.username, name)
	return nil
}

// RenameMailbox 重命名邮箱
func (u *CustomUser) RenameMailbox(existingName, newName string) error {
	return errors.New("不支持重命名邮箱")
}

// Logout 登出
func (u *CustomUser) Logout() error {
	log.Printf("IMAP用户登出: %s", u.username)
	return nil
}

// CustomMailbox 自定义邮箱
type CustomMailbox struct {
	name    string
	user    *CustomUser
	storage *MailStorage
	folder  *model.Folder // 新增
}

// Name 返回邮箱名称
func (mb *CustomMailbox) Name() string {
	return mb.name
}

// Info 返回邮箱信息
func (mb *CustomMailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{
		Attributes: []string{imap.NoInferiorsAttr},
		Delimiter:  "/",
		Name:       mb.name,
	}, nil
}

// Status 返回邮箱状态
func (mb *CustomMailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	// 获取邮件列表
	mails, err := mb.storage.GetMails(mb.user.username, mb.name, 1000)
	if err != nil {
		return nil, err
	}

	status := &imap.MailboxStatus{
		Name: mb.name,
	}

	for _, item := range items {
		switch item {
		case imap.StatusMessages:
			status.Messages = uint32(len(mails))
		case imap.StatusRecent:
			status.Recent = 0 // 简化实现，没有新邮件
		case imap.StatusUnseen:
			unseen := 0
			for _, mail := range mails {
				if !mail.IsRead {
					unseen++
				}
			}
			status.Unseen = uint32(unseen)
		case imap.StatusUidNext:
			status.UidNext = uint32(len(mails) + 1)
		case imap.StatusUidValidity:
			status.UidValidity = 1
		}
	}

	return status, nil
}

// SetSubscribed 设置订阅状态
func (mb *CustomMailbox) SetSubscribed(subscribed bool) error {
	return nil // 简化实现
}

// Check 检查邮箱
func (mb *CustomMailbox) Check() error {
	return nil
}

// ListMessages 列出消息
func (mb *CustomMailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	// 获取邮件列表
	mails, err := mb.storage.GetMails(mb.user.username, mb.name, 1000)
	if err != nil {
		return err
	}

	// 转换邮件为IMAP消息
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uidNum := uint32(mail.ID)

		// 检查是否在序列集中
		if uid {
			if !seqSet.Contains(uidNum) {
				continue
			}
		} else {
			if !seqSet.Contains(seqNum) {
				continue
			}
		}

		msg := &imap.Message{
			SeqNum: seqNum,
			Uid:    uidNum,
		}

		// 构建邮件内容
		body := mb.buildEmailBody(mail)
		bodyReader := strings.NewReader(body)

		// 解析邮件（暂时跳过解析，直接使用原始数据）
		_, err = message.Read(bodyReader)
		if err != nil {
			log.Printf("解析邮件失败: %v", err)
			continue
		}

		// 设置信封
		msg.Envelope = &imap.Envelope{
			Date:      mail.Received,
			Subject:   mail.Subject,
			From:      mb.parseAddressList(mail.From),
			To:        mb.parseAddressListList(mail.To),
			Cc:        mb.parseAddressListList(mail.Cc),
			Bcc:       mb.parseAddressListList(mail.Bcc),
			MessageId: mail.MessageID,
		}

		// 设置标志
		flags := []string{}
		if mail.IsRead {
			flags = append(flags, imap.SeenFlag)
		}
		msg.Flags = flags

		// 设置大小
		msg.Size = uint32(len(body))

		// 设置邮件体
		msg.Body = map[*imap.BodySectionName]imap.Literal{
			{}: strings.NewReader(body),
		}

		ch <- msg
	}

	return nil
}

// SearchMessages 搜索消息
func (mb *CustomMailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	// 简化实现，返回所有消息
	mails, err := mb.storage.GetMails(mb.user.username, mb.name, 1000)
	if err != nil {
		return nil, err
	}

	var results []uint32
	for i := range mails {
		if uid {
			results = append(results, uint32(mails[i].ID))
		} else {
			results = append(results, uint32(i+1))
		}
	}

	return results, nil
}

// CreateMessage 创建消息 (APPEND)
func (mb *CustomMailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	// 1. 读取邮件内容
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, body); err != nil {
		log.Printf("读取邮件正文失败: %v", err)
		return err
	}
	rawBody := buf.String()

	// 2. 解析邮件
	// 注意：这里的解析是简化的，实际应用中需要更健壮的MIME解析库
	// 暂时我们只将原始邮件内容存起来
	from, to, subject := parseEmailHeaders(rawBody)

	// 3. 存储邮件
	storedMail := &StoredMail{
		From:       from,
		To:         []string{to},
		Subject:    subject,
		Body:       rawBody,
		Received:   date,
		IsRead:     false, // 新邮件默认为未读
		FolderId:   mb.folder.Id,
		FolderName: mb.name,
		MailboxID:  mb.user.mailbox.Id,
		Username:   mb.user.Username(),
		MessageID:  fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), mb.user.storage.domain),
	}

	// 4. 调用storage层保存
	if err := mb.storage.SaveMail(storedMail); err != nil {
		log.Printf("保存邮件失败: %v", err)
		return err
	}

	log.Printf("成功追加邮件到 %s 文件夹，用户: %s", mb.name, mb.user.Username())
	return nil
}

// parseEmailHeaders 是一个简化的邮件头解析函数
func parseEmailHeaders(rawBody string) (from, to, subject string) {
	lines := strings.Split(rawBody, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "from: ") {
			from = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(strings.ToLower(line), "to: ") {
			to = strings.TrimSpace(line[4:])
		} else if strings.HasPrefix(strings.ToLower(line), "subject: ") {
			subject = strings.TrimSpace(line[9:])
		}
		// 遇到空行，说明邮件头结束
		if line == "" {
			break
		}
	}
	return
}

// UpdateMessagesFlags 更新消息标志
func (mb *CustomMailbox) UpdateMessagesFlags(uid bool, seqSet *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	// 获取邮件列表
	mails, err := mb.storage.GetMails(mb.user.username, mb.name, 1000)
	if err != nil {
		return err
	}

	// 更新标志
	for i, mail := range mails {
		seqNum := uint32(i + 1)
		uidNum := uint32(mail.ID)

		// 检查是否在序列集中
		if uid {
			if !seqSet.Contains(uidNum) {
				continue
			}
		} else {
			if !seqSet.Contains(seqNum) {
				continue
			}
		}

		// 检查是否包含Seen标志
		for _, flag := range flags {
			if flag == imap.SeenFlag {
				if op == imap.AddFlags {
					// 标记为已读
					if err := mb.storage.MarkAsRead(mb.user.username, mail.MessageID); err != nil {
						log.Printf("标记邮件已读失败: %v", err)
					}
				} else if op == imap.RemoveFlags {
					// 标记为未读
					if err := mb.storage.emailModel.MarkAsUnread(mail.ID); err != nil {
						log.Printf("标记邮件未读失败: %v", err)
					}
				}
				break
			}
		}
	}

	return nil
}

// CopyMessages 复制消息
func (mb *CustomMailbox) CopyMessages(uid bool, seqSet *imap.SeqSet, destName string) error {
	return errors.New("不支持复制消息")
}

// MoveMessages 移动消息
func (mb *CustomMailbox) MoveMessages(uid bool, seqSet *imap.SeqSet, destName string) error {
	return errors.New("不支持移动消息")
}

// Expunge 删除消息
func (mb *CustomMailbox) Expunge() error {
	return nil // 简化实现，不删除消息
}

// buildEmailBody 构建邮件体
func (mb *CustomMailbox) buildEmailBody(mail *StoredMail) string {
	// 构建简单的RFC822格式邮件
	body := fmt.Sprintf("From: %s\r\n", mail.From)
	body += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ", "))
	if len(mail.Cc) > 0 {
		body += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ", "))
	}
	body += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	body += fmt.Sprintf("Date: %s\r\n", mail.Received.Format(time.RFC1123Z))
	body += fmt.Sprintf("Message-ID: %s\r\n", mail.MessageID)
	body += fmt.Sprintf("Content-Type: %s\r\n", mail.ContentType)
	body += "\r\n"
	body += mail.Body

	return body
}

// parseAddressList 解析地址列表
func (mb *CustomMailbox) parseAddressList(email string) []*imap.Address {
	if email == "" {
		return nil
	}

	// 简单解析，假设格式为 "name@domain.com"
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil
	}

	return []*imap.Address{
		{
			MailboxName: parts[0],
			HostName:    parts[1],
		},
	}
}

// parseAddressListList 解析地址列表列表
func (mb *CustomMailbox) parseAddressListList(emails []string) []*imap.Address {
	var addresses []*imap.Address
	for _, email := range emails {
		addresses = append(addresses, mb.parseAddressList(email)...)
	}
	return addresses
}
