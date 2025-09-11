package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Folder 邮箱文件夹模型
type Folder struct {
	Id        int64          `gorm:"column:id;primaryKey;autoIncrement;comment:文件夹ID"`
	MailboxId int64          `gorm:"column:mailbox_id;type:bigint;not null;index:idx_mailbox_id_name_parent_id;comment:所属邮箱ID"`
	Name      string         `gorm:"column:name;type:varchar(255);not null;index:idx_mailbox_id_name_parent_id;comment:文件夹名称"`
	ParentId  *int64         `gorm:"column:parent_id;type:bigint;index:idx_mailbox_id_name_parent_id;comment:父文件夹ID"`
	IsSystem  bool           `gorm:"column:is_system;type:boolean;not null;default:false;comment:是否为系统预设文件夹"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:删除时间"`
}

// TableName Folder 表名
func (*Folder) TableName() string {
	return "folders"
}

// FolderModel 文件夹模型操作
type FolderModel struct {
	db *gorm.DB
}

// NewFolderModel 创建 FolderModel 实例
func NewFolderModel(db *gorm.DB) *FolderModel {
	return &FolderModel{
		db: db,
	}
}

// Create 创建新文件夹
func (m *FolderModel) Create(folder *Folder) error {
	return m.db.Create(folder).Error
}

// GetById 根据ID获取文件夹
func (m *FolderModel) GetById(id int64) (*Folder, error) {
	var folder Folder
	err := m.db.Where("id = ?", id).First(&folder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &folder, nil
}

// GetByMailboxIdAndName 根据邮箱ID和名称获取文件夹
func (m *FolderModel) GetByMailboxIdAndName(mailboxId int64, name string, parentId *int64) (*Folder, error) {
	var folder Folder
	query := m.db.Where("mailbox_id = ? AND name = ?", mailboxId, name)
	if parentId == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentId)
	}
	err := query.First(&folder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &folder, nil
}

// GetByMailboxId 获取邮箱下的所有文件夹
func (m *FolderModel) GetByMailboxId(mailboxId int64) ([]*Folder, error) {
	var folders []*Folder
	err := m.db.Where("mailbox_id = ?", mailboxId).Find(&folders).Error
	if err != nil {
		return nil, err
	}
	return folders, nil
}

// Update 更新文件夹
func (m *FolderModel) Update(folder *Folder) error {
	return m.db.Save(folder).Error
}

// Delete 删除文件夹
func (m *FolderModel) Delete(id int64) error {
	return m.db.Delete(&Folder{}, id).Error
}

// SoftDelete 软删除文件夹
func (m *FolderModel) SoftDelete(id int64) error {
	return m.db.Model(&Folder{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
