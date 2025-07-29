package types

import "time"

// ApiKeyCreateReq 创建API密钥请求
type ApiKeyCreateReq struct {
	Name        string     `json:"name" binding:"required,max=100"` // 密钥名称
	Permissions string     `json:"permissions" binding:"required"`  // 权限（JSON格式）
	Status      int        `json:"status" binding:"oneof=0 1"`      // 状态：0禁用 1启用
	ExpiresAt   *time.Time `json:"expiresAt"`                       // 过期时间
}

// ApiKeyUpdateReq 更新API密钥请求
type ApiKeyUpdateReq struct {
	Name        string     `json:"name" binding:"required,max=100"` // 密钥名称
	Permissions string     `json:"permissions" binding:"required"`  // 权限（JSON格式）
	Status      int        `json:"status" binding:"oneof=0 1"`      // 状态：0禁用 1启用
	ExpiresAt   *time.Time `json:"expiresAt"`                       // 过期时间
}

// ApiKeyListReq API密钥列表请求
type ApiKeyListReq struct {
	Name           string    `json:"name" form:"name"`                     // 密钥名称（模糊搜索）
	Status         *int      `json:"status" form:"status"`                 // 状态
	Permissions    string    `json:"permissions" form:"permissions"`       // 权限
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// ApiKeyResp API密钥响应
type ApiKeyResp struct {
	Id          int64      `json:"id"`          // 密钥ID
	UserId      int64      `json:"userId"`      // 用户ID
	Name        string     `json:"name"`        // 密钥名称
	Key         string     `json:"key"`         // 密钥（已脱敏）
	Permissions string     `json:"permissions"` // 权限
	Status      int        `json:"status"`      // 状态
	LastUsedAt  *time.Time `json:"lastUsedAt"`  // 最后使用时间
	ExpiresAt   *time.Time `json:"expiresAt"`   // 过期时间
	CreatedAt   time.Time  `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`   // 更新时间
}

// ApiKeyCreateResp 创建API密钥响应
type ApiKeyCreateResp struct {
	Id          int64      `json:"id"`          // 密钥ID
	UserId      int64      `json:"userId"`      // 用户ID
	Name        string     `json:"name"`        // 密钥名称
	Key         string     `json:"key"`         // 密钥（完整，仅创建时返回）
	Permissions string     `json:"permissions"` // 权限
	Status      int        `json:"status"`      // 状态
	ExpiresAt   *time.Time `json:"expiresAt"`   // 过期时间
	CreatedAt   time.Time  `json:"createdAt"`   // 创建时间
}
