package svc

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rankgice/new-email/internal/config"
	"github.com/rankgice/new-email/internal/model"
	"github.com/rankgice/new-email/internal/service"
	"github.com/rankgice/new-email/pkg/auth"
	"gorm.io/gorm/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB

	// 服务管理器
	ServiceManager *service.ServiceManager

	minioClient *minio.Client
	// Model层实例
	UserModel            *model.UserModel
	AdminModel           *model.AdminModel
	DomainModel          *model.DomainModel
	MailboxModel         *model.MailboxModel
	EmailModel           *model.EmailModel
	EmailAttachmentModel *model.EmailAttachmentModel
	ApiKeyModel          *model.ApiKeyModel
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

	// 初始化minio
	minioClient, err := initMinio(c)
	if err != nil {
		log.Fatalln("初始化minio失败", "error", err.Error())
	}

	// 初始化服务管理器
	serviceManager := initServiceManager(c)

	// 初始化Model层实例
	return &ServiceContext{
		Config:         c,
		DB:             db,
		ServiceManager: serviceManager,

		minioClient: minioClient,
		// 初始化所有Model实例
		UserModel:            model.NewUserModel(db),
		AdminModel:           model.NewAdminModel(db),
		DomainModel:          model.NewDomainModel(db),
		MailboxModel:         model.NewMailboxModel(db),
		EmailModel:           model.NewEmailModel(db),
		EmailAttachmentModel: model.NewEmailAttachmentModel(db),
		ApiKeyModel:          model.NewApiKeyModel(db),
	}
}

// initMinio 初始化minio
func initMinio(c config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyId, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSl,
	})
	if err != nil {
		return nil, err
	}
	// 初始化桶
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	for _, bucketName := range c.Minio.Buckets {
		if exists, err := minioClient.BucketExists(ctx, bucketName); err != nil {
			return nil, err
		} else if !exists {
			if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
				return nil, err
			}
		}
	}
	return minioClient, nil
}

// initDatabase 初始化数据库连接
func initDatabase(c config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

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

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印所有 SQL
	})

	if err != nil {
		log.Fatalln("连接SQLite数据库失败", "error", err.Error())
	}
	log.Printf("✅ SQLite数据库连接成功: %s", dbPath)

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

	// 如果Redis未启用，则将CacheService设置为nil
	if !c.Redis.Enabled {
		manager.Cache = nil
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
		&model.Folder{},
		&model.Email{},
		&model.EmailAttachment{},
		&model.ApiKey{},
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

	// 创建默认域名（如果不存在）
	var domainCount int64
	if err := db.Model(&model.Domain{}).Where("name = ?", "email.host").Count(&domainCount).Error; err != nil {
		return err
	}

	var defaultDomainId int64
	if domainCount == 0 {
		defaultDomain := &model.Domain{
			Name:        "email.host",
			Status:      1, // 启用
			DnsVerified: 1, // 假设已验证
			SpfRecord:   "v=spf1 mx a -all",
			DmarcRecord: "v=DMARC1; p=quarantine; rua=mailto:dmarc@email.host",
		}

		if err := db.Create(defaultDomain).Error; err != nil {
			return fmt.Errorf("创建默认域名失败: %v", err)
		}

		defaultDomainId = defaultDomain.Id
		log.Printf("✅ 创建默认域名成功: %s", defaultDomain.Name)
	} else {
		// 获取现有域名ID
		var domain model.Domain
		if err := db.Where("name = ?", "email.host").First(&domain).Error; err != nil {
			return err
		}
		defaultDomainId = domain.Id
	}

	// 创建测试邮箱（如果不存在）
	var mailboxCount int64
	if err := db.Model(&model.Mailbox{}).Where("email = ?", "test@email.host").Count(&mailboxCount).Error; err != nil {
		return err
	}

	if mailboxCount == 0 {
		// 加密测试密码
		hashedPassword, err := auth.HashPassword("test123")
		if err != nil {
			return fmt.Errorf("密码加密失败: %v", err)
		}

		testMailbox := &model.Mailbox{
			UserId:      1, // 假设用户ID为1
			DomainId:    defaultDomainId,
			Email:       "test@email.host",
			Password:    hashedPassword,
			Type:        "imap",
			Status:      1, // 启用
			AutoReceive: true,
		}

		if err := db.Create(testMailbox).Error; err != nil {
			return fmt.Errorf("创建测试邮箱失败: %v", err)
		}

		log.Printf("✅ 创建测试邮箱成功: %s (密码: test123)", testMailbox.Email)
	}

	log.Println("✅ 默认数据初始化完成")
	return nil
}
