package model

import (
	"errors"
	"new-email/internal/types"
	"time"

	"gorm.io/gorm"
)

// VerificationRule 验证码规则模型
type VerificationRule struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"` // 规则ID
	UserId      uint           `gorm:"index" json:"user_id"`               // 创建人ID，0表示公共规则
	Name        string         `gorm:"size:100;not null" json:"name"`      // 规则名称
	Pattern     string         `gorm:"type:text;not null" json:"pattern"`  // 正则表达式
	Description string         `gorm:"type:text" json:"description"`       // 规则描述
	IsGlobal    bool           `gorm:"default:false" json:"is_global"`     // 是否全局规则
	Status      int            `gorm:"default:1" json:"status"`            // 状态：1启用 0禁用
	Priority    int            `gorm:"default:0" json:"priority"`          // 优先级，数字越大优先级越高
	CreatedAt   time.Time      `json:"created_at"`                         // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                         // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                     // 软删除时间
}

// TableName 指定表名
func (VerificationRule) TableName() string {
	return "verification_rules"
}

// VerificationRuleModel 验证码规则模型
type VerificationRuleModel struct {
	db *gorm.DB
}

// NewVerificationRuleModel 创建验证码规则模型
func NewVerificationRuleModel(db *gorm.DB) *VerificationRuleModel {
	return &VerificationRuleModel{
		db: db,
	}
}

// Create 创建验证码规则
func (m *VerificationRuleModel) Create(rule *VerificationRule) error {
	return m.db.Create(rule).Error
}

// Update 更新验证码规则
func (m *VerificationRuleModel) Update(tx *gorm.DB, rule *VerificationRule) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(rule).Error
}

// MapUpdate 更新验证码规则（使用map）
func (m *VerificationRuleModel) MapUpdate(tx *gorm.DB, id uint, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&VerificationRule{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存验证码规则
func (m *VerificationRuleModel) Save(tx *gorm.DB, rule *VerificationRule) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(rule).Error
}

// Delete 删除验证码规则
func (m *VerificationRuleModel) Delete(rule *VerificationRule) error {
	return m.db.Delete(rule).Error
}

// GetById 根据ID获取验证码规则
func (m *VerificationRuleModel) GetById(id uint) (*VerificationRule, error) {
	var rule VerificationRule
	if err := m.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

// List 获取验证码规则列表
func (m *VerificationRuleModel) List(params types.VerificationRuleListReq) ([]*VerificationRule, int64, error) {
	var rules []*VerificationRule
	var total int64

	db := m.db.Model(&VerificationRule{})

	// 添加查询条件
	if params.UserId != nil {
		db = db.Where("user_id = ?", *params.UserId)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.IsGlobal != nil {
		db = db.Where("is_global = ?", *params.IsGlobal)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if !params.CreatedAtStart.IsZero() {
		db = db.Where("created_at >= ?", params.CreatedAtStart)
	}
	if !params.CreatedAtEnd.IsZero() {
		db = db.Where("created_at <= ?", params.CreatedAtEnd)
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

// BatchDelete 批量删除验证码规则
func (m *VerificationRuleModel) BatchDelete(ids []uint) error {
	return m.db.Where("id IN ?", ids).Delete(&VerificationRule{}).Error
}

// BatchUpdateStatus 批量更新验证码规则状态
func (m *VerificationRuleModel) BatchUpdateStatus(ids []uint, status int) error {
	return m.db.Model(&VerificationRule{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetActiveRules 获取活跃规则
func (m *VerificationRuleModel) GetActiveRules() ([]*VerificationRule, error) {
	var rules []*VerificationRule
	if err := m.db.Where("status = ?", 1).Order("priority DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetGlobalRules 获取全局规则
func (m *VerificationRuleModel) GetGlobalRules() ([]*VerificationRule, error) {
	var rules []*VerificationRule
	if err := m.db.Where("is_global = ? AND status = ?", true, 1).Order("priority DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetUserRules 获取用户规则
func (m *VerificationRuleModel) GetUserRules(userId uint) ([]*VerificationRule, error) {
	var rules []*VerificationRule
	if err := m.db.Where("user_id = ? AND status = ?", userId, 1).Order("priority DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetAvailableRules 获取用户可用的规则（全局规则 + 用户规则）
func (m *VerificationRuleModel) GetAvailableRules(userId uint) ([]*VerificationRule, error) {
	var rules []*VerificationRule
	if err := m.db.Where("(is_global = ? OR user_id = ?) AND status = ?", true, userId, 1).
		Order("priority DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// CountRules 统计规则数量
func (m *VerificationRuleModel) CountRules() (int64, error) {
	var count int64
	if err := m.db.Model(&VerificationRule{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountGlobalRules 统计全局规则数量
func (m *VerificationRuleModel) CountGlobalRules() (int64, error) {
	var count int64
	if err := m.db.Model(&VerificationRule{}).Where("is_global = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountUserRules 统计用户规则数量
func (m *VerificationRuleModel) CountUserRules(userId uint) (int64, error) {
	var count int64
	if err := m.db.Model(&VerificationRule{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CheckNameExists 检查规则名称是否存在
func (m *VerificationRuleModel) CheckNameExists(name string, userId uint, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&VerificationRule{}).Where("name = ? AND user_id = ?", name, userId)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdatePriority 更新规则优先级
func (m *VerificationRuleModel) UpdatePriority(id uint, priority int) error {
	return m.db.Model(&VerificationRule{}).Where("id = ?", id).Update("priority", priority).Error
}

// GetMaxPriority 获取最大优先级
func (m *VerificationRuleModel) GetMaxPriority(userId uint) (int, error) {
	var maxPriority int
	if err := m.db.Model(&VerificationRule{}).Where("user_id = ?", userId).
		Select("COALESCE(MAX(priority), 0)").Scan(&maxPriority).Error; err != nil {
		return 0, err
	}
	return maxPriority, nil
}
