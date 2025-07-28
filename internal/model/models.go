package model

import (
	"gorm.io/gorm"
	"time"
)

// 以下是其他模型的简化定义，用于确保项目能够编译

// EmailAttachment 邮件附件模型
type EmailAttachment struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailId   uint      `gorm:"not null;index" json:"email_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	MimeType  string    `gorm:"size:100" json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailAttachment) TableName() string { return "email_attachment" }

type EmailAttachmentModel struct{ db *gorm.DB }

func NewEmailAttachmentModel(db *gorm.DB) *EmailAttachmentModel { return &EmailAttachmentModel{db: db} }

// EmailTemplate 邮件模板模型
type EmailTemplate struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Subject     string         `gorm:"size:500" json:"subject"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailTemplate) TableName() string { return "email_template" }

type EmailTemplateModel struct{ db *gorm.DB }

func NewEmailTemplateModel(db *gorm.DB) *EmailTemplateModel { return &EmailTemplateModel{db: db} }

// EmailSignature 邮件签名模型
type EmailSignature struct {
	Id        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailSignature) TableName() string { return "email_signature" }

type EmailSignatureModel struct{ db *gorm.DB }

func NewEmailSignatureModel(db *gorm.DB) *EmailSignatureModel { return &EmailSignatureModel{db: db} }

// UserVerificationRule 用户验证码规则关联模型
type UserVerificationRule struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"not null;index" json:"user_id"`
	RuleId    uint      `gorm:"not null;index" json:"rule_id"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserVerificationRule) TableName() string { return "user_verification_rule" }

type UserVerificationRuleModel struct{ db *gorm.DB }

func NewUserVerificationRuleModel(db *gorm.DB) *UserVerificationRuleModel {
	return &UserVerificationRuleModel{db: db}
}

// ForwardRule 转发规则模型
type ForwardRule struct {
	Id             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId         uint           `gorm:"not null;index" json:"user_id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	FromPattern    string         `gorm:"size:255" json:"from_pattern"`
	SubjectPattern string         `gorm:"size:255" json:"subject_pattern"`
	ContentPattern string         `gorm:"type:text" json:"content_pattern"`
	ForwardTo      string         `gorm:"size:255;not null" json:"forward_to"`
	KeepOriginal   bool           `gorm:"default:true" json:"keep_original"`
	Status         int            `gorm:"default:1" json:"status"`
	Priority       int            `gorm:"default:0" json:"priority"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ForwardRule) TableName() string { return "forward_rule" }

type ForwardRuleModel struct{ db *gorm.DB }

func NewForwardRuleModel(db *gorm.DB) *ForwardRuleModel { return &ForwardRuleModel{db: db} }

// AntiSpamRule 反垃圾规则模型
type AntiSpamRule struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	RuleType    string         `gorm:"size:20;not null" json:"rule_type"`
	Pattern     string         `gorm:"type:text;not null" json:"pattern"`
	Action      string         `gorm:"size:20;not null" json:"action"`
	Description string         `gorm:"type:text" json:"description"`
	IsGlobal    bool           `gorm:"default:false" json:"is_global"`
	Status      int            `gorm:"default:1" json:"status"`
	Priority    int            `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AntiSpamRule) TableName() string { return "anti_spam_rule" }

type AntiSpamRuleModel struct{ db *gorm.DB }

func NewAntiSpamRuleModel(db *gorm.DB) *AntiSpamRuleModel { return &AntiSpamRuleModel{db: db} }

// OperationLog 操作日志模型
type OperationLog struct {
	Id         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId     uint      `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:100;not null" json:"action"`
	Resource   string    `gorm:"size:100" json:"resource"`
	ResourceId uint      `json:"resource_id"`
	Method     string    `gorm:"size:10" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	Ip         string    `gorm:"size:45" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Request    string    `gorm:"type:text" json:"request"`
	Response   string    `gorm:"type:text" json:"response"`
	Status     int       `json:"status"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

func (OperationLog) TableName() string { return "operation_log" }

type OperationLogModel struct{ db *gorm.DB }

func NewOperationLogModel(db *gorm.DB) *OperationLogModel { return &OperationLogModel{db: db} }

// EmailLog 邮件日志模型
type EmailLog struct {
	Id         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MailboxId  uint       `gorm:"not null;index" json:"mailbox_id"`
	EmailId    uint       `gorm:"index" json:"email_id"`
	Type       string     `gorm:"size:20;not null" json:"type"`
	Status     string     `gorm:"size:20;not null" json:"status"`
	Message    string     `gorm:"type:text" json:"message"`
	ErrorCode  string     `gorm:"size:50" json:"error_code"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (EmailLog) TableName() string { return "email_log" }

type EmailLogModel struct{ db *gorm.DB }

func NewEmailLogModel(db *gorm.DB) *EmailLogModel { return &EmailLogModel{db: db} }

// VerificationCode 验证码记录模型
type VerificationCode struct {
	Id          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailId     uint       `gorm:"not null;index" json:"email_id"`
	RuleId      uint       `gorm:"not null;index" json:"rule_id"`
	Code        string     `gorm:"size:50;not null" json:"code"`
	Source      string     `gorm:"size:100" json:"source"`
	ExtractedAt time.Time  `json:"extracted_at"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (VerificationCode) TableName() string { return "verification_code" }

type VerificationCodeModel struct{ db *gorm.DB }

func NewVerificationCodeModel(db *gorm.DB) *VerificationCodeModel {
	return &VerificationCodeModel{db: db}
}

// EmailDraft 草稿邮件模型
type EmailDraft struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint           `gorm:"not null;index" json:"user_id"`
	MailboxId   uint           `gorm:"not null;index" json:"mailbox_id"`
	Subject     string         `gorm:"size:500" json:"subject"`
	ToEmails    string         `gorm:"type:text" json:"to_emails"`
	CcEmails    string         `gorm:"type:text" json:"cc_emails"`
	BccEmails   string         `gorm:"type:text" json:"bcc_emails"`
	ContentType string         `gorm:"size:20;default:html" json:"content_type"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	Status      string         `gorm:"size:20;default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailDraft) TableName() string { return "email_draft" }

type EmailDraftModel struct{ db *gorm.DB }

func NewEmailDraftModel(db *gorm.DB) *EmailDraftModel { return &EmailDraftModel{db: db} }
