package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// VerificationCodeHandler 验证码处理器
type VerificationCodeHandler struct {
	svcCtx *svc.ServiceContext
}

// NewVerificationCodeHandler 创建验证码处理器
func NewVerificationCodeHandler(svcCtx *svc.ServiceContext) *VerificationCodeHandler {
	return &VerificationCodeHandler{
		svcCtx: svcCtx,
	}
}

// List 验证码列表
func (h *VerificationCodeHandler) List(c *gin.Context) {
	// TODO: 实现验证码列表查询
	c.JSON(http.StatusOK, result.SimpleResult("验证码列表"))
}

// GetById 获取验证码详情
func (h *VerificationCodeHandler) GetById(c *gin.Context) {
	// TODO: 实现获取验证码详情
	c.JSON(http.StatusOK, result.SimpleResult("获取验证码"))
}

// MarkUsed 标记验证码已使用
func (h *VerificationCodeHandler) MarkUsed(c *gin.Context) {
	// TODO: 实现标记验证码已使用
	c.JSON(http.StatusOK, result.SimpleResult("标记已使用"))
}
