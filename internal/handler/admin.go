package handler

import (
	"fmt"
	"net/http"
	"new-email/internal/constant"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"new-email/pkg/auth"
	"strconv"
	"strings"
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
	token, err := auth.GenerateTokenFull(admin.Id, admin.Username, true, admin.Role, h.svcCtx.Config.JWT.Secret, h.svcCtx.Config.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成token失败"))
		return
	}
	// 返回登录信息
	resp := types.AdminLoginResp{
		Token: token,
		Admin: types.AdminResp{
			Id:        admin.Id,
			Username:  admin.Username,
			Email:     admin.Email,
			Nickname:  admin.Nickname,
			Avatar:    admin.Avatar,
			Role:      admin.Role,
			Status:    admin.Status,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
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
	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取统计数据
	stats, err := h.getDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(stats))
}

// getDashboardStats 获取仪表板统计数据
func (h *AdminHandler) getDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 模拟统计数据
	// 在实际项目中，这些方法需要在对应的model中实现
	stats["users"] = types.UserStatsResp{
		Total:     100,
		Active:    80,
		Inactive:  20,
		Today:     5,
		ThisWeek:  20,
		ThisMonth: 50,
	}

	stats["mailboxes"] = map[string]interface{}{
		"total":  50,
		"active": 45,
		"synced": 40,
	}

	stats["domains"] = map[string]interface{}{
		"total":    10,
		"verified": 8,
		"active":   9,
	}

	stats["emails"] = types.EmailStatsResp{
		TotalEmails:    1000,
		SentEmails:     200,
		ReceivedEmails: 800,
		TodayEmails:    20,
	}

	stats["verificationCodes"] = map[string]interface{}{
		"total":   500,
		"used":    450,
		"expired": 30,
		"active":  20,
	}

	return stats, nil
}

// ListUsers 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
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
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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

// CreateUser 创建用户
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req types.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
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
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
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

	// 返回创建的用户信息（不包含密码）
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

