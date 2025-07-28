package types

import "time"

// TemplateCreateReq 创建模板请求
type TemplateCreateReq struct {
	Name        string `json:"name" binding:"required,max=100"`    // 模板名称
	Category    string `json:"category" binding:"required,max=50"` // 模板分类
	Subject     string `json:"subject" binding:"required,max=200"` // 邮件主题
	Content     string `json:"content" binding:"required"`         // 模板内容
	Variables   string `json:"variables"`                          // 变量定义（JSON格式）
	Description string `json:"description" binding:"max=500"`      // 模板描述
	Status      int    `json:"status" binding:"oneof=0 1"`         // 状态：0禁用 1启用
}

// TemplateUpdateReq 更新模板请求
type TemplateUpdateReq struct {
	Name        string `json:"name" binding:"required,max=100"`    // 模板名称
	Category    string `json:"category" binding:"required,max=50"` // 模板分类
	Subject     string `json:"subject" binding:"required,max=200"` // 邮件主题
	Content     string `json:"content" binding:"required"`         // 模板内容
	Variables   string `json:"variables"`                          // 变量定义（JSON格式）
	Description string `json:"description" binding:"max=500"`      // 模板描述
	Status      int    `json:"status" binding:"oneof=0 1"`         // 状态：0禁用 1启用
}

// TemplateListReq 模板列表请求
type TemplateListReq struct {
	Name           string    `json:"name" form:"name"`                     // 模板名称（模糊搜索）
	Category       string    `json:"category" form:"category"`             // 模板分类
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// TemplateResp 模板响应
type TemplateResp struct {
	Id          uint      `json:"id"`          // 模板ID
	UserId      uint      `json:"userId"`      // 用户ID
	Name        string    `json:"name"`        // 模板名称
	Category    string    `json:"category"`    // 模板分类
	Subject     string    `json:"subject"`     // 邮件主题
	Content     string    `json:"content"`     // 模板内容
	Variables   string    `json:"variables"`   // 变量定义
	Description string    `json:"description"` // 模板描述
	Status      int       `json:"status"`      // 状态
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// TemplateCopyReq 复制模板请求
type TemplateCopyReq struct {
	Name string `json:"name" binding:"required,max=100"` // 新模板名称
}

// TemplateCategoryResp 模板分类响应
type TemplateCategoryResp struct {
	Category string `json:"category"` // 分类名称
	Count    int64  `json:"count"`    // 模板数量
}
