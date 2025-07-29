package model

import (
	"gorm.io/gorm"
	"time"
)

// UserVerificationRule 用户验证码规则关联模型
type UserVerificationRule struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"not null;index" json:"user_id"`
	RuleId    int64     `gorm:"not null;index" json:"rule_id"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserVerificationRule) TableName() string { return "user_verification_rule" }

type UserVerificationRuleModel struct{ db *gorm.DB }

func NewUserVerificationRuleModel(db *gorm.DB) *UserVerificationRuleModel {
	return &UserVerificationRuleModel{db: db}
}

// Create 创建用户验证码规则关联
func (m *UserVerificationRuleModel) Create(userRule *UserVerificationRule) error {
	return m.db.Create(userRule).Error
}

// Delete 删除用户验证码规则关联
func (m *UserVerificationRuleModel) Delete(userRule *UserVerificationRule) error {
	return m.db.Delete(userRule).Error
}

// GetById 根据ID获取用户验证码规则关联
func (m *UserVerificationRuleModel) GetById(id int64) (*UserVerificationRule, error) {
	var userRule UserVerificationRule
	if err := m.db.First(&userRule, id).Error; err != nil {
		return nil, err
	}
	return &userRule, nil
}

// List 获取用户验证码规则关联列表
func (m *UserVerificationRuleModel) List(params UserVerificationRuleListParams) ([]*UserVerificationRule, int64, error) {
	var userRules []*UserVerificationRule
	var total int64

	db := m.db.Model(&UserVerificationRule{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.RuleId != 0 {
		db = db.Where("rule_id = ?", params.RuleId)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Order("created_at DESC").Find(&userRules).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(userRules))
	}

	return userRules, total, nil
}

// GetByUserId 根据用户ID获取规则关联列表
func (m *UserVerificationRuleModel) GetByUserId(userId int64) ([]*UserVerificationRule, error) {
	userRules, _, err := m.List(UserVerificationRuleListParams{UserId: userId})
	return userRules, err
}

// GetByRuleId 根据规则ID获取用户关联列表
func (m *UserVerificationRuleModel) GetByRuleId(ruleId int64) ([]*UserVerificationRule, error) {
	userRules, _, err := m.List(UserVerificationRuleListParams{RuleId: ruleId})
	return userRules, err
}
