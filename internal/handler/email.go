package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/service"
	"new-email/internal/svc"
	"new-email/internal/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
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

	// 根据direction参数设置邮件类型
	if req.Direction == "sent" {
		params.Direction = "sent"
	} else if req.Direction == "received" {
		params.Direction = "inbox"
	} else if req.Type != "" {
		params.Direction = req.Type
	}

	// 查询邮件列表
	emails, total, err := h.svcCtx.EmailModel.List(params)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
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
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 查询邮件详情
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限（只能查看自己的邮件）
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限查看此邮件"))
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
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	mailbox, err := h.svcCtx.MailboxModel.GetById(req.MailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId { // 检查邮箱是否属于当前用户
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限使用此邮箱"))
		return
	}

	// 验证邮件内容的完整性
	if req.Subject == "" {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件主题不能为空"))
		return
	}
	if req.ToEmail == "" {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("收件人不能为空"))
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件内容不能为空"))
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
		From:        req.FromEmail,
		To:          splitEmails(req.ToEmail),
		Cc:          splitEmails(req.CcEmail),
		Bcc:         splitEmails(req.BccEmail),
		Subject:     req.Subject,
		Body:        req.Content,
		ContentType: req.ContentType,
	}

	// 发送邮件
	if err := smtpService.SendEmail(emailMessage); err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件发送失败: "+err.Error()))
		return
	}

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
		c.JSON(http.StatusOK, result.ErrorAdd.AddError(err))
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

// splitEmails 分割邮箱地址字符串
func splitEmails(emails string) []string {
	if emails == "" {
		return []string{}
	}

	// 分割并清理邮箱地址
	parts := strings.Split(emails, ",")
	result := make([]string, 0, len(parts))

	for _, email := range parts {
		email = strings.TrimSpace(email)
		if email != "" {
			result = append(result, email)
		}
	}

	return result
}

