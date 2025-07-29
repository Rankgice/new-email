package types

import "time"

// SignatureCreateReq 创建签名请求
type SignatureCreateReq struct {
	Name      string `json:"name" binding:"required,max=100"`     // 签名名称
	Content   string `json:"content" binding:"required,max=5000"` // 签名内容
	IsDefault bool   `json:"isDefault"`                           // 是否默认签名
	Status    int    `json:"status" binding:"oneof=0 1"`          // 状态：0禁用 1启用
}

// SignatureUpdateReq 更新签名请求
type SignatureUpdateReq struct {
	Name      string `json:"name" binding:"required,max=100"`     // 签名名称
	Content   string `json:"content" binding:"required,max=5000"` // 签名内容
	IsDefault bool   `json:"isDefault"`                           // 是否默认签名
	Status    int    `json:"status" binding:"oneof=0 1"`          // 状态：0禁用 1启用
}

// SignatureListReq 签名列表请求
type SignatureListReq struct {
	Name           string    `json:"name" form:"name"`                     // 签名名称（模糊搜索）
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// SignatureResp 签名响应
type SignatureResp struct {
	Id        int64     `json:"id"`        // 签名ID
	UserId    int64     `json:"userId"`    // 用户ID
	Name      string    `json:"name"`      // 签名名称
	Content   string    `json:"content"`   // 签名内容
	IsDefault bool      `json:"isDefault"` // 是否默认签名
	Status    int       `json:"status"`    // 状态
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间
}
