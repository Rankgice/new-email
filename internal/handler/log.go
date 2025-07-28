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
		UserId:   0, // 默认值
		Action:   req.Action,
		Resource: req.Resource,
		Method:   req.Method,
		Status:   req.Status, // 直接使用*int类型
		Ip:       req.Ip,
	}

	// 处理可选参数
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
			ErrorMsg:   "", // OperationLog模型中没有ErrorMsg字段
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
			Status:    statusInt, // 使用转换后的int值
			FromEmail: "",        // EmailLog模型中没有FromEmail字段
			ToEmail:   "",        // EmailLog模型中没有ToEmail字段
			Subject:   "",        // EmailLog模型中没有Subject字段
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
