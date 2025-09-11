package model

import (
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"gorm.io/gorm"
)

// Email 邮件模型
type Email struct {
	Id          int64          `gorm:"primaryKey;autoIncrement" json:"id"`                                         // 邮件ID
	UserId      int64          `gorm:"not null;index" json:"user_id"`                                              // 用户ID
	MailboxId   int64          `gorm:"not null;index" json:"mailbox_id"`                                           // 邮箱ID
	MessageId   string         `gorm:"size:255;index" json:"message_id"`                                           // 邮件消息ID
	Subject     string         `gorm:"size:500" json:"subject"`                                                    // 邮件主题
	FromEmail   string         `gorm:"size:100;index" json:"from_email"`                                           // 发件人邮箱
	FromName    string         `gorm:"size:100" json:"from_name"`                                                  // 发件人姓名
	ToEmails    []string       `gorm:"type:json;serializer:json" json:"to_emails"`                                 // 收件人列表（JSON格式）
	CcEmails    []string       `gorm:"type:json;serializer:json" json:"cc_emails"`                                 // 抄送列表（JSON格式）
	BccEmails   []string       `gorm:"type:json;serializer:json" json:"bcc_emails"`                                // 密送列表（JSON格式）
	ReplyTo     string         `gorm:"size:100" json:"reply_to"`                                                   // 回复地址
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`                                   // 内容类型：html text
	Content     string         `gorm:"type:longtext" json:"content"`                                               // 邮件内容
	IsRead      bool           `gorm:"default:false" json:"is_read"`                                               // 是否已读
	IsStarred   bool           `gorm:"default:false" json:"is_starred"`                                            // 是否标星
	Direction   string         `gorm:"size:10;not null" json:"direction"`                                          // 方向：sent发送 received接收
	FolderId    int64          `gorm:"column:folder_id;type:bigint;not null;index:idx_folder_id" json:"folder_id"` // 文件夹ID
	SentAt      *time.Time     `json:"sent_at"`                                                                    // 发送时间
	ReceivedAt  *time.Time     `json:"received_at"`                                                                // 接收时间
	CreatedAt   time.Time      `json:"created_at"`                                                                 // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                                                                 // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                                             // 软删除时间
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
func (m *EmailModel) GetById(id int64) (*Email, error) {
	var email Email
	if err := m.db.First(&email, id).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

// List 获取邮件列表
func (m *EmailModel) List(params EmailListParams) ([]*Email, int64, error) {
	var emails []*Email
	var total int64

	db := m.db.Model(&Email{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.MailboxId != 0 {
		db = db.Where("mailbox_id = ?", params.MailboxId)
	}
	if params.MessageId != "" {
		db = db.Where("message_id = ?", params.MessageId)
	}
	if params.Subject != "" {
		db = db.Where("subject LIKE ?", "%"+params.Subject+"%")
	}
	if params.FromEmail != "" {
		db = db.Where("from_email LIKE ?", "%"+params.FromEmail+"%")
	}
	if params.ToEmails != "" {
		db = db.Where("to_emails LIKE ?", "%"+params.ToEmails+"%")
	}
	if params.Direction != "" {
		db = db.Where("direction = ?", params.Direction)
	}
	if params.IsRead != nil {
		db = db.Where("is_read = ?", *params.IsRead)
	}
	if params.IsStarred != nil {
		db = db.Where("is_starred = ?", *params.IsStarred)
	}
	if params.ContentType != "" {
		db = db.Where("content_type = ?", params.ContentType)
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

	if err := db.Order("created_at DESC").Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(emails))
	}

	return emails, total, nil
}

// Update 更新邮件
func (m *EmailModel) Update(email *Email) error {
	return m.db.Updates(email).Error
}

// MapUpdate 使用map更新邮件
func (m *EmailModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Email{}).Where("id = ?", id).Updates(data).Error
}

// Delete 删除邮件
func (m *EmailModel) Delete(email *Email) error {
	return m.db.Delete(email).Error
}

// BatchDelete 批量删除邮件
func (m *EmailModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&Email{}).Error
}

// GetStatistics 获取邮件统计信息
func (m *EmailModel) GetStatistics() (map[string]interface{}, error) {
	var total, sent, received, today int64

	// 总邮件数
	if err := m.db.Model(&Email{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 发送邮件数
	if err := m.db.Model(&Email{}).Where("direction = ?", "sent").Count(&sent).Error; err != nil {
		return nil, err
	}

	// 接收邮件数
	if err := m.db.Model(&Email{}).Where("direction = ?", "received").Count(&received).Error; err != nil {
		return nil, err
	}

	// 今日邮件数
	todayStart := time.Now().Truncate(24 * time.Hour)
	if err := m.db.Model(&Email{}).Where("created_at >= ?", todayStart).Count(&today).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalEmails":    total,
		"sentEmails":     sent,
		"receivedEmails": received,
		"todayEmails":    today,
	}, nil
}

// MarkAsRead 标记邮件为已读
func (m *EmailModel) MarkAsRead(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Update("is_read", true).Error
}

// CountByMailboxId 根据邮箱ID统计邮件数量
func (m *EmailModel) CountByMailboxId(mailboxId int64) (int, error) {
	var count int64
	err := m.db.Model(&Email{}).Where("mailbox_id = ?", mailboxId).Count(&count).Error
	return int(count), err
}

// GetByFolderId 根据文件夹ID获取邮件列表
func (m *EmailModel) GetByFolderId(folderId int64, limit int) ([]*Email, error) {
	var emails []*Email
	query := m.db.Where("folder_id = ?", folderId).
		Order("received_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&emails).Error
	return emails, err
}

// GetByUserId 根据用户ID获取邮件列表
func (m *EmailModel) GetByUserId(userId int64, limit int) ([]*Email, error) {
	var emails []*Email
	query := m.db.Where("user_id = ?", userId).
		Order("received_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&emails).Error
	return emails, err
}

// MarkAsUnread 标记邮件为未读
func (m *EmailModel) MarkAsUnread(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Update("is_read", false).Error
}

// MarkAsStarred 标记邮件为星标
func (m *EmailModel) MarkAsStarred(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Update("is_starred", true).Error
}

// UnmarkAsStarred 取消邮件星标
func (m *EmailModel) UnmarkAsStarred(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Update("is_starred", false).Error
}

// Count 获取邮件总数
func (m *EmailModel) Count() (int64, error) {
	var count int64
	if err := m.db.Model(&Email{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByDate 获取指定日期的邮件数
func (m *EmailModel) CountByDate(date string) (int64, error) {
	var count int64
	if err := m.db.Model(&Email{}).Where("DATE(created_at) = ?", date).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByDirection 根据方向获取邮件数
func (m *EmailModel) CountByDirection(direction string) (int64, error) {
	var count int64
	if err := m.db.Model(&Email{}).Where("direction = ?", direction).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByUserId 获取用户的邮件数
func (m *EmailModel) CountByUserId(userId int64) (int64, error) {
	var count int64
	if err := m.db.Model(&Email{}).
		Where("user_id = ?", userId).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Search 根据IMAP搜索条件搜索邮件
func (m *EmailModel) Search(mailboxId int64, folderId int64, criteria *imap.SearchCriteria) ([]*Email, error) {
	db := m.db.Model(&Email{}).
		Where("mailbox_id = ?", mailboxId).
		Where("folder_id = ?", folderId)

	// 处理各种搜索条件
	if len(criteria.Header) > 0 {
		for key, values := range criteria.Header { // 遍历Header map
			for _, value := range values {
				switch strings.ToLower(key) {
				case "from":
					db = db.Where("from_email LIKE ?", "%"+value+"%")
				case "to":
					db = db.Where("to_emails LIKE ?", "%\""+value+"\"%")
				case "subject":
					db = db.Where("subject LIKE ?", "%"+value+"%")
				case "body":
					db = db.Where("content LIKE ?", "%"+value+"%")
				case "text":
					db = db.Where("subject LIKE ? OR content LIKE ? OR from_email LIKE ? OR to_emails LIKE ?",
						"%"+value+"%", "%"+value+"%", "%"+value+"%", "%\""+value+"\"%")
				}
			}
		}
	}

	// 处理 WithFlags (例如 \Seen, \Flagged)
	for _, flag := range criteria.WithFlags {
		switch flag {
		case imap.SeenFlag:
			db = db.Where("is_read = ?", true)
		case imap.FlaggedFlag:
			db = db.Where("is_starred = ?", true)
			// TODO: 其他标志，如 \Answered, \Draft, \Deleted, \Recent
		}
	}

	// 处理 WithoutFlags (例如 \Unseen, \Unflagged)
	for _, flag := range criteria.WithoutFlags {
		switch flag {
		case imap.SeenFlag: // \Unseen
			db = db.Where("is_read = ?", false)
		case imap.FlaggedFlag: // \Unflagged
			db = db.Where("is_starred = ?", false)
		}
	}

	if !criteria.Since.IsZero() {
		db = db.Where("received_at >= ?", criteria.Since)
	}
	if !criteria.Before.IsZero() {
		db = db.Where("received_at < ?", criteria.Before)
	}
	if criteria.Larger != 0 {
		db = db.Where("size > ?", criteria.Larger)
	}
	if criteria.Smaller != 0 {
		db = db.Where("size < ?", criteria.Smaller)
	}
	if criteria.Uid != nil && len(criteria.Uid.Set) > 0 {
		var uids []int64
		for _, seqRange := range criteria.Uid.Set { // 遍历 SeqRange
			for i := seqRange.Start; i <= seqRange.Stop; i++ {
				uids = append(uids, int64(i))
			}
		}
		if len(uids) > 0 {
			db = db.Where("id IN (?)", uids)
		}
	}
	// TODO: 更多IMAP搜索条件，如 ON, SENTBEFORE, SENTSINCE等

	var emails []*Email
	if err := db.Order("received_at DESC").Find(&emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}
