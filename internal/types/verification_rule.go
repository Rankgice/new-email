package types

import "time"

// VerificationRuleCreateReq 创建验证码规则请求
type VerificationRuleCreateReq struct {
	Name        string `json:"name" binding:"required"`    // 规则名称
	Pattern     string `json:"pattern" binding:"required"` // 正则表达式
	Description string `json:"description"`                // 规则描述
	IsGlobal    bool   `json:"isGlobal"`                   // 是否全局规则
	Status      int    `json:"status" binding:"oneof=0 1"` // 状态
	Priority    int    `json:"priority"`                   // 优先级
}

// VerificationRuleUpdateReq 更新验证码规则请求
type VerificationRuleUpdateReq struct {
	Id          uint   `json:"id" binding:"required"`      // 规则ID
	Name        string `json:"name" binding:"required"`    // 规则名称
	Pattern     string `json:"pattern" binding:"required"` // 正则表达式
	Description string `json:"description"`                // 规则描述
	IsGlobal    bool   `json:"isGlobal"`                   // 是否全局规则
	Status      int    `json:"status" binding:"oneof=0 1"` // 状态
	Priority    int    `json:"priority"`                   // 优先级
}

// VerificationRuleListReq 验证码规则列表请求
type VerificationRuleListReq struct {
	UserId         *uint     `json:"userId" form:"userId"`                 // 用户ID
	Name           string    `json:"name" form:"name"`                     // 规则名称（模糊搜索）
	IsGlobal       *bool     `json:"isGlobal" form:"isGlobal"`             // 是否全局规则
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// VerificationRuleResp 验证码规则响应
type VerificationRuleResp struct {
	Id          uint      `json:"id"`          // 规则ID
	UserId      uint      `json:"userId"`      // 创建人ID
	Name        string    `json:"name"`        // 规则名称
	Pattern     string    `json:"pattern"`     // 正则表达式
	Description string    `json:"description"` // 规则描述
	IsGlobal    bool      `json:"isGlobal"`    // 是否全局规则
	Status      int       `json:"status"`      // 状态
	Priority    int       `json:"priority"`    // 优先级
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}
