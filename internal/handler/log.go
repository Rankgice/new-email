package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

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
	// TODO: 实现操作日志列表查询
	c.JSON(http.StatusOK, result.SimpleResult("操作日志列表"))
}

// ListEmailLogs 邮件日志列表
func (h *LogHandler) ListEmailLogs(c *gin.Context) {
	// TODO: 实现邮件日志列表查询
	c.JSON(http.StatusOK, result.SimpleResult("邮件日志列表"))
}
