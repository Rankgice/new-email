package types

import "time"

// DraftCreateReq 创建草稿请求
type DraftCreateReq struct {
	MailboxId   uint   `json:"mailboxId"`                             // 邮箱ID
	Subject     string `json:"subject" binding:"required,max=200"`    // 邮件主题
	FromEmail   string `json:"fromEmail" binding:"required,email"`    // 发件人邮箱
	ToEmail     string `json:"toEmail" binding:"required"`            // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content" binding:"required"`            // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
}

// DraftUpdateReq 更新草稿请求
type DraftUpdateReq struct {
	MailboxId   uint   `json:"mailboxId"`                             // 邮箱ID
	Subject     string `json:"subject" binding:"required,max=200"`    // 邮件主题
	FromEmail   string `json:"fromEmail" binding:"required,email"`    // 发件人邮箱
	ToEmail     string `json:"toEmail" binding:"required"`            // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content" binding:"required"`            // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
}

// DraftListReq 草稿列表请求
type DraftListReq struct {
	MailboxId      *uint     `json:"mailboxId" form:"mailboxId"`           // 邮箱ID
	Subject        string    `json:"subject" form:"subject"`               // 邮件主题（模糊搜索）
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 收件人邮箱
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// DraftResp 草稿响应
type DraftResp struct {
	Id          uint      `json:"id"`          // 草稿ID
	UserId      uint      `json:"userId"`      // 用户ID
	MailboxId   uint      `json:"mailboxId"`   // 邮箱ID
	Subject     string    `json:"subject"`     // 邮件主题
	FromEmail   string    `json:"fromEmail"`   // 发件人邮箱
	ToEmail     string    `json:"toEmail"`     // 收件人邮箱
	CcEmail     string    `json:"ccEmail"`     // 抄送邮箱
	BccEmail    string    `json:"bccEmail"`    // 密送邮箱
	Content     string    `json:"content"`     // 邮件内容
	ContentType string    `json:"contentType"` // 内容类型
	Attachments string    `json:"attachments"` // 附件信息
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// DraftSendResp 发送草稿响应
type DraftSendResp struct {
	Success bool      `json:"success"` // 是否成功
	Message string    `json:"message"` // 消息
	EmailId uint      `json:"emailId"` // 邮件ID
	SentAt  time.Time `json:"sentAt"`  // 发送时间
}
