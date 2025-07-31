package types

import "time"

// DraftCreateReq 创建草稿请求
type DraftCreateReq struct {
	MailboxId   int64  `json:"mailboxId"`                             // 邮箱ID
	Subject     string `json:"subject"`                               // 邮件主题
	FromEmail   string `json:"fromEmail"`                             // 发件人邮箱（可选，从邮箱信息自动获取）
	ToEmail     string `json:"toEmail"`                               // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content"`                               // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
}

// DraftUpdateReq 更新草稿请求
type DraftUpdateReq struct {
	MailboxId   int64  `json:"mailboxId"`                             // 邮箱ID
	Subject     string `json:"subject"`                               // 邮件主题
	FromEmail   string `json:"fromEmail"`                             // 发件人邮箱（可选，从邮箱信息自动获取）
	ToEmail     string `json:"toEmail"`                               // 收件人邮箱（多个用逗号分隔）
	CcEmail     string `json:"ccEmail"`                               // 抄送邮箱（多个用逗号分隔）
	BccEmail    string `json:"bccEmail"`                              // 密送邮箱（多个用逗号分隔）
	Content     string `json:"content"`                               // 邮件内容
	ContentType string `json:"contentType" binding:"oneof=text html"` // 内容类型：text, html
	Attachments string `json:"attachments"`                           // 附件信息（JSON格式）
}

// DraftListReq 草稿列表请求
type DraftListReq struct {
	MailboxId      *int64    `json:"mailboxId" form:"mailboxId"`           // 邮箱ID
	Subject        string    `json:"subject" form:"subject"`               // 邮件主题（模糊搜索）
	ToEmail        string    `json:"toEmail" form:"toEmail"`               // 收件人邮箱
	CreatedAtStart time.Time `json:"createdAtStart" form:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   time.Time `json:"createdAtEnd" form:"createdAtEnd"`     // 创建时间结束
	PageReq
}

// DraftResp 草稿响应
type DraftResp struct {
	Id          int64     `json:"id"`          // 草稿ID
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
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// DraftSendResp 发送草稿响应
type DraftSendResp struct {
	Success bool      `json:"success"` // 是否成功
	Message string    `json:"message"` // 消息
	EmailId int64     `json:"emailId"` // 邮件ID
	SentAt  time.Time `json:"sentAt"`  // 发送时间
}

// DraftAutoSaveReq 草稿自动保存请求
type DraftAutoSaveReq struct {
	DraftId     *int64 `json:"draftId"`     // 草稿ID（可选，用于更新现有草稿）
	MailboxId   int64  `json:"mailboxId"`   // 邮箱ID
	Subject     string `json:"subject"`     // 邮件主题
	ToEmail     string `json:"toEmail"`     // 收件人
	CcEmail     string `json:"ccEmail"`     // 抄送
	BccEmail    string `json:"bccEmail"`    // 密送
	Content     string `json:"content"`     // 邮件内容
	ContentType string `json:"contentType"` // 内容类型：text, html
}

// DraftAutoSaveResp 草稿自动保存响应
type DraftAutoSaveResp struct {
	DraftId int64     `json:"draftId"` // 草稿ID
	Success bool      `json:"success"` // 是否成功
	Message string    `json:"message"` // 消息
	SavedAt time.Time `json:"savedAt"` // 保存时间
	IsNew   bool      `json:"isNew"`   // 是否为新创建的草稿
}
