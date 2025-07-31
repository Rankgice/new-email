package types

import "time"

// VerificationCodeListReq 验证码列表请求
type VerificationCodeListReq struct {
	EmailId        *int64    `json:"emailId" form:"emailId"`               // 邮件ID
	Code           string    `json:"code" form:"code"`                     // 验证码（模糊搜索）
	Source         string    `json:"source" form:"source"`                 // 来源
	IsUsed         *bool     `json:"isUsed" form:"isUsed"`                 // 是否已使用
	IsExpired      *bool     `json:"isExpired" form:"isExpired"`           // 是否已过期
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// VerificationCodeResp 验证码响应
type VerificationCodeResp struct {
	Id        int64      `json:"id"`        // 验证码ID
	UserId    int64      `json:"userId"`    // 用户ID
	EmailId   int64      `json:"emailId"`   // 邮件ID
	Code      string     `json:"code"`      // 验证码
	Source    string     `json:"source"`    // 来源
	IsUsed    bool       `json:"isUsed"`    // 是否已使用
	IsExpired bool       `json:"isExpired"` // 是否已过期
	UsedAt    *time.Time `json:"usedAt"`    // 使用时间
	ExpiresAt time.Time  `json:"expiresAt"` // 过期时间
	CreatedAt time.Time  `json:"createdAt"` // 创建时间
}

// VerificationCodeLatestReq 获取最新验证码请求
type VerificationCodeLatestReq struct {
	Source string `json:"source" form:"source" binding:"required"` // 来源
}

// VerificationCodeExtractReq 验证码提取请求
type VerificationCodeExtractReq struct {
	EmailId int64 `json:"emailId" binding:"required"` // 邮件ID
}

// VerificationCodeExtractResp 验证码提取响应
type VerificationCodeExtractResp struct {
	EmailId     int64                    `json:"emailId"`     // 邮件ID
	Subject     string                   `json:"subject"`     // 邮件主题
	FromEmail   string                   `json:"fromEmail"`   // 发件人
	ExtractedAt time.Time                `json:"extractedAt"` // 提取时间
	Codes       []VerificationCodeResult `json:"codes"`       // 提取到的验证码
}

// VerificationCodeResult 验证码提取结果
type VerificationCodeResult struct {
	Code        string `json:"code"`        // 验证码
	Type        string `json:"type"`        // 验证码类型
	Context     string `json:"context"`     // 上下文信息
	Confidence  int    `json:"confidence"`  // 置信度 (0-100)
	Position    int    `json:"position"`    // 在邮件中的位置
	Length      int    `json:"length"`      // 验证码长度
	Pattern     string `json:"pattern"`     // 匹配的模式
	Description string `json:"description"` // 描述信息
}

// VerificationCodeMarkUsedReq 标记验证码已使用请求
type VerificationCodeMarkUsedReq struct {
	Used bool `json:"used"` // 是否已使用
}

// VerificationCodeStatsResp 验证码统计响应
type VerificationCodeStatsResp struct {
	TotalCodes  int64                        `json:"totalCodes"`  // 总验证码数
	UsedCodes   int64                        `json:"usedCodes"`   // 已使用数
	UnusedCodes int64                        `json:"unusedCodes"` // 未使用数
	TodayCodes  int64                        `json:"todayCodes"`  // 今日新增
	TypeStats   []VerificationCodeTypeStat   `json:"typeStats"`   // 类型统计
	SourceStats []VerificationCodeSourceStat `json:"sourceStats"` // 来源统计
}

// VerificationCodeTypeStat 验证码类型统计
type VerificationCodeTypeStat struct {
	Type  string `json:"type"`  // 类型
	Count int64  `json:"count"` // 数量
}

// VerificationCodeSourceStat 验证码来源统计
type VerificationCodeSourceStat struct {
	FromEmail string `json:"fromEmail"` // 发件人
	Count     int64  `json:"count"`     // 数量
}

// VerificationCodeBatchExtractReq 批量提取验证码请求
type VerificationCodeBatchExtractReq struct {
	EmailIds []int64 `json:"emailIds" binding:"required,min=1"` // 邮件ID列表
}

// VerificationCodeBatchExtractResp 批量提取验证码响应
type VerificationCodeBatchExtractResp struct {
	TotalEmails     int                           `json:"totalEmails"`     // 总邮件数
	ProcessedEmails int                           `json:"processedEmails"` // 已处理邮件数
	ExtractedCodes  int                           `json:"extractedCodes"`  // 提取到的验证码数
	Results         []VerificationCodeExtractResp `json:"results"`         // 提取结果
	Errors          []string                      `json:"errors"`          // 错误信息
}
