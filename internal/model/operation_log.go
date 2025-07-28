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
