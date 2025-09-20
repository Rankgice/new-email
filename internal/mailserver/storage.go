package mailserver

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/rankgice/new-email/internal/model"
	"github.com/rankgice/new-email/pkg/auth"
	"gorm.io/gorm"
)

// MailStorage 邮件存储
type MailStorage struct {
	emailModel   *model.EmailModel
	mailboxModel *model.MailboxModel
	domainModel  *model.DomainModel
	folderModel  *model.FolderModel // 新增
	domain       string
}

// StoredMail 存储的邮件
type StoredMail struct {
	ID          int64     `json:"id"`
	MessageID   string    `json:"message_id"`
	From        string    `json:"from"`
	To          []string  `json:"to"`
	Cc          []string  `json:"cc"`
	Bcc         []string  `json:"bcc"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	ContentType string    `json:"content_type"`
	Size        int       `json:"size"`
	Received    time.Time `json:"received"`
	IsRead      bool      `json:"is_read"`
	FolderId    int64     `json:"folder_id"`   // 文件夹ID
	FolderName  string    `json:"folder_name"` // 文件夹名称
	MailboxID   int64     `json:"mailbox_id"`
	Username    string    `json:"username"`
}

// NewMailStorage 创建邮件存储
func NewMailStorage(db *gorm.DB, domain string) *MailStorage {
	s := &MailStorage{
		emailModel:   model.NewEmailModel(db),
		mailboxModel: model.NewMailboxModel(db),
		domainModel:  model.NewDomainModel(db),
		folderModel:  model.NewFolderModel(db),
		domain:       domain,
	}
	// 确保系统文件夹存在
	s.ensureSystemFoldersExist(db)
	return s
}

// ensureSystemFoldersExist 确保每个邮箱都有默认的系统文件夹
func (s *MailStorage) ensureSystemFoldersExist(db *gorm.DB) {
	mailboxes, _, err := s.mailboxModel.List(model.MailboxListParams{})
	if err != nil {
		log.Printf("获取所有邮箱失败: %v", err)
		return
	}

	systemFolders := []string{"INBOX", "Sent", "Drafts", "Trash"}

	for _, mailbox := range mailboxes {
		for _, folderName := range systemFolders {
			_, err := s.getOrCreateFolder(mailbox.Id, folderName, nil, true)
			if err != nil {
				log.Printf("为邮箱 %s 创建系统文件夹 %s 失败: %v", mailbox.Email, folderName, err)
			}
		}
	}
}

// getOrCreateFolder 获取或创建文件夹
func (s *MailStorage) getOrCreateFolder(mailboxId int64, name string, parentId *int64, isSystem bool) (*model.Folder, error) {
	folder, err := s.folderModel.GetByMailboxIdAndName(mailboxId, name, parentId)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		folder = &model.Folder{
			MailboxId: mailboxId,
			Name:      name,
			ParentId:  parentId,
			IsSystem:  isSystem,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.folderModel.Create(folder); err != nil {
			return nil, err
		}
	}
	return folder, nil
}

// SaveMail (用于APPEND)
func (s *MailStorage) SaveMail(mail *StoredMail) error {
	// 1. 根据用户名查找邮箱
	mailbox, err := s.findMailboxByEmail(mail.Username)
	if err != nil {
		log.Printf("为APPEND查找邮箱失败 %s: %v", mail.Username, err)
		return err
	}
	if mailbox == nil {
		log.Printf("APPEND时邮箱不存在: %s", mail.Username)
		return fmt.Errorf("邮箱 %s 不存在", mail.Username)
	}

	// 2. 确定邮件方向
	direction := "received"
	// 如果发件人是自己，则认为是已发送邮件
	if mail.From == mail.Username {
		direction = "sent"
	}

	// 3. 获取或创建目标文件夹
	folder, err := s.getOrCreateFolder(mailbox.Id, mail.FolderName, nil, false) // 假设APPEND的文件夹不是系统文件夹
	if err != nil {
		log.Printf("为APPEND获取或创建文件夹失败 %s/%s: %v", mail.Username, mail.FolderName, err)
		return err
	}

	// 4. 创建邮件记录
	email := &model.Email{
		UserId:      mailbox.UserId,
		MailboxId:   mailbox.Id,
		MessageId:   mail.MessageID,
		Subject:     mail.Subject,
		FromEmail:   mail.From,
		ToEmails:    mail.To,
		CcEmails:    mail.Cc,
		BccEmails:   mail.Bcc,
		Content:     mail.Body, // 存储原始邮件体
		ContentType: mail.ContentType,
		IsRead:      mail.IsRead,
		IsStarred:   false,
		FolderId:    folder.Id, // 使用文件夹ID
		Direction:   direction,
		ReceivedAt:  &mail.Received,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 5. 保存到数据库
	if err := s.emailModel.Create(email); err != nil {
		log.Printf("APPEND存储邮件失败: %v", err)
		return err
	}

	log.Printf("✅ 邮件已通过APPEND存储到邮箱: %s, 文件夹: %s (ID: %d)", mail.Username, mail.FolderName, folder.Id)
	return nil
}

// StoreMail 存储邮件
func (s *MailStorage) StoreMail(mail *StoredMail) error {
	log.Printf("🎯 StoreMail: 开始存储邮件, From=%s, To=%v, Subject=%s", mail.From, mail.To, mail.Subject)

	// 查找目标邮箱
	for _, toAddr := range mail.To {
		log.Printf("🔍 处理收件人: %s", toAddr)
		mailbox, err := s.findMailboxByEmail(toAddr)
		if err != nil {
			log.Printf("❌ 查找邮箱失败 %s: %v", toAddr, err)
			continue
		}
		if mailbox == nil {
			log.Printf("❌ 邮箱不存在: %s", toAddr)
			continue
		}
		log.Printf("✅ 找到收件人邮箱: ID=%d, Email=%s, UserId=%d", mailbox.Id, mailbox.Email, mailbox.UserId)

		// 获取或创建INBOX文件夹
		inboxFolder, err := s.getOrCreateFolder(mailbox.Id, "INBOX", nil, true)
		if err != nil {
			log.Printf("为邮箱 %s 获取或创建INBOX文件夹失败: %v", toAddr, err)
			continue
		}

		// 创建邮件记录
		email := &model.Email{
			UserId:      mailbox.UserId, // 添加用户ID
			MailboxId:   mailbox.Id,
			MessageId:   mail.MessageID,
			Subject:     mail.Subject,
			FromEmail:   mail.From,
			ToEmails:    mail.To,
			CcEmails:    mail.Cc,
			BccEmails:   mail.Bcc,
			Content:     mail.Body,
			ContentType: mail.ContentType,
			IsRead:      false,
			IsStarred:   false,
			FolderId:    inboxFolder.Id, // 存储到INBOX文件夹
			Direction:   "received",
			ReceivedAt:  &mail.Received,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.emailModel.Create(email); err != nil {
			log.Printf("存储邮件失败: %v", err)
			return err
		}

		log.Printf("✅ 邮件已存储到邮箱: %s (ID: %d), 文件夹: %s (ID: %d)", toAddr, mailbox.Id, inboxFolder.Name, inboxFolder.Id)
	}

	return nil
}

// GetMails 获取邮件列表
func (s *MailStorage) GetMails(mailboxEmail string, folderName string, limit int) ([]*StoredMail, error) {
	log.Printf("🔍 GetMails: 查询邮箱=%s, 文件夹=%s", mailboxEmail, folderName)

	mailbox, err := s.findMailboxByEmail(mailboxEmail)
	if err != nil {
		log.Printf("❌ 查找邮箱失败: %v", err)
		return nil, err
	}
	if mailbox == nil {
		log.Printf("❌ 邮箱不存在: %s", mailboxEmail)
		return nil, fmt.Errorf("邮箱不存在: %s", mailboxEmail)
	}
	log.Printf("✅ 找到邮箱: ID=%d, Email=%s", mailbox.Id, mailbox.Email)

	folder, err := s.folderModel.GetByMailboxIdAndName(mailbox.Id, folderName, nil)
	if err != nil {
		log.Printf("❌ 查找文件夹失败: %v", err)
		return nil, err
	}
	if folder == nil {
		log.Printf("❌ 文件夹不存在: %s", folderName)
		return nil, fmt.Errorf("文件夹不存在: %s", folderName)
	}
	log.Printf("✅ 找到文件夹: ID=%d, Name=%s", folder.Id, folder.Name)

	emails, err := s.emailModel.GetByFolderId(folder.Id, limit)
	if err != nil {
		log.Printf("❌ 查询邮件失败: %v", err)
		return nil, err
	}
	log.Printf("📧 查询到 %d 封邮件", len(emails))

	var mails []*StoredMail
	for _, email := range emails {
		receivedAt := time.Now()
		if email.ReceivedAt != nil {
			receivedAt = *email.ReceivedAt
		}

		mail := &StoredMail{
			ID:          email.Id,
			MessageID:   email.MessageId,
			From:        email.FromEmail,
			To:          email.ToEmails,
			Cc:          email.CcEmails,
			Bcc:         email.BccEmails,
			Subject:     email.Subject,
			Body:        email.Content,
			ContentType: email.ContentType,
			Size:        len(email.Content), // 使用内容长度作为大小
			Received:    receivedAt,
			IsRead:      email.IsRead,
			FolderId:    email.FolderId,
			FolderName:  folder.Name, // 从获取到的文件夹对象中获取名称
			MailboxID:   email.MailboxId,
		}
		mails = append(mails, mail)
	}

	return mails, nil
}

// GetMail 获取单个邮件
func (s *MailStorage) GetMail(mailboxEmail string, messageID string) (*StoredMail, error) {
	mailbox, err := s.findMailboxByEmail(mailboxEmail)
	if err != nil {
		return nil, err
	}
	if mailbox == nil {
		return nil, fmt.Errorf("邮箱不存在: %s", mailboxEmail)
	}

	// 从MessageID中提取邮件ID
	var emailID int64
	if _, err := fmt.Sscanf(messageID, "<%d@localhost>", &emailID); err != nil {
		return nil, fmt.Errorf("无效的MessageID: %s", messageID)
	}

	email, err := s.emailModel.GetById(emailID)
	if err != nil {
		return nil, err
	}
	if email == nil || email.MailboxId != mailbox.Id {
		return nil, fmt.Errorf("邮件不存在")
	}

	receivedAt := time.Now()
	if email.ReceivedAt != nil {
		receivedAt = *email.ReceivedAt
	}

	folder, err := s.folderModel.GetById(email.FolderId)
	if err != nil {
		log.Printf("获取邮件文件夹失败 (ID: %d): %v", email.FolderId, err)
		return nil, err
	}
	folderName := "Unknown"
	if folder != nil {
		folderName = folder.Name
	}

	mail := &StoredMail{
		ID:          email.Id,
		MessageID:   messageID,
		From:        email.FromEmail,
		To:          email.ToEmails,
		Cc:          email.CcEmails,
		Bcc:         email.BccEmails,
		Subject:     email.Subject,
		Body:        email.Content,
		ContentType: email.ContentType,
		Size:        len(email.Content), // 使用内容长度作为大小
		Received:    receivedAt,
		IsRead:      email.IsRead,
		FolderId:    email.FolderId,
		FolderName:  folderName,
		MailboxID:   email.MailboxId,
	}

	return mail, nil
}

// SearchMails 根据IMAP搜索条件搜索邮件
func (s *MailStorage) SearchMails(mailboxId int64, folderId int64, criteria *imap.SearchCriteria) ([]*StoredMail, error) {
	emails, err := s.emailModel.Search(mailboxId, folderId, criteria)
	if err != nil {
		return nil, err
	}

	var mails []*StoredMail
	for _, email := range emails {
		receivedAt := time.Now()
		if email.ReceivedAt != nil {
			receivedAt = *email.ReceivedAt
		}

		folder, err := s.folderModel.GetById(email.FolderId)
		if err != nil {
			log.Printf("搜索邮件时获取文件夹失败 (ID: %d): %v", email.FolderId, err)
			continue
		}
		folderName := "Unknown"
		if folder != nil {
			folderName = folder.Name
		}

		mail := &StoredMail{
			ID:          email.Id,
			MessageID:   fmt.Sprintf("<%d@%s>", email.Id, "localhost"),
			From:        email.FromEmail,
			To:          email.ToEmails,
			Cc:          email.CcEmails,
			Bcc:         email.BccEmails,
			Subject:     email.Subject,
			Body:        email.Content,
			ContentType: email.ContentType,
			Size:        len(email.Content),
			Received:    receivedAt,
			IsRead:      email.IsRead,
			FolderId:    email.FolderId,
			FolderName:  folderName,
			MailboxID:   email.MailboxId,
		}
		mails = append(mails, mail)
	}

	return mails, nil
}

// MarkAsRead 标记邮件为已读
func (s *MailStorage) MarkAsRead(mailboxEmail string, messageID string) error {
	mail, err := s.GetMail(mailboxEmail, messageID)
	if err != nil {
		return err
	}

	return s.emailModel.MarkAsRead(mail.ID)
}

// ValidatePassword 验证邮箱密码
func (s *MailStorage) ValidatePassword(email, password string) bool {
	mailbox, err := s.findMailboxByEmail(email)
	if err != nil {
		log.Printf("验证凭据失败: %v", err)
		return false
	}
	if mailbox == nil {
		log.Printf("邮箱不存在: %s", email)
		return false
	}

	// 使用安全的密码验证
	if err := auth.CheckPassword(password, mailbox.Password); err != nil {
		log.Printf("密码验证失败 %s: %v", email, err)
		return false
	}

	log.Printf("✅ 邮箱凭据验证成功: %s", email)
	return true
}

// ValidateCredentials 验证邮箱凭据
func (s *MailStorage) ValidateCredentials(email, password string) bool {
	mailbox, err := s.findMailboxByEmail(email)
	if err != nil {
		log.Printf("验证凭据失败: %v", err)
		return false
	}
	if mailbox == nil {
		log.Printf("邮箱不存在: %s", email)
		return false
	}

	// 使用安全的密码验证
	if err := auth.CheckPassword(password, mailbox.Password); err != nil {
		log.Printf("密码验证失败 %s: %v", email, err)
		return false
	}

	log.Printf("✅ 邮箱凭据验证成功: %s", email)
	return true
}

// findMailboxByEmail 根据邮箱地址查找邮箱
func (s *MailStorage) findMailboxByEmail(email string) (*model.Mailbox, error) {
	return s.mailboxModel.GetByEmail(email)
}

// isMailboxExists 检查邮箱是否存在
func (s *MailStorage) isMailboxExists(email string) bool {
	mailbox, err := s.findMailboxByEmail(email)
	if err != nil {
		log.Printf("检查邮箱存在性时出错: %v", err)
		return false
	}
	return mailbox != nil
}
