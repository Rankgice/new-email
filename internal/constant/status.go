package constant

// 通用状态常量
const (
	StatusEnabled  = 1 // 启用
	StatusDisabled = 2 // 禁用
)

// 用户状态
const (
	UserStatusActive   = 1 // 激活
	UserStatusInactive = 2 // 未激活
	UserStatusLocked   = 3 // 锁定
)

// 管理员角色
const (
	AdminRoleSuper   = "admin"   // 超级管理员
	AdminRoleManager = "manager" // 普通管理员
)

// 邮箱类型
const (
	MailboxTypeSelf  = "self"  // 自建邮箱
	MailboxTypeThird = "third" // 第三方邮箱
)

// 邮箱提供商
const (
	ProviderGmail   = "gmail"
	ProviderOutlook = "outlook"
	ProviderQQ      = "qq"
	ProviderIMAP    = "imap"
)

// 邮件方向
const (
	EmailDirectionSent     = "sent"     // 发送
	EmailDirectionReceived = "received" // 接收
)

// 邮件内容类型
const (
	ContentTypeHTML = "html" // HTML格式
	ContentTypeText = "text" // 纯文本格式
)

// 规则类型
const (
	RuleTypeBlacklist = "blacklist" // 黑名单
	RuleTypeWhitelist = "whitelist" // 白名单
	RuleTypeKeyword   = "keyword"   // 关键词
	RuleTypeRegex     = "regex"     // 正则表达式
)

// 反垃圾规则动作
const (
	ActionBlock      = "block"      // 阻止
	ActionQuarantine = "quarantine" // 隔离
	ActionMark       = "mark"       // 标记
)

// 草稿状态
const (
	DraftStatusDraft     = "draft"     // 草稿
	DraftStatusScheduled = "scheduled" // 定时发送
	DraftStatusSent      = "sent"      // 已发送
)

// 日志类型
const (
	LogTypeSend    = "send"    // 发送
	LogTypeReceive = "receive" // 接收
	LogTypeSync    = "sync"    // 同步
)

// 日志状态
const (
	LogStatusSuccess = "success" // 成功
	LogStatusFailed  = "failed"  // 失败
	LogStatusPending = "pending" // 进行中
)

// 验证码使用状态
const (
	CodeStatusUnused = false // 未使用
	CodeStatusUsed   = true  // 已使用
)

// API权限
const (
	PermissionEmailRead  = "email:read"  // 邮件读取
	PermissionEmailWrite = "email:write" // 邮件写入
	PermissionCodeRead   = "code:read"   // 验证码读取
	PermissionUserRead   = "user:read"   // 用户读取
	PermissionUserWrite  = "user:write"  // 用户写入
)

// 分页默认值
const (
	DefaultPage     = 1   // 默认页码
	DefaultPageSize = 20  // 默认每页数量
	MaxPageSize     = 100 // 最大每页数量
)

// 文件大小限制
const (
	MaxAttachmentSize = 10 * 1024 * 1024 // 10MB
	MaxAvatarSize     = 2 * 1024 * 1024  // 2MB
)

// 缓存键前缀
const (
	CacheKeyUser    = "user:"
	CacheKeyAdmin   = "admin:"
	CacheKeyMailbox = "mailbox:"
	CacheKeyEmail   = "email:"
	CacheKeyRule    = "rule:"
	CacheKeyAPIKey  = "apikey:"
)

// 队列名称
const (
	QueueEmailSend    = "email:send"
	QueueEmailReceive = "email:receive"
	QueueEmailSync    = "email:sync"
)

// 默认值
const (
	DefaultSMTPPort = 587
	DefaultIMAPPort = 993
	DefaultTimeout  = 30 // 秒
)