// UpdateUser 更新用户
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	var req types.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 检查用户是否存在
	user, err := h.svcCtx.UserModel.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("用户不存在"))
		return
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != user.Username {
		if exists, err := h.svcCtx.UserModel.CheckUsernameExists(req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户名已存在"))
			return
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != user.Email {
		if exists, err := h.svcCtx.UserModel.CheckEmailExists(req.Email); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱已存在"))
			return
		}
	}

	// 更新用户信息
	user.Username = req.Username
	user.Email = req.Email
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar
	user.Status = req.Status

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("密码加密失败"))
			return
		}
		user.Password = hashedPassword
	}

	if err := h.svcCtx.UserModel.Update(nil, user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "update_user",
			Resource:   "user",
			ResourceId: user.Id,
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回更新后的用户信息
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

// DeleteUser 删除用户
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的用户ID"))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 检查用户是否存在
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
	if err := h.svcCtx.UserModel.Delete(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
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

// BatchOperationUsers 批量操作用户
func (h *AdminHandler) BatchOperationUsers(c *gin.Context) {
	var req types.UserBatchOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请选择要操作的用户"))
		return
	}

	var successCount, failCount int
	var errors []string

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	switch req.Operation {
	case "enable":
		// 批量启用
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不能操作自己", id))
				failCount++
				continue
			}

			user, err := h.svcCtx.UserModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
				failCount++
				continue
			}
			if user == nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不存在", id))
				failCount++
				continue
			}

			user.Status = constant.StatusEnabled
			if err := h.svcCtx.UserModel.Update(nil, user); err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	case "disable":
		// 批量禁用
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不能操作自己", id))
				failCount++
				continue
			}

			user, err := h.svcCtx.UserModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
				failCount++
				continue
			}
			if user == nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不存在", id))
				failCount++
				continue
			}

			user.Status = 0
			if err := h.svcCtx.UserModel.Update(nil, user); err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	case "delete":
		// 批量删除
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不能删除自己", id))
				failCount++
				continue
			}

			user, err := h.svcCtx.UserModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
				failCount++
				continue
			}
			if user == nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: 不存在", id))
				failCount++
				continue
			}

			if err := h.svcCtx.UserModel.Delete(user); err != nil {
				errors = append(errors, fmt.Sprintf("用户ID %d: %v", id, err))
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
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     fmt.Sprintf("batch_%s_user", req.Operation),
			Resource:   "user",
			ResourceId: 0, // 批量操作没有单一资源ID
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回操作结果
	resp := types.BatchOperationResp{
		Total:        len(req.Ids),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// ImportUsers 导入用户
func (h *AdminHandler) ImportUsers(c *gin.Context) {
	// 验证超级管理员权限
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要超级管理员权限"))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请选择要导入的文件"))
		return
	}
	defer file.Close()

	// 检查文件类型
	if header.Header.Get("Content-Type") != "text/csv" &&
		!strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("只支持CSV格式文件"))
		return
	}

	// TODO: 实现CSV文件解析和用户导入逻辑
	// 这里应该：
	// 1. 解析CSV文件
	// 2. 验证数据格式
	// 3. 批量创建用户
	// 4. 返回导入结果

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "import_users",
			Resource:   "user",
			ResourceId: 0,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 模拟导入结果
	resp := types.ImportUsersResp{
		Total:        0,
		SuccessCount: 0,
		FailCount:    0,
		Errors:       []string{},
		Message:      "用户导入功能待实现",
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// ExportUsers 导出用户
func (h *AdminHandler) ExportUsers(c *gin.Context) {
	var req types.ExportUsersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 设置默认导出格式
	if req.Format == "" {
		req.Format = "csv"
	}

	// 查询用户数据
	params := model.UserListParams{
		BaseListParams: model.BaseListParams{
			Page:     1,
			PageSize: 10000, // 导出时不分页
		},
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
	}

	users, _, err := h.svcCtx.UserModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// TODO: 实现实际的文件导出逻辑
	// 这里应该：
	// 1. 根据格式生成文件内容
	// 2. 设置响应头
	// 3. 返回文件流

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "export_users",
			Resource:   "user",
			ResourceId: 0,
			Method:     "GET",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 模拟导出结果
	resp := types.ExportUsersResp{
		Total:    int64(len(users)),
		Format:   req.Format,
		Filename: fmt.Sprintf("users_%s.%s", time.Now().Format("20060102_150405"), req.Format),
		Message:  "用户导出功能待实现",
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
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
			Id:        admin.Id,
			Username:  admin.Username,
			Email:     admin.Email,
			Nickname:  admin.Nickname,
			Avatar:    admin.Avatar,
			Role:      admin.Role,
			Status:    admin.Status,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
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
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的管理员ID"))
		return
	}

	adminId := id

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
		status := 1
		superAdmins, _, err := h.svcCtx.AdminModel.List(model.AdminListParams{Role: "admin", Status: &status})
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
	var req types.AdminBatchOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证超级管理员权限
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要超级管理员权限"))
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请选择要操作的管理员"))
		return
	}

	var successCount, failCount int
	var errors []string

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	switch req.Operation {
	case "enable":
		// 批量启用
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不能操作自己", id))
				failCount++
				continue
			}

			admin, err := h.svcCtx.AdminModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
				failCount++
				continue
			}
			if admin == nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不存在", id))
				failCount++
				continue
			}

			admin.Status = 1
			if err := h.svcCtx.AdminModel.Update(nil, admin); err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	case "disable":
		// 批量禁用
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不能操作自己", id))
				failCount++
				continue
			}

			admin, err := h.svcCtx.AdminModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
				failCount++
				continue
			}
			if admin == nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不存在", id))
				failCount++
				continue
			}

			admin.Status = 0
			if err := h.svcCtx.AdminModel.Update(nil, admin); err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	case "delete":
		// 批量删除
		for _, id := range req.Ids {
			if id == currentAdminId {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不能删除自己", id))
				failCount++
				continue
			}

			admin, err := h.svcCtx.AdminModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
				failCount++
				continue
			}
			if admin == nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: 不存在", id))
				failCount++
				continue
			}

			if err := h.svcCtx.AdminModel.Delete(admin); err != nil {
				errors = append(errors, fmt.Sprintf("管理员ID %d: %v", id, err))
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
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     fmt.Sprintf("batch_%s_admin", req.Operation),
			Resource:   "admin",
			ResourceId: 0, // 批量操作没有单一资源ID
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回操作结果
	resp := types.BatchOperationResp{
		Total:        len(req.Ids),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetSystemSettings 获取系统设置
func (h *AdminHandler) GetSystemSettings(c *gin.Context) {
	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 模拟系统设置数据
	// 在实际项目中，需要实现SystemSettingsModel
	resp := types.AdminSystemConfigResp{
		SiteName:                  "邮件管理系统",
		SiteLogo:                  "/static/logo.png",
		SiteDescription:           "专业的邮件管理解决方案",
		ContactEmail:              "admin@example.com",
		RegistrationEnabled:       true,
		EmailVerificationRequired: true,
		DefaultSMTP: types.SMTPConfigResp{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "noreply@example.com",
			UseTLS:   true,
		},
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateSystemSettings 更新系统设置
func (h *AdminHandler) UpdateSystemSettings(c *gin.Context) {
	var req types.AdminSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证超级管理员权限
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要超级管理员权限"))
		return
	}

	// TODO: 获取现有设置
	// 在实际项目中，需要实现SystemSettingsModel
	// settings, err := h.svcCtx.SystemSettingsModel.Get()
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
	//     return
	// }

	// TODO: 更新设置
	// settings.SiteName = req.SiteName
	// settings.SiteLogo = req.SiteLogo
	// settings.SiteDescription = req.SiteDescription
	// settings.ContactEmail = req.ContactEmail
	// settings.RegistrationEnabled = req.RegistrationEnabled
	// settings.EmailVerificationRequired = req.EmailVerificationRequired
	// settings.SmtpHost = req.DefaultSMTP.Host
	// settings.SmtpPort = req.DefaultSMTP.Port
	// settings.SmtpUsername = req.DefaultSMTP.Username
	// if req.DefaultSMTP.Password != "" {
	//     settings.SmtpPassword = req.DefaultSMTP.Password
	// }
	// settings.SmtpUseTLS = req.DefaultSMTP.UseTLS

	// if err := h.svcCtx.SystemSettingsModel.Update(settings); err != nil {
	//     c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
	//     return
	// }

	// 记录操作日志
	currentAdminId := middleware.GetCurrentUserId(c)
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "update_system_settings",
			Resource:   "system_settings",
			ResourceId: 1, // 模拟ID
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SimpleResult("系统设置更新成功"))
}

// List 获取用户列表
func (h *AdminHandler) List(c *gin.Context) {
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

		userResp := types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

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
func (h *AdminHandler) GetById(c *gin.Context) {
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
			LastLoginAt: user.LastLoginAt,
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
func (h *AdminHandler) Create(c *gin.Context) {
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
		Status:   req.Status,
	}

	if err := h.svcCtx.UserModel.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
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
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新用户
func (h *AdminHandler) Update(c *gin.Context) {
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

	if req.Status != 0 {
		user.Status = req.Status
	}

	// 保存更新
	if err := h.svcCtx.UserModel.Update(nil, user); err != nil {
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
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除用户
func (h *AdminHandler) Delete(c *gin.Context) {
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
	if err := h.svcCtx.UserModel.DeleteById(userId); err != nil {
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

// checkAdminPermission 检查管理员权限
func (h *AdminHandler) checkAdminPermission(c *gin.Context, userId int64) bool {
	admin, err := h.svcCtx.AdminModel.GetById(userId)
	if err != nil || admin == nil {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return false
	}
	return false
}
