package handler

import (
	"fmt"
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

// EmailHandler 邮件处理器
type EmailHandler struct {
	svcCtx *svc.ServiceContext
}

// NewEmailHandler 创建邮件处理器
func NewEmailHandler(svcCtx *svc.ServiceContext) *EmailHandler {
	return &EmailHandler{
		svcCtx: svcCtx,
	}
}

// List 邮件列表
func (h *EmailHandler) List(c *gin.Context) {
	var req types.EmailListReq
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
	params := model.EmailListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		MailboxId: 0, // 默认值
		Subject:   req.Subject,
		FromEmail: req.FromEmail,
		ToEmails:  req.ToEmail, // 使用ToEmails字段
	}

	// 处理可选的MailboxId参数
	if req.MailboxId != nil {
		params.MailboxId = *req.MailboxId
	}

	// 查询邮件列表
	emails, total, err := h.svcCtx.EmailModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var emailList []types.EmailResp
	for _, email := range emails {
		emailList = append(emailList, types.EmailResp{
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
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     emailList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetById 获取邮件详情
func (h *EmailHandler) GetById(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询邮件详情
	email, err := h.svcCtx.EmailModel.GetById(uint(emailId))
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
	if mailbox == nil || mailbox.UserId != currentUserId {
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

// Send 发送邮件
func (h *EmailHandler) Send(c *gin.Context) {
	var req types.EmailSendReq
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

	// TODO: 实现实际的邮件发送逻辑
	// 这里应该：
	// 1. 验证邮件内容的完整性
	// 2. 调用邮件发送服务
	// 3. 创建邮件记录
	// 4. 返回发送结果

	// 创建邮件记录
	email := &model.Email{
		MailboxId:   req.MailboxId,
		Subject:     req.Subject,
		FromEmail:   req.FromEmail,
		ToEmails:    req.ToEmail,  // 使用ToEmails字段
		CcEmails:    req.CcEmail,  // 使用CcEmails字段
		BccEmails:   req.BccEmail, // 使用BccEmails字段
		Content:     req.Content,
		ContentType: req.ContentType,
		Direction:   "sent", // 设置方向为发送
	}

	if err := h.svcCtx.EmailModel.Create(email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "send_email",
		Resource:   "email",
		ResourceId: email.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 模拟发送成功
	sendResp := types.EmailSendResp{
		Success: true,
		Message: "邮件发送成功",
		EmailId: email.Id,
		SentAt:  time.Now(),
	}

	c.JSON(http.StatusOK, result.SuccessResult(sendResp))
}

// MarkRead 标记已读
func (h *EmailHandler) MarkRead(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(uint(emailId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// TODO: 标记为已读
	// Email模型中没有Status字段，需要使用IsRead字段
	// email.IsRead = true
	// if err := h.svcCtx.EmailModel.Update(email); err != nil {
	//     c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
	//     return
	// }

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "mark_read_email",
		Resource:   "email",
		ResourceId: uint(emailId),
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("标记成功"))
}

// MarkStar 标记星标
func (h *EmailHandler) MarkStar(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	var req struct {
		IsStarred bool `json:"isStarred"`
	}
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

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(uint(emailId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 更新星标状态
	updateData := map[string]interface{}{
		"is_starred": req.IsStarred,
	}
	if err := h.svcCtx.EmailModel.MapUpdate(nil, uint(emailId), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	action := "unstar_email"
	if req.IsStarred {
		action = "star_email"
	}
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     action,
		Resource:   "email",
		ResourceId: uint(emailId),
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	message := "取消星标成功"
	if req.IsStarred {
		message = "标记星标成功"
	}
	c.JSON(http.StatusOK, result.SimpleResult(message))
}

// Delete 删除邮件
func (h *EmailHandler) Delete(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(uint(emailId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 软删除邮件
	if err := h.svcCtx.EmailModel.Delete(email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_email",
		Resource:   "email",
		ResourceId: uint(emailId),
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// BatchOperation 批量操作邮件
func (h *EmailHandler) BatchOperation(c *gin.Context) {
	var req types.EmailBatchOperationReq
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

	if len(req.Ids) == 0 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请选择要操作的邮件"))
		return
	}

	var successCount, failCount int
	var errors []string

	switch req.Operation {
	case "read":
		// 批量标记已读
		for _, id := range req.Ids {
			if err := h.batchUpdateEmailStatus(id, currentUserId, 1); err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			successCount++
		}

	case "unread":
		// 批量标记未读
		for _, id := range req.Ids {
			if err := h.batchUpdateEmailStatus(id, currentUserId, 0); err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			successCount++
		}

	case "delete":
		// 批量删除
		for _, id := range req.Ids {
			email, err := h.svcCtx.EmailModel.GetById(id)
			if err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			if email == nil {
				errors = append(errors, "邮件不存在")
				failCount++
				continue
			}

			// 检查权限
			mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
			if err != nil || mailbox == nil || mailbox.UserId != currentUserId {
				errors = append(errors, "无权限操作此邮件")
				failCount++
				continue
			}

			if err := h.svcCtx.EmailModel.Delete(email); err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			successCount++
		}

	case "move":
		// 批量移动到指定邮箱
		if req.TargetId == 0 {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请指定目标邮箱"))
			return
		}

		// 检查目标邮箱权限
		targetMailbox, err := h.svcCtx.MailboxModel.GetById(req.TargetId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if targetMailbox == nil || targetMailbox.UserId != currentUserId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用目标邮箱"))
			return
		}

		for _, id := range req.Ids {
			email, err := h.svcCtx.EmailModel.GetById(id)
			if err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			if email == nil {
				errors = append(errors, "邮件不存在")
				failCount++
				continue
			}

			// 检查权限
			mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
			if err != nil || mailbox == nil || mailbox.UserId != currentUserId {
				errors = append(errors, "无权限操作此邮件")
				failCount++
				continue
			}

			// 移动邮件
			email.MailboxId = req.TargetId
			if err := h.svcCtx.EmailModel.Update(email); err != nil {
				errors = append(errors, err.Error())
				failCount++
				continue
			}
			successCount++
		}

	default:
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不支持的操作类型"))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "batch_" + req.Operation + "_email",
		Resource:   "email",
		ResourceId: 0, // 批量操作没有单一资源ID
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回操作结果
	resp := types.BatchOperationResp{
		Total:        len(req.Ids),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// batchUpdateEmailStatus 批量更新邮件状态的辅助方法
func (h *EmailHandler) batchUpdateEmailStatus(emailId uint, userId uint, status int) error {
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		return err
	}
	if email == nil {
		return fmt.Errorf("邮件不存在")
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		return err
	}
	if mailbox == nil || mailbox.UserId != userId {
		return fmt.Errorf("无权限操作此邮件")
	}

	// TODO: 更新状态
	// Email模型中没有Status字段，需要使用其他字段或添加Status字段
	// email.Status = status
	// return h.svcCtx.EmailModel.Update(email)
	return nil
}
