package types

import "time"

// MailboxCreateReq 创建邮箱请求
type MailboxCreateReq struct {
	DomainId    int64  `json:"domainId" binding:"required"`    // 域名ID（自建邮箱）
	Email       string `json:"email" binding:"required,email"` // 邮箱地址
	Password    string `json:"password" binding:"required"`    // 邮箱密码
	AutoReceive bool   `json:"autoReceive"`                    // 是否自动收信
	Status      int    `json:"status" binding:"oneof=0 1"`     // 状态
}

// MailboxUpdateReq 更新邮箱请求
type MailboxUpdateReq struct {
	Id          int64  `json:"id" binding:"required"`      // 邮箱ID
	DomainId    int64  `json:"domainId"`                   // 域名ID（自建邮箱）
	Email       string `json:"email" binding:"email"`      // 邮箱地址
	Password    string `json:"password"`                   // 邮箱密码
	AutoReceive bool   `json:"autoReceive"`                // 是否自动收信
	Status      int    `json:"status" binding:"oneof=0 1"` // 状态
}

// MailboxListReq 邮箱列表请求
type MailboxListReq struct {
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	DomainId    int64  `json:"domainId" form:"domainId"`       // 域名ID
	Email       string `json:"email" form:"email"`             // 邮箱地址（模糊搜索）
	Status      *int   `json:"status" form:"status"`           // 状态
	AutoReceive *bool  `json:"autoReceive" form:"autoReceive"` // 自动收信
	TimeRangeReq
	PageReq
}

// MailboxResp 邮箱响应
type MailboxResp struct {
	Id          int64      `json:"id"`                   // 邮箱ID
	UserId      int64      `json:"userId"`               // 用户ID
	DomainId    int64      `json:"domainId"`             // 域名ID
	Email       string     `json:"email"`                // 邮箱地址
	AutoReceive bool       `json:"autoReceive"`          // 是否自动收信
	Status      int        `json:"status"`               // 状态
	LastSyncAt  *time.Time `json:"lastSyncAt,omitempty"` // 最后同步时间
	CreatedAt   time.Time  `json:"createdAt"`            // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`            // 更新时间
}

// MailboxSyncReq 同步邮箱请求
type MailboxSyncReq struct {
	Id        int64 `json:"id" binding:"required"`    // 邮箱ID
	ForceSync bool  `json:"forceSync"`                // 是否强制同步
	SyncDays  int   `json:"syncDays" binding:"min=1"` // 同步天数
}

// MailboxSyncResp 同步邮箱响应
type MailboxSyncResp struct {
	Success    bool      `json:"success"`    // 同步是否成功
	Message    string    `json:"message"`    // 同步结果信息
	SyncCount  int       `json:"syncCount"`  // 同步邮件数量
	ErrorCount int       `json:"errorCount"` // 错误数量
	LastSyncAt time.Time `json:"lastSyncAt"` // 最后同步时间
}

// MailboxStatsResp 邮箱统计响应
type MailboxStatsResp struct {
	TotalMailboxes  int64 `json:"totalMailboxes"`  // 总邮箱数
	ActiveMailboxes int64 `json:"activeMailboxes"` // 活跃邮箱数
}
