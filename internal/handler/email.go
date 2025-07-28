package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// EmailHandler 邮件处理器
type EmailHandler struct {
	svcCtx *svc.ServiceContext
}

// NewEmailHandler 创建邮件处理器
func NewEmailHandler(svcCtx *svc.ServiceContext) *EmailHandler {
	return &EmailHandler{
		svcCtx: svcCtx,
	}
}

// List 邮件列表
func (h *EmailHandler) List(c *gin.Context) {
	// TODO: 实现邮件列表查询
	c.JSON(http.StatusOK, result.SimpleResult("邮件列表接口"))
}

// GetById 获取邮件详情
func (h *EmailHandler) GetById(c *gin.Context) {
	// TODO: 实现获取邮件详情
	c.JSON(http.StatusOK, result.SimpleResult("获取邮件详情接口"))
}

// Send 发送邮件
func (h *EmailHandler) Send(c *gin.Context) {
	// TODO: 实现发送邮件
	c.JSON(http.StatusOK, result.SimpleResult("发送邮件接口"))
}

// MarkRead 标记已读
func (h *EmailHandler) MarkRead(c *gin.Context) {
	// TODO: 实现标记已读
	c.JSON(http.StatusOK, result.SimpleResult("标记已读接口"))
}

// MarkStar 标记星标
func (h *EmailHandler) MarkStar(c *gin.Context) {
	// TODO: 实现标记星标
	c.JSON(http.StatusOK, result.SimpleResult("标记星标接口"))
}

// Delete 删除邮件
func (h *EmailHandler) Delete(c *gin.Context) {
	// TODO: 实现删除邮件
	c.JSON(http.StatusOK, result.SimpleResult("删除邮件接口"))
}

// BatchOperation 批量操作邮件
func (h *EmailHandler) BatchOperation(c *gin.Context) {
	// TODO: 实现批量操作邮件
	c.JSON(http.StatusOK, result.SimpleResult("批量操作邮件接口"))
}
