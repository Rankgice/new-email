package model

import (
	"errors"
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

// Create 创建验证码记录
func (m *VerificationCodeModel) Create(code *VerificationCode) error {
	return m.db.Create(code).Error
}

// Update 更新验证码记录
func (m *VerificationCodeModel) Update(code *VerificationCode) error {
	return m.db.Updates(code).Error
}

// GetById 根据ID获取验证码记录
func (m *VerificationCodeModel) GetById(id uint) (*VerificationCode, error) {
	var code VerificationCode
	if err := m.db.First(&code, id).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

// List 获取验证码记录列表
func (m *VerificationCodeModel) List(params VerificationCodeListParams) ([]*VerificationCode, int64, error) {
	var codes []*VerificationCode
	var total int64

	db := m.db.Model(&VerificationCode{})

	// 添加查询条件
	if params.EmailId != 0 {
		db = db.Where("email_id = ?", params.EmailId)
	}
	if params.RuleId != 0 {
		db = db.Where("rule_id = ?", params.RuleId)
	}
	if params.Code != "" {
		db = db.Where("code = ?", params.Code)
	}
	if params.Source != "" {
		db = db.Where("source = ?", params.Source)
	}
	if params.IsUsed != nil {
		db = db.Where("is_used = ?", *params.IsUsed)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Order("extracted_at DESC").Find(&codes).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(codes))
	}

	return codes, total, nil
}

// GetByEmailId 根据邮件ID获取验证码记录
func (m *VerificationCodeModel) GetByEmailId(emailId uint) ([]*VerificationCode, error) {
	codes, _, err := m.List(VerificationCodeListParams{EmailId: emailId})
	return codes, err
}

// GetByRuleId 根据规则ID获取验证码记录
func (m *VerificationCodeModel) GetByRuleId(ruleId uint) ([]*VerificationCode, error) {
	codes, _, err := m.List(VerificationCodeListParams{RuleId: ruleId})
	return codes, err
}

// GetUnusedCodes 获取未使用的验证码
func (m *VerificationCodeModel) GetUnusedCodes() ([]*VerificationCode, error) {
	isUsed := false
	codes, _, err := m.List(VerificationCodeListParams{IsUsed: &isUsed})
	return codes, err
}

// MarkAsUsed 标记验证码为已使用
func (m *VerificationCodeModel) MarkAsUsed(id uint) error {
	now := time.Now()
	return m.db.Model(&VerificationCode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_used": true,
		"used_at": &now,
	}).Error
}

// GetLatestByTargetAndType 根据目标和类型获取最新验证码
func (m *VerificationCodeModel) GetLatestByTargetAndType(target, codeType string) (*VerificationCode, error) {
	var code VerificationCode
	if err := m.db.Where("source = ? AND code = ?", codeType, target).
		Order("created_at DESC").First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &code, nil
}

// MarkAsExpired 标记验证码为已过期
func (m *VerificationCodeModel) MarkAsExpired(id uint) error {
	return m.db.Model(&VerificationCode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_used": true, // 过期的验证码也标记为已使用
		"used_at": time.Now(),
	}).Error
}

// GetGlobalStatistics 获取全局验证码统计
func (m *VerificationCodeModel) GetGlobalStatistics() (map[string]interface{}, error) {
	var total, used, unused int64

	// 总数
	if err := m.db.Model(&VerificationCode{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 已使用数
	if err := m.db.Model(&VerificationCode{}).Where("is_used = ?", true).Count(&used).Error; err != nil {
		return nil, err
	}

	unused = total - used

	// 今日新增
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	if err := m.db.Model(&VerificationCode{}).Where("created_at >= ?", today).Count(&todayCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":  total,
		"used":   used,
		"unused": unused,
		"today":  todayCount,
	}, nil
}

// GetLatestBySource 根据来源获取最新验证码
func (m *VerificationCodeModel) GetLatestBySource(source string) (*VerificationCode, error) {
	var code VerificationCode
	if err := m.db.Where("source = ?", source).
		Order("created_at DESC").First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &code, nil
}

// GetStatistics 获取验证码统计信息（别名方法）
func (m *VerificationCodeModel) GetStatistics() (map[string]interface{}, error) {
	return m.GetGlobalStatistics()
}
