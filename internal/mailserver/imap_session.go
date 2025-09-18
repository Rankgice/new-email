package mailserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/rankgice/new-email/internal/model"
)

// IMAPSession 实现 imapserver.Session 接口
type IMAPSession struct {
	username       string
	mailbox        *model.Mailbox
	storage        *MailStorage
	selectedFolder *model.Folder
	authenticated  bool
	mailboxTracker *imapserver.MailboxTracker
}

// NewIMAPSession 创建新的 IMAP 会话
func NewIMAPSession(storage *MailStorage) *IMAPSession {
	return &IMAPSession{
		storage:       storage,
		authenticated: false,
	}
}

// Close 关闭会话
func (s *IMAPSession) Close() error {
	log.Printf("IMAP会话关闭: %s", s.username)
	return nil
}

// Login 用户登录
func (s *IMAPSession) Login(username, password string) error {
	log.Printf("IMAP登录尝试: %s", username)

	// 验证用户凭据
	if !s.storage.ValidateCredentials(username, password) {
		log.Printf("IMAP登录失败: %s", username)
		return imapserver.ErrAuthFailed
	}

	// 获取邮箱信息
	mailbox, err := s.storage.findMailboxByEmail(username)
	if err != nil {
		log.Printf("获取邮箱信息失败: %v", err)
		return err
	}
	if mailbox == nil {
		log.Printf("邮箱不存在: %s", username)
		return imapserver.ErrAuthFailed
	}

	s.username = username
	s.mailbox = mailbox
	s.authenticated = true

	log.Printf("IMAP登录成功: %s", username)
	return nil
}

// Select 选择邮箱
func (s *IMAPSession) Select(mailboxName string, options *imap.SelectOptions) (*imap.SelectData, error) {
	if !s.authenticated {
		return nil, errors.New("未认证")
	}

	log.Printf("选择邮箱: %s, 用户: %s", mailboxName, s.username)

	// 获取文件夹
	folder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, mailboxName, nil)
	if err != nil {
		log.Printf("获取文件夹失败: %v", err)
		return nil, err
	}
	if folder == nil {
		return nil, errors.New("邮箱不存在")
	}

	s.selectedFolder = folder

	// 获取邮件数量
	mails, err := s.storage.GetMails(s.username, mailboxName, 0)
	if err != nil {
		log.Printf("获取邮件数量失败: %v", err)
		return nil, err
	}

	numMessages := uint32(len(mails))
	s.mailboxTracker = imapserver.NewMailboxTracker(numMessages)

	// 计算未读邮件数量
	var numUnseen uint32
	for _, mail := range mails {
		if !mail.IsRead {
			numUnseen++
		}
	}

	selectData := &imap.SelectData{
		NumMessages: numMessages,
		UIDNext:     imap.UID(numMessages + 1), // 简化实现
		UIDValidity: 1,                         // 简化实现
		// NumUnseen 在 v2 中不再是 SelectData 的字段
	}

	return selectData, nil
}

// Create 创建邮箱
func (s *IMAPSession) Create(mailboxName string, options *imap.CreateOptions) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("创建邮箱: %s, 用户: %s", mailboxName, s.username)

	// 检查文件夹是否已存在
	existingFolder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, mailboxName, nil)
	if err != nil {
		return err
	}
	if existingFolder != nil {
		return errors.New("邮箱已存在")
	}

	// 创建新文件夹
	_, err = s.storage.getOrCreateFolder(s.mailbox.Id, mailboxName, nil, false)
	if err != nil {
		log.Printf("创建文件夹失败: %v", err)
		return err
	}

	log.Printf("成功创建邮箱: %s", mailboxName)
	return nil
}

// Delete 删除邮箱
func (s *IMAPSession) Delete(mailboxName string) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("删除邮箱: %s, 用户: %s", mailboxName, s.username)

	folder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, mailboxName, nil)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("邮箱不存在")
	}
	if folder.IsSystem {
		return errors.New("不能删除系统邮箱")
	}

	if err := s.storage.folderModel.Delete(folder.Id); err != nil {
		log.Printf("删除邮箱失败: %v", err)
		return err
	}

	log.Printf("成功删除邮箱: %s", mailboxName)
	return nil
}

