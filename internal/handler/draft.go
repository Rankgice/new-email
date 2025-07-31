package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/service"
	"new-email/internal/svc"
	"new-email/internal/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DraftHandler 草稿处理器
type DraftHandler struct {
	svcCtx *svc.ServiceContext
}

// NewDraftHandler 创建草稿处理器
func NewDraftHandler(svcCtx *svc.ServiceContext) *DraftHandler {
	return &DraftHandler{
		svcCtx: svcCtx,
	}
}

// List 草稿列表
func (h *DraftHandler) List(c *gin.Context) {
	var req types.DraftListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 转换为model参数
	params := model.EmailDraftListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId:    currentUserId,
		MailboxId: 0, // 默认值
		Subject:   req.Subject,
		ToEmails:  req.ToEmail, // 注意字段名是ToEmails
	}

	// 处理可选的MailboxId参数
	if req.MailboxId != nil {
		params.MailboxId = *req.MailboxId
	}

	// 查询草稿列表
	drafts, total, err := h.svcCtx.EmailDraftModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var draftList []types.DraftResp
	for _, draft := range drafts {
		draftList = append(draftList, types.DraftResp{
			Id:          draft.Id,
			UserId:      draft.UserId,
			MailboxId:   draft.MailboxId,
			Subject:     draft.Subject,
			FromEmail:   "",              // EmailDraft模型中没有FromEmail字段
			ToEmail:     draft.ToEmails,  // 使用ToEmails字段
			CcEmail:     draft.CcEmails,  // 使用CcEmails字段
			BccEmail:    draft.BccEmails, // 使用BccEmails字段
			Content:     draft.Content,
			ContentType: draft.ContentType,
			Attachments: "", // EmailDraft模型中没有Attachments字段
			CreatedAt:   draft.CreatedAt,
			UpdatedAt:   draft.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     draftList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建草稿
func (h *DraftHandler) Create(c *gin.Context) {
	var req types.DraftCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查邮箱是否属于当前用户
	if req.MailboxId > 0 {
		mailbox, err := h.svcCtx.MailboxModel.GetById(req.MailboxId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != currentUserId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
			return
		}
	}

	// 创建草稿
	draft := &model.EmailDraft{
		UserId:      currentUserId,
		MailboxId:   req.MailboxId,
		Subject:     req.Subject,
		ToEmails:    req.ToEmail,  // 使用ToEmails字段
		CcEmails:    req.CcEmail,  // 使用CcEmails字段
		BccEmails:   req.BccEmail, // 使用BccEmails字段
		Content:     req.Content,
		ContentType: req.ContentType,
		Status:      "draft", // 设置状态为草稿
	}

	if err := h.svcCtx.EmailDraftModel.Create(draft); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_draft",
		Resource:   "draft",
		ResourceId: draft.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的草稿信息
	resp := types.DraftResp{
		Id:          draft.Id,
		UserId:      draft.UserId,
		MailboxId:   draft.MailboxId,
		Subject:     draft.Subject,
		FromEmail:   "",              // EmailDraft模型中没有FromEmail字段
		ToEmail:     draft.ToEmails,  // 使用ToEmails字段
		CcEmail:     draft.CcEmails,  // 使用CcEmails字段
		BccEmail:    draft.BccEmails, // 使用BccEmails字段
		Content:     draft.Content,
		ContentType: draft.ContentType,
		Attachments: "", // EmailDraft模型中没有Attachments字段
		CreatedAt:   draft.CreatedAt,
		UpdatedAt:   draft.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新草稿
func (h *DraftHandler) Update(c *gin.Context) {
	// 获取草稿ID
	idStr := c.Param("id")
	draftId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的草稿ID"))
		return
	}

	var req types.DraftUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查草稿是否存在
	draft, err := h.svcCtx.EmailDraftModel.GetById(draftId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if draft == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("草稿不存在"))
		return
	}

	// 检查权限（只能更新自己的草稿）
	if draft.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此草稿"))
		return
	}

	// 检查邮箱是否属于当前用户
	if req.MailboxId > 0 {
		mailbox, err := h.svcCtx.MailboxModel.GetById(req.MailboxId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != currentUserId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
			return
		}
	}

	// 更新草稿信息
	draft.MailboxId = req.MailboxId
	draft.Subject = req.Subject
	draft.ToEmails = req.ToEmail   // 使用ToEmails字段
	draft.CcEmails = req.CcEmail   // 使用CcEmails字段
	draft.BccEmails = req.BccEmail // 使用BccEmails字段
	draft.Content = req.Content
	draft.ContentType = req.ContentType

	if err := h.svcCtx.EmailDraftModel.Update(draft); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_draft",
		Resource:   "draft",
		ResourceId: draft.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的草稿信息
	resp := types.DraftResp{
		Id:          draft.Id,
		UserId:      draft.UserId,
		MailboxId:   draft.MailboxId,
		Subject:     draft.Subject,
		FromEmail:   "",              // EmailDraft模型中没有FromEmail字段
		ToEmail:     draft.ToEmails,  // 使用ToEmails字段
		CcEmail:     draft.CcEmails,  // 使用CcEmails字段
		BccEmail:    draft.BccEmails, // 使用BccEmails字段
		Content:     draft.Content,
		ContentType: draft.ContentType,
		Attachments: "", // EmailDraft模型中没有Attachments字段
		CreatedAt:   draft.CreatedAt,
		UpdatedAt:   draft.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除草稿
func (h *DraftHandler) Delete(c *gin.Context) {
	// 获取草稿ID
	idStr := c.Param("id")
	draftId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的草稿ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查草稿是否存在
	draft, err := h.svcCtx.EmailDraftModel.GetById(draftId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if draft == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("草稿不存在"))
		return
	}

	// 检查权限（只能删除自己的草稿）
	if draft.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此草稿"))
		return
	}

	// 软删除草稿
	if err := h.svcCtx.EmailDraftModel.Delete(draft); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_draft",
		Resource:   "draft",
		ResourceId: draft.Id,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// Send 发送草稿
func (h *DraftHandler) Send(c *gin.Context) {
	// 获取草稿ID
	idStr := c.Param("id")
	draftId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的草稿ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查草稿是否存在
	draft, err := h.svcCtx.EmailDraftModel.GetById(draftId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if draft == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("草稿不存在"))
		return
	}

	// 检查权限（只能发送自己的草稿）
	if draft.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此草稿"))
		return
	}

	// 验证邮件内容的完整性
	if draft.Subject == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件主题不能为空"))
		return
	}
	if draft.ToEmails == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("收件人不能为空"))
		return
	}
	if draft.Content == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件内容不能为空"))
		return
	}

	// 获取邮箱信息
	mailbox, err := h.svcCtx.MailboxModel.GetById(draft.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
		return
	}
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
		return
	}

	// 调用邮件发送服务
	smtpConfig := service.SMTPConfig{
		Host:     h.svcCtx.Config.SMTP.Host,
		Port:     h.svcCtx.Config.SMTP.Port,
		Username: mailbox.Email,
		Password: mailbox.Password,
		UseTLS:   h.svcCtx.Config.SMTP.UseTLS,
	}

	smtpService := service.NewSMTPService(smtpConfig)

	// 构建邮件消息
	emailMessage := service.EmailMessage{
		From:        mailbox.Email,
		To:          splitEmails(draft.ToEmails),
		Cc:          splitEmails(draft.CcEmails),
		Bcc:         splitEmails(draft.BccEmails),
		Subject:     draft.Subject,
		Body:        draft.Content,
		ContentType: draft.ContentType,
	}

	// 发送邮件
	if err := smtpService.SendEmail(emailMessage); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("邮件发送失败: "+err.Error()))
		return
	}

	// 创建邮件记录
	email := &model.Email{
		MailboxId:   draft.MailboxId,
		Subject:     draft.Subject,
		FromEmail:   mailbox.Email,
		ToEmails:    draft.ToEmails,
		CcEmails:    draft.CcEmails,
		BccEmails:   draft.BccEmails,
		Content:     draft.Content,
		ContentType: draft.ContentType,
		Direction:   "sent",
		SentAt:      &[]time.Time{time.Now()}[0],
	}

	if err := h.svcCtx.EmailModel.Create(email); err != nil {
		// 邮件已发送，但记录创建失败，记录日志但不返回错误
		// TODO: 添加日志记录
	}

	// 删除草稿（发送成功后）
	if err := h.svcCtx.EmailDraftModel.Delete(draft); err != nil {
		// 草稿删除失败，记录日志但不影响发送结果
		// TODO: 添加日志记录
	}

	sendResp := types.DraftSendResp{
		Success: true,
		Message: "邮件发送成功",
		EmailId: email.Id,
		SentAt:  time.Now(),
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "send_draft",
		Resource:   "draft",
		ResourceId: draft.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SuccessResult(sendResp))
}

