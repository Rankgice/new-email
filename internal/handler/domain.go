package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// DomainHandler 域名处理器
type DomainHandler struct {
	svcCtx *svc.ServiceContext
}

// NewDomainHandler 创建域名处理器
func NewDomainHandler(svcCtx *svc.ServiceContext) *DomainHandler {
	return &DomainHandler{
		svcCtx: svcCtx,
	}
}

// List 域名列表
func (h *DomainHandler) List(c *gin.Context) {
	// TODO: 实现域名列表查询
	c.JSON(http.StatusOK, result.SimpleResult("域名列表"))
}

// Create 创建域名
func (h *DomainHandler) Create(c *gin.Context) {
	// TODO: 实现创建域名
	c.JSON(http.StatusOK, result.SimpleResult("创建域名"))
}

// Update 更新域名
func (h *DomainHandler) Update(c *gin.Context) {
	// TODO: 实现更新域名
	c.JSON(http.StatusOK, result.SimpleResult("更新域名"))
}

// Delete 删除域名
func (h *DomainHandler) Delete(c *gin.Context) {
	// TODO: 实现删除域名
	c.JSON(http.StatusOK, result.SimpleResult("删除域名"))
}

// Verify 验证域名
func (h *DomainHandler) Verify(c *gin.Context) {
	// TODO: 实现验证域名
	c.JSON(http.StatusOK, result.SimpleResult("验证域名"))
}

// BatchOperation 批量操作域名
func (h *DomainHandler) BatchOperation(c *gin.Context) {
	// TODO: 实现批量操作域名
	c.JSON(http.StatusOK, result.SimpleResult("批量操作域名"))
}
