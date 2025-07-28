package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailAttachment 邮件附件模型
type EmailAttachment struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailId   uint      `gorm:"not null;index" json:"email_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	MimeType  string    `gorm:"size:100" json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailAttachment) TableName() string { return "email_attachment" }

type EmailAttachmentModel struct{ db *gorm.DB }

func NewEmailAttachmentModel(db *gorm.DB) *EmailAttachmentModel { return &EmailAttachmentModel{db: db} }
