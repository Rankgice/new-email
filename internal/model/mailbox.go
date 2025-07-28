package model

import (
	"gorm.io/gorm"
	"time"
)

// Mailbox 邮箱模型
type Mailbox struct {
	Id           uint           `gorm:"primaryKey;autoIncrement" json:"id"`         // 邮箱ID
	UserId       uint           `gorm:"not null;index" json:"user_id"`              // 用户ID
	DomainId     uint           `gorm:"index" json:"domain_id"`                     // 域名ID（自建邮箱关联域名）
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"` // 邮箱地址
	Password     string         `gorm:"size:255;not null" json:"-"`                 // 邮箱密码（加密存储）
	Type         string         `gorm:"size:20;not null" json:"type"`               // 邮箱类型：self自建 third第三方
	Provider     string         `gorm:"size:50" json:"provider"`                    // 邮箱提供商：gmail outlook qq imap
	ImapHost     string         `gorm:"size:100" json:"imap_host"`                  // IMAP服务器地址
	ImapPort     int            `gorm:"default:993" json:"imap_port"`               // IMAP端口
	SmtpHost     string         `gorm:"size:100" json:"smtp_host"`                  // SMTP服务器地址
	SmtpPort     int            `gorm:"default:587" json:"smtp_port"`               // SMTP端口
	ClientId     string         `gorm:"size:255" json:"client_id"`                  // OAuth客户端ID
	RefreshToken string         `gorm:"size:500" json:"refresh_token"`              // OAuth刷新令牌
	Status       int            `gorm:"default:1" json:"status"`                    // 状态：1启用 0禁用
	AutoReceive  bool           `gorm:"default:true" json:"auto_receive"`           // 是否自动收信
	LastSyncAt   *time.Time     `json:"last_sync_at"`                               // 最后同步时间
	CreatedAt    time.Time      `json:"created_at"`                                 // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                                 // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                             // 软删除时间
}

// TableName 指定表名
func (Mailbox) TableName() string {
	return "mailbox"
}

// MailboxModel 邮箱模型
type MailboxModel struct {
	db *gorm.DB
}

// NewMailboxModel 创建邮箱模型
func NewMailboxModel(db *gorm.DB) *MailboxModel {
	return &MailboxModel{
		db: db,
	}
}

// Create 创建邮箱
func (m *MailboxModel) Create(mailbox *Mailbox) error {
	return m.db.Create(mailbox).Error
}

// GetById 根据ID获取邮箱
func (m *MailboxModel) GetById(id uint) (*Mailbox, error) {
	var mailbox Mailbox
	if err := m.db.First(&mailbox, id).Error; err != nil {
		return nil, err
	}
	return &mailbox, nil
}

// List 获取邮箱列表
func (m *MailboxModel) List(params MailboxListParams) ([]*Mailbox, int64, error) {
	var mailboxes []*Mailbox
	var total int64

	db := m.db.Model(&Mailbox{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.DomainId != 0 {
		db = db.Where("domain_id = ?", params.DomainId)
	}
	if params.Email != "" {
		db = db.Where("email LIKE ?", "%"+params.Email+"%")
	}
	if params.Type != "" {
		db = db.Where("type = ?", params.Type)
	}
	if params.Provider != "" {
		db = db.Where("provider = ?", params.Provider)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.AutoReceive != nil {
		db = db.Where("auto_receive = ?", *params.AutoReceive)
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

	if err := db.Order("created_at DESC").Find(&mailboxes).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(mailboxes))
	}

	return mailboxes, total, nil
}
