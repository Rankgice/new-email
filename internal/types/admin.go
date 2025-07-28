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
	Id       uint   `json:"id" binding:"required"`              // 管理员ID
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
	Id          uint       `json:"id"`                    // 管理员ID
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
	Token        string    `json:"token"`        // 访问令牌
	RefreshToken string    `json:"refreshToken"` // 刷新令牌
	ExpiresAt    time.Time `json:"expiresAt"`    // 过期时间
	Admin        AdminResp `json:"admin"`        // 管理员信息
}

// AdminStatsResp 管理员统计响应
type AdminStatsResp struct {
	TotalAdmins  int64 `json:"totalAdmins"`  // 总管理员数
	ActiveAdmins int64 `json:"activeAdmins"` // 活跃管理员数
	SuperAdmins  int64 `json:"superAdmins"`  // 超级管理员数
	Managers     int64 `json:"managers"`     // 普通管理员数
}

// AdminPermissionReq 管理员权限设置请求
type AdminPermissionReq struct {
	AdminId     uint     `json:"adminId" binding:"required"`     // 管理员ID
	Permissions []string `json:"permissions" binding:"required"` // 权限列表
}

// AdminPermissionResp 管理员权限响应
type AdminPermissionResp struct {
	AdminId     uint      `json:"adminId"`     // 管理员ID
	Permissions []string  `json:"permissions"` // 权限列表
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// AdminRoleReq 管理员角色设置请求
type AdminRoleReq struct {
	AdminId uint   `json:"adminId" binding:"required"`                  // 管理员ID
	Role    string `json:"role" binding:"required,oneof=admin manager"` // 角色
}

// AdminActivityReq 管理员活动记录请求
type AdminActivityReq struct {
	AdminId uint   `json:"adminId" form:"adminId"` // 管理员ID
	Action  string `json:"action" form:"action"`   // 操作类型
	TimeRangeReq
	PageReq
}

// AdminActivityResp 管理员活动记录响应
type AdminActivityResp struct {
	Id        uint      `json:"id"`        // 记录ID
	AdminId   uint      `json:"adminId"`   // 管理员ID
	Action    string    `json:"action"`    // 操作类型
	Resource  string    `json:"resource"`  // 操作资源
	IP        string    `json:"ip"`        // IP地址
	UserAgent string    `json:"userAgent"` // 用户代理
	CreatedAt time.Time `json:"createdAt"` // 创建时间
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

// AdminDashboardResp 管理员仪表板响应
type AdminDashboardResp struct {
	UserStats   UserStatsResp   `json:"userStats"`   // 用户统计
	AdminStats  AdminStatsResp  `json:"adminStats"`  // 管理员统计
	EmailStats  EmailStatsResp  `json:"emailStats"`  // 邮件统计
	SystemStats SystemStatsResp `json:"systemStats"` // 系统统计
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
	Ids       []uint `json:"ids" binding:"required,min=1"`                             // ID列表
	Operation string `json:"operation" binding:"required,oneof=enable disable delete"` // 操作类型
}

// AdminExportReq 管理员导出请求
type AdminExportReq struct {
	Format string   `json:"format" form:"format"` // 导出格式：csv, excel
	Fields []string `json:"fields"`               // 导出字段
	Ids    []uint   `json:"ids"`                  // 指定管理员ID（可选）
	AdminListReq
}
