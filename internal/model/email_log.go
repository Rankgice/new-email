package model

import (
	"gorm.io/gorm"
	"time"
)

// EmailLog 邮件日志模型
type EmailLog struct {
	Id         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MailboxId  uint       `gorm:"not null;index" json:"mailbox_id"`
	EmailId    uint       `gorm:"index" json:"email_id"`
	Type       string     `gorm:"size:20;not null" json:"type"`
	Status     string     `gorm:"size:20;not null" json:"status"`
	Message    string     `gorm:"type:text" json:"message"`
	ErrorCode  string     `gorm:"size:50" json:"error_code"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (EmailLog) TableName() string { return "email_log" }

type EmailLogModel struct{ db *gorm.DB }

func NewEmailLogModel(db *gorm.DB) *EmailLogModel { return &EmailLogModel{db: db} }

// Create 创建邮件日志
func (m *EmailLogModel) Create(log *EmailLog) error {
	return m.db.Create(log).Error
}

// GetById 根据ID获取邮件日志
func (m *EmailLogModel) GetById(id uint) (*EmailLog, error) {
	var log EmailLog
	if err := m.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// List 获取邮件日志列表
func (m *EmailLogModel) List(params EmailLogListParams) ([]*EmailLog, int64, error) {
	var logs []*EmailLog
	var total int64

	db := m.db.Model(&EmailLog{})

	// 添加查询条件
	if params.MailboxId != 0 {
		db = db.Where("mailbox_id = ?", params.MailboxId)
	}
	if params.EmailId != 0 {
		db = db.Where("email_id = ?", params.EmailId)
	}
	if params.Type != "" {
		db = db.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}
	if params.ErrorCode != "" {
		db = db.Where("error_code = ?", params.ErrorCode)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(logs))
	}

	return logs, total, nil
}

// GetByMailboxId 根据邮箱ID获取邮件日志
func (m *EmailLogModel) GetByMailboxId(mailboxId uint) ([]*EmailLog, error) {
	logs, _, err := m.List(EmailLogListParams{MailboxId: mailboxId})
	return logs, err
}

// GetByEmailId 根据邮件ID获取邮件日志
func (m *EmailLogModel) GetByEmailId(emailId uint) ([]*EmailLog, error) {
	logs, _, err := m.List(EmailLogListParams{EmailId: emailId})
	return logs, err
}
