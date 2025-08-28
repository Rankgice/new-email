package model

import (
	"errors"
	"github.com/rankgice/new-email/internal/constant"
	"time"

	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	Id          int64          `gorm:"primaryKey;autoIncrement" json:"id"`        // 域名ID
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"` // 域名
	Status      int            `gorm:"default:1" json:"status"`                   // 状态：1启用 2禁用
	DnsVerified int            `gorm:"default:1" json:"dns_verified"`             // DNS验证状态：1未验证 2已验证
	DkimRecord  string         `gorm:"type:text" json:"dkim_record"`              // DKIM记录
	SpfRecord   string         `gorm:"type:text" json:"spf_record"`               // SPF记录
	DmarcRecord string         `gorm:"type:text" json:"dmarc_record"`             // DMARC记录
	CreatedAt   time.Time      `json:"created_at"`                                // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                                // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                            // 软删除时间
}

// TableName 指定表名
func (Domain) TableName() string {
	return "domain"
}

// DomainModel 域名模型
type DomainModel struct {
	db *gorm.DB
}

// NewDomainModel 创建域名模型
func NewDomainModel(db *gorm.DB) *DomainModel {
	return &DomainModel{
		db: db,
	}
}

// Create 创建域名
func (m *DomainModel) Create(domain *Domain) error {
	return m.db.Create(domain).Error
}

// Update 更新域名
func (m *DomainModel) Update(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(domain).Error
}

// MapUpdate 更新域名（使用map）
func (m *DomainModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存域名
func (m *DomainModel) Save(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(domain).Error
}

// Delete 删除域名
func (m *DomainModel) Delete(domain *Domain) error {
	return m.db.Delete(domain).Error
}

// GetById 根据ID获取域名
func (m *DomainModel) GetById(id int64) (*Domain, error) {
	var domain Domain
	if err := m.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}

// GetByName 根据域名获取记录
func (m *DomainModel) GetByName(name string) (*Domain, error) {
	var domain Domain
	if err := m.db.Where("name = ?", name).First(&domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}

// List 获取域名列表
func (m *DomainModel) List(params DomainListParams) ([]*Domain, int64, error) {
	var domains []*Domain
	var total int64

	db := m.db.Model(&Domain{})

	// 添加查询条件
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}
	if params.DnsVerified != nil {
		db = db.Where("dns_verified = ?", *params.DnsVerified)
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

	if err := db.Order("created_at DESC").Find(&domains).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(domains))
	}

	return domains, total, nil
}

// BatchDelete 批量删除域名
func (m *DomainModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&Domain{}).Error
}

// BatchUpdateStatus 批量更新域名状态
func (m *DomainModel) BatchUpdateStatus(ids []int64, status int) error {
	return m.db.Model(&Domain{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetActiveDomains 获取活跃域名
func (m *DomainModel) GetActiveDomains() ([]*Domain, error) {
	var domains []*Domain
	if err := m.db.Where("status = ?", constant.StatusEnabled).Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

// GetVerifiedDomains 获取已验证域名
func (m *DomainModel) GetVerifiedDomains() ([]*Domain, error) {
	var domains []*Domain
	if err := m.db.Where("status = ? AND dns_verified = ?", constant.StatusEnabled, constant.VerifyStatusVerified).Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

// CountDomains 统计域名数量
func (m *DomainModel) CountDomains() (int64, error) {
	var count int64
	if err := m.db.Model(&Domain{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountVerifiedDomains 统计已验证域名数量
func (m *DomainModel) CountVerifiedDomains() (int64, error) {
	var count int64
	if err := m.db.Model(&Domain{}).Where("dns_verified = ?", constant.VerifyStatusVerified).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CheckNameExists 检查域名是否存在
func (m *DomainModel) CheckNameExists(name string, excludeId ...int64) (bool, error) {
	var count int64
	db := m.db.Model(&Domain{}).Where("name = ?", name)

	if len(excludeId) > 0 {
		db = db.Where("id != ?", excludeId[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateDNSVerification 更新DNS验证状态
func (m *DomainModel) UpdateDNSVerification(id int64, verified bool) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Update("dns_verified", verified).Error
}

// UpdateDKIMRecord 更新DKIM记录
func (m *DomainModel) UpdateDKIMRecord(id int64, dkimRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Update("dkim_record", dkimRecord).Error
}

// UpdateSPFRecord 更新SPF记录
func (m *DomainModel) UpdateSPFRecord(id int64, spfRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Update("spf_record", spfRecord).Error
}

// UpdateDMARCRecord 更新DMARC记录
func (m *DomainModel) UpdateDMARCRecord(id int64, dmarcRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Update("dmarc_record", dmarcRecord).Error
}

// GetStatistics 获取域名统计信息
func (m *DomainModel) GetStatistics() (map[string]interface{}, error) {
	var total, verified, active int64

	// 总域名数
	if err := m.db.Model(&Domain{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 已验证域名数
	if err := m.db.Model(&Domain{}).Where("dns_verified = ?", constant.VerifyStatusVerified).Count(&verified).Error; err != nil {
		return nil, err
	}

	// 活跃域名数
	if err := m.db.Model(&Domain{}).Where("status = ?", constant.StatusEnabled).Count(&active).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":    total,
		"verified": verified,
		"active":   active,
	}, nil
}
