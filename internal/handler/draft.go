package handler

import (
	"net/http"
	"new-email/internal/result"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// DraftHandler 草稿处理器
type DraftHandler struct {
	svcCtx *svc.ServiceContext
}

// NewDraftHandler 创建草稿处理器
func NewDraftHandler(svcCtx *svc.ServiceContext) *DraftHandler {
	return &DraftHandler{
		svcCtx: svcCtx,
	}
}

// List 草稿列表
func (h *DraftHandler) List(c *gin.Context) {
	// TODO: 实现草稿列表查询
	c.JSON(http.StatusOK, result.SimpleResult("草稿列表"))
}

// Create 创建草稿
func (h *DraftHandler) Create(c *gin.Context) {
	// TODO: 实现创建草稿
	c.JSON(http.StatusOK, result.SimpleResult("创建草稿"))
}

// Update 更新草稿
func (h *DraftHandler) Update(c *gin.Context) {
	// TODO: 实现更新草稿
	c.JSON(http.StatusOK, result.SimpleResult("更新草稿"))
}

// Delete 删除草稿
func (h *DraftHandler) Delete(c *gin.Context) {
	// TODO: 实现删除草稿
	c.JSON(http.StatusOK, result.SimpleResult("删除草稿"))
}

// Send 发送草稿
func (h *DraftHandler) Send(c *gin.Context) {
	// TODO: 实现发送草稿
	c.JSON(http.StatusOK, result.SimpleResult("发送草稿"))
}
