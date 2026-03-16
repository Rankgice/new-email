package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rankgice/new-email/internal/model"
	"github.com/rankgice/new-email/internal/service"
	"github.com/rankgice/new-email/internal/svc"
	"github.com/rankgice/new-email/internal/types"
	"gorm.io/gorm"
)

func isHashedMailboxPassword(password string) bool {
	return strings.HasPrefix(password, "$2a$") ||
		strings.HasPrefix(password, "$2b$") ||
		strings.HasPrefix(password, "$2y$") ||
		strings.HasPrefix(password, "$argon2")
}

func resolveMailboxCredentials(mailbox *model.Mailbox, fallbackUsername, fallbackPassword string) (string, string, error) {
	username := mailbox.Email
	if username == "" {
		username = fallbackUsername
	}

	password := mailbox.Password
	if password == "" {
		password = fallbackPassword
	}

	if isHashedMailboxPassword(password) {
		if mailbox.Email == fallbackUsername && fallbackPassword != "" {
			return fallbackUsername, fallbackPassword, nil
		}
		return "", "", fmt.Errorf("mailbox password is stored as a hash and cannot be used for external server login")
	}

	if username == "" || password == "" {
		return "", "", fmt.Errorf("mailbox credentials are incomplete")
	}

	return username, password, nil
}

func buildSMTPConfig(svcCtx *svc.ServiceContext, mailbox *model.Mailbox) (service.SMTPConfig, error) {
	username, password, err := resolveMailboxCredentials(mailbox, svcCtx.Config.SMTP.Username, svcCtx.Config.SMTP.Password)
	if err != nil {
		return service.SMTPConfig{}, err
	}

	return service.SMTPConfig{
		Host:     svcCtx.Config.SMTP.Host,
		Port:     svcCtx.Config.SMTP.Port,
		Username: username,
		Password: password,
		UseTLS:   svcCtx.Config.SMTP.UseTLS,
	}, nil
}

func buildIMAPConfig(svcCtx *svc.ServiceContext, mailbox *model.Mailbox) (service.IMAPConfig, error) {
	username, password, err := resolveMailboxCredentials(mailbox, svcCtx.Config.IMAP.Username, svcCtx.Config.IMAP.Password)
	if err != nil {
		return service.IMAPConfig{}, err
	}

	return service.IMAPConfig{
		Host:     svcCtx.Config.IMAP.Host,
		Port:     svcCtx.Config.IMAP.Port,
		Username: username,
		Password: password,
		UseTLS:   svcCtx.Config.IMAP.UseTLS,
	}, nil
}

func normalizeEmailContentType(contentType string) string {
	normalized := strings.ToLower(strings.TrimSpace(contentType))
	switch normalized {
	case "", "text", "plain", "text/plain":
		return "text/plain"
	case "html", "text/html":
		return "text/html"
	default:
		return contentType
	}
}

func normalizeMessageID(raw string) string {
	normalized := strings.TrimSpace(raw)
	for len(normalized) >= 2 && strings.HasPrefix(normalized, "<") && strings.HasSuffix(normalized, ">") {
		normalized = strings.TrimSpace(normalized[1 : len(normalized)-1])
	}
	return normalized
}

func normalizedIMAPMessageID(imapEmail *service.IMAPEmail) string {
	if imapEmail == nil {
		return ""
	}
	if normalized := normalizeMessageID(imapEmail.MessageID); normalized != "" {
		return normalized
	}
	if imapEmail.UID > 0 {
		return fmt.Sprintf("imap-uid:%d", imapEmail.UID)
	}
	if !imapEmail.Date.IsZero() {
		return fmt.Sprintf("imap:%s:%s:%s", imapEmail.From, imapEmail.Subject, imapEmail.Date.UTC().Format(time.RFC3339Nano))
	}
	return fmt.Sprintf("imap:%s:%s", imapEmail.From, imapEmail.Subject)
}

func findExistingEmailByNormalizedMessageID(svcCtx *svc.ServiceContext, mailboxID int64, normalizedMessageID string) (*model.Email, error) {
	if normalizedMessageID == "" {
		return nil, nil
	}

	candidates := []string{normalizedMessageID}
	wrapped := normalizedMessageID
	for i := 0; i < 8; i++ {
		wrapped = fmt.Sprintf("<%s>", wrapped)
		candidates = append(candidates, wrapped)
	}

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}

		existing, err := svcCtx.EmailModel.GetByMailboxIdAndMessageId(mailboxID, candidate)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	return nil, nil
}

func hasAPIPermission(rawPermissions string, required ...string) bool {
	trimmed := strings.TrimSpace(rawPermissions)
	if trimmed == "" {
		return false
	}

	requiredSet := make(map[string]struct{}, len(required))
	for _, permission := range required {
		requiredSet[strings.ToLower(strings.TrimSpace(permission))] = struct{}{}
	}

	matches := func(permission string) bool {
		normalized := strings.ToLower(strings.Trim(strings.TrimSpace(permission), `"[]`))
		if normalized == "" {
			return false
		}
		if normalized == "*" || normalized == "all" || normalized == "admin" {
			return true
		}
		_, ok := requiredSet[normalized]
		return ok
	}

	if strings.HasPrefix(trimmed, "[") {
		var permissions []string
		if err := json.Unmarshal([]byte(trimmed), &permissions); err == nil {
			for _, permission := range permissions {
				if matches(permission) {
					return true
				}
			}
		}
	}

	for _, permission := range strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\n' || r == '\t'
	}) {
		if matches(permission) {
			return true
		}
	}

	for permission := range requiredSet {
		if strings.Contains(strings.ToLower(trimmed), permission) {
			return true
		}
	}

	return false
}

