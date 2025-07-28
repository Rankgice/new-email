package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminLogHandler 管理员日志处理器
type AdminLogHandler struct {
	svcCtx *svc.ServiceContext
}

// NewAdminLogHandler 创建管理员日志处理器
func NewAdminLogHandler(svcCtx *svc.ServiceContext) *AdminLogHandler {
	return &AdminLogHandler{
		svcCtx: svcCtx,
	}
}

// ListOperationLogs 管理员操作日志列表
func (h *AdminLogHandler) ListOperationLogs(c *gin.Context) {
	var req types.LogListReq
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

	// 转换为model参数（管理员可以查看所有日志）
	params := model.OperationLogListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		UserId:   0, // 默认值
		Action:   req.Action,
		Resource: req.Resource,
		Method:   req.Method,
		Status:   req.Status,
		Ip:       req.Ip,
	}

	// 处理可选的UserId参数
	if req.UserId != nil {
		params.UserId = *req.UserId
	}

	// 查询操作日志列表
	logs, total, err := h.svcCtx.OperationLogModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
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
			ErrorMsg:   "", // OperationLog模型中没有ErrorMsg字段，使用空字符串
			CreatedAt:  log.CreatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     logList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// ListEmailLogs 管理员邮件日志列表
func (h *AdminLogHandler) ListEmailLogs(c *gin.Context) {
	var req types.EmailLogListReq
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

	// 转换为model参数（管理员可以查看所有邮件日志）
	params := model.EmailLogListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		EmailId:   0, // 默认值
		MailboxId: 0, // 默认值
		Type:      req.Type,
		Status:    "", // 默认值
	}

	// 处理可选参数
	if req.EmailId != nil {
		params.EmailId = *req.EmailId
	}
	if req.MailboxId != nil {
		params.MailboxId = *req.MailboxId
	}
	if req.Status != nil {
		// 将int转换为string，这里需要根据实际业务逻辑调整
		// params.Status = strconv.Itoa(*req.Status)
	}

	// 查询邮件日志列表
	logs, total, err := h.svcCtx.EmailLogModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var logList []types.EmailLogResp
	for _, log := range logs {
		// 将string类型的Status转换为int
		statusInt := 0
		if log.Status == "success" {
			statusInt = 1
		} else if log.Status == "failed" {
			statusInt = 2
		}

		logList = append(logList, types.EmailLogResp{
			Id:        log.Id,
			EmailId:   log.EmailId,
			MailboxId: log.MailboxId,
			Type:      log.Type,
			Status:    statusInt,
			FromEmail: "", // EmailLog模型中没有FromEmail字段
			ToEmail:   "", // EmailLog模型中没有ToEmail字段
			Subject:   "", // EmailLog模型中没有Subject字段
			ErrorMsg:  log.ErrorMsg,
			CreatedAt: log.CreatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     logList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// ListSystemLogs 系统日志列表
func (h *AdminLogHandler) ListSystemLogs(c *gin.Context) {
	var req types.LogListReq
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

	// 系统日志通常是文件日志，这里模拟返回系统运行日志
	// 实际项目中可能需要读取日志文件或从日志收集系统获取

	// 模拟系统日志数据
	systemLogs := []types.SystemLogResp{
		{
			Id:        1,
			Level:     "INFO",
			Message:   "系统启动成功",
			Module:    "system",
			CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			Id:        2,
			Level:     "WARN",
			Message:   "数据库连接池达到最大连接数",
			Module:    "database",
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
		{
			Id:        3,
			Level:     "ERROR",
			Message:   "邮件发送失败：SMTP连接超时",
			Module:    "email",
			CreatedAt: time.Now().Add(-15 * time.Minute),
		},
		{
			Id:        4,
			Level:     "INFO",
			Message:   "定时任务执行完成",
			Module:    "scheduler",
			CreatedAt: time.Now().Add(-5 * time.Minute),
		},
	}

	// 简单的过滤逻辑
	var filteredLogs []types.SystemLogResp
	for _, log := range systemLogs {
		// 按级别过滤
		if req.Level != "" && log.Level != req.Level {
			continue
		}
		// 按模块过滤
		if req.Module != "" && log.Module != req.Module {
			continue
		}
		// 按时间范围过滤
		if !req.CreatedAtStart.IsZero() && log.CreatedAt.Before(req.CreatedAtStart) {
			continue
		}
		if !req.CreatedAtEnd.IsZero() && log.CreatedAt.After(req.CreatedAtEnd) {
			continue
		}
		filteredLogs = append(filteredLogs, log)
	}

	// 简单分页
	total := int64(len(filteredLogs))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start > len(filteredLogs) {
		filteredLogs = []types.SystemLogResp{}
	} else {
		if end > len(filteredLogs) {
			end = len(filteredLogs)
		}
		filteredLogs = filteredLogs[start:end]
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     filteredLogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
