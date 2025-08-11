package service

import (
	"fmt"
	"io"
	"log"
	"mime"
	"strings"
	"time"

	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

// MessageParser 邮件解析器
type MessageParser struct{}

// ParsedMessage 解析后的邮件
type ParsedMessage struct {
	MessageID   string             `json:"message_id"`
	Subject     string             `json:"subject"`
	From        string             `json:"from"`
	To          []string           `json:"to"`
	Cc          []string           `json:"cc"`
	Bcc         []string           `json:"bcc"`
	Date        time.Time          `json:"date"`
	ContentType string             `json:"content_type"`
	Body        string             `json:"body"`
	HTMLBody    string             `json:"html_body"`
	TextBody    string             `json:"text_body"`
	Attachments []ParsedAttachment `json:"attachments"`
	Headers     map[string]string  `json:"headers"`
}

// ParsedAttachment 解析后的附件
type ParsedAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Data        []byte `json:"data"`
}

// NewMessageParser 创建邮件解析器
func NewMessageParser() *MessageParser {
	return &MessageParser{}
}

// ParseMessage 解析邮件消息
func (p *MessageParser) ParseMessage(reader io.Reader) (*ParsedMessage, error) {
	// 创建邮件读取器
	mr, err := mail.CreateReader(reader)
	if err != nil {
		return nil, fmt.Errorf("创建邮件读取器失败: %v", err)
	}

	parsed := &ParsedMessage{
		Headers:     make(map[string]string),
		Attachments: []ParsedAttachment{},
	}

	// 解析邮件头
	if err := p.parseHeadersFromReader(mr, parsed); err != nil {
		log.Printf("解析邮件头失败: %v", err)
	}

	// 解析邮件体
	if err := p.parseBodyFromReader(mr, parsed); err != nil {
		log.Printf("解析邮件体失败: %v", err)
	}

	return parsed, nil
}

// parseHeadersFromReader 从邮件读取器解析邮件头
func (p *MessageParser) parseHeadersFromReader(mr *mail.Reader, parsed *ParsedMessage) error {
	header := mr.Header

	// 基本信息
	parsed.MessageID = header.Get("Message-ID")
	parsed.Subject = header.Get("Subject")
	parsed.ContentType = header.Get("Content-Type")

	// 发件人
	if from, err := header.AddressList("From"); err == nil && len(from) > 0 {
		parsed.From = from[0].String()
	}

	// 收件人
	if to, err := header.AddressList("To"); err == nil {
		for _, addr := range to {
			parsed.To = append(parsed.To, addr.String())
		}
	}

	// 抄送
	if cc, err := header.AddressList("Cc"); err == nil {
		for _, addr := range cc {
			parsed.Cc = append(parsed.Cc, addr.String())
		}
	}

	// 密送
	if bcc, err := header.AddressList("Bcc"); err == nil {
		for _, addr := range bcc {
			parsed.Bcc = append(parsed.Bcc, addr.String())
		}
	}

	// 日期
	if date, err := header.Date(); err == nil {
		parsed.Date = date
	} else {
		parsed.Date = time.Now()
	}

	// 保存所有头信息
	fields := header.Fields()
	for fields.Next() {
		parsed.Headers[fields.Key()] = fields.Value()
	}

	return nil
}

// parseBodyFromReader 从邮件读取器解析邮件体
func (p *MessageParser) parseBodyFromReader(mr *mail.Reader, parsed *ParsedMessage) error {
	// 遍历邮件的所有部分
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("读取邮件部分失败: %v", err)
			continue
		}

		if err := p.parsePartFromReader(part, parsed); err != nil {
			log.Printf("解析邮件部分失败: %v", err)
		}
	}

	// 如果没有设置主体，使用文本体
	if parsed.Body == "" {
		if parsed.HTMLBody != "" {
			parsed.Body = parsed.HTMLBody
			parsed.ContentType = "text/html"
		} else if parsed.TextBody != "" {
			parsed.Body = parsed.TextBody
			parsed.ContentType = "text/plain"
		}
	}

	return nil
}

// parsePartFromReader 解析邮件部分
func (p *MessageParser) parsePartFromReader(part *mail.Part, parsed *ParsedMessage) error {
	contentType := part.Header.Get("Content-Type")
	disposition := part.Header.Get("Content-Disposition")

	// 读取内容
	content, err := io.ReadAll(part.Body)
	if err != nil {
		return err
	}

	// 判断是否是附件
	if strings.HasPrefix(disposition, "attachment") || strings.HasPrefix(disposition, "inline") {
		return p.parseAttachmentFromPart(part, content, parsed)
	}

	// 解析邮件正文
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = "text/plain"
	}

	contentStr := string(content)

	switch mediaType {
	case "text/plain":
		if parsed.TextBody == "" {
			parsed.TextBody = contentStr
		}
	case "text/html":
		if parsed.HTMLBody == "" {
			parsed.HTMLBody = contentStr
		}
	default:
		// 其他类型作为附件处理
		return p.parseAttachmentFromPart(part, content, parsed)
	}

	return nil
}

// parseAttachmentFromPart 解析附件
func (p *MessageParser) parseAttachmentFromPart(part *mail.Part, content []byte, parsed *ParsedMessage) error {
	attachment := ParsedAttachment{
		ContentType: part.Header.Get("Content-Type"),
		Size:        int64(len(content)),
		Data:        content,
	}

	// 获取文件名
	disposition := part.Header.Get("Content-Disposition")
	if disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			if filename, ok := params["filename"]; ok {
				attachment.Filename = filename
			}
		}
	}

	// 如果没有从Content-Disposition获取到文件名，尝试从Content-Type获取
	if attachment.Filename == "" {
		contentType := part.Header.Get("Content-Type")
		if contentType != "" {
			_, params, err := mime.ParseMediaType(contentType)
			if err == nil {
				if name, ok := params["name"]; ok {
					attachment.Filename = name
				}
			}
		}
	}

	// 如果还是没有文件名，生成一个默认的
	if attachment.Filename == "" {
		attachment.Filename = fmt.Sprintf("attachment_%d", len(parsed.Attachments)+1)
	}

	parsed.Attachments = append(parsed.Attachments, attachment)
	return nil
}

// ExtractTextContent 提取邮件的纯文本内容
func (p *MessageParser) ExtractTextContent(parsed *ParsedMessage) string {
	if parsed.TextBody != "" {
		return parsed.TextBody
	}
	if parsed.HTMLBody != "" {
		// 简单的HTML标签移除（实际项目中可能需要更复杂的HTML解析）
		text := parsed.HTMLBody
		text = strings.ReplaceAll(text, "<br>", "\n")
		text = strings.ReplaceAll(text, "<br/>", "\n")
		text = strings.ReplaceAll(text, "<br />", "\n")
		text = strings.ReplaceAll(text, "</p>", "\n")
		text = strings.ReplaceAll(text, "</div>", "\n")

		// 移除HTML标签（简单实现）
		for strings.Contains(text, "<") && strings.Contains(text, ">") {
			start := strings.Index(text, "<")
			end := strings.Index(text[start:], ">")
			if end == -1 {
				break
			}
			text = text[:start] + text[start+end+1:]
		}

		return strings.TrimSpace(text)
	}
	return parsed.Body
}
