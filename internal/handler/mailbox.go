package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// MailboxHandler 邮箱处理器
type MailboxHandler struct {
	svcCtx *svc.ServiceContext
}

// NewMailboxHandler 创建邮箱处理器
func NewMailboxHandler(svcCtx *svc.ServiceContext) *MailboxHandler {
	return &MailboxHandler{
		svcCtx: svcCtx,
	}
}

// List 邮箱列表
func (h *MailboxHandler) List(c *gin.Context) {
	// TODO: 实现邮箱列表查询
	c.JSON(http.StatusOK, result.SimpleResult("邮箱列表接口"))
}

// Create 创建邮箱
func (h *MailboxHandler) Create(c *gin.Context) {
	// TODO: 实现创建邮箱
	c.JSON(http.StatusOK, result.SimpleResult("创建邮箱接口"))
}

// Update 更新邮箱
func (h *MailboxHandler) Update(c *gin.Context) {
	// TODO: 实现更新邮箱
	c.JSON(http.StatusOK, result.SimpleResult("更新邮箱接口"))
}

// Delete 删除邮箱
func (h *MailboxHandler) Delete(c *gin.Context) {
	// TODO: 实现删除邮箱
	c.JSON(http.StatusOK, result.SimpleResult("删除邮箱接口"))
}

// Sync 同步邮箱
func (h *MailboxHandler) Sync(c *gin.Context) {
	// TODO: 实现同步邮箱
	c.JSON(http.StatusOK, result.SimpleResult("同步邮箱接口"))
}

// TestConnection 测试邮箱连接
func (h *MailboxHandler) TestConnection(c *gin.Context) {
	// TODO: 实现测试邮箱连接
	c.JSON(http.StatusOK, result.SimpleResult("测试邮箱连接接口"))
}
