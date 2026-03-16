package handler

import (
	"github.com/rankgice/new-email/internal/middleware"
	"github.com/rankgice/new-email/internal/model"
	"github.com/rankgice/new-email/internal/result"
	"github.com/rankgice/new-email/internal/service"
	"github.com/rankgice/new-email/internal/svc"
	"github.com/rankgice/new-email/internal/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ApiHandler API处理器
type ApiHandler struct {
	svcCtx *svc.ServiceContext
}

// NewApiHandler 创建API处理器
func NewApiHandler(svcCtx *svc.ServiceContext) *ApiHandler {
	return &ApiHandler{
		svcCtx: svcCtx,
	}
}

// GetEmail API获取邮件
func (h *ApiHandler) GetEmail(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 通过API密钥验证获取用户ID
	// TODO: 实现GetApiUserId方法
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
		return
	}
	permissions, _ := c.Get("permissions")
	permissionString, _ := permissions.(string)
	if !hasAPIPermission(permissionString, "emails:send", "email:send") {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("API密钥缺少发送邮件权限"))
		return
	}

	// 查询邮件详情
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限（只能查看自己的邮件）
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != userId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此邮件"))
		return
	}

	// 返回邮件详情
	resp := types.EmailResp{
		Id:          email.Id,
		MailboxId:   email.MailboxId,
		Subject:     email.Subject,
		FromEmail:   email.FromEmail,
		ToEmail:     email.ToEmails,  // 使用ToEmails字段
		CcEmail:     email.CcEmails,  // 使用CcEmails字段
		BccEmail:    email.BccEmails, // 使用BccEmails字段
		Content:     email.Content,
		ContentType: email.ContentType,
		Attachments: "", // Email模型中没有Attachments字段
		Status:      0,  // Email模型中没有Status字段
		Type:        "", // Email模型中没有Type字段
		CreatedAt:   email.CreatedAt,
		UpdatedAt:   email.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

func (h *ApiHandler) SendEmail(c *gin.Context) {
	var req types.EmailSendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
		return
	}

	var mailbox *model.Mailbox
	if req.MailboxId > 0 {
		var err error
		mailbox, err = h.svcCtx.MailboxModel.GetById(req.MailboxId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != userId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
			return
		}
	} else {
		mailboxes, err := h.svcCtx.MailboxModel.GetActiveMailboxes(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if len(mailboxes) == 0 {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("当前用户没有可用的发信邮箱"))
			return
		}
		mailbox = mailboxes[0]
	}

	smtpConfig, err := buildSMTPConfig(h.svcCtx, mailbox)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱凭据不可用于发信"))
		return
	}

	attachments, err := decodeEmailAttachments(req.Attachments)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("附件解析失败"))
		return
	}
	normalizedReq := req
	normalizedReq.ContentType = normalizeEmailContentType(req.ContentType)

	smtpService := service.NewSMTPService(smtpConfig)
	if err := smtpService.SendEmail(service.EmailMessage{
		From:        mailbox.Email,
		To:          normalizedReq.ToEmail,
		Cc:          normalizedReq.CcEmail,
		Bcc:         normalizedReq.BccEmail,
		Subject:     normalizedReq.Subject,
		Body:        normalizedReq.Content,
		ContentType: normalizedReq.ContentType,
		Attachments: attachments,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("邮件发送失败"))
		return
	}

	sentAt := time.Now()
	emailRecord, err := persistSentEmailRecord(h.svcCtx, userId, mailbox, &normalizedReq, sentAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	sendResp := types.EmailSendResp{
		Success: true,
		Message: "邮件发送成功",
		EmailId: emailRecord.Id,
		SentAt:  sentAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(sendResp))
}
