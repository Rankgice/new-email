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

// UserHandler 用户处理器
type UserHandler struct {
	svcCtx *svc.ServiceContext
}

// NewUserHandler 创建用户处理器
func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		svcCtx: svcCtx,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req types.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查用户名是否已存在
	if exists, err := h.svcCtx.UserModel.CheckUsernameExists(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已存在"))
		return
	}

	// 检查邮箱是否已存在
	if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已存在"))
		return
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(req.Password)
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
		Status:   1, // 默认启用
	}

	if err := h.svcCtx.UserModel.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
	}))
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req types.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 查找用户（支持用户名或邮箱登录）
	var user *model.User
	var err error

	// 先尝试用户名查找
	user, err = h.svcCtx.UserModel.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 如果用户名未找到，尝试邮箱查找
	if user == nil {
		user, err = h.svcCtx.UserModel.GetByEmail(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, result.ErrorUserNotFound)
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusUnauthorized, result.ErrorUserDisabled)
		return
	}

	// 验证密码
	if valid, err := auth.VerifyPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码验证失败"))
		return
	} else if !valid {
		c.JSON(http.StatusUnauthorized, result.ErrorPasswordWrong)
		return
	}

	// 生成JWT token
	token, err := auth.GenerateUserToken(user.Id, user.Username, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成token失败"))
		return
	}

	// 生成刷新token
	refreshToken, err := auth.GenerateUserToken(user.Id, user.Username, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.RefreshExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成刷新token失败"))
		return
	}

	// 更新最后登录时间
	if err := h.svcCtx.UserModel.UpdateLastLogin(user.Id); err != nil {
		// 记录日志但不影响登录
		// log.Printf("更新用户最后登录时间失败: %v", err)
	}

	// 返回登录信息
	resp := types.UserLoginResp{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(h.svcCtx.Config.JWT.ExpireHours) * time.Hour),
		User: types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// RefreshToken 刷新token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 刷新token
	newToken, err := auth.RefreshToken(req.RefreshToken, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"token":     newToken,
		"expiresAt": time.Now().Add(time.Duration(h.svcCtx.Config.JWT.ExpireHours) * time.Hour),
	}))
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	resp := types.UserResp{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	var req types.UserProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" {
		if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email, userId); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已被其他用户使用"))
			return
		}
	}

	// 更新用户信息
	updateData := map[string]interface{}{}
	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}

	if len(updateData) > 0 {
		if err := h.svcCtx.UserModel.MapUpdate(nil, userId, updateData); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	var req types.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取用户信息
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 验证旧密码
	if valid, err := auth.VerifyPassword(req.OldPassword, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码验证失败"))
		return
	} else if !valid {
		c.JSON(http.StatusBadRequest, result.ErrorPasswordWrong)
		return
	}

	// 加密新密码
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码加密失败"))
		return
	}

	// 更新密码
	if err := h.svcCtx.UserModel.MapUpdate(nil, userId, map[string]interface{}{
		"password": hashedPassword,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("密码修改成功"))
}

