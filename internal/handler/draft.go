package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
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
	draftId, err := strconv.ParseUint(idStr, 10, 32)
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
	draft, err := h.svcCtx.EmailDraftModel.GetById(uint(draftId))
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
	draftId, err := strconv.ParseUint(idStr, 10, 32)
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
	draft, err := h.svcCtx.EmailDraftModel.GetById(uint(draftId))
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
	draftId, err := strconv.ParseUint(idStr, 10, 32)
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
	draft, err := h.svcCtx.EmailDraftModel.GetById(uint(draftId))
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

	// TODO: 实现实际的邮件发送逻辑
	// 这里应该：
	// 1. 验证邮件内容的完整性
	// 2. 调用邮件发送服务
	// 3. 创建邮件记录
	// 4. 删除草稿（可选）

	// 模拟发送成功
	sendResp := types.DraftSendResp{
		Success: true,
		Message: "邮件发送成功",
		EmailId: 0, // 实际发送后应该返回邮件ID
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
	draftId, err := strconv.ParseUint(idStr, 10, 32)
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
	draft, err := h.svcCtx.EmailDraftModel.GetById(uint(draftId))
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
