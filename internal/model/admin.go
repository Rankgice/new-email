package model

import (
	"errors"
	"new-email/internal/types"
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`           // 管理员ID
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"` // 管理员用户名
	Email       string         `gorm:"uniqueIndex;size:100;not null" json:"email"`   // 管理员邮箱
	Password    string         `gorm:"size:255;not null" json:"-"`                   // 密码（加密存储）
	Nickname    string         `gorm:"size:50" json:"nickname"`                      // 昵称
	Avatar      string         `gorm:"size:255" json:"avatar"`                       // 头像URL
	Role        string         `gorm:"size:20;default:admin" json:"role"`            // 角色：admin超级管理员 manager普通管理员
	Status      int            `gorm:"default:1" json:"status"`                      // 状态：1启用 0禁用
	LastLoginAt *time.Time     `json:"last_login_at"`                                // 最后登录时间
	CreatedAt   time.Time      `json:"created_at"`                                   // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                                   // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除时间
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admin"
}

// AdminModel 管理员模型
type AdminModel struct {
	db *gorm.DB
}

// NewAdminModel 创建管理员模型
func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{
		db: db,
	}
}

// Create 创建管理员
func (m *AdminModel) Create(admin *Admin) error {
	return m.db.Create(admin).Error
}

// Update 更新管理员
func (m *AdminModel) Update(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(admin).Error
}

// MapUpdate 更新管理员（使用map）
func (m *AdminModel) MapUpdate(tx *gorm.DB, id uint, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存管理员
func (m *AdminModel) Save(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(admin).Error
}

// Delete 删除管理员
func (m *AdminModel) Delete(admin *Admin) error {
	return m.db.Delete(admin).Error
}

// GetById 根据ID获取管理员
func (m *AdminModel) GetById(id uint) (*Admin, error) {
	var admin Admin
	if err := m.db.First(&admin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// GetByUsername 根据用户名获取管理员
func (m *AdminModel) GetByUsername(username string) (*Admin, error) {
	var admin Admin
	if err := m.db.Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// GetByEmail 根据邮箱获取管理员
func (m *AdminModel) GetByEmail(email string) (*Admin, error) {
	var admin Admin
	if err := m.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// List 获取管理员列表
func (m *AdminModel) List(params types.AdminListReq) ([]*Admin, int64, error) {
	var admins []*Admin
	var total int64

	db := m.db.Model(&Admin{})

	// 添加查询条件
	if params.Username != "" {
		db = db.Where("username LIKE ?", "%"+params.Username+"%")
	}
	if params.Email != "" {
		db = db.Where("email LIKE ?", "%"+params.Email+"%")
	}
	if params.Role != "" {
		db = db.Where("role = ?", params.Role)
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

	if err := db.Order("created_at DESC").Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(admins))
	}

	return admins, total, nil
}

// BatchDelete 批量删除管理员
func (m *AdminModel) BatchDelete(ids []uint) error {
	return m.db.Where("id IN ?", ids).Delete(&Admin{}).Error
}

// BatchUpdateStatus 批量更新管理员状态
func (m *AdminModel) BatchUpdateStatus(ids []uint, status int) error {
	return m.db.Model(&Admin{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetActiveAdmins 获取活跃管理员
func (m *AdminModel) GetActiveAdmins() ([]*Admin, error) {
	var admins []*Admin
	if err := m.db.Where("status = ?", 1).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// CountAdmins 统计管理员数量
func (m *AdminModel) CountAdmins() (int64, error) {
	var count int64
	if err := m.db.Model(&Admin{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountActiveAdmins 统计活跃管理员数量
func (m *AdminModel) CountActiveAdmins() (int64, error) {
	var count int64
	if err := m.db.Model(&Admin{}).Where("status = ?", 1).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateLastLogin 更新最后登录时间
func (m *AdminModel) UpdateLastLogin(id uint) error {
	return m.db.Model(&Admin{}).Where("id = ?", id).Update("last_login_at", time.Now()).Error
}

// CheckUsernameExists 检查用户名是否存在
func (m *AdminModel) CheckUsernameExists(username string, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&Admin{}).Where("username = ?", username)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否存在
func (m *AdminModel) CheckEmailExists(email string, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&Admin{}).Where("email = ?", email)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetSuperAdmins 获取超级管理员列表
func (m *AdminModel) GetSuperAdmins() ([]*Admin, error) {
	var admins []*Admin
	if err := m.db.Where("role = ? AND status = ?", "admin", 1).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// HasSuperAdmin 检查是否存在超级管理员
func (m *AdminModel) HasSuperAdmin() (bool, error) {
	var count int64
	if err := m.db.Model(&Admin{}).Where("role = ? AND status = ?", "admin", 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
