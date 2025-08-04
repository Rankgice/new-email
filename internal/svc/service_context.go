package svc

import (
	"fmt"
	"log"
	"new-email/internal/config"
	"new-email/internal/model"
	"new-email/internal/service"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB

	// 服务管理器
	ServiceManager *service.ServiceManager

	// Model层实例
	UserModel                 *model.UserModel
	AdminModel                *model.AdminModel
	DomainModel               *model.DomainModel
	MailboxModel              *model.MailboxModel
	EmailModel                *model.EmailModel
	EmailAttachmentModel      *model.EmailAttachmentModel
	EmailTemplateModel        *model.EmailTemplateModel
	EmailSignatureModel       *model.EmailSignatureModel
	VerificationRuleModel     *model.VerificationRuleModel
	UserVerificationRuleModel *model.UserVerificationRuleModel
	ForwardRuleModel          *model.ForwardRuleModel
	AntiSpamRuleModel         *model.AntiSpamRuleModel
	OperationLogModel         *model.OperationLogModel
	EmailLogModel             *model.EmailLogModel
	ApiKeyModel               *model.ApiKeyModel
	VerificationCodeModel     *model.VerificationCodeModel
	EmailDraftModel           *model.EmailDraftModel
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db := initDatabase(c)

	// 自动迁移数据表
	if err := autoMigrate(db); err != nil {
		log.Fatalln("数据表迁移失败", "error", err.Error())
	}

	// 初始化默认数据
	if err := initDefaultData(db, c); err != nil {
		log.Printf("初始化默认数据失败: %v", err)
	}

	// 初始化服务管理器
	serviceManager := initServiceManager(c)

	// 初始化Model层实例
	return &ServiceContext{
		Config:         c,
		DB:             db,
		ServiceManager: serviceManager,

		// 初始化所有Model实例
		UserModel:                 model.NewUserModel(db),
		AdminModel:                model.NewAdminModel(db),
		DomainModel:               model.NewDomainModel(db),
		MailboxModel:              model.NewMailboxModel(db),
		EmailModel:                model.NewEmailModel(db),
		EmailAttachmentModel:      model.NewEmailAttachmentModel(db),
		EmailTemplateModel:        model.NewEmailTemplateModel(db),
		EmailSignatureModel:       model.NewEmailSignatureModel(db),
		VerificationRuleModel:     model.NewVerificationRuleModel(db),
		UserVerificationRuleModel: model.NewUserVerificationRuleModel(db),
		ForwardRuleModel:          model.NewForwardRuleModel(db),
		AntiSpamRuleModel:         model.NewAntiSpamRuleModel(db),
		OperationLogModel:         model.NewOperationLogModel(db),
		EmailLogModel:             model.NewEmailLogModel(db),
		ApiKeyModel:               model.NewApiKeyModel(db),
		VerificationCodeModel:     model.NewVerificationCodeModel(db),
		EmailDraftModel:           model.NewEmailDraftModel(db),
	}
}

// initDatabase 初始化数据库连接
func initDatabase(c config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	switch c.Database.Type {
	case "sqlite":
		// SQLite数据库
		dbPath := c.Database.SQLite.Path
		if dbPath == "" {
			dbPath = "./data/email.db"
		}

		// 确保数据库目录存在
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalln("创建数据库目录失败", "error", err.Error())
		}

		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			log.Fatalln("连接SQLite数据库失败", "error", err.Error())
		}
		log.Printf("✅ SQLite数据库连接成功: %s", dbPath)

	case "mysql":
		// MySQL数据库
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.Database.MySQL.Username,
			c.Database.MySQL.Password,
			c.Database.MySQL.Host,
			c.Database.MySQL.Port,
			c.Database.MySQL.Database,
			c.Database.MySQL.Charset)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalln("连接MySQL数据库失败", "error", err.Error())
		}

		// 设置连接池
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalln("获取数据库连接池失败", "error", err.Error())
		}
		sqlDB.SetMaxIdleConns(c.Database.MySQL.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.Database.MySQL.MaxOpenConns)

		log.Printf("✅ MySQL数据库连接成功: %s:%d/%s",
			c.Database.MySQL.Host, c.Database.MySQL.Port, c.Database.MySQL.Database)

	default:
		log.Fatalln("不支持的数据库类型:", c.Database.Type)
	}

	return db
}

