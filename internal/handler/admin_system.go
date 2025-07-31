package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/service"
	"new-email/internal/svc"
	"new-email/internal/types"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminSystemHandler 管理员系统处理器
type AdminSystemHandler struct {
	svcCtx      *svc.ServiceContext
	systemStats *service.SystemStatsService
}

// NewAdminSystemHandler 创建管理员系统处理器
func NewAdminSystemHandler(svcCtx *svc.ServiceContext) *AdminSystemHandler {
	return &AdminSystemHandler{
		svcCtx:      svcCtx,
		systemStats: service.NewSystemStatsService(svcCtx.DB),
	}
}

// Login 管理员登录
func (h *AdminSystemHandler) Login(c *gin.Context) {
	var req types.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 查找管理员用户
	user, err := h.svcCtx.UserModel.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户名或密码错误"))
		return
	}

	// 检查是否为管理员
	if user.Role != "admin" && user.Role != "manager" {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无管理员权限"))
		return
	}

	// 验证密码
	if !h.svcCtx.UserModel.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户名或密码错误"))
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("账户已被禁用"))
		return
	}

	// 生成JWT令牌
	token, expiresAt, err := h.svcCtx.JwtService.GenerateToken(user.Id, user.Username, true, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成令牌失败"))
		return
	}

	// 更新最后登录信息
	user.LastLoginAt = &[]time.Time{time.Now()}[0]
	user.LastLoginIp = c.ClientIP()
	h.svcCtx.UserModel.Update(user)

	// 记录登录日志
	log := &model.OperationLog{
		UserId:     user.Id,
		Action:     "admin_login",
		Resource:   "admin",
		ResourceId: user.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回登录结果
	resp := types.AdminLoginResp{
		Token:     token,
		ExpiresAt: expiresAt,
		Admin: types.AdminInfo{
			Id:       user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Dashboard 管理员仪表板
func (h *AdminSystemHandler) Dashboard(c *gin.Context) {
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

	// 获取系统统计信息
	stats, err := h.systemStats.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取最近用户
	recentUsers, err := h.getRecentUsers(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取最近日志
	recentLogs, err := h.getRecentLogs(20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 获取系统警告
	systemAlerts := h.getSystemAlerts()

	// 返回仪表板数据
	resp := types.AdminDashboardResp{
		Stats:        *stats,
		RecentUsers:  recentUsers,
		RecentLogs:   recentLogs,
		SystemAlerts: systemAlerts,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetSystemStats 获取系统统计信息
func (h *AdminSystemHandler) GetSystemStats(c *gin.Context) {
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

	// 获取系统统计信息
	stats, err := h.systemStats.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(stats))
}

// GetSystemSettings 获取系统设置
func (h *AdminSystemHandler) GetSystemSettings(c *gin.Context) {
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

	// 获取系统设置
	settings, err := h.getSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(settings))
}

// UpdateSystemSettings 更新系统设置
func (h *AdminSystemHandler) UpdateSystemSettings(c *gin.Context) {
	var req types.AdminSystemSettingsReq
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

	// 更新系统设置
	if err := h.updateSystemSettings(&req); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_system_settings",
		Resource:   "system",
		ResourceId: 0,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 获取更新后的设置
	settings, err := h.getSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(settings))
}

// checkAdminPermission 检查管理员权限
func (h *AdminSystemHandler) checkAdminPermission(c *gin.Context, userId int64) bool {
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

// getRecentUsers 获取最近用户
func (h *AdminSystemHandler) getRecentUsers(limit int) ([]types.UserResp, error) {
	params := model.UserListParams{
		BaseListParams: model.BaseListParams{
			Page:     1,
			PageSize: limit,
		},
	}

	users, _, err := h.svcCtx.UserModel.List(params)
	if err != nil {
		return nil, err
	}

	var userList []types.UserResp
	for _, user := range users {
		userList = append(userList, types.UserResp{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Status:      user.Status,
			Role:        user.Role,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	return userList, nil
}

// getRecentLogs 获取最近日志
func (h *AdminSystemHandler) getRecentLogs(limit int) ([]types.OperationLogResp, error) {
	params := model.OperationLogListParams{
		BaseListParams: model.BaseListParams{
			Page:     1,
			PageSize: limit,
		},
	}

	logs, _, err := h.svcCtx.OperationLogModel.List(params)
	if err != nil {
		return nil, err
	}

	var logList []types.OperationLogResp
	for _, log := range logs {
		logList = append(logList, types.OperationLogResp{
			Id:         log.Id,
			UserId:     log.UserId,
			Action:     log.Action,
			Resource:   log.Resource,
			ResourceId: log.ResourceId,
			Method:     log.Method,
			Path:       log.Path,
			Ip:         log.Ip,
			UserAgent:  log.UserAgent,
			Status:     log.Status,
			CreatedAt:  log.CreatedAt,
		})
	}

	return logList, nil
}

// getSystemAlerts 获取系统警告
func (h *AdminSystemHandler) getSystemAlerts() []types.SystemAlert {
	var alerts []types.SystemAlert

	// 检查系统资源使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memUsage := float64(m.Alloc) / float64(m.Sys) * 100
	if memUsage > 80 {
		alerts = append(alerts, types.SystemAlert{
			Id:        1,
			Type:      "resource",
			Level:     "warning",
			Title:     "内存使用率过高",
			Message:   "当前内存使用率为 " + strconv.FormatFloat(memUsage, 'f', 1, 64) + "%",
			Status:    "active",
			CreatedAt: time.Now(),
		})
	}

	// 检查磁盘空间（这里简化处理）
	// 实际项目中应该检查真实的磁盘使用情况

	return alerts
}

// getSystemSettings 获取系统设置
func (h *AdminSystemHandler) getSystemSettings() (*types.AdminSystemSettingsResp, error) {
	// 这里应该从配置文件或数据库中读取系统设置
	// 简化处理，返回默认设置
	settings := &types.AdminSystemSettingsResp{
		SiteName:        "邮件管理系统",
		SiteDescription: "专业的邮件管理解决方案",
		SiteLogo:        "/static/logo.png",
		AllowRegister:   true,
		RequireInvite:   false,
		DefaultQuota:    1024 * 1024 * 1024, // 1GB
		MaxMailboxes:    10,
		UpdatedAt:       time.Now(),
	}

	return settings, nil
}

// updateSystemSettings 更新系统设置
func (h *AdminSystemHandler) updateSystemSettings(req *types.AdminSystemSettingsReq) error {
	// 这里应该将设置保存到配置文件或数据库中
	// 简化处理，直接返回成功
	return nil
}
