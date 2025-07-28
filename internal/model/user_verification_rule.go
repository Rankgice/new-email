package model

import (
	"gorm.io/gorm"
	"time"
)

// UserVerificationRule 用户验证码规则关联模型
type UserVerificationRule struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"not null;index" json:"user_id"`
	RuleId    uint      `gorm:"not null;index" json:"rule_id"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserVerificationRule) TableName() string { return "user_verification_rule" }

type UserVerificationRuleModel struct{ db *gorm.DB }

func NewUserVerificationRuleModel(db *gorm.DB) *UserVerificationRuleModel {
	return &UserVerificationRuleModel{db: db}
}
