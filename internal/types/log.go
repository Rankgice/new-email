package types

import "time"

// LogListReq 日志列表请求
type LogListReq struct {
	UserId         *uint     `json:"userId" form:"userId"`                 // 用户ID
	Action         string    `json:"action" form:"action"`                 // 操作类型
	Resource       string    `json:"resource" form:"resource"`             // 资源类型
	Method         string    `json:"method" form:"method"`                 // 请求方法
	Status         *int      `json:"status" form:"status"`                 // 状态码
	Ip             string    `json:"ip" form:"ip"`                         // IP地址
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// OperationLogResp 操作日志响应
type OperationLogResp struct {
	Id         uint      `json:"id"`         // 日志ID
	UserId     uint      `json:"userId"`     // 用户ID
	Action     string    `json:"action"`     // 操作类型
	Resource   string    `json:"resource"`   // 资源类型
	ResourceId uint      `json:"resourceId"` // 资源ID
	Method     string    `json:"method"`     // 请求方法
	Path       string    `json:"path"`       // 请求路径
	Ip         string    `json:"ip"`         // IP地址
	UserAgent  string    `json:"userAgent"`  // 用户代理
	Status     int       `json:"status"`     // 状态码
	ErrorMsg   string    `json:"errorMsg"`   // 错误信息
	CreatedAt  time.Time `json:"createdAt"`  // 创建时间
}

// EmailLogListReq 邮件日志列表请求
type EmailLogListReq struct {
	EmailId        *uint     `json:"emailId" form:"emailId"`               // 邮件ID
	MailboxId      *uint     `json:"mailboxId" form:"mailboxId"`           // 邮箱ID
	Type           string    `json:"type" form:"type"`                     // 日志类型
	Status         *int      `json:"status" form:"status"`                 // 状态
	FromEmail      string    `json:"fromEmail" form:"fromEmail"`           // 发件人
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 收件人
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// EmailLogResp 邮件日志响应
type EmailLogResp struct {
	Id        uint      `json:"id"`        // 日志ID
	EmailId   uint      `json:"emailId"`   // 邮件ID
	MailboxId uint      `json:"mailboxId"` // 邮箱ID
	Type      string    `json:"type"`      // 日志类型
	Status    int       `json:"status"`    // 状态
	FromEmail string    `json:"fromEmail"` // 发件人
	ToEmail   string    `json:"toEmail"`   // 收件人
	Subject   string    `json:"subject"`   // 邮件主题
	ErrorMsg  string    `json:"errorMsg"`  // 错误信息
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}
