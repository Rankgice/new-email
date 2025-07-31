package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"new-email/pkg/auth"
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

	// 检查域名是否存在且属于当前用户
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
	// 创建邮箱
	mailbox := &model.Mailbox{
		UserId:      currentUserId,
		DomainId:    req.DomainId,
		Email:       req.Email,
		Password:    req.Password,
		AutoReceive: req.AutoReceive,
		Status:      req.Status,
	}

	if err := h.svcCtx.MailboxModel.Create(mailbox); err != nil {
		c.JSON(http.StatusOK, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_mailbox",
		Resource:   "mailbox",
		ResourceId: mailbox.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

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
		// 加密新密码
		encryptedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("密码加密失败"))
			return
		}
		updateData["password"] = encryptedPassword
	}
	updateData["auto_receive"] = req.AutoReceive
	updateData["status"] = req.Status

	// 更新邮箱信息
	if err := h.svcCtx.MailboxModel.MapUpdate(nil, req.Id, updateData); err != nil {
		c.JSON(http.StatusOK, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_mailbox",
		Resource:   "mailbox",
		ResourceId: req.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

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

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_mailbox",
		Resource:   "mailbox",
		ResourceId: mailboxId,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// Sync 同步邮箱
func (h *MailboxHandler) Sync(c *gin.Context) {
	var req types.MailboxSyncReq
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

	// 检查权限（只能同步自己的邮箱）
	if mailbox.UserId != currentUserId {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	// TODO: 实现实际的邮箱同步逻辑
	// 这里应该：
	// 1. 连接IMAP服务器
	// 2. 获取邮件列表
	// 3. 同步邮件到数据库
	// 4. 更新最后同步时间

	// 模拟同步结果
	resp := types.MailboxSyncResp{
		Success:    true,
		Message:    "同步成功",
		SyncCount:  10,
		ErrorCount: 0,
		LastSyncAt: time.Now(),
	}

	// 更新最后同步时间
	now := time.Now()
	h.svcCtx.MailboxModel.MapUpdate(nil, req.Id, map[string]interface{}{
		"last_sync_at": &now,
	})

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "sync_mailbox",
		Resource:   "mailbox",
		ResourceId: req.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

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
