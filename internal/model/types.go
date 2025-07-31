package model

import "time"

// BaseListParams 基础列表查询参数
type BaseListParams struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页数量
}

// BaseTimeRangeParams 基础时间范围参数
type BaseTimeRangeParams struct {
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	UpdatedAtStart time.Time `json:"updatedAtStart" form:"updatedAtStart"` // 更新时间开始
	UpdatedAtEnd   time.Time `json:"updatedAtEnd" form:"updatedAtEnd"`     // 更新时间结束
}

// UserListParams 用户列表查询参数
type UserListParams struct {
	BaseListParams
	BaseTimeRangeParams
	Username string `json:"username" form:"username"` // 用户名
	Email    string `json:"email" form:"email"`       // 邮箱
	Status   *int   `json:"status" form:"status"`     // 状态
}

// AdminListParams 管理员列表查询参数
type AdminListParams struct {
	BaseListParams
	BaseTimeRangeParams
	Username string `json:"username" form:"username"` // 用户名
	Email    string `json:"email" form:"email"`       // 邮箱
	Role     string `json:"role" form:"role"`         // 角色
	Status   *int   `json:"status" form:"status"`     // 状态
}

// DomainListParams 域名列表查询参数
type DomainListParams struct {
	BaseListParams
	BaseTimeRangeParams
	Name        string `json:"name" form:"name"`               // 域名
	Status      *int   `json:"status" form:"status"`           // 状态
	DnsVerified *bool  `json:"dnsVerified" form:"dnsVerified"` // DNS验证状态
}

// MailboxListParams 邮箱列表查询参数
type MailboxListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	DomainId    int64  `json:"domainId" form:"domainId"`       // 域名ID
	Email       string `json:"email" form:"email"`             // 邮箱地址
	Status      *int   `json:"status" form:"status"`           // 状态
	AutoReceive *bool  `json:"autoReceive" form:"autoReceive"` // 自动收信
}

// EmailListParams 邮件列表查询参数
type EmailListParams struct {
	BaseListParams
	BaseTimeRangeParams
	MailboxId   int64  `json:"mailboxId" form:"mailboxId"`     // 邮箱ID
	MessageId   string `json:"messageId" form:"messageId"`     // 消息ID
	Subject     string `json:"subject" form:"subject"`         // 主题
	FromEmail   string `json:"fromEmail" form:"fromEmail"`     // 发件人
	ToEmails    string `json:"toEmails" form:"toEmails"`       // 收件人
	Direction   string `json:"direction" form:"direction"`     // 方向
	IsRead      *bool  `json:"isRead" form:"isRead"`           // 是否已读
	IsStarred   *bool  `json:"isStarred" form:"isStarred"`     // 是否标星
	ContentType string `json:"contentType" form:"contentType"` // 内容类型
}

// EmailAttachmentListParams 邮件附件列表查询参数
type EmailAttachmentListParams struct {
	BaseListParams
	EmailId  int64  `json:"emailId" form:"emailId"`   // 邮件ID
	Filename string `json:"filename" form:"filename"` // 文件名
	MimeType string `json:"mimeType" form:"mimeType"` // MIME类型
}

// EmailTemplateListParams 邮件模板列表查询参数
type EmailTemplateListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	Name        string `json:"name" form:"name"`               // 模板名称
	Category    string `json:"category" form:"category"`       // 模板分类
	Subject     string `json:"subject" form:"subject"`         // 主题
	ContentType string `json:"contentType" form:"contentType"` // 内容类型
	IsDefault   *bool  `json:"isDefault" form:"isDefault"`     // 是否默认
	Status      *int   `json:"status" form:"status"`           // 状态
}

// EmailSignatureListParams 邮件签名列表查询参数
type EmailSignatureListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId    int64  `json:"userId" form:"userId"`       // 用户ID
	Name      string `json:"name" form:"name"`           // 签名名称
	IsDefault *bool  `json:"isDefault" form:"isDefault"` // 是否默认
	Status    *int   `json:"status" form:"status"`       // 状态
}

// VerificationRuleListParams 验证码规则列表查询参数
type VerificationRuleListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	Name        string `json:"name" form:"name"`               // 规则名称
	Pattern     string `json:"pattern" form:"pattern"`         // 匹配模式
	Source      string `json:"source" form:"source"`           // 来源
	Status      *int   `json:"status" form:"status"`           // 状态
	IsGlobal    *bool  `json:"isGlobal" form:"isGlobal"`       // 是否全局
	PatternType string `json:"patternType" form:"patternType"` // 模式类型
}

