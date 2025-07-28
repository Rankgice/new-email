package model

import (
	"gorm.io/gorm"
	"time"
)

// Email 邮件模型
type Email struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`       // 邮件ID
	MailboxId   uint           `gorm:"not null;index" json:"mailbox_id"`         // 邮箱ID
	MessageId   string         `gorm:"size:255;index" json:"message_id"`         // 邮件消息ID
	Subject     string         `gorm:"size:500" json:"subject"`                  // 邮件主题
	FromEmail   string         `gorm:"size:100;index" json:"from_email"`         // 发件人邮箱
	FromName    string         `gorm:"size:100" json:"from_name"`                // 发件人姓名
	ToEmails    string         `gorm:"type:text" json:"to_emails"`               // 收件人列表（JSON格式）
	CcEmails    string         `gorm:"type:text" json:"cc_emails"`               // 抄送列表（JSON格式）
	BccEmails   string         `gorm:"type:text" json:"bcc_emails"`              // 密送列表（JSON格式）
	ReplyTo     string         `gorm:"size:100" json:"reply_to"`                 // 回复地址
	ContentType string         `gorm:"size:20;default:html" json:"content_type"` // 内容类型：html text
	Content     string         `gorm:"type:longtext" json:"content"`             // 邮件内容
	IsRead      bool           `gorm:"default:false" json:"is_read"`             // 是否已读
	IsStarred   bool           `gorm:"default:false" json:"is_starred"`          // 是否标星
	Direction   string         `gorm:"size:10;not null" json:"direction"`        // 方向：sent发送 received接收
	SentAt      *time.Time     `json:"sent_at"`                                  // 发送时间
	ReceivedAt  *time.Time     `json:"received_at"`                              // 接收时间
	CreatedAt   time.Time      `json:"created_at"`                               // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                               // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                           // 软删除时间
}

// TableName 指定表名
func (Email) TableName() string {
	return "email"
}

// EmailModel 邮件模型
type EmailModel struct {
	db *gorm.DB
}

// NewEmailModel 创建邮件模型
func NewEmailModel(db *gorm.DB) *EmailModel {
	return &EmailModel{
		db: db,
	}
}

// Create 创建邮件
func (m *EmailModel) Create(email *Email) error {
	return m.db.Create(email).Error
}

// GetById 根据ID获取邮件
func (m *EmailModel) GetById(id uint) (*Email, error) {
	var email Email
	if err := m.db.First(&email, id).Error; err != nil {
		return nil, err
	}
	return &email, nil
}
