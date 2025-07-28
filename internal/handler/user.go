package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"new-email/pkg/auth"
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
