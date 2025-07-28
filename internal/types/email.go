package types

import "time"

// EmailCreateReq 创建邮件请求
type EmailCreateReq struct {
	MailboxId   uint   `json:"mailboxId" binding:"required"`          // 邮箱ID
	Subject     string `json:"subject" binding:"required,max=200"`    // 邮件主题
	FromEmail   string `json:"fromEmail" binding:"required,email"`    // 发件人邮箱
	ToEmail     string `json:"toEmail" binding:"required"`            // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content" binding:"required"`            // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
	Type        string `json:"type" binding:"oneof=inbox sent draft"` // 邮件类型
}

// EmailUpdateReq 更新邮件请求
type EmailUpdateReq struct {
	Subject     string `json:"subject" binding:"required,max=200"`    // 邮件主题
	FromEmail   string `json:"fromEmail" binding:"required,email"`    // 发件人邮箱
	ToEmail     string `json:"toEmail" binding:"required"`            // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content" binding:"required"`            // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
	Status      int    `json:"status" binding:"oneof=0 1 2"`          // 状态：0未读 1已读 2已删除
}

// EmailListReq 邮件列表请求
type EmailListReq struct {
	MailboxId      *uint     `json:"mailboxId" form:"mailboxId"`           // 邮箱ID
	Subject        string    `json:"subject" form:"subject"`               // 邮件主题（模糊搜索）
	FromEmail      string    `json:"fromEmail" form:"fromEmail"`           // 发件人邮箱
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 收件人邮箱
	Status         *int      `json:"status" form:"status"`                 // 状态
	Type           string    `json:"type" form:"type"`                     // 邮件类型
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// EmailResp 邮件响应
type EmailResp struct {
	Id          uint      `json:"id"`          // 邮件ID
	MailboxId   uint      `json:"mailboxId"`   // 邮箱ID
	Subject     string    `json:"subject"`     // 邮件主题
	FromEmail   string    `json:"fromEmail"`   // 发件人邮箱
	ToEmail     string    `json:"toEmail"`     // 收件人邮箱
	CcEmail     string    `json:"ccEmail"`     // 抄送邮箱
	BccEmail    string    `json:"bccEmail"`    // 密送邮箱
	Content     string    `json:"content"`     // 邮件内容
	ContentType string    `json:"contentType"` // 内容类型
	Attachments string    `json:"attachments"` // 附件信息
	Status      int       `json:"status"`      // 状态
	Type        string    `json:"type"`        // 邮件类型
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// EmailSendReq 发送邮件请求
type EmailSendReq struct {
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

// EmailSendResp 发送邮件响应
type EmailSendResp struct {
	Success bool      `json:"success"` // 是否成功
	Message string    `json:"message"` // 消息
	EmailId uint      `json:"emailId"` // 邮件ID
	SentAt  time.Time `json:"sentAt"`  // 发送时间
}

// EmailStatsResp 邮件统计响应
type EmailStatsResp struct {
	Total     int64 `json:"total"`     // 总邮件数
	Unread    int64 `json:"unread"`    // 未读邮件数
	Read      int64 `json:"read"`      // 已读邮件数
	Sent      int64 `json:"sent"`      // 已发送邮件数
	Draft     int64 `json:"draft"`     // 草稿邮件数
	Today     int64 `json:"today"`     // 今日邮件数
	ThisWeek  int64 `json:"thisWeek"`  // 本周邮件数
	ThisMonth int64 `json:"thisMonth"` // 本月邮件数
}

// EmailBatchOperationReq 邮件批量操作请求
type EmailBatchOperationReq struct {
	Ids       []uint `json:"ids" binding:"required,min=1"`                                    // 邮件ID列表
	Operation string `json:"operation" binding:"required,oneof=read unread delete move copy"` // 操作类型
	TargetId  uint   `json:"targetId"`                                                        // 目标ID（移动或复制时使用）
}

// EmailSearchReq 邮件搜索请求
type EmailSearchReq struct {
	Keyword        string    `json:"keyword" form:"keyword" binding:"required,min=1"` // 搜索关键词
	MailboxId      *uint     `json:"mailboxId" form:"mailboxId"`                      // 邮箱ID
	SearchIn       string    `json:"searchIn" form:"searchIn"`                        // 搜索范围：subject, content, from, to, all
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"`            // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`                // 创建时间结束
	PageReq
}
