package mailserver

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"new-email/internal/model"
)

// MailStorage 邮件存储
type MailStorage struct {
	emailModel   *model.EmailModel
	mailboxModel *model.MailboxModel
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
	Folder      string    `json:"folder"`
	MailboxID   int64     `json:"mailbox_id"`
}

// NewMailStorage 创建邮件存储
func NewMailStorage(db *gorm.DB) *MailStorage {
	return &MailStorage{
		emailModel:   model.NewEmailModel(db),
		mailboxModel: model.NewMailboxModel(db),
	}
}

// StoreMail 存储邮件
func (s *MailStorage) StoreMail(mail *StoredMail) error {
	// 查找目标邮箱
	for _, toAddr := range mail.To {
		mailbox, err := s.findMailboxByEmail(toAddr)
		if err != nil {
			log.Printf("查找邮箱失败 %s: %v", toAddr, err)
			continue
		}
		if mailbox == nil {
			log.Printf("邮箱不存在: %s", toAddr)
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
			Direction:   "received",
			ReceivedAt:  &mail.Received,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.emailModel.Create(email); err != nil {
			log.Printf("存储邮件失败: %v", err)
			return err
		}

		log.Printf("✅ 邮件已存储到邮箱: %s (ID: %d)", toAddr, mailbox.Id)
	}

	return nil
}

// GetMails 获取邮件列表
func (s *MailStorage) GetMails(mailboxEmail string, folder string, limit int) ([]*StoredMail, error) {
	mailbox, err := s.findMailboxByEmail(mailboxEmail)
	if err != nil {
		return nil, err
	}
	if mailbox == nil {
		return nil, fmt.Errorf("邮箱不存在: %s", mailboxEmail)
	}

	emails, err := s.emailModel.GetByMailboxId(mailbox.Id, limit)
	if err != nil {
		return nil, err
	}

	var mails []*StoredMail
	for _, email := range emails {
		receivedAt := time.Now()
		if email.ReceivedAt != nil {
			receivedAt = *email.ReceivedAt
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
			Size:        len(email.Content), // 使用内容长度作为大小
			Received:    receivedAt,
			IsRead:      email.IsRead,
			Folder:      "INBOX",
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
		Folder:      "INBOX",
		MailboxID:   email.MailboxId,
	}

	return mail, nil
}

// MarkAsRead 标记邮件为已读
func (s *MailStorage) MarkAsRead(mailboxEmail string, messageID string) error {
	mail, err := s.GetMail(mailboxEmail, messageID)
	if err != nil {
		return err
	}

	return s.emailModel.MarkAsRead(mail.ID)
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

	// 简单的密码验证（实际项目中应该使用加密）
	return mailbox.Password == password
}

// findMailboxByEmail 根据邮箱地址查找邮箱
func (s *MailStorage) findMailboxByEmail(email string) (*model.Mailbox, error) {
	return s.mailboxModel.GetByEmail(email)
}

// GetMailboxes 获取邮箱列表
func (s *MailStorage) GetMailboxes(email string) ([]string, error) {
	// 对于简单实现，只返回INBOX
	return []string{"INBOX"}, nil
}
