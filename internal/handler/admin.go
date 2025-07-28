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

// AdminHandler 管理员处理器
type AdminHandler struct {
	svcCtx *svc.ServiceContext
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(svcCtx *svc.ServiceContext) *AdminHandler {
	return &AdminHandler{
		svcCtx: svcCtx,
	}
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req types.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 查找管理员（支持用户名或邮箱登录）
	var admin *model.Admin
	var err error

	// 先尝试用户名查找
	admin, err = h.svcCtx.AdminModel.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 如果用户名未找到，尝试邮箱查找
	if admin == nil {
		admin, err = h.svcCtx.AdminModel.GetByEmail(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
	}

	if admin == nil {
		c.JSON(http.StatusUnauthorized, result.ErrorUserNotFound)
		return
	}

	// 检查管理员状态
	if admin.Status != 1 {
		c.JSON(http.StatusUnauthorized, result.ErrorUserDisabled)
		return
	}

	// 验证密码
	if valid, err := auth.VerifyPassword(req.Password, admin.Password); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码验证失败"))
		return
	} else if !valid {
		c.JSON(http.StatusUnauthorized, result.ErrorPasswordWrong)
		return
	}

	// 生成JWT token
	token, err := auth.GenerateAdminToken(admin.Id, admin.Username, admin.Role, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成token失败"))
		return
	}

	// 生成刷新token
	refreshToken, err := auth.GenerateAdminToken(admin.Id, admin.Username, admin.Role, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.RefreshExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成刷新token失败"))
		return
	}

	// 更新最后登录时间
	if err := h.svcCtx.AdminModel.UpdateLastLogin(admin.Id); err != nil {
		// 记录日志但不影响登录
	}

	// 返回登录信息
	resp := types.AdminLoginResp{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(h.svcCtx.Config.JWT.ExpireHours) * time.Hour),
		Admin: types.AdminResp{
			Id:          admin.Id,
			Username:    admin.Username,
			Email:       admin.Email,
			Nickname:    admin.Nickname,
			Avatar:      admin.Avatar,
			Role:        admin.Role,
			Status:      admin.Status,
			LastLoginAt: admin.LastLoginAt,
			CreatedAt:   admin.CreatedAt,
			UpdatedAt:   admin.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// RefreshToken 刷新token
func (h *AdminHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 刷新token
	newToken, err := auth.RefreshAdminToken(req.RefreshToken, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"token":     newToken,
		"expiresAt": time.Now().Add(time.Duration(h.svcCtx.Config.JWT.ExpireHours) * time.Hour),
	}))
}

// GetProfile 获取管理员资料
func (h *AdminHandler) GetProfile(c *gin.Context) {
	adminId := middleware.GetCurrentAdminId(c)
	if adminId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	admin, err := h.svcCtx.AdminModel.GetById(adminId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if admin == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	resp := types.AdminResp{
		Id:          admin.Id,
		Username:    admin.Username,
		Email:       admin.Email,
		Nickname:    admin.Nickname,
		Avatar:      admin.Avatar,
		Role:        admin.Role,
		Status:      admin.Status,
		LastLoginAt: admin.LastLoginAt,
		CreatedAt:   admin.CreatedAt,
		UpdatedAt:   admin.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateProfile 更新管理员资料
func (h *AdminHandler) UpdateProfile(c *gin.Context) {
	adminId := middleware.GetCurrentAdminId(c)
	if adminId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	var req types.AdminProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查邮箱是否已被其他管理员使用
	if req.Email != "" {
		if exists, err := h.svcCtx.AdminModel.CheckEmailExists(req.Email, adminId); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已被其他管理员使用"))
			return
		}
	}

	// 更新管理员信息
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
		if err := h.svcCtx.AdminModel.MapUpdate(nil, adminId, updateData); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// ChangePassword 修改密码
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	adminId := middleware.GetCurrentAdminId(c)
	if adminId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	var req types.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取管理员信息
	admin, err := h.svcCtx.AdminModel.GetById(adminId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if admin == nil {
		c.JSON(http.StatusNotFound, result.ErrorUserNotFound)
		return
	}

	// 验证旧密码
	if valid, err := auth.VerifyPassword(req.OldPassword, admin.Password); err != nil {
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
	if err := h.svcCtx.AdminModel.MapUpdate(nil, adminId, map[string]interface{}{
		"password": hashedPassword,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("密码修改成功"))
}

// Dashboard 管理员仪表板
func (h *AdminHandler) Dashboard(c *gin.Context) {
	// TODO: 实现仪表板数据统计
	c.JSON(http.StatusOK, result.SimpleResult("管理员仪表板接口"))
}

// ListUsers 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// TODO: 实现用户列表查询
	c.JSON(http.StatusOK, result.SimpleResult("用户列表接口"))
}

// CreateUser 创建用户
func (h *AdminHandler) CreateUser(c *gin.Context) {
	// TODO: 实现创建用户
	c.JSON(http.StatusOK, result.SimpleResult("创建用户接口"))
}

// UpdateUser 更新用户
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	// TODO: 实现更新用户
	c.JSON(http.StatusOK, result.SimpleResult("更新用户接口"))
}

// DeleteUser 删除用户
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// TODO: 实现删除用户
	c.JSON(http.StatusOK, result.SimpleResult("删除用户接口"))
}

// BatchOperationUsers 批量操作用户
func (h *AdminHandler) BatchOperationUsers(c *gin.Context) {
	// TODO: 实现批量操作用户
	c.JSON(http.StatusOK, result.SimpleResult("批量操作用户接口"))
}

// ImportUsers 导入用户
func (h *AdminHandler) ImportUsers(c *gin.Context) {
	// TODO: 实现导入用户
	c.JSON(http.StatusOK, result.SimpleResult("导入用户接口"))
}

// ExportUsers 导出用户
func (h *AdminHandler) ExportUsers(c *gin.Context) {
	// TODO: 实现导出用户
	c.JSON(http.StatusOK, result.SimpleResult("导出用户接口"))
}

// ListAdmins 管理员列表
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	var req types.AdminListReq
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
	params := model.AdminListParams{
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
		Role:     req.Role,
		Status:   req.Status,
	}

	// 查询管理员列表
	admins, total, err := h.svcCtx.AdminModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var adminList []types.AdminResp
	for _, admin := range admins {
		adminList = append(adminList, types.AdminResp{
			Id:          admin.Id,
			Username:    admin.Username,
			Email:       admin.Email,
			Nickname:    admin.Nickname,
			Avatar:      admin.Avatar,
			Role:        admin.Role,
			Status:      admin.Status,
			LastLoginAt: admin.LastLoginAt,
			CreatedAt:   admin.CreatedAt,
			UpdatedAt:   admin.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     adminList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// CreateAdmin 创建管理员
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req types.AdminCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查用户名是否已存在
	if exists, err := h.svcCtx.AdminModel.CheckUsernameExists(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已存在"))
		return
	}

	// 检查邮箱是否已存在
	if exists, err := h.svcCtx.AdminModel.CheckEmailExists(req.Email); err != nil {
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

	// 创建管理员
	admin := &model.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Role:     req.Role,
		Status:   req.Status,
	}

	if err := h.svcCtx.AdminModel.Create(admin); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "create_admin",
			Resource:   "admin",
			ResourceId: admin.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	resp := types.AdminResp{
		Id:        admin.Id,
		Username:  admin.Username,
		Email:     admin.Email,
		Nickname:  admin.Nickname,
		Avatar:    admin.Avatar,
		Role:      admin.Role,
		Status:    admin.Status,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateAdmin 更新管理员
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	var req types.AdminUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查管理员是否存在
	admin, err := h.svcCtx.AdminModel.GetById(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if admin == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("管理员不存在"))
		return
	}

	// 检查用户名是否已被其他管理员使用
	if req.Username != "" && req.Username != admin.Username {
		if exists, err := h.svcCtx.AdminModel.CheckUsernameExists(req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已被其他管理员使用"))
			return
		}
	}

	// 检查邮箱是否已被其他管理员使用
	if req.Email != "" && req.Email != admin.Email {
		if exists, err := h.svcCtx.AdminModel.CheckEmailExists(req.Email, req.Id); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已被其他管理员使用"))
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
	if req.Role != "" {
		updateData["role"] = req.Role
	}
	updateData["status"] = req.Status

	// 更新管理员信息
	if err := h.svcCtx.AdminModel.MapUpdate(nil, req.Id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "update_admin",
			Resource:   "admin",
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

// DeleteAdmin 删除管理员
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的管理员ID"))
		return
	}

	adminId := uint(id)

	// 检查管理员是否存在
	admin, err := h.svcCtx.AdminModel.GetById(adminId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if admin == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("管理员不存在"))
		return
	}

	// 检查是否为当前登录管理员（不能删除自己）
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId == adminId {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不能删除自己"))
		return
	}

	// 检查是否为超级管理员（可选：保护超级管理员不被删除）
	if admin.Role == "admin" {
		// 检查是否还有其他超级管理员
		superAdmins, err := h.svcCtx.AdminModel.GetSuperAdmins()
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if len(superAdmins) <= 1 {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不能删除最后一个超级管理员"))
			return
		}
	}

	// 软删除管理员
	if err := h.svcCtx.AdminModel.Delete(admin); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "delete_admin",
			Resource:   "admin",
			ResourceId: adminId,
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

// BatchOperationAdmins 批量操作管理员
func (h *AdminHandler) BatchOperationAdmins(c *gin.Context) {
	// TODO: 实现批量操作管理员
	c.JSON(http.StatusOK, result.SimpleResult("批量操作管理员接口"))
}

// GetSystemSettings 获取系统设置
func (h *AdminHandler) GetSystemSettings(c *gin.Context) {
	// TODO: 实现获取系统设置
	c.JSON(http.StatusOK, result.SimpleResult("获取系统设置接口"))
}

// UpdateSystemSettings 更新系统设置
func (h *AdminHandler) UpdateSystemSettings(c *gin.Context) {
	// TODO: 实现更新系统设置
	c.JSON(http.StatusOK, result.SimpleResult("更新系统设置接口"))
}
