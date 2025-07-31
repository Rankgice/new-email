package fixtures

import (
	"log"
	"new-email/internal/config"
	"new-email/internal/model"
	"new-email/internal/svc"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB 测试数据库结构
type TestDB struct {
	DB     *gorm.DB
	SvcCtx *svc.ServiceContext
	Config config.Config
}

// NewTestDB 创建测试数据库
func NewTestDB() *TestDB {
	// 创建测试配置
	c := config.Config{
		App: config.AppConfig{
			Name:    "邮件管理系统-测试",
			Version: "1.0.0-test",
			Debug:   true,
		},
		Web: config.WebConfig{
			Port:         8888,
			Mode:         "debug",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: config.DatabaseConfig{
			Type: "sqlite",
			SQLite: config.SQLiteConfig{
				Path: ":memory:",
			},
		},
		JWT: config.JWTConfig{
			Secret:             "test-jwt-secret-key-for-testing-only",
			ExpireHours:        1,
			RefreshExpireHours: 24,
		},
	}

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 测试时静默日志
	})
	if err != nil {
		log.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移数据表
	err = db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.Domain{},
		&model.Mailbox{},
		&model.Email{},
		&model.EmailAttachment{},
		&model.EmailTemplate{},
		&model.EmailSignature{},
		&model.VerificationRule{},
		&model.UserVerificationRule{},
		&model.ForwardRule{},
		&model.AntiSpamRule{},
		&model.OperationLog{},
		&model.EmailLog{},
		&model.ApiKey{},
		&model.VerificationCode{},
		&model.EmailDraft{},
	)
	if err != nil {
		log.Fatalf("数据表迁移失败: %v", err)
	}

	// 创建服务上下文（使用测试数据库）
	svcCtx := &svc.ServiceContext{
		Config: c,
		DB:     db,
	}

	// 初始化Model实例
	svcCtx.UserModel = model.NewUserModel(db)
	svcCtx.AdminModel = model.NewAdminModel(db)
	svcCtx.DomainModel = model.NewDomainModel(db)
	svcCtx.MailboxModel = model.NewMailboxModel(db)
	svcCtx.EmailModel = model.NewEmailModel(db)
	svcCtx.EmailAttachmentModel = model.NewEmailAttachmentModel(db)
	svcCtx.EmailTemplateModel = model.NewEmailTemplateModel(db)
	svcCtx.EmailSignatureModel = model.NewEmailSignatureModel(db)
	svcCtx.VerificationRuleModel = model.NewVerificationRuleModel(db)
	svcCtx.UserVerificationRuleModel = model.NewUserVerificationRuleModel(db)
	svcCtx.ForwardRuleModel = model.NewForwardRuleModel(db)
	svcCtx.AntiSpamRuleModel = model.NewAntiSpamRuleModel(db)
	svcCtx.OperationLogModel = model.NewOperationLogModel(db)
	svcCtx.EmailLogModel = model.NewEmailLogModel(db)
	svcCtx.ApiKeyModel = model.NewApiKeyModel(db)
	svcCtx.VerificationCodeModel = model.NewVerificationCodeModel(db)
	svcCtx.EmailDraftModel = model.NewEmailDraftModel(db)

	return &TestDB{
		DB:     db,
		SvcCtx: svcCtx,
		Config: c,
	}
}

// Close 关闭测试数据库
func (tdb *TestDB) Close() {
	sqlDB, err := tdb.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// Clean 清理测试数据
func (tdb *TestDB) Clean() {
	// 清理所有表的数据
	tables := []string{
		"users", "admins", "domains", "mailboxes", "emails",
		"email_attachments", "email_templates", "email_signatures",
		"verification_rules", "user_verification_rules", "forward_rules",
		"anti_spam_rules", "operation_logs", "email_logs",
		"api_keys", "verification_codes", "email_drafts",
	}

	for _, table := range tables {
		tdb.DB.Exec("DELETE FROM " + table)
	}
}

// SeedTestData 插入测试数据
func (tdb *TestDB) SeedTestData() {
	// 创建测试管理员
	admin := &model.Admin{
		Username: "testadmin",
		Email:    "testadmin@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Nickname: "测试管理员",
		Role:     "admin",
		Status:   1,
	}
	tdb.DB.Create(admin)

	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Nickname: "测试用户",
		Status:   1,
	}
	tdb.DB.Create(user)

	// 创建测试域名
	domain := &model.Domain{
		Name:        "test.com",
		Status:      1,
		DnsVerified: 2, // 2表示已验证
	}
	tdb.DB.Create(domain)

	// 创建测试邮箱
	mailbox := &model.Mailbox{
		UserId:      user.Id,
		DomainId:    domain.Id,
		Email:       "test@test.com",
		Password:    "encrypted_password",
		Status:      1,
		AutoReceive: true,
	}
	tdb.DB.Create(mailbox)

	// 创建测试验证码规则
	rule := &model.VerificationRule{
		UserId:      0, // 公共规则
		Name:        "测试验证码规则",
		Pattern:     `\b\d{6}\b`,
		Description: "6位数字验证码",
		IsGlobal:    true,
		Status:      1,
		Priority:    1,
	}
	tdb.DB.Create(rule)
}
