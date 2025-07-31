package service

import (
	"html"
	"new-email/internal/types"
	"regexp"
	"strconv"
	"strings"
)

// VerificationCodeExtractor 验证码提取器
type VerificationCodeExtractor struct {
	patterns []CodePattern
}

// CodePattern 验证码模式
type CodePattern struct {
	Name        string         // 模式名称
	Pattern     *regexp.Regexp // 正则表达式
	Type        string         // 验证码类型
	Description string         // 描述
	Confidence  int            // 基础置信度
}

// NewVerificationCodeExtractor 创建验证码提取器
func NewVerificationCodeExtractor() *VerificationCodeExtractor {
	extractor := &VerificationCodeExtractor{}
	extractor.initPatterns()
	return extractor
}

// initPatterns 初始化验证码匹配模式
func (e *VerificationCodeExtractor) initPatterns() {
	e.patterns = []CodePattern{
		// 6位数字验证码（最常见）
		{
			Name:        "6位数字验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:验证码|verification\s*code|code|验证|码)[：:\s]*([0-9]{6})`),
			Type:        "numeric_6",
			Description: "6位数字验证码",
			Confidence:  95,
		},
		// 4位数字验证码
		{
			Name:        "4位数字验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:验证码|verification\s*code|code|验证|码)[：:\s]*([0-9]{4})`),
			Type:        "numeric_4",
			Description: "4位数字验证码",
			Confidence:  90,
		},
		// 8位数字验证码
		{
			Name:        "8位数字验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:验证码|verification\s*code|code|验证|码)[：:\s]*([0-9]{8})`),
			Type:        "numeric_8",
			Description: "8位数字验证码",
			Confidence:  85,
		},
		// 混合字母数字验证码（6位）
		{
			Name:        "6位字母数字验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:验证码|verification\s*code|code|验证|码)[：:\s]*([A-Z0-9]{6})`),
			Type:        "alphanumeric_6",
			Description: "6位字母数字验证码",
			Confidence:  80,
		},
		// 混合字母数字验证码（4位）
		{
			Name:        "4位字母数字验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:验证码|verification\s*code|code|验证|码)[：:\s]*([A-Z0-9]{4})`),
			Type:        "alphanumeric_4",
			Description: "4位字母数字验证码",
			Confidence:  75,
		},
		// 纯数字验证码（通用模式）
		{
			Name:        "纯数字验证码",
			Pattern:     regexp.MustCompile(`\b([0-9]{4,8})\b`),
			Type:        "numeric_general",
			Description: "4-8位纯数字验证码",
			Confidence:  60,
		},
		// 短信验证码特殊格式
		{
			Name:        "短信验证码格式",
			Pattern:     regexp.MustCompile(`(?i)(?:您的|your)\s*(?:验证码|verification\s*code)\s*(?:是|is)[：:\s]*([0-9A-Z]{4,8})`),
			Type:        "sms_format",
			Description: "短信验证码格式",
			Confidence:  90,
		},
		// 登录验证码
		{
			Name:        "登录验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:登录|login)\s*(?:验证码|verification\s*code)[：:\s]*([0-9A-Z]{4,8})`),
			Type:        "login_code",
			Description: "登录验证码",
			Confidence:  85,
		},
		// 注册验证码
		{
			Name:        "注册验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:注册|register|registration)\s*(?:验证码|verification\s*code)[：:\s]*([0-9A-Z]{4,8})`),
			Type:        "register_code",
			Description: "注册验证码",
			Confidence:  85,
		},
		// 重置密码验证码
		{
			Name:        "重置密码验证码",
			Pattern:     regexp.MustCompile(`(?i)(?:重置|reset|找回|forgot)\s*(?:密码|password)\s*(?:验证码|verification\s*code)[：:\s]*([0-9A-Z]{4,8})`),
			Type:        "reset_password_code",
			Description: "重置密码验证码",
			Confidence:  85,
		},
		// 动态密码/OTP
		{
			Name:        "动态密码",
			Pattern:     regexp.MustCompile(`(?i)(?:动态密码|otp|one.time.password)[：:\s]*([0-9]{6,8})`),
			Type:        "otp",
			Description: "动态密码/OTP",
			Confidence:  90,
		},
		// 安全码
		{
			Name:        "安全码",
			Pattern:     regexp.MustCompile(`(?i)(?:安全码|security\s*code)[：:\s]*([0-9A-Z]{4,8})`),
			Type:        "security_code",
			Description: "安全码",
			Confidence:  80,
		},
	}
}

