package handler

import (
	"fmt"
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

// MailboxHandler 邮箱处理器
type MailboxHandler struct {
	svcCtx *svc.ServiceContext
}

// NewMailboxHandler 创建邮箱处理器
func NewMailboxHandler(svcCtx *svc.ServiceContext) *MailboxHandler {
	return &MailboxHandler{
		svcCtx: svcCtx,
	}
}

// List 邮箱列表
func (h *MailboxHandler) List(c *gin.Context) {
	var req types.MailboxListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 如果不是管理员，只能查看自己的邮箱
	if req.UserId == 0 {
		req.UserId = currentUserId
	} else if req.UserId != currentUserId {
		// 这里可以添加管理员权限检查
		// 暂时只允许查看自己的邮箱
		req.UserId = currentUserId
	}

	// 转换为model参数
	params := model.MailboxListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.StartTime,
			CreatedAtEnd:   req.EndTime,
		},
		UserId:      req.UserId,
		DomainId:    req.DomainId,
		Email:       req.Email,
		Status:      req.Status,
		AutoReceive: req.AutoReceive,
	}

	// 查询邮箱列表
	mailboxes, total, err := h.svcCtx.MailboxModel.List(params)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var mailboxList []types.MailboxResp
	for _, mailbox := range mailboxes {
		mailboxList = append(mailboxList, types.MailboxResp{
			Id:          mailbox.Id,
			UserId:      mailbox.UserId,
			DomainId:    mailbox.DomainId,
			Email:       mailbox.Email,
			AutoReceive: mailbox.AutoReceive,
			Status:      mailbox.Status,
			LastSyncAt:  mailbox.LastSyncAt,
			CreatedAt:   mailbox.CreatedAt,
			UpdatedAt:   mailbox.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     mailboxList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建邮箱
func (h *MailboxHandler) Create(c *gin.Context) {
	var req types.MailboxCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 检查邮箱是否已存在
	if exists, err := h.svcCtx.MailboxModel.CheckEmailExists(req.Email); err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱已存在"))
		return
	}

	// 如果没有提供域名ID，使用默认域名
	domainId := req.DomainId
	if domainId == 0 {
		// 获取第一个可用的域名作为默认域名
		domains, err := h.svcCtx.DomainModel.GetActiveDomains()
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		}
		if len(domains) == 0 {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("系统中没有可用的域名，请联系管理员"))
			return
		}
		domainId = domains[0].Id
	}

	// 检查域名是否存在且启用
	domain, err := h.svcCtx.DomainModel.GetById(domainId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if domain == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("域名不存在"))
		return
	}
	if domain.Status != 1 {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("域名未启用"))
		return
	}
	// 创建邮箱
	mailbox := &model.Mailbox{
		UserId:      currentUserId,
		DomainId:    domainId,
		Email:       req.Email,
		Password:    req.Password,
		AutoReceive: req.AutoReceive,
		Status:      req.Status,
	}

	if err := h.svcCtx.MailboxModel.Create(mailbox); err != nil {
		c.JSON(http.StatusOK, result.ErrorAdd.AddError(err))
		return
	}

	// 返回创建的邮箱信息
	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"id":    mailbox.Id,
		"email": mailbox.Email,
	}))
}

