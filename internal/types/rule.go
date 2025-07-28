package types

import "time"

// VerificationRuleCreateReq 创建验证码规则请求
type VerificationRuleCreateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	Source      string `json:"source" binding:"required,max=50"` // 来源
	Pattern     string `json:"pattern" binding:"required"`       // 匹配模式
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// VerificationRuleUpdateReq 更新验证码规则请求
type VerificationRuleUpdateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	Source      string `json:"source" binding:"required,max=50"` // 来源
	Pattern     string `json:"pattern" binding:"required"`       // 匹配模式
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// VerificationRuleListReq 验证码规则列表请求
type VerificationRuleListReq struct {
	Name           string    `json:"name" form:"name"`                     // 规则名称（模糊搜索）
	Source         string    `json:"source" form:"source"`                 // 来源
	Status         *int      `json:"status" form:"status"`                 // 状态
	Priority       *int      `json:"priority" form:"priority"`             // 优先级
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// VerificationRuleResp 验证码规则响应
type VerificationRuleResp struct {
	Id          uint      `json:"id"`          // 规则ID
	Name        string    `json:"name"`        // 规则名称
	Source      string    `json:"source"`      // 来源
	Pattern     string    `json:"pattern"`     // 匹配模式
	Description string    `json:"description"` // 规则描述
	Priority    int       `json:"priority"`    // 优先级
	Status      int       `json:"status"`      // 状态
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// ForwardRuleCreateReq 创建转发规则请求
type ForwardRuleCreateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	FromPattern string `json:"fromPattern" binding:"required"`   // 发件人匹配模式
	ToEmail     string `json:"toEmail" binding:"required,email"` // 转发目标邮箱
	Conditions  string `json:"conditions"`                       // 转发条件（JSON格式）
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// ForwardRuleUpdateReq 更新转发规则请求
type ForwardRuleUpdateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	FromPattern string `json:"fromPattern" binding:"required"`   // 发件人匹配模式
	ToEmail     string `json:"toEmail" binding:"required,email"` // 转发目标邮箱
	Conditions  string `json:"conditions"`                       // 转发条件（JSON格式）
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// ForwardRuleListReq 转发规则列表请求
type ForwardRuleListReq struct {
	Name           string    `json:"name" form:"name"`                     // 规则名称（模糊搜索）
	FromPattern    string    `json:"fromPattern" form:"fromPattern"`       // 发件人匹配模式
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 转发目标邮箱
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// ForwardRuleResp 转发规则响应
type ForwardRuleResp struct {
	Id          uint      `json:"id"`          // 规则ID
	UserId      uint      `json:"userId"`      // 用户ID
	Name        string    `json:"name"`        // 规则名称
	FromPattern string    `json:"fromPattern"` // 发件人匹配模式
	ToEmail     string    `json:"toEmail"`     // 转发目标邮箱
	Conditions  string    `json:"conditions"`  // 转发条件
	Description string    `json:"description"` // 规则描述
	Priority    int       `json:"priority"`    // 优先级
	Status      int       `json:"status"`      // 状态
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// AntiSpamRuleCreateReq 创建反垃圾规则请求
type AntiSpamRuleCreateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	Type        string `json:"type" binding:"required,max=20"`   // 规则类型
	Pattern     string `json:"pattern" binding:"required"`       // 匹配模式
	Action      string `json:"action" binding:"required,max=20"` // 处理动作
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// AntiSpamRuleUpdateReq 更新反垃圾规则请求
type AntiSpamRuleUpdateReq struct {
	Name        string `json:"name" binding:"required,max=100"`  // 规则名称
	Type        string `json:"type" binding:"required,max=20"`   // 规则类型
	Pattern     string `json:"pattern" binding:"required"`       // 匹配模式
	Action      string `json:"action" binding:"required,max=20"` // 处理动作
	Description string `json:"description" binding:"max=500"`    // 规则描述
	Priority    int    `json:"priority" binding:"min=1,max=100"` // 优先级
	Status      int    `json:"status" binding:"oneof=0 1"`       // 状态：0禁用 1启用
}

// AntiSpamRuleListReq 反垃圾规则列表请求
type AntiSpamRuleListReq struct {
	Name           string    `json:"name" form:"name"`                     // 规则名称（模糊搜索）
	Type           string    `json:"type" form:"type"`                     // 规则类型
	Action         string    `json:"action" form:"action"`                 // 处理动作
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// AntiSpamRuleResp 反垃圾规则响应
type AntiSpamRuleResp struct {
	Id          uint      `json:"id"`          // 规则ID
	Name        string    `json:"name"`        // 规则名称
	Type        string    `json:"type"`        // 规则类型
	Pattern     string    `json:"pattern"`     // 匹配模式
	Action      string    `json:"action"`      // 处理动作
	Description string    `json:"description"` // 规则描述
	Priority    int       `json:"priority"`    // 优先级
	Status      int       `json:"status"`      // 状态
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}
