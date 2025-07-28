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
	var req types.TemplateListReq
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
	params := model.EmailTemplateListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId:   &currentUserId,
		Name:     req.Name,
		Category: req.Category,
		Status:   req.Status,
	}

	// 查询模板列表
	templates, total, err := h.svcCtx.EmailTemplateModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var templateList []types.TemplateResp
	for _, template := range templates {
		templateList = append(templateList, types.TemplateResp{
			Id:          template.Id,
			UserId:      template.UserId,
			Name:        template.Name,
			Category:    template.Category,
			Subject:     template.Subject,
			Content:     template.Content,
			Variables:   template.Variables,
			Description: template.Description,
			Status:      template.Status,
			CreatedAt:   template.CreatedAt,
			UpdatedAt:   template.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     templateList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建模板
func (h *TemplateHandler) Create(c *gin.Context) {
	var req types.TemplateCreateReq
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

	// 创建模板
	template := &model.EmailTemplate{
		UserId:      currentUserId,
		Name:        req.Name,
		Category:    req.Category,
		Subject:     req.Subject,
		Content:     req.Content,
		Variables:   req.Variables,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := h.svcCtx.EmailTemplateModel.Create(template); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_template",
		Resource:   "template",
		ResourceId: template.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的模板信息
	resp := types.TemplateResp{
		Id:          template.Id,
		UserId:      template.UserId,
		Name:        template.Name,
		Category:    template.Category,
		Subject:     template.Subject,
		Content:     template.Content,
		Variables:   template.Variables,
		Description: template.Description,
		Status:      template.Status,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新模板
func (h *TemplateHandler) Update(c *gin.Context) {
	// 获取模板ID
	idStr := c.Param("id")
	templateId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的模板ID"))
		return
	}

	var req types.TemplateUpdateReq
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

	// 检查模板是否存在
	template, err := h.svcCtx.EmailTemplateModel.GetById(uint(templateId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("模板不存在"))
		return
	}

	// 检查权限（只能更新自己的模板）
	if template.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此模板"))
		return
	}

	// 更新模板信息
	template.Name = req.Name
	template.Category = req.Category
	template.Subject = req.Subject
	template.Content = req.Content
	template.Variables = req.Variables
	template.Description = req.Description
	template.Status = req.Status

	if err := h.svcCtx.EmailTemplateModel.Update(nil, template); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_template",
		Resource:   "template",
		ResourceId: template.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的模板信息
	resp := types.TemplateResp{
		Id:          template.Id,
		UserId:      template.UserId,
		Name:        template.Name,
		Category:    template.Category,
		Subject:     template.Subject,
		Content:     template.Content,
		Variables:   template.Variables,
		Description: template.Description,
		Status:      template.Status,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除模板
func (h *TemplateHandler) Delete(c *gin.Context) {
	// 获取模板ID
	idStr := c.Param("id")
	templateId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的模板ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查模板是否存在
	template, err := h.svcCtx.EmailTemplateModel.GetById(uint(templateId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("模板不存在"))
		return
	}

	// 检查权限（只能删除自己的模板）
	if template.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此模板"))
		return
	}

	// 软删除模板
	if err := h.svcCtx.EmailTemplateModel.Delete(template); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_template",
		Resource:   "template",
		ResourceId: template.Id,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// GetById 获取模板详情
func (h *TemplateHandler) GetById(c *gin.Context) {
	// 获取模板ID
	idStr := c.Param("id")
	templateId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的模板ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询模板详情
	template, err := h.svcCtx.EmailTemplateModel.GetById(uint(templateId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("模板不存在"))
		return
	}

	// 检查权限（只能查看自己的模板）
	if template.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此模板"))
		return
	}

	// 返回模板详情
	resp := types.TemplateResp{
		Id:          template.Id,
		UserId:      template.UserId,
		Name:        template.Name,
		Category:    template.Category,
		Subject:     template.Subject,
		Content:     template.Content,
		Variables:   template.Variables,
		Description: template.Description,
		Status:      template.Status,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetCategories 获取模板分类
func (h *TemplateHandler) GetCategories(c *gin.Context) {
	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 获取用户的模板分类
	categories, err := h.svcCtx.EmailTemplateModel.GetCategoriesByUserId(currentUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(categories))
}

// Copy 复制模板
func (h *TemplateHandler) Copy(c *gin.Context) {
	// 获取模板ID
	idStr := c.Param("id")
	templateId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的模板ID"))
		return
	}

	var req types.TemplateCopyReq
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

	// 查询原模板
	originalTemplate, err := h.svcCtx.EmailTemplateModel.GetById(uint(templateId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if originalTemplate == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("模板不存在"))
		return
	}

	// 检查权限（只能复制自己的模板）
	if originalTemplate.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此模板"))
		return
	}

	// 创建新模板
	newTemplate := &model.EmailTemplate{
		UserId:      currentUserId,
		Name:        req.Name,
		Category:    originalTemplate.Category,
		Subject:     originalTemplate.Subject,
		Content:     originalTemplate.Content,
		Variables:   originalTemplate.Variables,
		Description: originalTemplate.Description,
		Status:      1, // 默认启用
	}

	if err := h.svcCtx.EmailTemplateModel.Create(newTemplate); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "copy_template",
		Resource:   "template",
		ResourceId: newTemplate.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回新模板信息
	resp := types.TemplateResp{
		Id:          newTemplate.Id,
		UserId:      newTemplate.UserId,
		Name:        newTemplate.Name,
		Category:    newTemplate.Category,
		Subject:     newTemplate.Subject,
		Content:     newTemplate.Content,
		Variables:   newTemplate.Variables,
		Description: newTemplate.Description,
		Status:      newTemplate.Status,
		CreatedAt:   newTemplate.CreatedAt,
		UpdatedAt:   newTemplate.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Delete 删除模板
func (h *TemplateHandler) Delete(c *gin.Context) {
	// TODO: 实现删除模板
	c.JSON(http.StatusOK, result.SimpleResult("删除模板"))
}
