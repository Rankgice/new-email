package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// ApiHandler API处理器
type ApiHandler struct {
	svcCtx *svc.ServiceContext
}

// NewApiHandler 创建API处理器
func NewApiHandler(svcCtx *svc.ServiceContext) *ApiHandler {
	return &ApiHandler{
		svcCtx: svcCtx,
	}
}

// ListEmails API邮件列表
func (h *ApiHandler) ListEmails(c *gin.Context) {
	// TODO: 实现API邮件列表查询
	c.JSON(http.StatusOK, result.SimpleResult("API邮件列表"))
}

// GetEmail API获取邮件
func (h *ApiHandler) GetEmail(c *gin.Context) {
	// TODO: 实现API获取邮件详情
	c.JSON(http.StatusOK, result.SimpleResult("API获取邮件"))
}

// SendEmail API发送邮件
func (h *ApiHandler) SendEmail(c *gin.Context) {
	// TODO: 实现API发送邮件
	c.JSON(http.StatusOK, result.SimpleResult("API发送邮件"))
}

// ListVerificationCodes API验证码列表
func (h *ApiHandler) ListVerificationCodes(c *gin.Context) {
	// TODO: 实现API验证码列表查询
	c.JSON(http.StatusOK, result.SimpleResult("API验证码列表"))
}

// GetVerificationCode API获取验证码
func (h *ApiHandler) GetVerificationCode(c *gin.Context) {
	// TODO: 实现API获取验证码详情
	c.JSON(http.StatusOK, result.SimpleResult("API获取验证码"))
}
