package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailDraft 草稿邮件模型
type EmailDraft struct {
	Id          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64          `gorm:"not null;index" json:"user_id"`
	MailboxId   int64          `gorm:"not null;index" json:"mailbox_id"`
	Subject     string         `gorm:"size:500" json:"subject"`
	ToEmails    []string       `gorm:"type:json;serializer:json" json:"to_emails"`
	CcEmails    []string       `gorm:"type:json;serializer:json" json:"cc_emails"`
	BccEmails   []string       `gorm:"type:json;serializer:json" json:"bcc_emails"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	Status      string         `gorm:"size:20;default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailDraft) TableName() string { return "email_draft" }

type EmailDraftModel struct{ db *gorm.DB }

func NewEmailDraftModel(db *gorm.DB) *EmailDraftModel { return &EmailDraftModel{db: db} }

// Create 创建草稿邮件
func (m *EmailDraftModel) Create(draft *EmailDraft) error {
	return m.db.Create(draft).Error
}

// Update 更新草稿邮件
func (m *EmailDraftModel) Update(draft *EmailDraft) error {
	return m.db.Updates(draft).Error
}

// Delete 删除草稿邮件
func (m *EmailDraftModel) Delete(draft *EmailDraft) error {
	return m.db.Delete(draft).Error
}

// GetById 根据ID获取草稿邮件
func (m *EmailDraftModel) GetById(id int64) (*EmailDraft, error) {
	var draft EmailDraft
	if err := m.db.First(&draft, id).Error; err != nil {
		return nil, err
	}
	return &draft, nil
}

// List 获取草稿邮件列表
func (m *EmailDraftModel) List(params EmailDraftListParams) ([]*EmailDraft, int64, error) {
	var drafts []*EmailDraft
	var total int64

	db := m.db.Model(&EmailDraft{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.MailboxId != 0 {
		db = db.Where("mailbox_id = ?", params.MailboxId)
	}
	if params.Subject != "" {
		db = db.Where("subject LIKE ?", "%"+params.Subject+"%")
	}
	if params.ToEmails != "" {
		db = db.Where("to_emails LIKE ?", "%"+params.ToEmails+"%")
	}
	if params.ContentType != "" {
		db = db.Where("content_type = ?", params.ContentType)
	}
	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
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

	if err := db.Order("updated_at DESC").Find(&drafts).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(drafts))
	}

	return drafts, total, nil
}

// GetByUserId 根据用户ID获取草稿邮件列表
func (m *EmailDraftModel) GetByUserId(userId int64) ([]*EmailDraft, error) {
	drafts, _, err := m.List(EmailDraftListParams{UserId: userId})
	return drafts, err
}

// GetByMailboxId 根据邮箱ID获取草稿邮件列表
func (m *EmailDraftModel) GetByMailboxId(mailboxId int64) ([]*EmailDraft, error) {
	drafts, _, err := m.List(EmailDraftListParams{MailboxId: mailboxId})
	return drafts, err
}

// GetDraftsByStatus 根据状态获取草稿邮件列表
func (m *EmailDraftModel) GetDraftsByStatus(status string) ([]*EmailDraft, error) {
	drafts, _, err := m.List(EmailDraftListParams{Status: status})
	return drafts, err
}
