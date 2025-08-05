package types

import "time"

// AdminCreateReq 创建管理员请求
type AdminCreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`    // 用户名
	Email    string `json:"email" binding:"required,email"`              // 邮箱
	Password string `json:"password" binding:"required,min=6"`           // 密码
	Nickname string `json:"nickname"`                                    // 昵称
	Avatar   string `json:"avatar"`                                      // 头像
	Role     string `json:"role" binding:"required,oneof=admin manager"` // 角色
	Status   int    `json:"status" binding:"oneof=0 1"`                  // 状态
}

// AdminUpdateReq 更新管理员请求
type AdminUpdateReq struct {
	Id       int64  `json:"id" binding:"required"`              // 管理员ID
	Username string `json:"username" binding:"min=3,max=50"`    // 用户名
	Email    string `json:"email" binding:"email"`              // 邮箱
	Nickname string `json:"nickname"`                           // 昵称
	Avatar   string `json:"avatar"`                             // 头像
	Role     string `json:"role" binding:"oneof=admin manager"` // 角色
	Status   int    `json:"status" binding:"oneof=0 1"`         // 状态
}

// AdminListReq 管理员列表请求
type AdminListReq struct {
	Username       string    `json:"username" form:"username"`             // 用户名（模糊搜索）
	Email          string    `json:"email" form:"email"`                   // 邮箱（模糊搜索）
	Role           string    `json:"role" form:"role"`                     // 角色
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// AdminResp 管理员响应
type AdminResp struct {
	Id          int64      `json:"id"`                    // 管理员ID
	Username    string     `json:"username"`              // 用户名
	Email       string     `json:"email"`                 // 邮箱
	Nickname    string     `json:"nickname"`              // 昵称
	Avatar      string     `json:"avatar"`                // 头像
	Role        string     `json:"role"`                  // 角色
	Status      int        `json:"status"`                // 状态
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"` // 最后登录时间
	CreatedAt   time.Time  `json:"createdAt"`             // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`             // 更新时间
}

// AdminProfileReq 管理员资料更新请求
type AdminProfileReq struct {
	Nickname string `json:"nickname"`              // 昵称
	Avatar   string `json:"avatar"`                // 头像
	Email    string `json:"email" binding:"email"` // 邮箱
}

// AdminLoginReq 管理员登录请求
type AdminLoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"` // 密码
}

// AdminLoginResp 管理员登录响应
type AdminLoginResp struct {
	Token string    `json:"token"` // 访问令牌
	Admin AdminResp `json:"admin"` // 管理员信息
}

// AdminStatsResp 管理员统计响应
type AdminStatsResp struct {
	TotalAdmins  int64 `json:"totalAdmins"`  // 总管理员数
	ActiveAdmins int64 `json:"activeAdmins"` // 活跃管理员数
	SuperAdmins  int64 `json:"superAdmins"`  // 超级管理员数
	Managers     int64 `json:"managers"`     // 普通管理员数
}

// AdminSystemConfigReq 系统配置请求
type AdminSystemConfigReq struct {
	SiteName                  string        `json:"siteName"`                  // 网站名称
	SiteLogo                  string        `json:"siteLogo"`                  // 网站Logo
	SiteDescription           string        `json:"siteDescription"`           // 网站描述
	ContactEmail              string        `json:"contactEmail"`              // 联系邮箱
	DefaultSMTP               SMTPConfigReq `json:"defaultSMTP"`               // 默认SMTP配置
	RegistrationEnabled       bool          `json:"registrationEnabled"`       // 是否开放注册
	EmailVerificationRequired bool          `json:"emailVerificationRequired"` // 是否需要邮箱验证
}

// SMTPConfigReq SMTP配置请求
type SMTPConfigReq struct {
	Host     string `json:"host" binding:"required"`     // SMTP服务器
	Port     int    `json:"port" binding:"required"`     // 端口
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
	UseTLS   bool   `json:"useTLS"`                      // 是否使用TLS
}

// AdminSystemConfigResp 系统配置响应
type AdminSystemConfigResp struct {
	SiteName                  string         `json:"siteName"`                  // 网站名称
	SiteLogo                  string         `json:"siteLogo"`                  // 网站Logo
	SiteDescription           string         `json:"siteDescription"`           // 网站描述
	ContactEmail              string         `json:"contactEmail"`              // 联系邮箱
	DefaultSMTP               SMTPConfigResp `json:"defaultSMTP"`               // 默认SMTP配置
	RegistrationEnabled       bool           `json:"registrationEnabled"`       // 是否开放注册
	EmailVerificationRequired bool           `json:"emailVerificationRequired"` // 是否需要邮箱验证
	UpdatedAt                 time.Time      `json:"updatedAt"`                 // 更新时间
}

// SMTPConfigResp SMTP配置响应
type SMTPConfigResp struct {
	Host     string `json:"host"`     // SMTP服务器
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	UseTLS   bool   `json:"useTLS"`   // 是否使用TLS
}

