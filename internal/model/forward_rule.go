package model

import (
	"gorm.io/gorm"
	"time"
)

// ForwardRule 转发规则模型
type ForwardRule struct {
	Id             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId         int64          `gorm:"not null;index" json:"user_id"`
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

// Create 创建转发规则
func (m *ForwardRuleModel) Create(rule *ForwardRule) error {
	return m.db.Create(rule).Error
}

// Update 更新转发规则
func (m *ForwardRuleModel) Update(rule *ForwardRule) error {
	return m.db.Updates(rule).Error
}

// Delete 删除转发规则
func (m *ForwardRuleModel) Delete(rule *ForwardRule) error {
	return m.db.Delete(rule).Error
}

// GetById 根据ID获取转发规则
func (m *ForwardRuleModel) GetById(id int64) (*ForwardRule, error) {
	var rule ForwardRule
	if err := m.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// List 获取转发规则列表
func (m *ForwardRuleModel) List(params ForwardRuleListParams) ([]*ForwardRule, int64, error) {
	var rules []*ForwardRule
	var total int64

	db := m.db.Model(&ForwardRule{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.FromPattern != "" {
		db = db.Where("from_pattern LIKE ?", "%"+params.FromPattern+"%")
	}
	if params.SubjectPattern != "" {
		db = db.Where("subject_pattern LIKE ?", "%"+params.SubjectPattern+"%")
	}
	if params.ForwardTo != "" {
		db = db.Where("forward_to LIKE ?", "%"+params.ForwardTo+"%")
	}
	if params.KeepOriginal != nil {
		db = db.Where("keep_original = ?", *params.KeepOriginal)
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

// GetByUserId 根据用户ID获取转发规则列表
func (m *ForwardRuleModel) GetByUserId(userId int64) ([]*ForwardRule, error) {
	rules, _, err := m.List(ForwardRuleListParams{UserId: userId})
	return rules, err
}

// GetActiveRules 获取活跃的转发规则
func (m *ForwardRuleModel) GetActiveRules(userId int64) ([]*ForwardRule, error) {
	status := 1
	rules, _, err := m.List(ForwardRuleListParams{UserId: userId, Status: &status})
	return rules, err
}
