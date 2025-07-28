package types

import "time"

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Password string `json:"password" binding:"required,min=6"`        // 密码
	Nickname string `json:"nickname"`                                 // 昵称
	Avatar   string `json:"avatar"`                                   // 头像
	Status   int    `json:"status" binding:"oneof=0 1"`               // 状态
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	Id       uint   `json:"id" binding:"required"`           // 用户ID
	Username string `json:"username" binding:"min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"email"`           // 邮箱
	Nickname string `json:"nickname"`                        // 昵称
	Avatar   string `json:"avatar"`                          // 头像
	Status   int    `json:"status" binding:"oneof=0 1"`      // 状态
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
	Id          uint       `json:"id"`                    // 用户ID
	Username    string     `json:"username"`              // 用户名
	Email       string     `json:"email"`                 // 邮箱
	Nickname    string     `json:"nickname"`              // 昵称
	Avatar      string     `json:"avatar"`                // 头像
	Status      int        `json:"status"`                // 状态
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"` // 最后登录时间
	CreatedAt   time.Time  `json:"createdAt"`             // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`             // 更新时间
}

// UserProfileReq 用户资料更新请求
type UserProfileReq struct {
	Nickname string `json:"nickname"`              // 昵称
	Avatar   string `json:"avatar"`                // 头像
	Email    string `json:"email" binding:"email"` // 邮箱
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Password string `json:"password" binding:"required,min=6"`        // 密码
	Nickname string `json:"nickname"`                                 // 昵称
	Code     string `json:"code"`                                     // 验证码（如果需要）
}

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"` // 密码
}

// UserLoginResp 用户登录响应
type UserLoginResp struct {
	Token        string    `json:"token"`        // 访问令牌
	RefreshToken string    `json:"refreshToken"` // 刷新令牌
	ExpiresAt    time.Time `json:"expiresAt"`    // 过期时间
	User         UserResp  `json:"user"`         // 用户信息
}

// UserStatsResp 用户统计响应
type UserStatsResp struct {
	TotalUsers  int64 `json:"totalUsers"`  // 总用户数
	ActiveUsers int64 `json:"activeUsers"` // 活跃用户数
	NewUsers    int64 `json:"newUsers"`    // 新增用户数（今日）
	OnlineUsers int64 `json:"onlineUsers"` // 在线用户数
}

// UserBatchCreateReq 批量创建用户请求
type UserBatchCreateReq struct {
	Users []UserCreateReq `json:"users" binding:"required,min=1"` // 用户列表
}

// UserBatchImportReq 批量导入用户请求
type UserBatchImportReq struct {
	File   string `json:"file" binding:"required"`   // 文件路径
	Format string `json:"format" binding:"required"` // 文件格式：csv, excel
}

// UserExportReq 用户导出请求
type UserExportReq struct {
	Format string   `json:"format" form:"format"` // 导出格式：csv, excel
	Fields []string `json:"fields"`               // 导出字段
	Ids    []uint   `json:"ids"`                  // 指定用户ID（可选）
	UserListReq
}

// UserSearchReq 用户搜索请求
type UserSearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"` // 搜索关键词
	Type    string `json:"type" form:"type"`                          // 搜索类型：username, email, nickname
	PageReq
}

// UserActivityReq 用户活动记录请求
type UserActivityReq struct {
	UserId uint   `json:"userId" form:"userId"` // 用户ID
	Action string `json:"action" form:"action"` // 操作类型
	TimeRangeReq
	PageReq
}

// UserActivityResp 用户活动记录响应
type UserActivityResp struct {
	Id        uint      `json:"id"`        // 记录ID
	UserId    uint      `json:"userId"`    // 用户ID
	Action    string    `json:"action"`    // 操作类型
	Resource  string    `json:"resource"`  // 操作资源
	IP        string    `json:"ip"`        // IP地址
	UserAgent string    `json:"userAgent"` // 用户代理
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

// UserPreferenceReq 用户偏好设置请求
type UserPreferenceReq struct {
	Language string            `json:"language"` // 语言
	Timezone string            `json:"timezone"` // 时区
	Theme    string            `json:"theme"`    // 主题
	Settings map[string]string `json:"settings"` // 其他设置
}

// UserPreferenceResp 用户偏好设置响应
type UserPreferenceResp struct {
	Language  string            `json:"language"`  // 语言
	Timezone  string            `json:"timezone"`  // 时区
	Theme     string            `json:"theme"`     // 主题
	Settings  map[string]string `json:"settings"`  // 其他设置
	UpdatedAt time.Time         `json:"updatedAt"` // 更新时间
}

// UserSecurityReq 用户安全设置请求
type UserSecurityReq struct {
	EnableTwoFactor bool   `json:"enableTwoFactor"` // 启用双因子认证
	TwoFactorSecret string `json:"twoFactorSecret"` // 双因子认证密钥
}

// UserSecurityResp 用户安全设置响应
type UserSecurityResp struct {
	EnableTwoFactor    bool      `json:"enableTwoFactor"`    // 启用双因子认证
	LastPasswordChange time.Time `json:"lastPasswordChange"` // 最后修改密码时间
	LoginAttempts      int       `json:"loginAttempts"`      // 登录尝试次数
	IsLocked           bool      `json:"isLocked"`           // 是否锁定
}

// UserNotificationReq 用户通知设置请求
type UserNotificationReq struct {
	EmailNotification bool `json:"emailNotification"` // 邮件通知
	SmsNotification   bool `json:"smsNotification"`   // 短信通知
	PushNotification  bool `json:"pushNotification"`  // 推送通知
}

// UserNotificationResp 用户通知设置响应
type UserNotificationResp struct {
	EmailNotification bool      `json:"emailNotification"` // 邮件通知
	SmsNotification   bool      `json:"smsNotification"`   // 短信通知
	PushNotification  bool      `json:"pushNotification"`  // 推送通知
	UpdatedAt         time.Time `json:"updatedAt"`         // 更新时间
}
