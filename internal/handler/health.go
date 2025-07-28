package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	svcCtx *svc.ServiceContext
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(svcCtx *svc.ServiceContext) *HealthHandler {
	return &HealthHandler{
		svcCtx: svcCtx,
	}
}

// Health 健康检查
func (h *HealthHandler) Health(c *gin.Context) {
	// 检查数据库连接
	dbStatus := "ok"
	sqlDB, err := h.svcCtx.DB.DB()
	if err != nil {
		dbStatus = "error: " + err.Error()
	} else {
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "error: " + err.Error()
		}
	}

	// 构建响应
	resp := types.HealthResp{
		Status:  "ok",
		Version: h.svcCtx.Config.App.Version,
		Uptime:  time.Since(time.Now()).String(), // 这里应该记录实际启动时间
		Services: map[string]string{
			"database": dbStatus,
			"app":      "ok",
		},
	}

	// 如果有服务异常，整体状态为error
	for _, status := range resp.Services {
		if status != "ok" {
			resp.Status = "error"
			break
		}
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