// initServiceManager 初始化服务管理器
func initServiceManager(c config.Config) *service.ServiceManager {
	// 构建服务配置
	serviceConfig := service.ServiceConfig{
		SMTP: service.SMTPConfig{
			Host:     c.SMTP.Host,
			Port:     c.SMTP.Port,
			Username: c.SMTP.Username,
			Password: c.SMTP.Password,
			UseTLS:   c.SMTP.UseTLS,
		},
		IMAP: service.IMAPConfig{
			Host:     c.IMAP.Host,
			Port:     c.IMAP.Port,
			Username: c.IMAP.Username,
			Password: c.IMAP.Password,
			UseTLS:   c.IMAP.UseTLS,
		},
		SMS: service.SMSConfig{
			Provider:  c.SMS.Provider,
			AccessKey: c.SMS.AccessKey,
			SecretKey: c.SMS.SecretKey,
			SignName:  c.SMS.SignName,
			Region:    c.SMS.Region,
		},
		Storage: service.StorageConfig{
			Type:      c.Storage.Type,
			BasePath:  c.Storage.BasePath,
			MaxSize:   c.Storage.MaxSize,
			AllowExts: c.Storage.AllowExts,
			CDNDomain: c.Storage.CDNDomain,
		},
		Cache: service.CacheConfig{
			Host:     c.Redis.Host,
			Port:     c.Redis.Port,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
			PoolSize: c.Redis.PoolSize,
		},
	}

	// 创建服务管理器
	manager := service.NewServiceManager(serviceConfig)

	// 测试所有连接
	results := manager.TestAllConnections()
	for serviceName, err := range results {
		if err != nil {
			log.Printf("服务 %s 连接失败: %v", serviceName, err)
		} else {
			log.Printf("服务 %s 连接成功", serviceName)
		}
	}

	return manager
}

// autoMigrate 自动迁移数据表结构
func autoMigrate(db *gorm.DB) error {
	log.Println("🔄 开始数据库迁移...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
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
		return err
	}

	log.Println("✅ 数据库迁移完成")
	return nil
}

// initDefaultData 初始化默认数据
func initDefaultData(db *gorm.DB, c config.Config) error {
	log.Println("🔄 初始化默认数据...")

	// 检查是否已存在超级管理员
	var count int64
	if err := db.Model(&model.Admin{}).Where("role = ? AND status = ?", "admin", 1).Count(&count).Error; err != nil {
		return err
	}

	// 如果不存在超级管理员，创建默认管理员
	if count == 0 {
		defaultAdmin := &model.Admin{
			Username: c.System.DefaultAdmin.Username,
			Email:    c.System.DefaultAdmin.Email,
			Password: c.System.DefaultAdmin.Password, // 注意：实际使用时需要加密
			Nickname: c.System.DefaultAdmin.Nickname,
			Role:     "admin",
			Status:   1,
		}

		if err := db.Create(defaultAdmin).Error; err != nil {
			return fmt.Errorf("创建默认管理员失败: %v", err)
		}

		log.Printf("✅ 创建默认管理员成功: %s", defaultAdmin.Username)
	}

	// 初始化默认验证码规则
	if err := initDefaultVerificationRules(db); err != nil {
		return err
	}

	log.Println("✅ 默认数据初始化完成")
	return nil
}

// initDefaultVerificationRules 初始化默认验证码规则
func initDefaultVerificationRules(db *gorm.DB) error {
	// 检查是否已存在全局规则
	var count int64
	if err := db.Model(&model.VerificationRule{}).Where("is_global = ?", true).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // 已存在全局规则，跳过初始化
	}

	// 默认验证码规则
	defaultRules := []*model.VerificationRule{
		{
			UserId:      0, // 公共规则
			Name:        "通用数字验证码",
			Pattern:     `\b\d{4,8}\b`,
			Description: "匹配4-8位数字验证码",
			IsGlobal:    true,
			Status:      1,
			Priority:    1,
		},
		{
			UserId:      0,
			Name:        "通用字母数字验证码",
			Pattern:     `\b[A-Za-z0-9]{4,8}\b`,
			Description: "匹配4-8位字母数字验证码",
			IsGlobal:    true,
			Status:      1,
			Priority:    2,
		},
		{
			UserId:      0,
			Name:        "邮箱验证码",
			Pattern:     `验证码[：:]\s*([A-Za-z0-9]{4,8})`,
			Description: "匹配邮箱中的验证码格式",
			IsGlobal:    true,
			Status:      1,
			Priority:    3,
		},
	}

	for _, rule := range defaultRules {
		if err := db.Create(rule).Error; err != nil {
			return fmt.Errorf("创建默认验证码规则失败: %v", err)
		}
	}

	log.Println("✅ 默认验证码规则初始化完成")
	return nil
}
