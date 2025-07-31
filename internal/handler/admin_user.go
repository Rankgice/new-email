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

// AdminUserHandler 管理员用户管理处理器
type AdminUserHandler struct {
	svcCtx *svc.ServiceContext
}

// NewAdminUserHandler 创建管理员用户管理处理器
func NewAdminUserHandler(svcCtx *svc.ServiceContext) *AdminUserHandler {
	return &AdminUserHandler{
		svcCtx: svcCtx,
	}
}

// List 获取用户列表
func (h *AdminUserHandler) List(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
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
	params := model.UserListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		Role:     req.Role,
	}

	// 查询用户列表
	users, total, err := h.svcCtx.UserModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var userList []types.UserResp
	for _, user := range users {
		// 获取用户的邮箱数量
		mailboxCount, _ := h.svcCtx.MailboxModel.CountByUserId(user.Id)

		// 获取用户的邮件数量
		emailCount, _ := h.svcCtx.EmailModel.CountByUserId(user.Id)

		userResp := types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Status:      user.Status,
			Role:        user.Role,
			LastLoginAt: user.LastLoginAt,
			LastLoginIp: user.LastLoginIp,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

		// 添加统计信息（如果需要）
		userResp.MailboxCount = mailboxCount
		userResp.EmailCount = emailCount

		userList = append(userList, userResp)
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     userList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetById 获取用户详情
func (h *AdminUserHandler) GetById(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
		return
	}

	// 查询用户详情
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("用户不存在"))
		return
	}

	// 获取用户的邮箱列表
	mailboxes, _, err := h.svcCtx.MailboxModel.List(model.MailboxListParams{
		BaseListParams: model.BaseListParams{Page: 1, PageSize: 100},
		UserId:         userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取用户的最近邮件
	recentEmails, _, err := h.svcCtx.EmailModel.List(model.EmailListParams{
		BaseListParams: model.BaseListParams{Page: 1, PageSize: 10},
		UserId:         &userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换邮箱信息
	var mailboxList []types.MailboxResp
	for _, mailbox := range mailboxes {
		mailboxList = append(mailboxList, types.MailboxResp{
			Id:        mailbox.Id,
			UserId:    mailbox.UserId,
			Email:     mailbox.Email,
			Password:  "", // 不返回密码
			ImapHost:  mailbox.ImapHost,
			ImapPort:  mailbox.ImapPort,
			SmtpHost:  mailbox.SmtpHost,
			SmtpPort:  mailbox.SmtpPort,
			Status:    mailbox.Status,
			CreatedAt: mailbox.CreatedAt,
			UpdatedAt: mailbox.UpdatedAt,
		})
	}

	// 转换邮件信息
	var emailList []types.EmailResp
	for _, email := range recentEmails {
		emailList = append(emailList, types.EmailResp{
			Id:        email.Id,
			MailboxId: email.MailboxId,
			Subject:   email.Subject,
			FromEmail: email.FromEmail,
			ToEmails:  email.ToEmails,
			Direction: email.Direction,
			IsRead:    email.IsRead,
			CreatedAt: email.CreatedAt,
		})
	}

	// 构建用户详情响应
	userDetail := map[string]interface{}{
		"user": types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Status:      user.Status,
			Role:        user.Role,
			LastLoginAt: user.LastLoginAt,
			LastLoginIp: user.LastLoginIp,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		"mailboxes":    mailboxList,
		"recentEmails": emailList,
		"statistics": map[string]interface{}{
			"mailboxCount": len(mailboxList),
			"emailCount":   len(emailList),
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(userDetail))
}

// Create 创建用户
func (h *AdminUserHandler) Create(c *gin.Context) {
	var req types.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := h.svcCtx.UserModel.GetByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已存在"))
		return
	}

	// 检查邮箱是否已存在
	existingUser, _ = h.svcCtx.UserModel.GetByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已存在"))
		return
	}

	// 加密密码
	hashedPassword, err := h.svcCtx.UserModel.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码加密失败"))
		return
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     req.Role,
		Status:   req.Status,
	}

	if err := h.svcCtx.UserModel.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorCreate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "admin_create_user",
		Resource:   "user",
		ResourceId: user.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的用户信息
	resp := types.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新用户
func (h *AdminUserHandler) Update(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	var req types.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
		return
	}

	// 查询用户
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("用户不存在"))
		return
	}

	// 更新用户信息
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		existingUser, _ := h.svcCtx.UserModel.GetByEmail(req.Email)
		if existingUser != nil && existingUser.Id != userId {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已被其他用户使用"))
			return
		}
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != 0 {
		user.Status = req.Status
	}

	// 保存更新
	if err := h.svcCtx.UserModel.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "admin_update_user",
		Resource:   "user",
		ResourceId: user.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的用户信息
	resp := types.UserResp{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Status:      user.Status,
		Role:        user.Role,
		LastLoginAt: user.LastLoginAt,
		LastLoginIp: user.LastLoginIp,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除用户
func (h *AdminUserHandler) Delete(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
		return
	}

	// 不能删除自己
	if userId == currentUserId {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不能删除自己"))
		return
	}

	// 查询用户
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("用户不存在"))
		return
	}

	// 软删除用户
	if err := h.svcCtx.UserModel.Delete(userId); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "admin_delete_user",
		Resource:   "user",
		ResourceId: userId,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// ResetPassword 重置用户密码
func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	var req types.UserPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前管理员ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查管理员权限
	if !h.checkAdminPermission(c, currentUserId) {
		return
	}

	// 查询用户
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("用户不存在"))
		return
	}

	// 加密新密码
	hashedPassword, err := h.svcCtx.UserModel.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码加密失败"))
		return
	}

	// 更新密码
	user.Password = hashedPassword
	if err := h.svcCtx.UserModel.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "admin_reset_password",
		Resource:   "user",
		ResourceId: userId,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("密码重置成功"))
}

// checkAdminPermission 检查管理员权限
func (h *AdminUserHandler) checkAdminPermission(c *gin.Context, userId int64) bool {
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return false
	}

	if user.Role != "admin" && user.Role != "manager" {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无管理员权限"))
		return false
	}

	return true
}