// ExtractFromEmail 从邮件中提取验证码
func (e *VerificationCodeExtractor) ExtractFromEmail(subject, content string) []types.VerificationCodeResult {
	var results []types.VerificationCodeResult

	// 清理HTML标签
	cleanContent := e.cleanHTML(content)

	// 合并主题和内容进行匹配
	fullText := subject + "\n" + cleanContent

	// 使用每个模式进行匹配
	for _, pattern := range e.patterns {
		matches := pattern.Pattern.FindAllStringSubmatch(fullText, -1)

		for _, match := range matches {
			if len(match) > 1 {
				code := strings.TrimSpace(match[1])
				if code != "" {
					// 计算位置
					position := strings.Index(fullText, match[0])

					// 提取上下文
					context := e.extractContext(fullText, position, len(match[0]))

					// 计算置信度
					confidence := e.calculateConfidence(pattern, code, context, subject, cleanContent)

					result := types.VerificationCodeResult{
						Code:        code,
						Type:        pattern.Type,
						Context:     context,
						Confidence:  confidence,
						Position:    position,
						Length:      len(code),
						Pattern:     pattern.Name,
						Description: pattern.Description,
					}

					results = append(results, result)
				}
			}
		}
	}

	// 去重和排序
	results = e.deduplicateAndSort(results)

	return results
}

// cleanHTML 清理HTML标签
func (e *VerificationCodeExtractor) cleanHTML(content string) string {
	// 解码HTML实体
	content = html.UnescapeString(content)

	// 移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, " ")

	// 清理多余的空白字符
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	return strings.TrimSpace(content)
}

// extractContext 提取验证码的上下文信息
func (e *VerificationCodeExtractor) extractContext(text string, position, matchLength int) string {
	contextLength := 50 // 前后各50个字符

	start := position - contextLength
	if start < 0 {
		start = 0
	}

	end := position + matchLength + contextLength
	if end > len(text) {
		end = len(text)
	}

	context := text[start:end]

	// 清理换行符和多余空格
	context = strings.ReplaceAll(context, "\n", " ")
	context = regexp.MustCompile(`\s+`).ReplaceAllString(context, " ")

	return strings.TrimSpace(context)
}

// calculateConfidence 计算置信度
func (e *VerificationCodeExtractor) calculateConfidence(pattern CodePattern, code, context, subject, content string) int {
	confidence := pattern.Confidence

	// 根据验证码特征调整置信度
	if len(code) == 6 && isNumeric(code) {
		confidence += 5 // 6位数字验证码最常见
	}

	// 根据上下文关键词调整置信度
	contextLower := strings.ToLower(context)
	if strings.Contains(contextLower, "验证码") || strings.Contains(contextLower, "verification") {
		confidence += 10
	}
	if strings.Contains(contextLower, "登录") || strings.Contains(contextLower, "login") {
		confidence += 5
	}
	if strings.Contains(contextLower, "注册") || strings.Contains(contextLower, "register") {
		confidence += 5
	}

	// 根据邮件主题调整置信度
	subjectLower := strings.ToLower(subject)
	if strings.Contains(subjectLower, "验证码") || strings.Contains(subjectLower, "verification") {
		confidence += 15
	}

	// 根据发件人域名调整置信度（常见的服务提供商）
	commonDomains := []string{"noreply", "no-reply", "service", "support", "notification", "auth"}
	for _, domain := range commonDomains {
		if strings.Contains(strings.ToLower(content), domain) {
			confidence += 5
			break
		}
	}

	// 确保置信度在合理范围内
	if confidence > 100 {
		confidence = 100
	}
	if confidence < 0 {
		confidence = 0
	}

	return confidence
}

// deduplicateAndSort 去重和排序
func (e *VerificationCodeExtractor) deduplicateAndSort(results []types.VerificationCodeResult) []types.VerificationCodeResult {
	// 使用map去重
	seen := make(map[string]bool)
	var unique []types.VerificationCodeResult

	for _, result := range results {
		key := result.Code + "|" + result.Type
		if !seen[key] {
			seen[key] = true
			unique = append(unique, result)
		}
	}

	// 按置信度排序（降序）
	for i := 0; i < len(unique)-1; i++ {
		for j := i + 1; j < len(unique); j++ {
			if unique[i].Confidence < unique[j].Confidence {
				unique[i], unique[j] = unique[j], unique[i]
			}
		}
	}

	return unique
}

// isNumeric 检查字符串是否为纯数字
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// GetSupportedTypes 获取支持的验证码类型
func (e *VerificationCodeExtractor) GetSupportedTypes() []string {
	var types []string
	seen := make(map[string]bool)

	for _, pattern := range e.patterns {
		if !seen[pattern.Type] {
			types = append(types, pattern.Type)
			seen[pattern.Type] = true
		}
	}

	return types
}

// ValidateCode 验证验证码格式
func (e *VerificationCodeExtractor) ValidateCode(code string) (bool, string) {
	if len(code) < 4 || len(code) > 8 {
		return false, "验证码长度应在4-8位之间"
	}

	// 检查是否包含有效字符
	validChars := regexp.MustCompile(`^[0-9A-Za-z]+$`)
	if !validChars.MatchString(code) {
		return false, "验证码只能包含数字和字母"
	}

	return true, ""
}
