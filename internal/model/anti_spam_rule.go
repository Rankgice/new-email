package model

import (
	"gorm.io/gorm"
	"time"
)

// AntiSpamRule 反垃圾规则模型
type AntiSpamRule struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	RuleType    string         `gorm:"size:20;not null" json:"rule_type"`
	Pattern     string         `gorm:"type:text;not null" json:"pattern"`
	Action      string         `gorm:"size:20;not null" json:"action"`
	Description string         `gorm:"type:text" json:"description"`
	IsGlobal    bool           `gorm:"default:false" json:"is_global"`
	Status      int            `gorm:"default:1" json:"status"`
	Priority    int            `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AntiSpamRule) TableName() string { return "anti_spam_rule" }

type AntiSpamRuleModel struct{ db *gorm.DB }

func NewAntiSpamRuleModel(db *gorm.DB) *AntiSpamRuleModel { return &AntiSpamRuleModel{db: db} }