func ensureMailboxFolder(svcCtx *svc.ServiceContext, mailboxId int64, name string, isSystem bool) (*model.Folder, error) {
	folderModel := model.NewFolderModel(svcCtx.DB)
	folder, err := folderModel.GetByMailboxIdAndName(mailboxId, name, nil)
	if err != nil {
		return nil, err
	}
	if folder != nil {
		return folder, nil
	}

	folder = &model.Folder{
		MailboxId: mailboxId,
		Name:      name,
		IsSystem:  isSystem,
	}
	if err := folderModel.Create(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func decodeEmailAttachments(attachments []types.AttachmentData) ([]service.EmailAttachment, error) {
	decoded := make([]service.EmailAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		data, err := base64.StdEncoding.DecodeString(attachment.Data)
		if err != nil {
			return nil, fmt.Errorf("decode attachment %s: %w", attachment.Filename, err)
		}

		decoded = append(decoded, service.EmailAttachment{
			Filename:    attachment.Filename,
			ContentType: attachment.ContentType,
			Data:        data,
		})
	}

	return decoded, nil
}

func persistEmailAttachments(svcCtx *svc.ServiceContext, emailId int64, attachments []types.AttachmentData) error {
	for _, attachment := range attachments {
		if err := svcCtx.EmailAttachmentModel.Create(&model.EmailAttachment{
			EmailId:   emailId,
			Filename:  attachment.Filename,
			FilePath:  "",
			FileSize:  attachment.Size,
			MimeType:  attachment.ContentType,
			CreatedAt: time.Now(),
		}); err != nil {
			return err
		}
	}

	return nil
}

func persistSentEmailRecord(svcCtx *svc.ServiceContext, userId int64, mailbox *model.Mailbox, req *types.EmailSendReq, sentAt time.Time) (*model.Email, error) {
	var emailRecord *model.Email
	err := svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		folderModel := model.NewFolderModel(tx)
		sentFolder, err := folderModel.GetByMailboxIdAndName(mailbox.Id, "Sent", nil)
		if err != nil {
			return err
		}
		if sentFolder == nil {
			sentFolder = &model.Folder{
				MailboxId: mailbox.Id,
				Name:      "Sent",
				IsSystem:  true,
			}
			if err := folderModel.Create(sentFolder); err != nil {
				return err
			}
		}

		emailRecord = &model.Email{
			UserId:      userId,
			MailboxId:   mailbox.Id,
			FolderId:    sentFolder.Id,
			Subject:     req.Subject,
			FromEmail:   mailbox.Email,
			ToEmails:    req.ToEmail,
			CcEmails:    req.CcEmail,
			BccEmails:   req.BccEmail,
			Content:     req.Content,
			ContentType: normalizeEmailContentType(req.ContentType),
			Direction:   "sent",
			IsRead:      true,
			SentAt:      &sentAt,
		}

		if err := model.NewEmailModel(tx).Create(emailRecord); err != nil {
			return err
		}

		attachmentModel := model.NewEmailAttachmentModel(tx)
		for _, attachment := range req.Attachments {
			if err := attachmentModel.Create(&model.EmailAttachment{
				EmailId:   emailRecord.Id,
				Filename:  attachment.Filename,
				FilePath:  "",
				FileSize:  attachment.Size,
				MimeType:  attachment.ContentType,
				CreatedAt: time.Now(),
			}); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return emailRecord, nil
}

func persistReceivedEmail(svcCtx *svc.ServiceContext, mailbox *model.Mailbox, imapEmail *service.IMAPEmail, content string) (*model.Email, error) {
	inboxFolder, err := ensureMailboxFolder(svcCtx, mailbox.Id, "INBOX", true)
	if err != nil {
		return nil, err
	}

	contentType := normalizeEmailContentType(imapEmail.ContentType)

	receivedAt := imapEmail.Date
	if receivedAt.IsZero() {
		receivedAt = time.Now()
	}

	email := &model.Email{
		UserId:      mailbox.UserId,
		MailboxId:   mailbox.Id,
		FolderId:    inboxFolder.Id,
		MessageId:   normalizedIMAPMessageID(imapEmail),
		Subject:     imapEmail.Subject,
		FromEmail:   imapEmail.From,
		ToEmails:    imapEmail.To,
		CcEmails:    imapEmail.Cc,
		BccEmails:   imapEmail.Bcc,
		Content:     content,
		ContentType: contentType,
		IsRead:      imapEmail.IsRead,
		Direction:   "received",
		ReceivedAt:  &receivedAt,
	}

	if err := svcCtx.EmailModel.Create(email); err != nil {
		return nil, err
	}

	return email, nil
}
