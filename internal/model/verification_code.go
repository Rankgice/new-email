package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// VerificationCodeStatsResp 验证码统计响应
type VerificationCodeStatsResp struct {
	TotalCodes  int64                        `json:"totalCodes"`
	UsedCodes   int64                        `json:"usedCodes"`
	UnusedCodes int64                        `json:"unusedCodes"`
	TodayCodes  int64                        `json:"todayCodes"`
	TypeStats   []VerificationCodeTypeStat   `json:"typeStats"`
	SourceStats []VerificationCodeSourceStat `json:"sourceStats"`
}

// VerificationCodeTypeStat 验证码类型统计
type VerificationCodeTypeStat struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

// VerificationCodeSourceStat 验证码来源统计
type VerificationCodeSourceStat struct {
	FromEmail string `json:"fromEmail"`
	Count     int64  `json:"count"`
}

// VerificationCode 验证码记录模型
type VerificationCode struct {
	Id          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64      `gorm:"not null;index" json:"user_id"`
	EmailId     int64      `gorm:"not null;index" json:"email_id"`
	Code        string     `gorm:"size:50;not null" json:"code"`
	Source      string     `gorm:"size:100" json:"source"`
	Type        string     `gorm:"size:50" json:"type"`
	Context     string     `gorm:"type:text" json:"context"`
	Confidence  int        `gorm:"default:0" json:"confidence"`
	Pattern     string     `gorm:"size:100" json:"pattern"`
	Description string     `gorm:"size:200" json:"description"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
func (m *VerificationCodeModel) GetById(id int64) (*VerificationCode, error) {
	var code VerificationCode
	if err := m.db.First(&code, id).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

// GetByEmailAndCode 根据邮件ID和验证码获取记录
func (m *VerificationCodeModel) GetByEmailAndCode(emailId int64, code string) (*VerificationCode, error) {
	var verificationCode VerificationCode
	if err := m.db.Where("email_id = ? AND code = ?", emailId, code).First(&verificationCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &verificationCode, nil
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
func (m *VerificationCodeModel) GetByEmailId(emailId int64) ([]*VerificationCode, error) {
	codes, _, err := m.List(VerificationCodeListParams{EmailId: emailId})
	return codes, err
}

// GetByRuleId 根据规则ID获取验证码记录
func (m *VerificationCodeModel) GetByRuleId(ruleId int64) ([]*VerificationCode, error) {
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
func (m *VerificationCodeModel) MarkAsUsed(id int64) error {
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
func (m *VerificationCodeModel) MarkAsExpired(id int64) error {
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

// FindByCode 根据验证码查找记录
func (m *VerificationCodeModel) FindByCode(code string) (*VerificationCode, error) {
	var verificationCode VerificationCode
	if err := m.db.Where("code = ?", code).
		Order("created_at DESC").First(&verificationCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("验证码不存在")
		}
		return nil, err
	}
	return &verificationCode, nil
}

// GetStatistics 获取验证码统计信息（别名方法）
func (m *VerificationCodeModel) GetStatistics() (map[string]interface{}, error) {
	return m.GetGlobalStatistics()
}

// GetStats 获取用户验证码统计信息
func (m *VerificationCodeModel) GetStats(userId int64) (*VerificationCodeStatsResp, error) {
	var stats VerificationCodeStatsResp

	// 总验证码数
	m.db.Model(&VerificationCode{}).Where("user_id = ?", userId).Count(&stats.TotalCodes)

	// 已使用数
	m.db.Model(&VerificationCode{}).Where("user_id = ? AND is_used = ?", userId, true).Count(&stats.UsedCodes)

	// 未使用数
	stats.UnusedCodes = stats.TotalCodes - stats.UsedCodes

	// 今日新增
	today := time.Now().Format("2006-01-02")
	m.db.Model(&VerificationCode{}).Where("user_id = ? AND DATE(created_at) = ?", userId, today).Count(&stats.TodayCodes)

	// 类型统计
	var typeStats []VerificationCodeTypeStat
	m.db.Model(&VerificationCode{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ?", userId).
		Group("type").
		Scan(&typeStats)
	stats.TypeStats = typeStats

	// 来源统计
	var sourceStats []VerificationCodeSourceStat
	m.db.Model(&VerificationCode{}).
		Select("source as from_email, COUNT(*) as count").
		Where("user_id = ?", userId).
		Group("source").
		Order("count DESC").
		Limit(10).
		Scan(&sourceStats)
	stats.SourceStats = sourceStats

	return &stats, nil
}

// GetLatestBySource 获取指定来源的最新验证码
func (m *VerificationCodeModel) GetLatestBySource(userId int64, source string) (*VerificationCode, error) {
	var code VerificationCode
	if err := m.db.Where("user_id = ? AND source = ?", userId, source).
		Order("created_at DESC").First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &code, nil
}