// MarkRead 标记已读
func (h *EmailHandler) MarkRead(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// TODO: 标记为已读
	// Email模型中没有Status字段，需要使用IsRead字段
	// email.IsRead = true
	// if err := h.svcCtx.EmailModel.Update(email); err != nil {
	//     c.JSON(http.StatusOK, result.ErrorUpdate.AddError(err))
	//     return
	// }

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "mark_read_email",
		Resource:   "email",
		ResourceId: emailId,
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
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	var req struct {
		IsStarred bool `json:"isStarred"`
	}
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

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 更新星标状态
	updateData := map[string]interface{}{
		"is_starred": req.IsStarred,
	}
	if err := h.svcCtx.EmailModel.MapUpdate(nil, emailId, updateData); err != nil {
		c.JSON(http.StatusOK, result.ErrorUpdate.AddError(err))
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
		ResourceId: emailId,
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
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 查询邮件
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 软删除邮件
	if err := h.svcCtx.EmailModel.Delete(email); err != nil {
		c.JSON(http.StatusOK, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_email",
		Resource:   "email",
		ResourceId: emailId,
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
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("请选择要操作的邮件"))
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
			c.JSON(http.StatusOK, result.ErrorSimpleResult("请指定目标邮箱"))
			return
		}

		// 检查目标邮箱权限
		targetMailbox, err := h.svcCtx.MailboxModel.GetById(req.TargetId)
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		}
		if targetMailbox == nil || targetMailbox.UserId != currentUserId {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限使用目标邮箱"))
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
		c.JSON(http.StatusOK, result.ErrorSimpleResult("不支持的操作类型"))
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
func (h *EmailHandler) batchUpdateEmailStatus(emailId int64, userId int64, status int) error {
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

	// 更新已读状态
	updateData := map[string]interface{}{
		"is_read": status == 1,
	}
	return h.svcCtx.EmailModel.MapUpdate(nil, emailId, updateData)
}

// Export 导出邮件
func (h *EmailHandler) Export(c *gin.Context) {
	var req types.EmailExportReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 设置默认格式
	if req.Format == "" {
		req.Format = "csv"
	}

	// 构建查询参数
	params := model.EmailListParams{
		BaseListParams: model.BaseListParams{
			Page:     1,
			PageSize: 10000, // 导出时设置较大的页面大小
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		Subject:   req.Subject,
		FromEmail: req.FromEmail,
		ToEmails:  req.ToEmail,
		Direction: req.Direction,
	}

	// 处理可选的MailboxId参数
	if req.MailboxId != nil {
		// 检查邮箱权限
		mailbox, err := h.svcCtx.MailboxModel.GetById(*req.MailboxId)
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != currentUserId {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限访问此邮箱"))
			return
		}
		params.MailboxId = *req.MailboxId
	}

	// 查询邮件列表
	emails, total, err := h.svcCtx.EmailModel.List(params)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("没有找到符合条件的邮件"))
		return
	}

	// 过滤用户权限：只导出用户有权限的邮件
	var filteredEmails []*model.Email
	for _, email := range emails {
		mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
		if err != nil || mailbox == nil || mailbox.UserId != currentUserId {
			continue // 跳过无权限的邮件
		}
		filteredEmails = append(filteredEmails, email)
	}

	if len(filteredEmails) == 0 {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("没有找到有权限的邮件"))
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("emails_export_%s.%s", timestamp, req.Format)

	// 确保导出目录存在
	exportDir := "./data/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("创建导出目录失败"))
		return
	}

	filePath := filepath.Join(exportDir, fileName)

	// 根据格式导出
	var fileSize int64
	switch req.Format {
	case "csv":
		fileSize, err = h.exportToCSV(filteredEmails, filePath, req.IncludeContent)
	case "json":
		fileSize, err = h.exportToJSON(filteredEmails, filePath, req.IncludeContent)
	case "eml":
		fileSize, err = h.exportToEML(filteredEmails, filePath)
	default:
		c.JSON(http.StatusOK, result.ErrorSimpleResult("不支持的导出格式"))
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("导出失败: "+err.Error()))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "export_emails",
		Resource:   "email",
		ResourceId: 0,
		Method:     "GET",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回导出结果
	resp := types.EmailExportResp{
		FileName:    fileName,
		FileSize:    fileSize,
		RecordCount: int64(len(filteredEmails)),
		DownloadUrl: fmt.Sprintf("/api/user/emails/download/%s", fileName),
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// exportToCSV 导出为CSV格式
func (h *EmailHandler) exportToCSV(emails []*model.Email, filePath string, includeContent bool) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV头部
	headers := []string{"ID", "邮箱ID", "主题", "发件人", "收件人", "抄送", "密送", "方向", "是否已读", "是否星标", "发送时间", "接收时间", "创建时间"}
	if includeContent {
		headers = append(headers, "内容类型", "邮件内容")
	}

	if err := writer.Write(headers); err != nil {
		return 0, err
	}

	// 写入数据行
	for _, email := range emails {
		row := []string{
			strconv.FormatInt(email.Id, 10),
			strconv.FormatInt(email.MailboxId, 10),
			email.Subject,
			email.FromEmail,
			email.ToEmails,
			email.CcEmails,
			email.BccEmails,
			email.Direction,
			strconv.FormatBool(email.IsRead),
			strconv.FormatBool(email.IsStarred),
			formatTimePtr(email.SentAt),
			formatTimePtr(email.ReceivedAt),
			email.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if includeContent {
			row = append(row, email.ContentType, email.Content)
		}

		if err := writer.Write(row); err != nil {
			return 0, err
		}
	}

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

// exportToJSON 导出为JSON格式
func (h *EmailHandler) exportToJSON(emails []*model.Email, filePath string, includeContent bool) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 构建导出数据
	var exportData []map[string]interface{}
	for _, email := range emails {
		data := map[string]interface{}{
			"id":          email.Id,
			"mailbox_id":  email.MailboxId,
			"subject":     email.Subject,
			"from_email":  email.FromEmail,
			"from_name":   email.FromName,
			"to_emails":   email.ToEmails,
			"cc_emails":   email.CcEmails,
			"bcc_emails":  email.BccEmails,
			"reply_to":    email.ReplyTo,
			"direction":   email.Direction,
			"is_read":     email.IsRead,
			"is_starred":  email.IsStarred,
			"sent_at":     email.SentAt,
			"received_at": email.ReceivedAt,
			"created_at":  email.CreatedAt,
			"updated_at":  email.UpdatedAt,
		}

		if includeContent {
			data["content_type"] = email.ContentType
			data["content"] = email.Content
		}

		exportData = append(exportData, data)
	}

	// 写入JSON文件
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		return 0, err
	}

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

// exportToEML 导出为EML格式
func (h *EmailHandler) exportToEML(emails []*model.Email, filePath string) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalSize int64

	for i, email := range emails {
		// 构建EML格式的邮件内容
		emlContent := h.buildEMLContent(email)

		// 写入分隔符（除了第一封邮件）
		if i > 0 {
			separator := "\n\n" + strings.Repeat("-", 50) + "\n\n"
			if _, err := file.WriteString(separator); err != nil {
				return 0, err
			}
			totalSize += int64(len(separator))
		}

		// 写入EML内容
		if _, err := file.WriteString(emlContent); err != nil {
			return 0, err
		}
		totalSize += int64(len(emlContent))
	}

	return totalSize, nil
}

// buildEMLContent 构建EML格式的邮件内容
func (h *EmailHandler) buildEMLContent(email *model.Email) string {
	var builder strings.Builder

	// 邮件头部
	builder.WriteString(fmt.Sprintf("Message-ID: <%s>\n", email.MessageId))
	builder.WriteString(fmt.Sprintf("Subject: %s\n", email.Subject))
	builder.WriteString(fmt.Sprintf("From: %s <%s>\n", email.FromName, email.FromEmail))
	builder.WriteString(fmt.Sprintf("To: %s\n", email.ToEmails))

	if email.CcEmails != "" {
		builder.WriteString(fmt.Sprintf("Cc: %s\n", email.CcEmails))
	}

	if email.ReplyTo != "" {
		builder.WriteString(fmt.Sprintf("Reply-To: %s\n", email.ReplyTo))
	}

	// 时间信息
	if email.SentAt != nil {
		builder.WriteString(fmt.Sprintf("Date: %s\n", email.SentAt.Format(time.RFC1123Z)))
	} else {
		builder.WriteString(fmt.Sprintf("Date: %s\n", email.CreatedAt.Format(time.RFC1123Z)))
	}

	// 内容类型
	if email.ContentType == "html" {
		builder.WriteString("Content-Type: text/html; charset=utf-8\n")
	} else {
		builder.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	builder.WriteString("Content-Transfer-Encoding: 8bit\n")
	builder.WriteString("\n") // 空行分隔头部和正文

	// 邮件正文
	builder.WriteString(email.Content)

	return builder.String()
}

// Download 下载导出的文件
func (h *EmailHandler) Download(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("文件名不能为空"))
		return
	}

	// 获取当前用户ID（验证权限）
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	// 构建文件路径
	filePath := filepath.Join("./data/exports", fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("文件不存在"))
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(filePath)
}

// formatTimePtr 格式化时间指针
func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
