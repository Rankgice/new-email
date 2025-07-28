package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ApiKey API密钥模型
type ApiKey struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`      // 密钥ID
	UserId      uint           `gorm:"not null;index" json:"user_id"`           // 用户ID
	Name        string         `gorm:"size:100;not null" json:"name"`           // 密钥名称
	Key         string         `gorm:"uniqueIndex;size:64;not null" json:"key"` // API密钥
	Secret      string         `gorm:"size:128;not null" json:"-"`              // 密钥秘钥（加密存储）
	Permissions string         `gorm:"type:text" json:"permissions"`            // 权限列表（JSON格式）
	Status      int            `gorm:"default:1" json:"status"`                 // 状态：1启用 0禁用
	LastUsedAt  *time.Time     `json:"last_used_at"`                            // 最后使用时间
	ExpiresAt   *time.Time     `json:"expires_at"`                              // 过期时间
	CreatedAt   time.Time      `json:"created_at"`                              // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                              // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间
}

// TableName 指定表名
func (ApiKey) TableName() string {
	return "api_key"
}

// ApiKeyModel API密钥模型
type ApiKeyModel struct {
	db *gorm.DB
}

// NewApiKeyModel 创建API密钥模型
func NewApiKeyModel(db *gorm.DB) *ApiKeyModel {
	return &ApiKeyModel{
		db: db,
	}
}

// Create 创建API密钥
func (m *ApiKeyModel) Create(apiKey *ApiKey) error {
	return m.db.Create(apiKey).Error
}

// Update 更新API密钥
func (m *ApiKeyModel) Update(tx *gorm.DB, apiKey *ApiKey) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(apiKey).Error
}

// MapUpdate 更新API密钥（使用map）
func (m *ApiKeyModel) MapUpdate(tx *gorm.DB, id uint, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&ApiKey{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存API密钥
func (m *ApiKeyModel) Save(tx *gorm.DB, apiKey *ApiKey) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(apiKey).Error
}

// Delete 删除API密钥
func (m *ApiKeyModel) Delete(apiKey *ApiKey) error {
	return m.db.Delete(apiKey).Error
}

// GetById 根据ID获取API密钥
func (m *ApiKeyModel) GetById(id uint) (*ApiKey, error) {
	var apiKey ApiKey
	if err := m.db.First(&apiKey, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}

// GetByKey 根据密钥获取记录
func (m *ApiKeyModel) GetByKey(key string) (*ApiKey, error) {
	var apiKey ApiKey
	if err := m.db.Where("key = ?", key).First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}

// List 获取API密钥列表
func (m *ApiKeyModel) List(params ApiKeyListParams) ([]*ApiKey, int64, error) {
	var apiKeys []*ApiKey
	var total int64

	db := m.db.Model(&ApiKey{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Permissions != "" {
		db = db.Where("permissions LIKE ?", "%"+params.Permissions+"%")
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

	if err := db.Order("created_at DESC").Find(&apiKeys).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(apiKeys))
	}

	return apiKeys, total, nil
}

// GetByUserId 根据用户ID获取API密钥列表 (使用List方法替代)
func (m *ApiKeyModel) GetByUserId(userId uint) ([]*ApiKey, error) {
	apiKeys, _, err := m.List(ApiKeyListParams{UserId: userId})
	return apiKeys, err
}

// UpdateLastUsed 更新最后使用时间
func (m *ApiKeyModel) UpdateLastUsed(id uint) error {
	return m.db.Model(&ApiKey{}).Where("id = ?", id).Update("last_used_at", time.Now()).Error
}

// CheckKeyExists 检查密钥是否存在
func (m *ApiKeyModel) CheckKeyExists(key string, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&ApiKey{}).Where("key = ?", key)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// BatchDelete 批量删除API密钥
func (m *ApiKeyModel) BatchDelete(ids []uint) error {
	return m.db.Where("id IN ?", ids).Delete(&ApiKey{}).Error
}

// BatchUpdateStatus 批量更新API密钥状态
func (m *ApiKeyModel) BatchUpdateStatus(ids []uint, status int) error {
	return m.db.Model(&ApiKey{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetActiveKeys 获取活跃API密钥 (使用List方法替代)
func (m *ApiKeyModel) GetActiveKeys() ([]*ApiKey, error) {
	status := 1
	apiKeys, _, err := m.List(ApiKeyListParams{Status: &status})
	return apiKeys, err
}

// CountKeys 统计API密钥数量
func (m *ApiKeyModel) CountKeys() (int64, error) {
	var count int64
	if err := m.db.Model(&ApiKey{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountUserKeys 统计用户API密钥数量
func (m *ApiKeyModel) CountUserKeys(userId uint) (int64, error) {
	var count int64
	if err := m.db.Model(&ApiKey{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
