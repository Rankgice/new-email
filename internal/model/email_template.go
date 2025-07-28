package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailTemplate 邮件模板模型
type EmailTemplate struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Subject     string         `gorm:"size:500" json:"subject"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailTemplate) TableName() string { return "email_template" }

type EmailTemplateModel struct{ db *gorm.DB }

func NewEmailTemplateModel(db *gorm.DB) *EmailTemplateModel { return &EmailTemplateModel{db: db} }
