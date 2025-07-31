package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailTemplate 邮件模板模型
type EmailTemplate struct {
	Id          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64          `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Category    string         `gorm:"size:50;default:通用模板" json:"category"`
	Subject     string         `gorm:"size:500" json:"subject"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	Variables   string         `gorm:"type:text" json:"variables"`
	Description string         `gorm:"size:500" json:"description"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailTemplate) TableName() string { return "email_template" }

type EmailTemplateModel struct{ db *gorm.DB }

func NewEmailTemplateModel(db *gorm.DB) *EmailTemplateModel { return &EmailTemplateModel{db: db} }

// Create 创建邮件模板
func (m *EmailTemplateModel) Create(template *EmailTemplate) error {
	return m.db.Create(template).Error
}

// Update 更新邮件模板
func (m *EmailTemplateModel) Update(template *EmailTemplate) error {
	return m.db.Updates(template).Error
}

// Delete 删除邮件模板
func (m *EmailTemplateModel) Delete(template *EmailTemplate) error {
	return m.db.Delete(template).Error
}

// GetById 根据ID获取邮件模板
func (m *EmailTemplateModel) GetById(id int64) (*EmailTemplate, error) {
	var template EmailTemplate
	if err := m.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// List 获取邮件模板列表
func (m *EmailTemplateModel) List(params EmailTemplateListParams) ([]*EmailTemplate, int64, error) {
	var templates []*EmailTemplate
	var total int64

	db := m.db.Model(&EmailTemplate{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Category != "" {
		db = db.Where("category = ?", params.Category)
	}
	if params.Subject != "" {
		db = db.Where("subject LIKE ?", "%"+params.Subject+"%")
	}
	if params.ContentType != "" {
		db = db.Where("content_type = ?", params.ContentType)
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

	if err := db.Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(templates))
	}

	return templates, total, nil
}

// GetByUserId 根据用户ID获取模板列表
func (m *EmailTemplateModel) GetByUserId(userId int64) ([]*EmailTemplate, error) {
	templates, _, err := m.List(EmailTemplateListParams{UserId: userId})
	return templates, err
}

// GetDefaultTemplates 获取默认模板列表
func (m *EmailTemplateModel) GetDefaultTemplates(userId int64) ([]*EmailTemplate, error) {
	isDefault := true
	templates, _, err := m.List(EmailTemplateListParams{UserId: userId, IsDefault: &isDefault})
	return templates, err
}

// GetCategoriesByUserId 获取用户的模板分类列表
func (m *EmailTemplateModel) GetCategoriesByUserId(userId int64) ([]string, error) {
	var categories []string

	// 从数据库中查询用户的所有分类
	err := m.db.Model(&EmailTemplate{}).
		Where("user_id = ?", userId).
		Distinct("category").
		Pluck("category", &categories).Error

	if err != nil {
		return nil, err
	}

	// 如果没有分类，返回默认分类
	if len(categories) == 0 {
		categories = []string{
			"通用模板",
			"营销邮件",
			"通知邮件",
			"系统邮件",
		}
	}

	return categories, nil
}
