package service

import (
	"fmt"
	"log"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	SMTP          *SMTPService
	IMAP          *IMAPService
	SMS           *SMSService
	Storage       *StorageService
	Cache         *CacheService
	MessageParser *MessageParser
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	SMTP    SMTPConfig    `json:"smtp"`
	IMAP    IMAPConfig    `json:"imap"`
	SMS     SMSConfig     `json:"sms"`
	Storage StorageConfig `json:"storage"`
	Cache   CacheConfig   `json:"cache"`
}

// NewServiceManager 创建服务管理器
func NewServiceManager(config ServiceConfig) *ServiceManager {
	manager := &ServiceManager{}

	// 初始化SMTP服务
	if config.SMTP.Host != "" {
		manager.SMTP = NewSMTPService(config.SMTP)
		log.Printf("SMTP服务已初始化: %s:%d", config.SMTP.Host, config.SMTP.Port)
	}

	// 初始化IMAP服务
	if config.IMAP.Host != "" {
		manager.IMAP = NewIMAPService(config.IMAP)
		log.Printf("IMAP服务已初始化: %s:%d", config.IMAP.Host, config.IMAP.Port)
	}

	// 初始化SMS服务
	if config.SMS.Provider != "" {
		manager.SMS = NewSMSService(config.SMS)
		log.Printf("SMS服务已初始化: %s", config.SMS.Provider)
	}

	// 初始化存储服务
	if config.Storage.BasePath != "" {
		manager.Storage = NewStorageService(config.Storage)
		log.Printf("存储服务已初始化: %s", config.Storage.BasePath)
	}

	// 初始化缓存服务
	if config.Cache.Host != "" {
		manager.Cache = NewCacheService(config.Cache)
		log.Printf("缓存服务已初始化: %s:%d", config.Cache.Host, config.Cache.Port)
	}

	// 初始化邮件解析器
	manager.MessageParser = NewMessageParser()
	log.Printf("邮件解析器已初始化")

	return manager
}

// TestAllConnections 测试所有服务连接
func (m *ServiceManager) TestAllConnections() map[string]error {
	results := make(map[string]error)

	// 测试SMTP连接
	if m.SMTP != nil {
		if err := m.SMTP.TestConnection(); err != nil {
			results["smtp"] = err
			log.Printf("SMTP连接测试失败: %v", err)
		} else {
			results["smtp"] = nil
			log.Println("SMTP连接测试成功")
		}
	}

	// 测试IMAP连接
	if m.IMAP != nil {
		if err := m.IMAP.TestConnection(); err != nil {
			results["imap"] = err
			log.Printf("IMAP连接测试失败: %v", err)
		} else {
			results["imap"] = nil
			log.Println("IMAP连接测试成功")
		}
	}

	// 测试SMS连接
	if m.SMS != nil {
		// SMS测试可能会产生费用，这里只做配置检查
		config := m.SMS.GetSMSConfig()
		if config.Provider == "" {
			results["sms"] = fmt.Errorf("SMS服务商未配置")
		} else {
			results["sms"] = nil
			log.Println("SMS配置检查通过")
		}
	}

	// 测试存储服务
	if m.Storage != nil {
		if stats, err := m.Storage.GetStorageStats(); err != nil {
			results["storage"] = err
			log.Printf("存储服务测试失败: %v", err)
		} else {
			results["storage"] = nil
			log.Printf("存储服务测试成功: %v", stats)
		}
	}

	// 测试缓存连接
	if m.Cache != nil {
		if err := m.Cache.TestConnection(); err != nil {
			results["cache"] = err
			log.Printf("缓存连接测试失败: %v", err)
		} else {
			results["cache"] = nil
			log.Println("缓存连接测试成功")
		}
	}

	return results
}

// GetServiceStatus 获取服务状态
func (m *ServiceManager) GetServiceStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// SMTP状态
	if m.SMTP != nil {
		config := m.SMTP.GetSMTPConfig()
		status["smtp"] = map[string]interface{}{
			"enabled": true,
			"host":    config.Host,
			"port":    config.Port,
			"useTLS":  config.UseTLS,
		}
	} else {
		status["smtp"] = map[string]interface{}{
			"enabled": false,
		}
	}

	// IMAP状态
	if m.IMAP != nil {
		status["imap"] = map[string]interface{}{
			"enabled": true,
		}
	} else {
		status["imap"] = map[string]interface{}{
			"enabled": false,
		}
	}

	// SMS状态
	if m.SMS != nil {
		config := m.SMS.GetSMSConfig()
		status["sms"] = map[string]interface{}{
			"enabled":  true,
			"provider": config.Provider,
		}
	} else {
		status["sms"] = map[string]interface{}{
			"enabled": false,
		}
	}

	// 存储状态
	if m.Storage != nil {
		stats, _ := m.Storage.GetStorageStats()
		status["storage"] = map[string]interface{}{
			"enabled": true,
			"stats":   stats,
		}
	} else {
		status["storage"] = map[string]interface{}{
			"enabled": false,
		}
	}

	// 缓存状态
	if m.Cache != nil {
		stats, _ := m.Cache.GetCacheStats()
		status["cache"] = map[string]interface{}{
			"enabled": true,
			"stats":   stats,
		}
	} else {
		status["cache"] = map[string]interface{}{
			"enabled": false,
		}
	}

	return status
}

// Close 关闭所有服务
func (m *ServiceManager) Close() error {
	var errors []error

	// 关闭IMAP连接
	if m.IMAP != nil {
		if err := m.IMAP.Disconnect(); err != nil {
			errors = append(errors, fmt.Errorf("关闭IMAP连接失败: %v", err))
		}
	}

	// 关闭缓存连接
	if m.Cache != nil {
		if err := m.Cache.Close(); err != nil {
			errors = append(errors, fmt.Errorf("关闭缓存连接失败: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("关闭服务时发生错误: %v", errors)
	}

	return nil
}

// SendVerificationEmail 发送验证码邮件
func (m *ServiceManager) SendVerificationEmail(to, code string) error {
	if m.SMTP == nil {
		return fmt.Errorf("SMTP服务未配置")
	}

	message := EmailMessage{
		From:        "noreply@example.com", // 应该从配置中获取
		To:          []string{to},
		Subject:     "邮箱验证码",
		Body:        fmt.Sprintf("您的验证码是: %s，有效期5分钟。", code),
		ContentType: "text/plain",
	}

	return m.SMTP.SendEmail(message)
}

// FetchEmails 获取邮件
func (m *ServiceManager) FetchEmails(mailbox string, limit uint32) ([]*IMAPEmail, error) {
	if m.IMAP == nil {
		return nil, fmt.Errorf("IMAP服务未配置")
	}

	if err := m.IMAP.Connect(); err != nil {
		return nil, err
	}
	defer m.IMAP.Disconnect()

	return m.IMAP.FetchEmails(mailbox, limit)
}

// SendVerificationSMS 发送验证码短信
func (m *ServiceManager) SendVerificationSMS(phone, code string) (*SMSResponse, error) {
	if m.SMS == nil {
		return nil, fmt.Errorf("SMS服务未配置")
	}
	return m.SMS.SendVerificationCode(phone, code)
}
