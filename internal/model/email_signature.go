package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailSignature 邮件签名模型
type EmailSignature struct {
	Id        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailSignature) TableName() string { return "email_signature" }

type EmailSignatureModel struct{ db *gorm.DB }

func NewEmailSignatureModel(db *gorm.DB) *EmailSignatureModel { return &EmailSignatureModel{db: db} }
