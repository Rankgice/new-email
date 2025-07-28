package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// CommonHandler 通用处理器
type CommonHandler struct {
	svcCtx *svc.ServiceContext
}

// NewCommonHandler 创建通用处理器
func NewCommonHandler(svcCtx *svc.ServiceContext) *CommonHandler {
	return &CommonHandler{
		svcCtx: svcCtx,
	}
}

// SendCode 发送验证码
func (h *CommonHandler) SendCode(c *gin.Context) {
	// TODO: 实现发送验证码功能
	c.JSON(http.StatusOK, result.SimpleResult("发送验证码接口"))
}

// VerifyCode 验证验证码
func (h *CommonHandler) VerifyCode(c *gin.Context) {
	// TODO: 实现验证验证码功能
	c.JSON(http.StatusOK, result.SimpleResult("验证验证码接口"))
}
