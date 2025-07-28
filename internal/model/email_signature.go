package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailSignature 邮件签名模型
type EmailSignature struct {
	Id        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailSignature) TableName() string { return "email_signature" }

type EmailSignatureModel struct{ db *gorm.DB }

func NewEmailSignatureModel(db *gorm.DB) *EmailSignatureModel { return &EmailSignatureModel{db: db} }

// Create 创建邮件签名
func (m *EmailSignatureModel) Create(signature *EmailSignature) error {
	return m.db.Create(signature).Error
}

// Update 更新邮件签名
func (m *EmailSignatureModel) Update(signature *EmailSignature) error {
	return m.db.Updates(signature).Error
}

// Delete 删除邮件签名
func (m *EmailSignatureModel) Delete(signature *EmailSignature) error {
	return m.db.Delete(signature).Error
}

// GetById 根据ID获取邮件签名
func (m *EmailSignatureModel) GetById(id uint) (*EmailSignature, error) {
	var signature EmailSignature
	if err := m.db.First(&signature, id).Error; err != nil {
		return nil, err
	}
	return &signature, nil
}

// List 获取邮件签名列表
func (m *EmailSignatureModel) List(params EmailSignatureListParams) ([]*EmailSignature, int64, error) {
	var signatures []*EmailSignature
	var total int64

	db := m.db.Model(&EmailSignature{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.IsDefault != nil {
		db = db.Where("is_default = ?", *params.IsDefault)
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

	if err := db.Order("created_at DESC").Find(&signatures).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(signatures))
	}

	return signatures, total, nil
}

// GetByUserId 根据用户ID获取签名列表
func (m *EmailSignatureModel) GetByUserId(userId uint) ([]*EmailSignature, error) {
	signatures, _, err := m.List(EmailSignatureListParams{UserId: userId})
	return signatures, err
}

// GetDefaultSignature 获取默认签名
func (m *EmailSignatureModel) GetDefaultSignature(userId uint) (*EmailSignature, error) {
	isDefault := true
	signatures, _, err := m.List(EmailSignatureListParams{UserId: userId, IsDefault: &isDefault})
	if err != nil {
		return nil, err
	}
	if len(signatures) > 0 {
		return signatures[0], nil
	}
	return nil, nil
}

// ClearDefaultByUserId 清除用户的默认签名
func (m *EmailSignatureModel) ClearDefaultByUserId(userId uint) error {
	return m.db.Model(&EmailSignature{}).Where("user_id = ?", userId).Update("is_default", false).Error
}