// EmailStatsResp 邮件统计响应
type EmailStatsResp struct {
	TotalEmails    int64 `json:"totalEmails"`    // 总邮件数
	SentEmails     int64 `json:"sentEmails"`     // 发送邮件数
	ReceivedEmails int64 `json:"receivedEmails"` // 接收邮件数
	TodayEmails    int64 `json:"todayEmails"`    // 今日邮件数
}

// SystemStatsResp 系统统计响应
type SystemStatsResp struct {
	TotalMailboxes  int64 `json:"totalMailboxes"`  // 总邮箱数
	ActiveMailboxes int64 `json:"activeMailboxes"` // 活跃邮箱数
	TotalDomains    int64 `json:"totalDomains"`    // 总域名数
	VerifiedDomains int64 `json:"verifiedDomains"` // 已验证域名数
}

// AdminBatchOperationReq 批量操作请求
type AdminBatchOperationReq struct {
	Ids       []int64 `json:"ids" binding:"required,min=1"`                             // ID列表
	Operation string  `json:"operation" binding:"required,oneof=enable disable delete"` // 操作类型
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`  // 用户名
	Email    string `json:"email" binding:"required,email"`            // 邮箱
	Password string `json:"password" binding:"required,min=6,max=100"` // 密码
	Nickname string `json:"nickname" binding:"max=100"`                // 昵称
	Avatar   string `json:"avatar" binding:"max=500"`                  // 头像URL
	Status   int    `json:"status" binding:"oneof=0 1"`                // 状态：0禁用 1启用
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Password string `json:"password" binding:"max=100"`               // 密码（可选）
	Nickname string `json:"nickname" binding:"max=100"`               // 昵称
	Avatar   string `json:"avatar" binding:"max=500"`                 // 头像URL
	Status   int    `json:"status" binding:"oneof=0 1"`               // 状态：0禁用 1启用
}

// UserListReq 用户列表请求
type UserListReq struct {
	Username       string    `json:"username" form:"username"`             // 用户名（模糊搜索）
	Email          string    `json:"email" form:"email"`                   // 邮箱（模糊搜索）
	Status         *int      `json:"status" form:"status"`                 // 状态
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// UserResp 用户响应
type UserResp struct {
	Id          int64     `json:"id"`          // 用户ID
	Username    string    `json:"username"`    // 用户名
	Email       string    `json:"email"`       // 邮箱
	Nickname    string    `json:"nickname"`    // 昵称
	Avatar      string    `json:"avatar"`      // 头像URL
	Status      int       `json:"status"`      // 状态
	LastLoginAt time.Time `json:"lastLoginAt"` // 最后登录时间
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// UserStatsResp 用户统计响应
type UserStatsResp struct {
	Total     int64 `json:"total"`     // 总用户数
	Active    int64 `json:"active"`    // 活跃用户数
	Inactive  int64 `json:"inactive"`  // 非活跃用户数
	Today     int64 `json:"today"`     // 今日新增用户数
	ThisWeek  int64 `json:"thisWeek"`  // 本周新增用户数
	ThisMonth int64 `json:"thisMonth"` // 本月新增用户数
}

// UserBatchOperationReq 用户批量操作请求
type UserBatchOperationReq struct {
	Ids       []int64 `json:"ids" binding:"required,min=1"`                             // 用户ID列表
	Operation string  `json:"operation" binding:"required,oneof=enable disable delete"` // 操作类型
}

// ImportUsersResp 导入用户响应
type ImportUsersResp struct {
	Total        int      `json:"total"`        // 总数
	SuccessCount int      `json:"successCount"` // 成功数
	FailCount    int      `json:"failCount"`    // 失败数
	Errors       []string `json:"errors"`       // 错误信息
	Message      string   `json:"message"`      // 消息
}

// ExportUsersReq 导出用户请求
type ExportUsersReq struct {
	Username       string    `json:"username" form:"username"`             // 用户名（模糊搜索）
	Email          string    `json:"email" form:"email"`                   // 邮箱（模糊搜索）
	Status         *int      `json:"status" form:"status"`                 // 状态
	Format         string    `json:"format" form:"format"`                 // 导出格式：csv, excel
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
}

// ExportUsersResp 导出用户响应
type ExportUsersResp struct {
	Total    int64  `json:"total"`    // 导出总数
	Format   string `json:"format"`   // 导出格式
	Filename string `json:"filename"` // 文件名
	Message  string `json:"message"`  // 消息
}

// AdminInfo 管理员信息
type AdminInfo struct {
	Id       int64  `json:"id"`       // 管理员ID
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Email    string `json:"email"`    // 邮箱
	Role     string `json:"role"`     // 角色
	Status   int    `json:"status"`   // 状态
}

// AdminSystemStatsResp 管理员系统统计响应
type AdminSystemStatsResp struct {
	UserStats    AdminUserStats    `json:"userStats"`    // 用户统计
	EmailStats   AdminEmailStats   `json:"emailStats"`   // 邮件统计
	MailboxStats AdminMailboxStats `json:"mailboxStats"` // 邮箱统计
	SystemStats  AdminSystemInfo   `json:"systemStats"`  // 系统统计
}

// AdminUserStats 用户统计
type AdminUserStats struct {
	TotalUsers  int64 `json:"totalUsers"`  // 总用户数
	ActiveUsers int64 `json:"activeUsers"` // 活跃用户数
	NewUsers    int64 `json:"newUsers"`    // 新用户数（今日）
	OnlineUsers int64 `json:"onlineUsers"` // 在线用户数
}

// AdminEmailStats 邮件统计
type AdminEmailStats struct {
	TotalEmails    int64 `json:"totalEmails"`    // 总邮件数
	TodayEmails    int64 `json:"todayEmails"`    // 今日邮件数
	SentEmails     int64 `json:"sentEmails"`     // 发送邮件数
	ReceivedEmails int64 `json:"receivedEmails"` // 接收邮件数
}

// AdminMailboxStats 邮箱统计
type AdminMailboxStats struct {
	TotalMailboxes  int64 `json:"totalMailboxes"`  // 总邮箱数
	ActiveMailboxes int64 `json:"activeMailboxes"` // 活跃邮箱数
	ImapMailboxes   int64 `json:"imapMailboxes"`   // IMAP邮箱数
	Pop3Mailboxes   int64 `json:"pop3Mailboxes"`   // POP3邮箱数
}

// AdminSystemInfo 系统信息
type AdminSystemInfo struct {
	Version   string    `json:"version"`   // 系统版本
	StartTime time.Time `json:"startTime"` // 启动时间
	Uptime    string    `json:"uptime"`    // 运行时间
	GoVersion string    `json:"goVersion"` // Go版本
	Platform  string    `json:"platform"`  // 平台信息
	CPUUsage  float64   `json:"cpuUsage"`  // CPU使用率
	MemUsage  float64   `json:"memUsage"`  // 内存使用率
	DiskUsage float64   `json:"diskUsage"` // 磁盘使用率
}

// AdminDashboardResp 管理员仪表板响应
type AdminDashboardResp struct {
	Stats        AdminSystemStatsResp `json:"stats"`        // 统计信息
	RecentUsers  []UserResp           `json:"recentUsers"`  // 最近用户
	RecentLogs   []OperationLogResp   `json:"recentLogs"`   // 最近日志
	SystemAlerts []SystemAlert        `json:"systemAlerts"` // 系统警告
}

// SystemAlert 系统警告
type SystemAlert struct {
	Id        int64     `json:"id"`        // 警告ID
	Type      string    `json:"type"`      // 警告类型
	Level     string    `json:"level"`     // 警告级别
	Title     string    `json:"title"`     // 警告标题
	Message   string    `json:"message"`   // 警告消息
	Status    string    `json:"status"`    // 状态
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

// AdminSystemSettingsReq 管理员系统设置请求
type AdminSystemSettingsReq struct {
	SiteName        string `json:"siteName"`        // 网站名称
	SiteDescription string `json:"siteDescription"` // 网站描述
	SiteLogo        string `json:"siteLogo"`        // 网站Logo
	AllowRegister   bool   `json:"allowRegister"`   // 允许注册
	RequireInvite   bool   `json:"requireInvite"`   // 需要邀请码
	DefaultQuota    int64  `json:"defaultQuota"`    // 默认配额
	MaxMailboxes    int    `json:"maxMailboxes"`    // 最大邮箱数
}

// AdminSystemSettingsResp 管理员系统设置响应
type AdminSystemSettingsResp struct {
	SiteName        string    `json:"siteName"`        // 网站名称
	SiteDescription string    `json:"siteDescription"` // 网站描述
	SiteLogo        string    `json:"siteLogo"`        // 网站Logo
	AllowRegister   bool      `json:"allowRegister"`   // 允许注册
	RequireInvite   bool      `json:"requireInvite"`   // 需要邀请码
	DefaultQuota    int64     `json:"defaultQuota"`    // 默认配额
	MaxMailboxes    int       `json:"maxMailboxes"`    // 最大邮箱数
	UpdatedAt       time.Time `json:"updatedAt"`       // 更新时间
}

// DailyStats 每日统计
type DailyStats struct {
	Date  string `json:"date"`  // 日期
	Count int64  `json:"count"` // 数量
}

// UserActivityStats 用户活动统计
type UserActivityStats struct {
	UserId       int64      `json:"userId"`       // 用户ID
	Username     string     `json:"username"`     // 用户名
	EmailCount   int64      `json:"emailCount"`   // 邮件数量
	MailboxCount int64      `json:"mailboxCount"` // 邮箱数量
	LastLoginAt  *time.Time `json:"lastLoginAt"`  // 最后登录时间
}

// SystemHealth 系统健康状态
type SystemHealth struct {
	Status    string    `json:"status"`    // 状态
	Score     float64   `json:"score"`     // 健康分数
	CPUUsage  float64   `json:"cpuUsage"`  // CPU使用率
	MemUsage  float64   `json:"memUsage"`  // 内存使用率
	DiskUsage float64   `json:"diskUsage"` // 磁盘使用率
	Uptime    string    `json:"uptime"`    // 运行时间
	CheckedAt time.Time `json:"checkedAt"` // 检查时间
}
