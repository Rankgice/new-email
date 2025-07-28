package types

import "time"

// VerificationCodeListReq 验证码列表请求
type VerificationCodeListReq struct {
	EmailId        *uint     `json:"emailId" form:"emailId"`               // 邮件ID
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
	Id        uint       `json:"id"`        // 验证码ID
	UserId    uint       `json:"userId"`    // 用户ID
	EmailId   uint       `json:"emailId"`   // 邮件ID
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

// VerificationCodeStatsResp 验证码统计响应
type VerificationCodeStatsResp struct {
	Total     int64 `json:"total"`     // 总数
	Used      int64 `json:"used"`      // 已使用数
	Unused    int64 `json:"unused"`    // 未使用数
	Expired   int64 `json:"expired"`   // 已过期数
	Today     int64 `json:"today"`     // 今日新增数
	ThisWeek  int64 `json:"thisWeek"`  // 本周新增数
	ThisMonth int64 `json:"thisMonth"` // 本月新增数
}
