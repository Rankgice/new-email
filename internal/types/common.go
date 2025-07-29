package types

import "time"

// PageReq 分页请求
type PageReq struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页数量
}

// PageResp 分页响应
type PageResp struct {
	List     interface{} `json:"list"`     // 数据列表
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页数量
	Total    int64       `json:"total"`    // 总数量
}

// TimeRangeReq 时间范围请求
type TimeRangeReq struct {
	StartTime time.Time `json:"startTime" form:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime" form:"endTime"`     // 结束时间
}

// UploadResp 上传响应
type UploadResp struct {
	Url      string `json:"url"`      // 文件URL
	Filename string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
	Type     string `json:"type"`     // 上传类型
}

// CaptchaResp 图形验证码响应
type CaptchaResp struct {
	CaptchaId string    `json:"captchaId"` // 验证码ID
	ImageData string    `json:"imageData"` // base64编码的图片数据
	ExpiresAt time.Time `json:"expiresAt"` // 过期时间
}

// CaptchaVerifyReq 图形验证码验证请求
type CaptchaVerifyReq struct {
	CaptchaId string `json:"captchaId" binding:"required"` // 验证码ID
	Code      string `json:"code" binding:"required"`      // 验证码
}

// CaptchaVerifyResp 图形验证码验证响应
type CaptchaVerifyResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}

// SendCodeReq 发送验证码请求
type SendCodeReq struct {
	Type          string `json:"type" binding:"required,oneof=email sms"` // 验证码类型：email, sms
	Target        string `json:"target" binding:"required"`               // 目标（邮箱或手机号）
	Length        int    `json:"length" binding:"min=4,max=8"`            // 验证码长度
	ExpireMinutes int    `json:"expireMinutes" binding:"min=1,max=60"`    // 过期时间（分钟）
}

// SendCodeResp 发送验证码响应
type SendCodeResp struct {
	Success   bool      `json:"success"`   // 是否成功
	Message   string    `json:"message"`   // 消息
	ExpiresAt time.Time `json:"expiresAt"` // 过期时间
	CodeId    int64     `json:"codeId"`    // 验证码ID
}

// VerifyCodeReq 验证验证码请求
type VerifyCodeReq struct {
	Code   string `json:"code" binding:"required"` // 验证码
	Type   string `json:"type" binding:"required"` // 验证码类型
	Target string `json:"target"`                  // 目标（邮箱或手机号）
	CodeId int64  `json:"codeId"`                  // 验证码ID（可选）
}

// VerifyCodeResp 验证验证码响应
type VerifyCodeResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	CodeId  int64  `json:"codeId"`  // 验证码ID
}

// HealthResp 健康检查响应
type HealthResp struct {
	Status   string            `json:"status"`   // 状态
	Version  string            `json:"version"`  // 版本
	Uptime   string            `json:"uptime"`   // 运行时间
	Services map[string]string `json:"services"` // 服务状态
}

// BatchOperationResp 批量操作响应
type BatchOperationResp struct {
	Total        int      `json:"total"`        // 总数
	SuccessCount int      `json:"successCount"` // 成功数
	FailCount    int      `json:"failCount"`    // 失败数
	Errors       []string `json:"errors"`       // 错误信息
}
