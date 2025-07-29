package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	Id        uint           `gorm:"primaryKey;autoIncrement" json:"id"`           // 用户ID
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"` // 用户名
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`   // 邮箱地址
	Password  string         `gorm:"size:255;not null" json:"-"`                   // 密码（加密存储）
	Nickname  string         `gorm:"size:50" json:"nickname"`                      // 昵称
	Avatar    string         `gorm:"size:255" json:"avatar"`                       // 头像URL
	Status    int            `gorm:"default:1" json:"status"`                      // 状态：1启用 0禁用
	CreatedAt time.Time      `json:"created_at"`                                   // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                   // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除时间
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// UserModel 用户模型
type UserModel struct {
	db *gorm.DB
}

// NewUserModel 创建用户模型
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

// Create 创建用户
func (m *UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}

// Update 更新用户
func (m *UserModel) Update(tx *gorm.DB, user *User) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(user).Error
}

// MapUpdate 更新用户（使用map）
func (m *UserModel) MapUpdate(tx *gorm.DB, id uint, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&User{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存用户
func (m *UserModel) Save(tx *gorm.DB, user *User) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(user).Error
}

// Delete 删除用户
func (m *UserModel) Delete(user *User) error {
	return m.db.Delete(user).Error
}

// GetById 根据ID获取用户
func (m *UserModel) GetById(id uint) (*User, error) {
	var user User
	if err := m.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	if err := m.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (m *UserModel) List(params UserListParams) ([]*User, int64, error) {
	var users []*User
	var total int64

	db := m.db.Model(&User{})

	// 添加查询条件
	if params.Username != "" {
		db = db.Where("username LIKE ?", "%"+params.Username+"%")
	}
	if params.Email != "" {
		db = db.Where("email LIKE ?", "%"+params.Email+"%")
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

	if err := db.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(users))
	}

	return users, total, nil
}

// BatchDelete 批量删除用户
func (m *UserModel) BatchDelete(ids []uint) error {
	return m.db.Where("id IN ?", ids).Delete(&User{}).Error
}

// BatchUpdateStatus 批量更新用户状态
func (m *UserModel) BatchUpdateStatus(ids []uint, status int) error {
	return m.db.Model(&User{}).Where("id IN ?", ids).Update("status", status).Error
}

// CheckUsernameExists 检查用户名是否存在
func (m *UserModel) CheckUsernameExists(username string, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&User{}).Where("username = ?", username)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否存在
func (m *UserModel) CheckEmailExists(email string, excludeId ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&User{}).Where("email = ?", email)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
