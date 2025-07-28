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
	ImapSsl      bool           `gorm:"default:true" json:"imap_ssl"`               // IMAP是否使用SSL
	SmtpHost     string         `gorm:"size:100" json:"smtp_host"`                  // SMTP服务器地址
	SmtpPort     int            `gorm:"default:587" json:"smtp_port"`               // SMTP端口
	SmtpSsl      bool           `gorm:"default:true" json:"smtp_ssl"`               // SMTP是否使用SSL
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

// Update 更新邮箱
func (m *MailboxModel) Update(mailbox *Mailbox) error {
	return m.db.Updates(mailbox).Error
}

// Delete 删除邮箱
func (m *MailboxModel) Delete(mailbox *Mailbox) error {
	return m.db.Delete(mailbox).Error
}

// MapUpdate 根据条件更新邮箱
func (m *MailboxModel) MapUpdate(tx *gorm.DB, id uint, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Mailbox{}).Where("id = ?", id).Updates(data).Error
}

// CheckEmailExists 检查邮箱是否存在
func (m *MailboxModel) CheckEmailExists(email string, excludeIds ...uint) (bool, error) {
	var count int64
	db := m.db.Model(&Mailbox{}).Where("email = ?", email)

	// 排除指定的ID
	if len(excludeIds) > 0 {
		db = db.Where("id NOT IN ?", excludeIds)
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByUserId 根据用户ID获取邮箱列表
func (m *MailboxModel) GetByUserId(userId uint) ([]*Mailbox, error) {
	mailboxes, _, err := m.List(MailboxListParams{UserId: userId})
	return mailboxes, err
}

// GetByEmail 根据邮箱地址获取邮箱
func (m *MailboxModel) GetByEmail(email string) (*Mailbox, error) {
	var mailbox Mailbox
	if err := m.db.Where("email = ?", email).First(&mailbox).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mailbox, nil
}

// GetActiveMailboxes 获取活跃邮箱列表
func (m *MailboxModel) GetActiveMailboxes(userId uint) ([]*Mailbox, error) {
	status := 1
	mailboxes, _, err := m.List(MailboxListParams{UserId: userId, Status: &status})
	return mailboxes, err
}

// GetStatistics 获取邮箱统计信息
func (m *MailboxModel) GetStatistics() (map[string]interface{}, error) {
	var total, active, self, third int64

	// 总邮箱数
	if err := m.db.Model(&Mailbox{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 活跃邮箱数
	if err := m.db.Model(&Mailbox{}).Where("status = ?", 1).Count(&active).Error; err != nil {
		return nil, err
	}

	// 自建邮箱数
	if err := m.db.Model(&Mailbox{}).Where("type = ?", "self").Count(&self).Error; err != nil {
		return nil, err
	}

	// 第三方邮箱数
	if err := m.db.Model(&Mailbox{}).Where("type = ?", "third").Count(&third).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":  total,
		"active": active,
		"self":   self,
		"third":  third,
	}, nil
}