// UserVerificationRuleListParams 用户验证码规则关联列表查询参数
type UserVerificationRuleListParams struct {
	BaseListParams
	UserId int64 `json:"userId" form:"userId"` // 用户ID
	RuleId int64 `json:"ruleId" form:"ruleId"` // 规则ID
	Status *int  `json:"status" form:"status"` // 状态
}

// ForwardRuleListParams 转发规则列表查询参数
type ForwardRuleListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId         int64  `json:"userId" form:"userId"`                 // 用户ID
	Name           string `json:"name" form:"name"`                     // 规则名称
	FromPattern    string `json:"fromPattern" form:"fromPattern"`       // 发件人模式
	SubjectPattern string `json:"subjectPattern" form:"subjectPattern"` // 主题模式
	ForwardTo      string `json:"forwardTo" form:"forwardTo"`           // 转发地址
	KeepOriginal   *bool  `json:"keepOriginal" form:"keepOriginal"`     // 保留原邮件
	Status         *int   `json:"status" form:"status"`                 // 状态
	Priority       *int   `json:"priority" form:"priority"`             // 优先级
}

// AntiSpamRuleListParams 反垃圾规则列表查询参数
type AntiSpamRuleListParams struct {
	BaseListParams
	BaseTimeRangeParams
	Name     string `json:"name" form:"name"`         // 规则名称
	RuleType string `json:"ruleType" form:"ruleType"` // 规则类型
	Pattern  string `json:"pattern" form:"pattern"`   // 匹配模式
	Action   string `json:"action" form:"action"`     // 动作
	IsGlobal *bool  `json:"isGlobal" form:"isGlobal"` // 是否全局
	Status   *int   `json:"status" form:"status"`     // 状态
	Priority *int   `json:"priority" form:"priority"` // 优先级
}

// OperationLogListParams 操作日志列表查询参数
type OperationLogListParams struct {
	BaseListParams
	UserId     int64  `json:"userId" form:"userId"`         // 用户ID
	Action     string `json:"action" form:"action"`         // 操作
	Resource   string `json:"resource" form:"resource"`     // 资源
	ResourceId int64  `json:"resourceId" form:"resourceId"` // 资源ID
	Method     string `json:"method" form:"method"`         // 请求方法
	Path       string `json:"path" form:"path"`             // 请求路径
	Ip         string `json:"ip" form:"ip"`                 // IP地址
	Status     *int   `json:"status" form:"status"`         // 状态码
}

// EmailLogListParams 邮件日志列表查询参数
type EmailLogListParams struct {
	BaseListParams
	MailboxId int64  `json:"mailboxId" form:"mailboxId"` // 邮箱ID
	EmailId   int64  `json:"emailId" form:"emailId"`     // 邮件ID
	Type      string `json:"type" form:"type"`           // 类型
	Status    string `json:"status" form:"status"`       // 状态
	ErrorCode string `json:"errorCode" form:"errorCode"` // 错误码
}

// ApiKeyListParams API密钥列表查询参数
type ApiKeyListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	Name        string `json:"name" form:"name"`               // 密钥名称
	Permissions string `json:"permissions" form:"permissions"` // 权限
	Status      *int   `json:"status" form:"status"`           // 状态
}

// VerificationCodeListParams 验证码记录列表查询参数
type VerificationCodeListParams struct {
	BaseListParams
	EmailId int64  `json:"emailId" form:"emailId"` // 邮件ID
	RuleId  int64  `json:"ruleId" form:"ruleId"`   // 规则ID
	Code    string `json:"code" form:"code"`       // 验证码
	Source  string `json:"source" form:"source"`   // 来源
	IsUsed  *bool  `json:"isUsed" form:"isUsed"`   // 是否已使用
}

// EmailDraftListParams 草稿邮件列表查询参数
type EmailDraftListParams struct {
	BaseListParams
	BaseTimeRangeParams
	UserId      int64  `json:"userId" form:"userId"`           // 用户ID
	MailboxId   int64  `json:"mailboxId" form:"mailboxId"`     // 邮箱ID
	Subject     string `json:"subject" form:"subject"`         // 主题
	ToEmails    string `json:"toEmails" form:"toEmails"`       // 收件人
	ContentType string `json:"contentType" form:"contentType"` // 内容类型
	Status      string `json:"status" form:"status"`           // 状态
}
