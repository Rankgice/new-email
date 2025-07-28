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

// Create 创建反垃圾规则
func (m *AntiSpamRuleModel) Create(rule *AntiSpamRule) error {
	return m.db.Create(rule).Error
}

// Update 更新反垃圾规则
func (m *AntiSpamRuleModel) Update(rule *AntiSpamRule) error {
	return m.db.Updates(rule).Error
}

// Delete 删除反垃圾规则
func (m *AntiSpamRuleModel) Delete(rule *AntiSpamRule) error {
	return m.db.Delete(rule).Error
}

// GetById 根据ID获取反垃圾规则
func (m *AntiSpamRuleModel) GetById(id uint) (*AntiSpamRule, error) {
	var rule AntiSpamRule
	if err := m.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// List 获取反垃圾规则列表
func (m *AntiSpamRuleModel) List(params AntiSpamRuleListParams) ([]*AntiSpamRule, int64, error) {
	var rules []*AntiSpamRule
	var total int64

	db := m.db.Model(&AntiSpamRule{})

	// 添加查询条件
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.RuleType != "" {
		db = db.Where("rule_type = ?", params.RuleType)
	}
	if params.Pattern != "" {
		db = db.Where("pattern LIKE ?", "%"+params.Pattern+"%")
	}
	if params.Action != "" {
		db = db.Where("action = ?", params.Action)
	}
	if params.IsGlobal != nil {
		db = db.Where("is_global = ?", *params.IsGlobal)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.Priority != nil {
		db = db.Where("priority = ?", *params.Priority)
	}
	if !params.CreatedAtStart.IsZero() {
		db = db.Where("created_at >= ?", params.CreatedAtStart)
	}
	if !params.CreatedAtEnd.IsZero() {
		db = db.Where("created_at <= ?", params.CreatedAtEnd)
	}
	if !params.UpdatedAtStart.IsZero() {
		db = db.Where("updated_at >= ?", params.UpdatedAtStart)
	}
	if !params.UpdatedAtEnd.IsZero() {
		db = db.Where("updated_at <= ?", params.UpdatedAtEnd)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Order("priority DESC, created_at DESC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(rules))
	}

	return rules, total, nil
}

// GetGlobalRules 获取全局反垃圾规则
func (m *AntiSpamRuleModel) GetGlobalRules() ([]*AntiSpamRule, error) {
	isGlobal := true
	status := 1
	rules, _, err := m.List(AntiSpamRuleListParams{IsGlobal: &isGlobal, Status: &status})
	return rules, err
}

// GetActiveRules 获取活跃的反垃圾规则
func (m *AntiSpamRuleModel) GetActiveRules() ([]*AntiSpamRule, error) {
	status := 1
	rules, _, err := m.List(AntiSpamRuleListParams{Status: &status})
	return rules, err
}
