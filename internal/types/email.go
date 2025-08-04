package types

import "time"

// EmailCreateReq 创建邮件请求
type EmailCreateReq struct {
	MailboxId   int64  `json:"mailboxId" binding:"required"`          // 邮箱ID
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

// EmailListReq 邮件列表请求
type EmailListReq struct {
	UserId         int64     `json:"userId" form:"userId"`                 // 用户ID
	MailboxId      int64     `json:"mailboxId" form:"mailboxId"`           // 邮箱ID
	Subject        string    `json:"subject" form:"subject"`               // 邮件主题（模糊搜索）
	FromEmail      string    `json:"fromEmail" form:"fromEmail"`           // 发件人邮箱
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 收件人邮箱
	Direction      string    `json:"direction" form:"direction"`           // 邮件方向：sent(已发送), received(已接收)
	Status         *int      `json:"status" form:"status"`                 // 状态
	Type           string    `json:"type" form:"type"`                     // 邮件类型
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// EmailResp 邮件响应
type EmailResp struct {
	Id          int64     `json:"id"`          // 邮件ID
	UserId      int64     `json:"userId"`      // 用户ID
	MailboxId   int64     `json:"mailboxId"`   // 邮箱ID
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
	MailboxId   int64  `json:"mailboxId"`                             // 邮箱ID
	Subject     string `json:"subject" binding:"required,max=200"`    // 邮件主题
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
	EmailId int64     `json:"emailId"` // 邮件ID
	SentAt  time.Time `json:"sentAt"`  // 发送时间
}

// EmailBatchOperationReq 邮件批量操作请求
type EmailBatchOperationReq struct {
	Ids       []int64 `json:"ids" binding:"required,min=1"`                                    // 邮件ID列表
	Operation string  `json:"operation" binding:"required,oneof=read unread delete move copy"` // 操作类型
	TargetId  int64   `json:"targetId"`                                                        // 目标ID（移动或复制时使用）
}

// EmailBatchOperationResp 邮件批量操作响应
type EmailBatchOperationResp struct {
	Success int `json:"success"` // 成功数量
	Failed  int `json:"failed"`  // 失败数量
}

// EmailExportReq 邮件导出请求
type EmailExportReq struct {
	MailboxId      *int64    `json:"mailboxId" form:"mailboxId"`                        // 邮箱ID
	Subject        string    `json:"subject" form:"subject"`                            // 邮件主题（模糊搜索）
	FromEmail      string    `json:"fromEmail" form:"fromEmail"`                        // 发件人邮箱
	ToEmail        string    `json:"toEmail" form:"toEmail"`                            // 收件人邮箱
	Direction      string    `json:"direction" form:"direction"`                        // 邮件方向：sent(已发送), received(已接收)
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"`              // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`                  // 创建时间结束
	Format         string    `json:"format" form:"format" binding:"oneof=csv json eml"` // 导出格式：csv, json, eml
	IncludeContent bool      `json:"includeContent" form:"includeContent"`              // 是否包含邮件内容
}

// EmailExportResp 邮件导出响应
type EmailExportResp struct {
	FileName    string `json:"fileName"`    // 文件名
	FileSize    int64  `json:"fileSize"`    // 文件大小
	RecordCount int64  `json:"recordCount"` // 记录数量
	DownloadUrl string `json:"downloadUrl"` // 下载链接
}
