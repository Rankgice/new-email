package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	App       AppConfig       `yaml:"app"`
	Web       WebConfig       `yaml:"web"`
	Database  DatabaseConfig  `yaml:"database"`
	JWT       JWTConfig       `yaml:"jwt"`
	Email     EmailConfig     `yaml:"email"`
	SMTP      SMTPConfig      `yaml:"smtp"`    // 新增SMTP配置
	IMAP      IMAPConfig      `yaml:"imap"`    // 新增IMAP配置
	SMS       SMSConfig       `yaml:"sms"`     // 新增SMS配置
	Storage   StorageConfig   `yaml:"storage"` // 新增存储配置
	Log       LogConfig       `yaml:"log"`
	Redis     RedisConfig     `yaml:"redis"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	System    SystemConfig    `yaml:"system"`
	Minio     MinioConfig     `yaml:"minio"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

// WebConfig Web服务配置
type WebConfig struct {
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `yaml:"type"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	MySQL  MySQLConfig  `yaml:"mysql"`
}

// SQLiteConfig SQLite配置
type SQLiteConfig struct {
	Path string `yaml:"path"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type MinioConfig struct {
	Endpoint        string   `yaml:"endpoint"`
	AccessKeyId     string   `yaml:"access_key_id"`
	SecretAccessKey string   `yaml:"secret_access_key"`
	UseSSl          bool     `yaml:"use_ssl"`
	Buckets         []string `yaml:"buckets"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `yaml:"secret"`
	ExpireHours        int    `yaml:"expire_hours"`
	RefreshExpireHours int    `yaml:"refresh_expire_hours"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	DefaultSMTP SMTPConfig       `yaml:"default_smtp"`
	Attachment  AttachmentConfig `yaml:"attachment"`
	Receive     ReceiveConfig    `yaml:"receive"`
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	UseTLS   bool   `yaml:"use_tls"`
}

// AttachmentConfig 附件配置
type AttachmentConfig struct {
	MaxSize      int64    `yaml:"max_size"`
	StoragePath  string   `yaml:"storage_path"`
	AllowedTypes []string `yaml:"allowed_types"`
}

// ReceiveConfig 收信配置
type ReceiveConfig struct {
	BatchSize    int `yaml:"batch_size"`
	SyncInterval int `yaml:"sync_interval"`
	MaxRetries   int `yaml:"max_retries"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
	Burst             int  `yaml:"burst"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	DefaultAdmin DefaultAdminConfig `yaml:"default_admin"`
	Registration RegistrationConfig `yaml:"registration"`
	Security     SecurityConfig     `yaml:"security"`
}

// DefaultAdminConfig 默认管理员配置
type DefaultAdminConfig struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Nickname string `yaml:"nickname"`
}

// RegistrationConfig 注册配置
type RegistrationConfig struct {
	Enabled                  bool `yaml:"enabled"`
	RequireEmailVerification bool `yaml:"require_email_verification"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	PasswordMinLength int `yaml:"password_min_length"`
	SessionTimeout    int `yaml:"session_timeout"`
	MaxLoginAttempts  int `yaml:"max_login_attempts"`
	LockoutDuration   int `yaml:"lockout_duration"`
}

// IMAPConfig IMAP配置
type IMAPConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	UseTLS      bool   `yaml:"use_tls"`
	TLSCertPath string `yaml:"tls_cert_path"` // TLS证书路径
	TLSKeyPath  string `yaml:"tls_key_path"`  // TLS密钥路径
}

// SMSConfig SMS配置
type SMSConfig struct {
	Provider  string `yaml:"provider"` // aliyun, tencent, twilio
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	SignName  string `yaml:"sign_name"`
	Region    string `yaml:"region"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string   `yaml:"type"` // local, oss, s3
	BasePath  string   `yaml:"base_path"`
	MaxSize   int64    `yaml:"max_size"`
	AllowExts []string `yaml:"allow_exts"`
	CDNDomain string   `yaml:"cdn_domain"`
}

// NewConfig 创建配置
func NewConfig(path string) Config {
	conf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("读取配置文件失败：", err)
	}

	var c Config
	if err := yaml.Unmarshal(conf, &c); err != nil {
		log.Fatal("解析配置文件失败：", err)
	}

	return c
}
