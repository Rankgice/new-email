package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailLog 邮件日志模型
type EmailLog struct {
	Id         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MailboxId  uint       `gorm:"not null;index" json:"mailbox_id"`
	EmailId    uint       `gorm:"index" json:"email_id"`
	Type       string     `gorm:"size:20;not null" json:"type"`
	Status     string     `gorm:"size:20;not null" json:"status"`
	Message    string     `gorm:"type:text" json:"message"`
	ErrorCode  string     `gorm:"size:50" json:"error_code"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (EmailLog) TableName() string { return "email_log" }

type EmailLogModel struct{ db *gorm.DB }

func NewEmailLogModel(db *gorm.DB) *EmailLogModel { return &EmailLogModel{db: db} }
