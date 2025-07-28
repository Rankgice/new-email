package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// SignatureHandler 签名处理器
type SignatureHandler struct {
	svcCtx *svc.ServiceContext
}

// NewSignatureHandler 创建签名处理器
func NewSignatureHandler(svcCtx *svc.ServiceContext) *SignatureHandler {
	return &SignatureHandler{
		svcCtx: svcCtx,
	}
}

// List 签名列表
func (h *SignatureHandler) List(c *gin.Context) {
	// TODO: 实现签名列表查询
	c.JSON(http.StatusOK, result.SimpleResult("签名列表"))
}

// Create 创建签名
func (h *SignatureHandler) Create(c *gin.Context) {
	// TODO: 实现创建签名
	c.JSON(http.StatusOK, result.SimpleResult("创建签名"))
}

// Update 更新签名
func (h *SignatureHandler) Update(c *gin.Context) {
	// TODO: 实现更新签名
	c.JSON(http.StatusOK, result.SimpleResult("更新签名"))
}

// Delete 删除签名
func (h *SignatureHandler) Delete(c *gin.Context) {
	// TODO: 实现删除签名
	c.JSON(http.StatusOK, result.SimpleResult("删除签名"))
}
