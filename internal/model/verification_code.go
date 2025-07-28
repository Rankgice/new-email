package model

import (
	"gorm.io/gorm"
	"time"
)

// VerificationCode 验证码记录模型
type VerificationCode struct {
	Id          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailId     uint       `gorm:"not null;index" json:"email_id"`
	RuleId      uint       `gorm:"not null;index" json:"rule_id"`
	Code        string     `gorm:"size:50;not null" json:"code"`
	Source      string     `gorm:"size:100" json:"source"`
	ExtractedAt time.Time  `json:"extracted_at"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (VerificationCode) TableName() string { return "verification_code" }

type VerificationCodeModel struct{ db *gorm.DB }

func NewVerificationCodeModel(db *gorm.DB) *VerificationCodeModel {
	return &VerificationCodeModel{db: db}
}