// Rename 重命名邮箱
func (s *IMAPSession) Rename(oldName, newName string, options *imap.RenameOptions) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("重命名邮箱: %s -> %s, 用户: %s", oldName, newName, s.username)

	// 获取原文件夹
	folder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, oldName, nil)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("原邮箱不存在")
	}
	if folder.IsSystem {
		return errors.New("不能重命名系统邮箱")
	}

	// 检查新名称是否已存在
	existingFolder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, newName, folder.ParentId)
	if err != nil {
		return err
	}
	if existingFolder != nil {
		return errors.New("新邮箱名称已存在")
	}

	// 更新文件夹名称
	folder.Name = newName
	folder.UpdatedAt = time.Now()
	if err := s.storage.folderModel.Update(folder); err != nil {
		log.Printf("重命名邮箱失败: %v", err)
		return err
	}

	log.Printf("成功重命名邮箱: %s -> %s", oldName, newName)
	return nil
}

// Subscribe 订阅邮箱
func (s *IMAPSession) Subscribe(mailboxName string) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("订阅邮箱: %s, 用户: %s (简化实现)", mailboxName, s.username)
	// 简化实现，不实际存储订阅状态
	return nil
}

// Unsubscribe 取消订阅邮箱
func (s *IMAPSession) Unsubscribe(mailboxName string) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("取消订阅邮箱: %s, 用户: %s (简化实现)", mailboxName, s.username)
	// 简化实现，不实际存储订阅状态
	return nil
}

// List 列出邮箱
func (s *IMAPSession) List(w *imapserver.ListWriter, ref string, patterns []string, options *imap.ListOptions) error {
	if !s.authenticated {
		return errors.New("未认证")
	}

	log.Printf("列出邮箱: ref=%s, patterns=%v, 用户: %s", ref, patterns, s.username)

	folders, err := s.storage.folderModel.GetByMailboxId(s.mailbox.Id)
	if err != nil {
		log.Printf("获取文件夹列表失败: %v", err)
		return err
	}

	for _, folder := range folders {
		// 简化的模式匹配，实际应该使用更复杂的匹配逻辑
		matched := false
		for _, pattern := range patterns {
			if pattern == "*" || pattern == folder.Name || strings.Contains(folder.Name, strings.Trim(pattern, "*")) {
				matched = true
				break
			}
		}

		if matched {
			listData := &imap.ListData{
				Attrs:   []imap.MailboxAttr{},
				Delim:   '/',
				Mailbox: folder.Name,
			}
			if err := w.WriteList(listData); err != nil {
				log.Printf("写入列表数据失败: %v", err)
				return err
			}
		}
	}

	return nil
}

// Status 获取邮箱状态
func (s *IMAPSession) Status(mailboxName string, options *imap.StatusOptions) (*imap.StatusData, error) {
	if !s.authenticated {
		return nil, errors.New("未认证")
	}

	log.Printf("获取邮箱状态: %s, 用户: %s", mailboxName, s.username)

	// 检查邮箱是否存在
	folder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, mailboxName, nil)
	if err != nil {
		log.Printf("获取文件夹失败: %v", err)
		return nil, err
	}
	if folder == nil {
		log.Printf("邮箱不存在: %s", mailboxName)
		return nil, errors.New("邮箱不存在")
	}

	// 获取邮件列表
	mails, err := s.storage.GetMails(s.username, mailboxName, 0)
	if err != nil {
		log.Printf("获取邮件失败: %v", err)
		// 如果获取邮件失败，返回空状态而不是错误
		mails = []*StoredMail{}
	}

	statusData := &imap.StatusData{
		Mailbox: mailboxName,
	}

	// 根据请求的选项有条件地设置字段
	if options != nil {
		if options.NumMessages {
			numMessages := uint32(len(mails))
			statusData.NumMessages = &numMessages
		}

		if options.NumUnseen {
			var numUnseen uint32
			for _, mail := range mails {
				if mail != nil && !mail.IsRead {
					numUnseen++
				}
			}
			statusData.NumUnseen = &numUnseen
		}

		if options.UIDNext {
			statusData.UIDNext = imap.UID(uint32(len(mails)) + 1)
		}

		if options.UIDValidity {
			statusData.UIDValidity = 1
		}

		if options.NumRecent {
			numRecent := uint32(0) // 简化实现，没有最近邮件
			statusData.NumRecent = &numRecent
		}

		if options.Size {
			var totalSize int64
			for _, mail := range mails {
				if mail != nil {
					totalSize += int64(len(mail.Body))
				}
			}
			statusData.Size = &totalSize
		}
	} else {
		// 如果 options 为 nil，提供基本状态信息
		log.Printf("警告: StatusOptions 为 nil，提供默认状态")
		numMessages := uint32(len(mails))
		statusData.NumMessages = &numMessages
		statusData.UIDNext = imap.UID(numMessages + 1)
		statusData.UIDValidity = 1
	}

	log.Printf("邮箱状态: %s - 请求选项: %+v", mailboxName, options)
	return statusData, nil
}

