package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	svcCtx *svc.ServiceContext
}

// NewTemplateHandler 创建模板处理器
func NewTemplateHandler(svcCtx *svc.ServiceContext) *TemplateHandler {
	return &TemplateHandler{
		svcCtx: svcCtx,
	}
}

// List 模板列表
func (h *TemplateHandler) List(c *gin.Context) {
	// TODO: 实现模板列表查询
	c.JSON(http.StatusOK, result.SimpleResult("模板列表"))
}

// Create 创建模板
func (h *TemplateHandler) Create(c *gin.Context) {
	// TODO: 实现创建模板
	c.JSON(http.StatusOK, result.SimpleResult("创建模板"))
}

// Update 更新模板
func (h *TemplateHandler) Update(c *gin.Context) {
	// TODO: 实现更新模板
	c.JSON(http.StatusOK, result.SimpleResult("更新模板"))
}

// Delete 删除模板
func (h *TemplateHandler) Delete(c *gin.Context) {
	// TODO: 实现删除模板
	c.JSON(http.StatusOK, result.SimpleResult("删除模板"))
}
