package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// ApiKeyHandler API密钥处理器
type ApiKeyHandler struct {
	svcCtx *svc.ServiceContext
}

// NewApiKeyHandler 创建API密钥处理器
func NewApiKeyHandler(svcCtx *svc.ServiceContext) *ApiKeyHandler {
	return &ApiKeyHandler{
		svcCtx: svcCtx,
	}
}

// List API密钥列表
func (h *ApiKeyHandler) List(c *gin.Context) {
	// TODO: 实现API密钥列表查询
	c.JSON(http.StatusOK, result.SimpleResult("API密钥列表"))
}

// Create 创建API密钥
func (h *ApiKeyHandler) Create(c *gin.Context) {
	// TODO: 实现创建API密钥
	c.JSON(http.StatusOK, result.SimpleResult("创建API密钥"))
}

// Update 更新API密钥
func (h *ApiKeyHandler) Update(c *gin.Context) {
	// TODO: 实现更新API密钥
	c.JSON(http.StatusOK, result.SimpleResult("更新API密钥"))
}

// Delete 删除API密钥
func (h *ApiKeyHandler) Delete(c *gin.Context) {
	// TODO: 实现删除API密钥
	c.JSON(http.StatusOK, result.SimpleResult("删除API密钥"))
}
