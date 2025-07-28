package model

import (
	"gorm.io/gorm"
	"time"
)

// ForwardRule 转发规则模型
type ForwardRule struct {
	Id             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId         uint           `gorm:"not null;index" json:"user_id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	FromPattern    string         `gorm:"size:255" json:"from_pattern"`
	SubjectPattern string         `gorm:"size:255" json:"subject_pattern"`
	ContentPattern string         `gorm:"type:text" json:"content_pattern"`
	ForwardTo      string         `gorm:"size:255;not null" json:"forward_to"`
	KeepOriginal   bool           `gorm:"default:true" json:"keep_original"`
	Status         int            `gorm:"default:1" json:"status"`
	Priority       int            `gorm:"default:0" json:"priority"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ForwardRule) TableName() string { return "forward_rule" }

type ForwardRuleModel struct{ db *gorm.DB }

func NewForwardRuleModel(db *gorm.DB) *ForwardRuleModel { return &ForwardRuleModel{db: db} }
