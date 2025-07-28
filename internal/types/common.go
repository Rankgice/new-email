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

// IdReq ID请求
type IdReq struct {
	Id uint `json:"id" form:"id" binding:"required"` // ID
}

// IdsReq 批量ID请求
type IdsReq struct {
	Ids []uint `json:"ids" binding:"required,min=1"` // ID列表
}

// StatusReq 状态请求
type StatusReq struct {
	Status int `json:"status" binding:"oneof=0 1"` // 状态：0禁用 1启用
}

// BatchStatusReq 批量状态请求
type BatchStatusReq struct {
	Ids    []uint `json:"ids" binding:"required,min=1"` // ID列表
	Status int    `json:"status" binding:"oneof=0 1"`   // 状态：0禁用 1启用
}

// TimeRangeReq 时间范围请求
type TimeRangeReq struct {
	StartTime time.Time `json:"startTime" form:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime" form:"endTime"`     // 结束时间
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResp 登录响应
type LoginResp struct {
	Token        string      `json:"token"`        // 访问令牌
	RefreshToken string      `json:"refreshToken"` // 刷新令牌
	ExpiresAt    time.Time   `json:"expiresAt"`    // 过期时间
	UserInfo     interface{} `json:"userInfo"`     // 用户信息
}

// RefreshTokenReq 刷新令牌请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required"` // 新密码
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱
	Code     string `json:"code" binding:"required"`        // 验证码
	Password string `json:"password" binding:"required"`    // 新密码
}

// UploadReq 上传请求
type UploadReq struct {
	Type string `json:"type" form:"type"` // 上传类型：avatar, attachment
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

// SearchReq 搜索请求
type SearchReq struct {
	Keyword  string `json:"keyword" form:"keyword"`   // 搜索关键词
	Category string `json:"category" form:"category"` // 搜索分类
	PageReq
}

// ExportReq 导出请求
type ExportReq struct {
	Format string   `json:"format" form:"format"` // 导出格式：csv, excel
	Fields []string `json:"fields"`               // 导出字段
	TimeRangeReq
}

// ImportReq 导入请求
type ImportReq struct {
	File   string `json:"file" binding:"required"`   // 文件路径
	Format string `json:"format" binding:"required"` // 文件格式
}

// ImportResp 导入响应
type ImportResp struct {
	Total   int      `json:"total"`            // 总数
	Success int      `json:"success"`          // 成功数
	Failed  int      `json:"failed"`           // 失败数
	Errors  []string `json:"errors,omitempty"` // 错误信息
}

// StatisticsReq 统计请求
type StatisticsReq struct {
	Type      string `json:"type" form:"type"`           // 统计类型
	Dimension string `json:"dimension" form:"dimension"` // 统计维度：day, week, month
	TimeRangeReq
}

// StatisticsResp 统计响应
type StatisticsResp struct {
	Labels []string      `json:"labels"` // 标签
	Data   []interface{} `json:"data"`   // 数据
}

// ConfigReq 配置请求
type ConfigReq struct {
	Key   string `json:"key" binding:"required"`   // 配置键
	Value string `json:"value" binding:"required"` // 配置值
}

// ConfigResp 配置响应
type ConfigResp struct {
	Key         string `json:"key"`         // 配置键
	Value       string `json:"value"`       // 配置值
	Description string `json:"description"` // 描述
	Type        string `json:"type"`        // 类型
}

// ErrorResp 错误响应
type ErrorResp struct {
	Code    int    `json:"code"`              // 错误码
	Message string `json:"message"`           // 错误信息
	Details string `json:"details,omitempty"` // 详细信息
}

// SuccessResp 成功响应
type SuccessResp struct {
	Message string      `json:"message"`        // 成功信息
	Data    interface{} `json:"data,omitempty"` // 数据
}

// ListResp 列表响应
type ListResp struct {
	List  interface{} `json:"list"`  // 列表数据
	Total int64       `json:"total"` // 总数
	PageReq
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
	CodeId    uint      `json:"codeId"`    // 验证码ID
}

// VerifyCodeReq 验证验证码请求
type VerifyCodeReq struct {
	Code   string `json:"code" binding:"required"` // 验证码
	Type   string `json:"type" binding:"required"` // 验证码类型
	Target string `json:"target"`                  // 目标（邮箱或手机号）
	CodeId uint   `json:"codeId"`                  // 验证码ID（可选）
}

// VerifyCodeResp 验证验证码响应
type VerifyCodeResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	CodeId  uint   `json:"codeId"`  // 验证码ID
}

// SystemInfoResp 系统信息响应
type SystemInfoResp struct {
	Version     string    `json:"version"`     // 系统版本
	Environment string    `json:"environment"` // 运行环境
	ServerTime  time.Time `json:"serverTime"`  // 服务器时间
	Timezone    string    `json:"timezone"`    // 时区
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
