package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"new-email/pkg/auth"
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
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 如果提供了验证码，则验证邮箱验证码
	if req.Code != "" {
		if err := h.verifyEmailCode(req.Email, req.Code); err != nil {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("验证码验证失败: "+err.Error()))
			return
		}
	}

	// 检查用户名是否已存在
	if exists, err := h.svcCtx.UserModel.CheckUsernameExists(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("用户名已存在"))
		return
	}

	// 检查邮箱是否已存在
	if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱已存在"))
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

// verifyEmailCode 验证邮箱验证码
func (h *UserHandler) verifyEmailCode(email, code string) error {
	// 查找未使用的验证码
	verificationCode, err := h.svcCtx.VerificationCodeModel.FindByCode(code)
	if err != nil {
		return fmt.Errorf("验证码不存在")
	}

	// 检查验证码是否已使用
	if verificationCode.IsUsed {
		return fmt.Errorf("验证码已使用")
	}

	// 检查验证码是否过期（这里需要根据实际的过期时间字段调整）
	// 假设验证码有效期为10分钟
	if time.Since(verificationCode.ExtractedAt) > 10*time.Minute {
		return fmt.Errorf("验证码已过期")
	}

	// 标记验证码为已使用
	verificationCode.IsUsed = true
	if err := h.svcCtx.VerificationCodeModel.Update(verificationCode); err != nil {
		return fmt.Errorf("更新验证码状态失败")
	}

	return nil
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req types.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 先尝试用户名查找
	user, err := h.svcCtx.UserModel.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
		return
	}

	// 如果用户名未找到，尝试邮箱查找
	if user == nil {
		user, err = h.svcCtx.UserModel.GetByEmail(req.Username)
		if err != nil {
			c.JSON(http.StatusOK, result.ErrorSelect.AddError(err))
			return
		}
	}

	if user == nil {
		c.JSON(http.StatusOK, result.ErrorUserNotFound)
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusOK, result.ErrorUserDisabled)
		return
	}

	// 验证密码
	if valid, err := auth.VerifyPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("密码验证失败"))
		return
	} else if !valid {
		c.JSON(http.StatusOK, result.ErrorPasswordWrong)
		return
	}

	// 生成JWT token
	token, err := auth.GenerateTokenFull(user.Id, user.Username, false, "", h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusOK, result.ErrorSimpleResult("生成token失败"))
		return
	}

	// 返回登录信息
	resp := types.UserLoginResp{
		Token: token,
		User: types.UserResp{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
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

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	var req types.UserProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" {
		if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email, userId); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusOK, result.ErrorSimpleResult("邮箱已被其他用户使用"))
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
		c.JSON(http.StatusOK, result.ErrorUnauthorized)
		return
	}

	var req types.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.ErrorBindingParam.AddError(err))
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
		c.JSON(http.StatusOK, result.ErrorPasswordWrong)
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