// Append 追加邮件到邮箱
func (s *IMAPSession) Append(mailboxName string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error) {
	if !s.authenticated {
		return nil, errors.New("未认证")
	}

	log.Printf("追加邮件到邮箱: %s, 用户: %s", mailboxName, s.username)

	// 获取目标文件夹
	folder, err := s.storage.folderModel.GetByMailboxIdAndName(s.mailbox.Id, mailboxName, nil)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, errors.New("邮箱不存在")
	}

	// 读取邮件内容
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		log.Printf("读取邮件内容失败: %v", err)
		return nil, err
	}
	rawBody := buf.String()

	// 解析邮件头部
	from, to, subject := parseEmailHeaders(rawBody)

	// 创建存储邮件对象
	storedMail := &StoredMail{
		From:       from,
		To:         []string{to},
		Subject:    subject,
		Body:       rawBody,
		Received:   time.Now(),
		IsRead:     false,
		FolderId:   folder.Id,
		FolderName: mailboxName,
		MailboxID:  s.mailbox.Id,
		Username:   s.username,
		MessageID:  fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), s.storage.domain),
	}

	// 如果有标志，设置已读状态
	if options != nil && options.Flags != nil {
		for _, flag := range options.Flags {
			if flag == imap.FlagSeen {
				storedMail.IsRead = true
				break
			}
		}
	}

	// 保存邮件
	if err := s.storage.SaveMail(storedMail); err != nil {
		log.Printf("保存邮件失败: %v", err)
		return nil, err
	}

	log.Printf("成功追加邮件到邮箱: %s", mailboxName)

	// 返回追加数据
	appendData := &imap.AppendData{
		UID: imap.UID(storedMail.ID),
	}

	return appendData, nil
}

// Poll 轮询更新
func (s *IMAPSession) Poll(w *imapserver.UpdateWriter, allowExpunge bool) error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("轮询邮箱更新: %s, 用户: %s", s.selectedFolder.Name, s.username)
	// 简化实现，不发送更新
	return nil
}

// Idle 空闲模式
func (s *IMAPSession) Idle(w *imapserver.UpdateWriter, stop <-chan struct{}) error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("进入空闲模式: %s, 用户: %s", s.selectedFolder.Name, s.username)

	// 等待停止信号
	<-stop
	log.Printf("退出空闲模式: %s, 用户: %s", s.selectedFolder.Name, s.username)
	return nil
}

// Unselect 取消选择邮箱
func (s *IMAPSession) Unselect() error {
	if !s.authenticated || s.selectedFolder == nil {
		return errors.New("未选择邮箱")
	}

	log.Printf("取消选择邮箱: %s, 用户: %s", s.selectedFolder.Name, s.username)
	s.selectedFolder = nil
	s.mailboxTracker = nil
	return nil
}

// parseEmailHeaders 解析邮件头部（简化版本）
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
