package types

import "time"

// DomainCreateReq 创建域名请求
type DomainCreateReq struct {
	Name string `json:"name" binding:"required"` // 域名
}

// DomainUpdateReq 更新域名请求
type DomainUpdateReq struct {
	Id     uint   `json:"id" binding:"required"`      // 域名ID
	Name   string `json:"name" binding:"required"`    // 域名
	Status int    `json:"status" binding:"oneof=0 1"` // 状态
}

// DomainListReq 域名列表请求
type DomainListReq struct {
	Name           string    `json:"name" form:"name"`                     // 域名（模糊搜索）
	Status         *int      `json:"status" form:"status"`                 // 状态
	DnsVerified    *bool     `json:"dnsVerified" form:"dnsVerified"`       // DNS验证状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// DomainResp 域名响应
type DomainResp struct {
	Id          uint      `json:"id"`          // 域名ID
	Name        string    `json:"name"`        // 域名
	Status      int       `json:"status"`      // 状态
	DnsVerified bool      `json:"dnsVerified"` // DNS验证状态
	DkimRecord  string    `json:"dkimRecord"`  // DKIM记录
	SpfRecord   string    `json:"spfRecord"`   // SPF记录
	DmarcRecord string    `json:"dmarcRecord"` // DMARC记录
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// DomainBatchOperationReq 域名批量操作请求
type DomainBatchOperationReq struct {
	Ids       []uint `json:"ids" binding:"required,min=1"`                                    // 域名ID列表
	Operation string `json:"operation" binding:"required,oneof=enable disable delete verify"` // 操作类型
}