// GetById 获取草稿详情
func (h *DraftHandler) GetById(c *gin.Context) {
	// 获取草稿ID
	idStr := c.Param("id")
	draftId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的草稿ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询草稿详情
	draft, err := h.svcCtx.EmailDraftModel.GetById(draftId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if draft == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("草稿不存在"))
		return
	}

	// 检查权限（只能查看自己的草稿）
	if draft.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此草稿"))
		return
	}

	// 返回草稿详情
	resp := types.DraftResp{
		Id:          draft.Id,
		UserId:      draft.UserId,
		MailboxId:   draft.MailboxId,
		Subject:     draft.Subject,
		FromEmail:   "",              // EmailDraft模型中没有FromEmail字段
		ToEmail:     draft.ToEmails,  // 使用ToEmails字段
		CcEmail:     draft.CcEmails,  // 使用CcEmails字段
		BccEmail:    draft.BccEmails, // 使用BccEmails字段
		Content:     draft.Content,
		ContentType: draft.ContentType,
		Attachments: "", // EmailDraft模型中没有Attachments字段
		CreatedAt:   draft.CreatedAt,
		UpdatedAt:   draft.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// AutoSave 自动保存草稿
func (h *DraftHandler) AutoSave(c *gin.Context) {
	var req types.DraftAutoSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查邮箱是否属于当前用户
	if req.MailboxId > 0 {
		mailbox, err := h.svcCtx.MailboxModel.GetById(req.MailboxId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != currentUserId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
			return
		}
	}

	var draft *model.EmailDraft
	var err error

	// 如果提供了草稿ID，则更新现有草稿
	if req.DraftId != nil && *req.DraftId > 0 {
		draft, err = h.svcCtx.EmailDraftModel.GetById(*req.DraftId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}

		// 检查权限
		if draft != nil && draft.UserId != currentUserId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此草稿"))
			return
		}
	}

	// 如果草稿不存在，创建新草稿
	if draft == nil {
		draft = &model.EmailDraft{
			UserId:      currentUserId,
			MailboxId:   req.MailboxId,
			Subject:     req.Subject,
			ToEmails:    req.ToEmail,
			CcEmails:    req.CcEmail,
			BccEmails:   req.BccEmail,
			Content:     req.Content,
			ContentType: req.ContentType,
			Status:      "draft",
		}

		if err := h.svcCtx.EmailDraftModel.Create(draft); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
			return
		}
	} else {
		// 更新现有草稿
		draft.MailboxId = req.MailboxId
		draft.Subject = req.Subject
		draft.ToEmails = req.ToEmail
		draft.CcEmails = req.CcEmail
		draft.BccEmails = req.BccEmail
		draft.Content = req.Content
		draft.ContentType = req.ContentType

		if err := h.svcCtx.EmailDraftModel.Update(draft); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
	}

	// 返回自动保存结果
	resp := types.DraftAutoSaveResp{
		DraftId: draft.Id,
		Success: true,
		Message: "草稿自动保存成功",
		SavedAt: time.Now(),
		IsNew:   req.DraftId == nil || *req.DraftId == 0,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
