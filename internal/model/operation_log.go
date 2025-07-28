package model

import (
	"gorm.io/gorm"
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	Id         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId     uint      `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:100;not null" json:"action"`
	Resource   string    `gorm:"size:100" json:"resource"`
	ResourceId uint      `json:"resource_id"`
	Method     string    `gorm:"size:10" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	Ip         string    `gorm:"size:45" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Request    string    `gorm:"type:text" json:"request"`
	Response   string    `gorm:"type:text" json:"response"`
	Status     int       `json:"status"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

func (OperationLog) TableName() string { return "operation_log" }

type OperationLogModel struct{ db *gorm.DB }

func NewOperationLogModel(db *gorm.DB) *OperationLogModel { return &OperationLogModel{db: db} }

// Create 创建操作日志
func (m *OperationLogModel) Create(log *OperationLog) error {
	return m.db.Create(log).Error
}

// GetById 根据ID获取操作日志
func (m *OperationLogModel) GetById(id uint) (*OperationLog, error) {
	var log OperationLog
	if err := m.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// List 获取操作日志列表
func (m *OperationLogModel) List(params OperationLogListParams) ([]*OperationLog, int64, error) {
	var logs []*OperationLog
	var total int64

	db := m.db.Model(&OperationLog{})

	// 添加查询条件
	if params.UserId != 0 {
		db = db.Where("user_id = ?", params.UserId)
	}
	if params.Action != "" {
		db = db.Where("action LIKE ?", "%"+params.Action+"%")
	}
	if params.Resource != "" {
		db = db.Where("resource = ?", params.Resource)
	}
	if params.ResourceId != 0 {
		db = db.Where("resource_id = ?", params.ResourceId)
	}
	if params.Method != "" {
		db = db.Where("method = ?", params.Method)
	}
	if params.Path != "" {
		db = db.Where("path LIKE ?", "%"+params.Path+"%")
	}
	if params.Ip != "" {
		db = db.Where("ip = ?", params.Ip)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
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

// GetByUserId 根据用户ID获取操作日志
func (m *OperationLogModel) GetByUserId(userId uint) ([]*OperationLog, error) {
	logs, _, err := m.List(OperationLogListParams{UserId: userId})
	return logs, err
}

// GetByResource 根据资源获取操作日志
func (m *OperationLogModel) GetByResource(resource string, resourceId uint) ([]*OperationLog, error) {
	logs, _, err := m.List(OperationLogListParams{Resource: resource, ResourceId: resourceId})
	return logs, err
}