// Update 更新邮箱
func (h *MailboxHandler) Update(c *gin.Context) {
	var req types.MailboxUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 检查邮箱是否存在
	mailbox, err := h.svcCtx.MailboxModel.GetById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限（只能更新自己的邮箱）
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	// 检查邮箱地址是否已被其他邮箱使用
	if req.Email != "" && req.Email != mailbox.Email {
		if exists, err := h.svcCtx.MailboxModel.CheckEmailExists(req.Email); err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱地址已被其他邮箱使用"))
			return
		}
	}

	// 如果更新域名，检查域名是否存在
	if req.DomainId > 0 && req.DomainId != mailbox.DomainId {
		domain, err := h.svcCtx.DomainModel.GetById(req.DomainId)
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		}
		if domain == nil {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("域名不存在"))
			return
		}
		if domain.Status != 1 {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("域名未启用"))
			return
		}
	}

	// 构建更新数据
	updateData := map[string]interface{}{}
	if req.DomainId > 0 {
		updateData["domain_id"] = req.DomainId
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Password != "" {
		updateData["password"] = req.Password
	}
	updateData["auto_receive"] = req.AutoReceive
	updateData["status"] = req.Status

	// 更新邮箱信息
	if err := h.svcCtx.MailboxModel.MapUpdate(nil, req.Id, updateData); err != nil {
		c.JSON(http.StatusOK, result.ErrorUpdate.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// Delete 删除邮箱
func (h *MailboxHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮箱ID"))
		return
	}

	mailboxId := int64(id)

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 检查邮箱是否存在
	mailbox, err := h.svcCtx.MailboxModel.GetById(mailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限（只能删除自己的邮箱）
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	// 检查是否有关联的邮件（可选：防止删除有邮件的邮箱）
	// 这里可以添加检查逻辑

	// 软删除邮箱
	if err := h.svcCtx.MailboxModel.Delete(mailbox); err != nil {
		c.JSON(http.StatusOK, result.ErrorDelete.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

func (h *MailboxHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮箱ID"))
		return
	}

	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	mailbox, err := h.svcCtx.MailboxModel.GetById(id)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱不存在"))
		return
	}
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	resp := types.MailboxConnectionTestResp{}

	imapConfig, err := buildIMAPConfig(h.svcCtx, mailbox)
	if err != nil {
		resp.ImapError = "IMAP配置不可用"
	} else {
		imapService := service.NewIMAPService(imapConfig)
		if err := imapService.TestConnection(); err != nil {
			resp.ImapError = "IMAP连接失败"
		} else {
			resp.ImapSuccess = true
		}
	}

	smtpConfig, err := buildSMTPConfig(h.svcCtx, mailbox)
	if err != nil {
		resp.SmtpError = "SMTP配置不可用"
	} else {
		smtpService := service.NewSMTPService(smtpConfig)
		if err := smtpService.TestConnection(); err != nil {
			resp.SmtpError = "SMTP连接失败"
		} else {
			resp.SmtpSuccess = true
		}
	}

	resp.Success = resp.ImapSuccess && resp.SmtpSuccess
	switch {
	case resp.Success:
		resp.Message = "连接测试成功"
	case resp.ImapSuccess || resp.SmtpSuccess:
		resp.Message = "连接测试部分成功"
	default:
		resp.Message = "连接测试失败"
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

func (h *MailboxHandler) Sync(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮箱ID"))
		return
	}

	var req struct {
		ForceSync bool `json:"forceSync"`
		SyncDays  int  `json:"syncDays"`
	}
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
			return
		}
	}
	if req.SyncDays < 0 {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("syncDays 不能小于 0"))
		return
	}

	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	mailbox, err := h.svcCtx.MailboxModel.GetById(id)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限（只能同步自己的邮箱）
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	imapConfig, err := buildIMAPConfig(h.svcCtx, mailbox)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱同步配置不可用"))
		return
	}

	imapService := service.NewIMAPService(imapConfig)
	if err := imapService.Connect(); err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("连接IMAP服务器失败"))
		return
	}
	defer imapService.Disconnect()

	limit := uint32(h.svcCtx.Config.Email.Receive.BatchSize)
	if limit == 0 {
		limit = 100
	}

	imapEmails, err := imapService.FetchEmails("INBOX", limit)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("获取邮件列表失败"))
		return
	}

	cutoff := time.Time{}
	if req.SyncDays > 0 {
		cutoff = time.Now().AddDate(0, 0, -req.SyncDays)
	}

	syncCount := 0
	errorCount := 0
	for _, imapEmail := range imapEmails {
		if imapEmail == nil {
			continue
		}
		if !cutoff.IsZero() && !imapEmail.Date.IsZero() && imapEmail.Date.Before(cutoff) {
			continue
		}

		messageID := normalizedIMAPMessageID(imapEmail)
		existing, err := findExistingEmailByNormalizedMessageID(h.svcCtx, mailbox.Id, messageID)
		if err != nil {
			errorCount++
			continue
		}
		if existing != nil && existing.MessageId != messageID && messageID != "" {
			if err := h.svcCtx.EmailModel.MapUpdate(nil, existing.Id, map[string]interface{}{
				"message_id": messageID,
				"updated_at": time.Now(),
			}); err != nil {
				errorCount++
				continue
			}
			existing.MessageId = messageID
		}
		if existing != nil && !req.ForceSync {
			continue
		}

		content := imapEmail.Body
		if content == "" && imapEmail.UID > 0 {
			if body, err := imapService.FetchEmailBody(imapEmail.UID); err == nil {
				content = body
			}
		}

		if existing != nil {
			updateData := map[string]interface{}{
				"message_id":   messageID,
				"subject":      imapEmail.Subject,
				"from_email":   imapEmail.From,
				"to_emails":    imapEmail.To,
				"cc_emails":    imapEmail.Cc,
				"bcc_emails":   imapEmail.Bcc,
				"content":      content,
				"content_type": normalizeEmailContentType(imapEmail.ContentType),
				"is_read":      imapEmail.IsRead,
				"received_at":  imapEmail.Date,
				"updated_at":   time.Now(),
			}
			if err := h.svcCtx.EmailModel.MapUpdate(nil, existing.Id, updateData); err != nil {
				errorCount++
				continue
			}
			syncCount++
			continue
		}

		if _, err := persistReceivedEmail(h.svcCtx, mailbox, imapEmail, content); err != nil {
			errorCount++
			continue
		}
		syncCount++
	}

	now := time.Now()
	_ = h.svcCtx.MailboxModel.MapUpdate(nil, id, map[string]interface{}{
		"last_sync_at": &now,
	})

	resp := types.MailboxSyncResp{
		Success:    errorCount == 0,
		SyncCount:  syncCount,
		ErrorCount: errorCount,
		LastSyncAt: now,
	}
	switch {
	case syncCount == 0 && errorCount == 0:
		resp.Message = "同步完成，暂无新邮件"
	case errorCount == 0:
		resp.Message = fmt.Sprintf("同步成功，共同步 %d 封邮件", syncCount)
	case syncCount > 0:
		resp.Success = true
		resp.Message = fmt.Sprintf("同步完成，成功 %d 封，失败 %d 封", syncCount, errorCount)
	default:
		resp.Message = fmt.Sprintf("同步失败，错误数量 %d", errorCount)
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetById 根据ID获取邮箱信息
func (h *MailboxHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮箱ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	mailbox, err := h.svcCtx.MailboxModel.GetById(id)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	if mailbox == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限（只能查看自己的邮箱）
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限查看此邮箱"))
		return
	}

	resp := types.MailboxResp{
		Id:          mailbox.Id,
		UserId:      mailbox.UserId,
		DomainId:    mailbox.DomainId,
		Email:       mailbox.Email,
		AutoReceive: mailbox.AutoReceive,
		Status:      mailbox.Status,
		LastSyncAt:  mailbox.LastSyncAt,
		CreatedAt:   mailbox.CreatedAt,
		UpdatedAt:   mailbox.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetStats 获取邮箱统计信息
func (h *MailboxHandler) GetStats(c *gin.Context) {
	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 获取用户的邮箱统计
	totalMailboxes, _, err := h.svcCtx.MailboxModel.List(model.MailboxListParams{
		UserId: currentUserId,
	})
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	// 获取活跃邮箱数
	status := 1
	activeMailboxes, _, err := h.svcCtx.MailboxModel.List(model.MailboxListParams{
		UserId: currentUserId,
		Status: &status,
	})
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	resp := types.MailboxStatsResp{
		TotalMailboxes:  int64(len(totalMailboxes)),
		ActiveMailboxes: int64(len(activeMailboxes)),
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
