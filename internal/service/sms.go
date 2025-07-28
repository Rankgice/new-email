package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SMSConfig SMS配置
type SMSConfig struct {
	Provider  string `json:"provider"`  // 服务商：aliyun, tencent, twilio
	AccessKey string `json:"accessKey"` // 访问密钥
	SecretKey string `json:"secretKey"` // 密钥
	SignName  string `json:"signName"`  // 签名
	Region    string `json:"region"`    // 地区
}

// SMSService SMS服务
type SMSService struct {
	config SMSConfig
}

// NewSMSService 创建SMS服务
func NewSMSService(config SMSConfig) *SMSService {
	return &SMSService{
		config: config,
	}
}

// SMSMessage SMS消息
type SMSMessage struct {
	Phone      string            `json:"phone"`      // 手机号
	Template   string            `json:"template"`   // 模板ID
	Params     map[string]string `json:"params"`     // 模板参数
	Content    string            `json:"content"`    // 短信内容（自定义内容）
	ScheduleAt *time.Time        `json:"scheduleAt"` // 定时发送时间
}

// SMSResponse SMS响应
type SMSResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId"`
	Error     string `json:"error"`
	Cost      int    `json:"cost"` // 费用（分）
}

// SendSMS 发送短信
func (s *SMSService) SendSMS(message SMSMessage) (*SMSResponse, error) {
	switch s.config.Provider {
	case "aliyun":
		return s.sendAliyunSMS(message)
	case "tencent":
		return s.sendTencentSMS(message)
	case "twilio":
		return s.sendTwilioSMS(message)
	default:
		return s.sendMockSMS(message)
	}
}

// sendMockSMS 发送模拟短信（用于测试）
func (s *SMSService) sendMockSMS(message SMSMessage) (*SMSResponse, error) {
	// 模拟发送延迟
	time.Sleep(100 * time.Millisecond)

	// 模拟成功响应
	return &SMSResponse{
		Success:   true,
		MessageID: fmt.Sprintf("mock_%d", time.Now().Unix()),
		Cost:      5, // 5分钱
	}, nil
}

// sendAliyunSMS 发送阿里云短信
func (s *SMSService) sendAliyunSMS(message SMSMessage) (*SMSResponse, error) {
	// 阿里云短信API实现
	// 这里是简化的实现，实际需要使用阿里云SDK

	url := "https://dysmsapi.aliyuncs.com/"

	// 构建请求参数
	params := map[string]interface{}{
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         s.config.Region,
		"PhoneNumbers":     message.Phone,
		"SignName":         s.config.SignName,
		"TemplateCode":     message.Template,
		"TemplateParam":    message.Params,
		"AccessKeyId":      s.config.AccessKey,
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Format":           "JSON",
	}

	// 发送HTTP请求
	resp, err := s.sendHTTPRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	if code, ok := result["Code"].(string); ok && code == "OK" {
		return &SMSResponse{
			Success:   true,
			MessageID: result["BizId"].(string),
			Cost:      5,
		}, nil
	} else {
		return &SMSResponse{
			Success: false,
			Error:   result["Message"].(string),
		}, nil
	}
}

// sendTencentSMS 发送腾讯云短信
func (s *SMSService) sendTencentSMS(message SMSMessage) (*SMSResponse, error) {
	// 腾讯云短信API实现
	// 这里是简化的实现，实际需要使用腾讯云SDK

	url := "https://sms.tencentcloudapi.com/"

	// 构建请求参数
	params := map[string]interface{}{
		"Action":           "SendSms",
		"Version":          "2021-01-11",
		"Region":           s.config.Region,
		"PhoneNumberSet":   []string{message.Phone},
		"SmsSdkAppId":      s.config.AccessKey,
		"SignName":         s.config.SignName,
		"TemplateId":       message.Template,
		"TemplateParamSet": s.mapToSlice(message.Params),
	}

	// 发送HTTP请求
	resp, err := s.sendHTTPRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	if response, ok := result["Response"].(map[string]interface{}); ok {
		if sendStatusSet, ok := response["SendStatusSet"].([]interface{}); ok && len(sendStatusSet) > 0 {
			status := sendStatusSet[0].(map[string]interface{})
			if code, ok := status["Code"].(string); ok && code == "Ok" {
				return &SMSResponse{
					Success:   true,
					MessageID: status["SerialNo"].(string),
					Cost:      5,
				}, nil
			} else {
				return &SMSResponse{
					Success: false,
					Error:   status["Message"].(string),
				}, nil
			}
		}
	}

	return &SMSResponse{
		Success: false,
		Error:   "未知错误",
	}, nil
}

// sendTwilioSMS 发送Twilio短信
func (s *SMSService) sendTwilioSMS(message SMSMessage) (*SMSResponse, error) {
	// Twilio短信API实现
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.config.AccessKey)

	// 构建请求参数
	params := map[string]interface{}{
		"From": s.config.SignName,
		"To":   message.Phone,
		"Body": message.Content,
	}

	// 发送HTTP请求
	resp, err := s.sendHTTPRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	if sid, ok := result["sid"].(string); ok {
		return &SMSResponse{
			Success:   true,
			MessageID: sid,
			Cost:      10, // Twilio相对较贵
		}, nil
	} else {
		return &SMSResponse{
			Success: false,
			Error:   result["message"].(string),
		}, nil
	}
}

// sendHTTPRequest 发送HTTP请求
func (s *SMSService) sendHTTPRequest(method, url string, params map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// mapToSlice 将map转换为slice
func (s *SMSService) mapToSlice(m map[string]string) []string {
	var result []string
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// TestConnection 测试SMS服务连接
func (s *SMSService) TestConnection() error {
	// 发送测试短信
	testMessage := SMSMessage{
		Phone:   "13800138000", // 测试号码
		Content: "SMS服务测试",
	}

	resp, err := s.SendSMS(testMessage)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("SMS测试失败: %s", resp.Error)
	}

	return nil
}

// GetSMSConfig 获取SMS配置
func (s *SMSService) GetSMSConfig() SMSConfig {
	// 返回配置副本，隐藏密钥
	config := s.config
	config.SecretKey = "***"
	return config
}

// SendVerificationCode 发送验证码
func (s *SMSService) SendVerificationCode(phone, code string) (*SMSResponse, error) {
	message := SMSMessage{
		Phone:    phone,
		Template: "SMS_VERIFICATION", // 验证码模板
		Params: map[string]string{
			"code": code,
		},
	}

	return s.SendSMS(message)
}

// SendNotification 发送通知短信
func (s *SMSService) SendNotification(phone, content string) (*SMSResponse, error) {
	message := SMSMessage{
		Phone:   phone,
		Content: content,
	}

	return s.SendSMS(message)
}
