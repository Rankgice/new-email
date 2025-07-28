package handler

import (
	"net/http"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"

	"github.com/gin-gonic/gin"
)

// LogHandler 日志处理器
type LogHandler struct {
	svcCtx *svc.ServiceContext
}

// NewLogHandler 创建日志处理器
func NewLogHandler(svcCtx *svc.ServiceContext) *LogHandler {
	return &LogHandler{
		svcCtx: svcCtx,
	}
}

// ListOperationLogs 操作日志列表
func (h *LogHandler) ListOperationLogs(c *gin.Context) {
	var req types.LogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
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
	params := model.OperationLogListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId:   req.UserId,
		Action:   req.Action,
		Resource: req.Resource,
		Method:   req.Method,
		Status:   req.Status,
		Ip:       req.Ip,
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
			ErrorMsg:   log.ErrorMsg,
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

// ListEmailLogs 邮件日志列表
func (h *LogHandler) ListEmailLogs(c *gin.Context) {
	var req types.EmailLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
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
	params := model.EmailLogListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		EmailId:   req.EmailId,
		MailboxId: req.MailboxId,
		Type:      req.Type,
		Status:    req.Status,
		FromEmail: req.FromEmail,
		ToEmail:   req.ToEmail,
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
		logList = append(logList, types.EmailLogResp{
			Id:        log.Id,
			EmailId:   log.EmailId,
			MailboxId: log.MailboxId,
			Type:      log.Type,
			Status:    log.Status,
			FromEmail: log.FromEmail,
			ToEmail:   log.ToEmail,
			Subject:   log.Subject,
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
