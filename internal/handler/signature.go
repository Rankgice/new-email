package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"strconv"

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
	var req types.SignatureListReq
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
	params := model.EmailSignatureListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId: &currentUserId,
		Name:   req.Name,
		Status: req.Status,
	}

	// 查询签名列表
	signatures, total, err := h.svcCtx.EmailSignatureModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var signatureList []types.SignatureResp
	for _, signature := range signatures {
		signatureList = append(signatureList, types.SignatureResp{
			Id:        signature.Id,
			UserId:    signature.UserId,
			Name:      signature.Name,
			Content:   signature.Content,
			IsDefault: signature.IsDefault,
			Status:    signature.Status,
			CreatedAt: signature.CreatedAt,
			UpdatedAt: signature.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     signatureList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建签名
func (h *SignatureHandler) Create(c *gin.Context) {
	var req types.SignatureCreateReq
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

	// 如果设置为默认签名，需要先取消其他默认签名
	if req.IsDefault {
		if err := h.svcCtx.EmailSignatureModel.ClearDefaultByUserId(currentUserId); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
	}

	// 创建签名
	signature := &model.EmailSignature{
		UserId:    currentUserId,
		Name:      req.Name,
		Content:   req.Content,
		IsDefault: req.IsDefault,
		Status:    req.Status,
	}

	if err := h.svcCtx.EmailSignatureModel.Create(signature); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_signature",
		Resource:   "signature",
		ResourceId: signature.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的签名信息
	resp := types.SignatureResp{
		Id:        signature.Id,
		UserId:    signature.UserId,
		Name:      signature.Name,
		Content:   signature.Content,
		IsDefault: signature.IsDefault,
		Status:    signature.Status,
		CreatedAt: signature.CreatedAt,
		UpdatedAt: signature.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新签名
func (h *SignatureHandler) Update(c *gin.Context) {
	// 获取签名ID
	idStr := c.Param("id")
	signatureId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的签名ID"))
		return
	}

	var req types.SignatureUpdateReq
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

	// 检查签名是否存在
	signature, err := h.svcCtx.EmailSignatureModel.GetById(uint(signatureId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if signature == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("签名不存在"))
		return
	}

	// 检查权限（只能更新自己的签名）
	if signature.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此签名"))
		return
	}

	// 如果设置为默认签名，需要先取消其他默认签名
	if req.IsDefault && !signature.IsDefault {
		if err := h.svcCtx.EmailSignatureModel.ClearDefaultByUserId(currentUserId); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
	}

	// 更新签名信息
	signature.Name = req.Name
	signature.Content = req.Content
	signature.IsDefault = req.IsDefault
	signature.Status = req.Status

	if err := h.svcCtx.EmailSignatureModel.Update(nil, signature); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_signature",
		Resource:   "signature",
		ResourceId: signature.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的签名信息
	resp := types.SignatureResp{
		Id:        signature.Id,
		UserId:    signature.UserId,
		Name:      signature.Name,
		Content:   signature.Content,
		IsDefault: signature.IsDefault,
		Status:    signature.Status,
		CreatedAt: signature.CreatedAt,
		UpdatedAt: signature.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除签名
func (h *SignatureHandler) Delete(c *gin.Context) {
	// 获取签名ID
	idStr := c.Param("id")
	signatureId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的签名ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查签名是否存在
	signature, err := h.svcCtx.EmailSignatureModel.GetById(uint(signatureId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if signature == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("签名不存在"))
		return
	}

	// 检查权限（只能删除自己的签名）
	if signature.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此签名"))
		return
	}

	// 软删除签名
	if err := h.svcCtx.EmailSignatureModel.Delete(signature); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_signature",
		Resource:   "signature",
		ResourceId: signature.Id,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// GetById 获取签名详情
func (h *SignatureHandler) GetById(c *gin.Context) {
	// 获取签名ID
	idStr := c.Param("id")
	signatureId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的签名ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询签名详情
	signature, err := h.svcCtx.EmailSignatureModel.GetById(uint(signatureId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if signature == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("签名不存在"))
		return
	}

	// 检查权限（只能查看自己的签名）
	if signature.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此签名"))
		return
	}

	// 返回签名详情
	resp := types.SignatureResp{
		Id:        signature.Id,
		UserId:    signature.UserId,
		Name:      signature.Name,
		Content:   signature.Content,
		IsDefault: signature.IsDefault,
		Status:    signature.Status,
		CreatedAt: signature.CreatedAt,
		UpdatedAt: signature.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// SetDefault 设置默认签名
func (h *SignatureHandler) SetDefault(c *gin.Context) {
	// 获取签名ID
	idStr := c.Param("id")
	signatureId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的签名ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查签名是否存在
	signature, err := h.svcCtx.EmailSignatureModel.GetById(uint(signatureId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if signature == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("签名不存在"))
		return
	}

	// 检查权限（只能设置自己的签名）
	if signature.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此签名"))
		return
	}

	// 先取消其他默认签名
	if err := h.svcCtx.EmailSignatureModel.ClearDefaultByUserId(currentUserId); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 设置为默认签名
	signature.IsDefault = true
	if err := h.svcCtx.EmailSignatureModel.Update(nil, signature); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "set_default_signature",
		Resource:   "signature",
		ResourceId: signature.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("设置成功"))
}

// Delete 删除签名
func (h *SignatureHandler) Delete(c *gin.Context) {
	// TODO: 实现删除签名
	c.JSON(http.StatusOK, result.SimpleResult("删除签名"))
}