// List 获取用户列表（管理员功能）
func (h *UserHandler) List(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
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
		userList = append(userList, types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
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

// GetById 根据ID获取用户信息（管理员功能）
func (h *UserHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	user, err := h.svcCtx.UserModel.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	resp := types.UserResp{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建用户（管理员功能）
func (h *UserHandler) Create(c *gin.Context) {
	var req types.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查用户名是否已存在
	if exists, err := h.svcCtx.UserModel.CheckUsernameExists(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已存在"))
		return
	}

	// 检查邮箱是否已存在
	if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已存在"))
		return
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(req.Password)
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
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if err := h.svcCtx.UserModel.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "create_user",
			Resource:   "user",
			ResourceId: user.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	resp := types.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新用户信息（管理员功能）
func (h *UserHandler) Update(c *gin.Context) {
	var req types.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查用户是否存在
	user, err := h.svcCtx.UserModel.GetById(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != "" && req.Username != user.Username {
		if exists, err := h.svcCtx.UserModel.CheckUsernameExists(req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已被其他用户使用"))
			return
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email, req.Id); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已被其他用户使用"))
			return
		}
	}

	// 构建更新数据
	updateData := map[string]interface{}{}
	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}
	updateData["status"] = req.Status

	// 更新用户信息
	if err := h.svcCtx.UserModel.MapUpdate(nil, req.Id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "update_user",
			Resource:   "user",
			ResourceId: req.Id,
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// Delete 删除用户（管理员功能）
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	userId := uint(id)

	// 检查用户是否存在
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 检查是否为当前登录用户（不能删除自己）
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == userId {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不能删除自己"))
		return
	}

	// 软删除用户
	if err := h.svcCtx.UserModel.Delete(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "delete_user",
			Resource:   "user",
			ResourceId: userId,
			Method:     "DELETE",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// ResetPassword 重置用户密码（管理员功能）
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	userId := uint(id)

	// 检查用户是否存在
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 生成默认密码（可以是随机密码或固定密码）
	defaultPassword := "123456" // 实际项目中应该生成随机密码
	hashedPassword, err := auth.HashPassword(defaultPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码加密失败"))
		return
	}

	// 更新密码
	if err := h.svcCtx.UserModel.MapUpdate(nil, userId, map[string]interface{}{
		"password": hashedPassword,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "reset_password",
			Resource:   "user",
			ResourceId: userId,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"message":         "密码重置成功",
		"defaultPassword": defaultPassword,
	}))
}

// ToggleStatus 切换用户状态（管理员功能）
func (h *UserHandler) ToggleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	userId := uint(id)

	// 检查用户是否存在
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 检查是否为当前登录用户（不能禁用自己）
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == userId {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不能禁用自己"))
		return
	}

	// 切换状态
	newStatus := 1
	if user.Status == 1 {
		newStatus = 0
	}

	if err := h.svcCtx.UserModel.MapUpdate(nil, userId, map[string]interface{}{
		"status": newStatus,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	action := "enable_user"
	if newStatus == 0 {
		action = "disable_user"
	}

	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     action,
			Resource:   "user",
			ResourceId: userId,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	statusText := "启用"
	if newStatus == 0 {
		statusText = "禁用"
	}

	c.JSON(http.StatusOK, result.SimpleResult("用户"+statusText+"成功"))
}

// GetStats 获取用户统计信息（管理员功能）
func (h *UserHandler) GetStats(c *gin.Context) {
	// 获取总用户数
	totalUsers, _, err := h.svcCtx.UserModel.List(model.UserListParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取活跃用户数
	status := 1
	activeUsers, _, err := h.svcCtx.UserModel.List(model.UserListParams{Status: &status})
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取今日新增用户数
	today := time.Now().Truncate(24 * time.Hour)
	newUsers, _, err := h.svcCtx.UserModel.List(model.UserListParams{
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: today,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	resp := types.UserStatsResp{
		TotalUsers:  int64(len(totalUsers)),
		ActiveUsers: int64(len(activeUsers)),
		NewUsers:    int64(len(newUsers)),
		OnlineUsers: 0, // 在线用户数需要额外的逻辑来实现
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Search 搜索用户
func (h *UserHandler) Search(c *gin.Context) {
	var req types.UserSearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 根据搜索类型构建查询参数
	params := model.UserListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}

	switch req.Type {
	case "username":
		params.Username = req.Keyword
	case "email":
		params.Email = req.Keyword
	case "nickname":
		// 这里需要在model层添加nickname搜索支持
		// 暂时使用username搜索
		params.Username = req.Keyword
	default:
		// 默认搜索用户名和邮箱
		params.Username = req.Keyword
		if params.Username == "" {
			params.Email = req.Keyword
		}
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
		userList = append(userList, types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
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

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	// 获取当前用户ID
	userId := middleware.GetCurrentUserId(c)

	// 记录操作日志
	if userId > 0 {
		log := &model.OperationLog{
			UserId:    userId,
			Action:    "logout",
			Resource:  "user",
			Method:    "POST",
			Path:      c.Request.URL.Path,
			Ip:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Status:    http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 在实际项目中，这里可以将token加入黑名单
	// 或者清除服务端的session信息

	c.JSON(http.StatusOK, result.SimpleResult("登出成功"))
}
