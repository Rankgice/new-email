package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailAttachment 邮件附件模型
type EmailAttachment struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailId   int64     `gorm:"not null;index" json:"email_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	MimeType  string    `gorm:"size:100" json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailAttachment) TableName() string { return "email_attachment" }

type EmailAttachmentModel struct{ db *gorm.DB }

func NewEmailAttachmentModel(db *gorm.DB) *EmailAttachmentModel { return &EmailAttachmentModel{db: db} }

// Create 创建邮件附件
func (m *EmailAttachmentModel) Create(attachment *EmailAttachment) error {
	return m.db.Create(attachment).Error
}

// GetById 根据ID获取邮件附件
func (m *EmailAttachmentModel) GetById(id int64) (*EmailAttachment, error) {
	var attachment EmailAttachment
	if err := m.db.First(&attachment, id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

// List 获取邮件附件列表
func (m *EmailAttachmentModel) List(params EmailAttachmentListParams) ([]*EmailAttachment, int64, error) {
	var attachments []*EmailAttachment
	var total int64

	db := m.db.Model(&EmailAttachment{})

	// 添加查询条件
	if params.EmailId != 0 {
		db = db.Where("email_id = ?", params.EmailId)
	}
	if params.Filename != "" {
		db = db.Where("filename LIKE ?", "%"+params.Filename+"%")
	}
	if params.MimeType != "" {
		db = db.Where("mime_type = ?", params.MimeType)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(attachments))
	}

	return attachments, total, nil
}

// Delete 删除邮件附件
func (m *EmailAttachmentModel) Delete(attachment *EmailAttachment) error {
	return m.db.Delete(attachment).Error
}

// GetByEmailId 根据邮件ID获取附件列表
func (m *EmailAttachmentModel) GetByEmailId(emailId int64) ([]*EmailAttachment, error) {
	attachments, _, err := m.List(EmailAttachmentListParams{EmailId: emailId})
	return attachments, err
}
