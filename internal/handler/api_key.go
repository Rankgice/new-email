package handler

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/rankgice/new-email/internal/middleware"
	"github.com/rankgice/new-email/internal/model"
	"github.com/rankgice/new-email/internal/result"
	"github.com/rankgice/new-email/internal/svc"
	"github.com/rankgice/new-email/internal/types"
	"net/http"
	"strconv"
	"time"

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
	var req types.ApiKeyListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 转换为model参数
	params := model.ApiKeyListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId:      currentUserId,
		Name:        req.Name,
		Status:      req.Status,
		Permissions: req.Permissions,
	}

	// 查询API密钥列表
	apiKeys, total, err := h.svcCtx.ApiKeyModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var apiKeyList []types.ApiKeyResp
	for _, apiKey := range apiKeys {
		// 隐藏密钥的完整内容，只显示前8位和后4位
		maskedKey := ""
		if len(apiKey.Key) > 12 {
			maskedKey = apiKey.Key[:8] + "****" + apiKey.Key[len(apiKey.Key)-4:]
		} else {
			maskedKey = "****"
		}

		apiKeyList = append(apiKeyList, types.ApiKeyResp{
			Id:          apiKey.Id,
			UserId:      apiKey.UserId,
			Name:        apiKey.Name,
			Key:         maskedKey,
			Permissions: apiKey.Permissions,
			Status:      apiKey.Status,
			LastUsedAt:  apiKey.LastUsedAt,
			ExpiresAt:   apiKey.ExpiresAt,
			CreatedAt:   apiKey.CreatedAt,
			UpdatedAt:   apiKey.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     apiKeyList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建API密钥
func (h *ApiKeyHandler) Create(c *gin.Context) {
	var req types.ApiKeyCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 生成API密钥
	apiKey, err := h.generateApiKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成API密钥失败"))
		return
	}

	// 处理过期时间
	var expiresAt *time.Time
	if req.ExpiresAt != nil && !req.ExpiresAt.IsZero() {
		expiresAt = req.ExpiresAt
	}

	// 创建API密钥记录
	apiKeyRecord := &model.ApiKey{
		UserId:      currentUserId,
		Name:        req.Name,
		Key:         apiKey,
		Permissions: req.Permissions,
		Status:      req.Status,
		ExpiresAt:   expiresAt,
	}

	if err := h.svcCtx.ApiKeyModel.Create(apiKeyRecord); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 返回创建的API密钥信息（包含完整密钥，仅此一次显示）
	resp := types.ApiKeyCreateResp{
		Id:          apiKeyRecord.Id,
		UserId:      apiKeyRecord.UserId,
		Name:        apiKeyRecord.Name,
		Key:         apiKey, // 完整密钥，仅创建时返回
		Permissions: apiKeyRecord.Permissions,
		Status:      apiKeyRecord.Status,
		ExpiresAt:   apiKeyRecord.ExpiresAt,
		CreatedAt:   apiKeyRecord.CreatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新API密钥
func (h *ApiKeyHandler) Update(c *gin.Context) {
	// 获取API密钥ID
	idStr := c.Param("id")
	apiKeyId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的API密钥ID"))
		return
	}

	var req types.ApiKeyUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查API密钥是否存在
	apiKey, err := h.svcCtx.ApiKeyModel.GetById(apiKeyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if apiKey == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("API密钥不存在"))
		return
	}

	// 检查权限（只能更新自己的API密钥）
	if apiKey.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此API密钥"))
		return
	}

	// 更新API密钥信息
	apiKey.Name = req.Name
	apiKey.Permissions = req.Permissions
	apiKey.Status = req.Status
	if req.ExpiresAt != nil && !req.ExpiresAt.IsZero() {
		apiKey.ExpiresAt = req.ExpiresAt
	}

	if err := h.svcCtx.ApiKeyModel.Update(nil, apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 返回更新后的API密钥信息（隐藏密钥）
	maskedKey := ""
	if len(apiKey.Key) > 12 {
		maskedKey = apiKey.Key[:8] + "****" + apiKey.Key[len(apiKey.Key)-4:]
	} else {
		maskedKey = "****"
	}

	resp := types.ApiKeyResp{
		Id:          apiKey.Id,
		UserId:      apiKey.UserId,
		Name:        apiKey.Name,
		Key:         maskedKey,
		Permissions: apiKey.Permissions,
		Status:      apiKey.Status,
		LastUsedAt:  apiKey.LastUsedAt,
		ExpiresAt:   apiKey.ExpiresAt,
		CreatedAt:   apiKey.CreatedAt,
		UpdatedAt:   apiKey.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除API密钥
func (h *ApiKeyHandler) Delete(c *gin.Context) {
	// 获取API密钥ID
	idStr := c.Param("id")
	apiKeyId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的API密钥ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查API密钥是否存在
	apiKey, err := h.svcCtx.ApiKeyModel.GetById(apiKeyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if apiKey == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("API密钥不存在"))
		return
	}

	// 检查权限（只能删除自己的API密钥）
	if apiKey.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此API密钥"))
		return
	}

	// 软删除API密钥
	if err := h.svcCtx.ApiKeyModel.Delete(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// GetById 获取API密钥详情
func (h *ApiKeyHandler) GetById(c *gin.Context) {
	// 获取API密钥ID
	idStr := c.Param("id")
	apiKeyId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的API密钥ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询API密钥详情
	apiKey, err := h.svcCtx.ApiKeyModel.GetById(apiKeyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if apiKey == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("API密钥不存在"))
		return
	}

	// 检查权限（只能查看自己的API密钥）
	if apiKey.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此API密钥"))
		return
	}

	// 返回API密钥详情（隐藏密钥）
	maskedKey := ""
	if len(apiKey.Key) > 12 {
		maskedKey = apiKey.Key[:8] + "****" + apiKey.Key[len(apiKey.Key)-4:]
	} else {
		maskedKey = "****"
	}

	resp := types.ApiKeyResp{
		Id:          apiKey.Id,
		UserId:      apiKey.UserId,
		Name:        apiKey.Name,
		Key:         maskedKey,
		Permissions: apiKey.Permissions,
		Status:      apiKey.Status,
		LastUsedAt:  apiKey.LastUsedAt,
		ExpiresAt:   apiKey.ExpiresAt,
		CreatedAt:   apiKey.CreatedAt,
		UpdatedAt:   apiKey.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// generateApiKey 生成API密钥
func (h *ApiKeyHandler) generateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ak_" + hex.EncodeToString(bytes), nil
}
