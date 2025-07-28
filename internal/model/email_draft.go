package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailDraft 草稿邮件模型
type EmailDraft struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint           `gorm:"not null;index" json:"user_id"`
	MailboxId   uint           `gorm:"not null;index" json:"mailbox_id"`
	Subject     string         `gorm:"size:500" json:"subject"`
	ToEmails    string         `gorm:"type:text" json:"to_emails"`
	CcEmails    string         `gorm:"type:text" json:"cc_emails"`
	BccEmails   string         `gorm:"type:text" json:"bcc_emails"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	Status      string         `gorm:"size:20;default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailDraft) TableName() string { return "email_draft" }

type EmailDraftModel struct{ db *gorm.DB }

func NewEmailDraftModel(db *gorm.DB) *EmailDraftModel { return &EmailDraftModel{db: db} }
