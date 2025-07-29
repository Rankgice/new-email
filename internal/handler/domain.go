package handler

import (
	"fmt"
	"net/http"
	"new-email/internal/constant"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"strconv"

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
	var req types.DomainListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 转换为model参数
	params := model.DomainListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		Name:        req.Name,
		Status:      req.Status,
		DnsVerified: req.DnsVerified,
	}

	// 查询域名列表
	domains, total, err := h.svcCtx.DomainModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var domainList []types.DomainResp
	for _, domain := range domains {
		domainList = append(domainList, types.DomainResp{
			Id:          domain.Id,
			Name:        domain.Name,
			Status:      domain.Status,
			DnsVerified: domain.DnsVerified,
			DkimRecord:  domain.DkimRecord,
			SpfRecord:   domain.SpfRecord,
			DmarcRecord: domain.DmarcRecord,
			CreatedAt:   domain.CreatedAt,
			UpdatedAt:   domain.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     domainList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Create 创建域名
func (h *DomainHandler) Create(c *gin.Context) {
	var req types.DomainCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查域名是否已存在
	if exists, err := h.svcCtx.DomainModel.CheckNameExists(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	} else if exists {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名已存在"))
		return
	}

	// 创建域名
	domain := &model.Domain{
		Name:   req.Name,
		Status: constant.StatusEnabled, // 默认启用
	}

	if err := h.svcCtx.DomainModel.Create(domain); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "create_domain",
			Resource:   "domain",
			ResourceId: domain.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	resp := types.DomainResp{
		Id:          domain.Id,
		Name:        domain.Name,
		Status:      domain.Status,
		DnsVerified: domain.DnsVerified,
		DkimRecord:  domain.DkimRecord,
		SpfRecord:   domain.SpfRecord,
		DmarcRecord: domain.DmarcRecord,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Update 更新域名
func (h *DomainHandler) Update(c *gin.Context) {
	var req types.DomainUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 检查域名是否存在
	domain, err := h.svcCtx.DomainModel.GetById(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if domain == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("域名不存在"))
		return
	}

	// 检查域名是否已被其他记录使用
	if req.Name != domain.Name {
		if exists, err := h.svcCtx.DomainModel.CheckNameExists(req.Name); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		} else if exists {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名已被其他记录使用"))
			return
		}
	}

	// 构建更新数据
	updateData := map[string]interface{}{
		"name":   req.Name,
		"status": req.Status,
	}

	// 更新域名信息
	if err := h.svcCtx.DomainModel.MapUpdate(nil, req.Id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "update_domain",
			Resource:   "domain",
			ResourceId: req.Id,
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// Delete 删除域名
func (h *DomainHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

	domainId := uint(id)

	// 检查域名是否存在
	domain, err := h.svcCtx.DomainModel.GetById(domainId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if domain == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("域名不存在"))
		return
	}

	// 检查是否有关联的邮箱（可选：防止删除正在使用的域名）
	// 这里可以添加检查逻辑

	// 软删除域名
	if err := h.svcCtx.DomainModel.Delete(domain); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "delete_domain",
			Resource:   "domain",
			ResourceId: domainId,
			Method:     "DELETE",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// Verify 验证域名DNS配置
func (h *DomainHandler) Verify(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

	domainId := uint(id)

	// 检查域名是否存在
	domain, err := h.svcCtx.DomainModel.GetById(domainId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if domain == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("域名不存在"))
		return
	}

	// TODO: 实现DNS验证逻辑
	// 这里应该检查域名的DNS记录是否正确配置
	// 包括MX记录、SPF记录、DKIM记录、DMARC记录等

	// 暂时模拟验证成功
	verified := true

	// 更新验证状态
	if err := h.svcCtx.DomainModel.MapUpdate(nil, domainId, map[string]interface{}{
		"dns_verified": verified,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "verify_domain",
			Resource:   "domain",
			ResourceId: domainId,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	message := "DNS验证成功"
	if !verified {
		message = "DNS验证失败"
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"verified": verified,
		"message":  message,
	}))
}

// GetById 根据ID获取域名信息
func (h *DomainHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

	domain, err := h.svcCtx.DomainModel.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if domain == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("域名不存在"))
		return
	}

	resp := types.DomainResp{
		Id:          domain.Id,
		Name:        domain.Name,
		Status:      domain.Status,
		DnsVerified: domain.DnsVerified,
		DkimRecord:  domain.DkimRecord,
		SpfRecord:   domain.SpfRecord,
		DmarcRecord: domain.DmarcRecord,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// BatchOperation 批量操作域名
func (h *DomainHandler) BatchOperation(c *gin.Context) {
	var req types.DomainBatchOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请选择要操作的域名"))
		return
	}

	var successCount, failCount int
	var errors []string

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	switch req.Operation {
	case "enable":
		// 批量启用
		err := h.svcCtx.DomainModel.BatchUpdateStatus(req.Ids, constant.StatusEnabled)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
		successCount = len(req.Ids)

	case "disable":
		// 批量禁用
		err := h.svcCtx.DomainModel.BatchUpdateStatus(req.Ids, constant.StatusDisabled)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
			return
		}
		successCount = len(req.Ids)

	case "delete":
		// 批量删除
		for _, id := range req.Ids {
			domain, err := h.svcCtx.DomainModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: %v", id, err))
				failCount++
				continue
			}
			if domain == nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: 不存在", id))
				failCount++
				continue
			}

			if err := h.svcCtx.DomainModel.Delete(domain); err != nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	case "verify":
		// 批量验证DNS
		for _, id := range req.Ids {
			domain, err := h.svcCtx.DomainModel.GetById(id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: %v", id, err))
				failCount++
				continue
			}
			if domain == nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: 不存在", id))
				failCount++
				continue
			}

			// TODO: 实现实际的DNS验证逻辑
			// 这里应该调用DNS验证服务
			verified := true // 模拟验证结果

			// 更新验证状态
			updateData := map[string]interface{}{
				"dns_verified": verified,
			}
			if err := h.svcCtx.DomainModel.MapUpdate(nil, id, updateData); err != nil {
				errors = append(errors, fmt.Sprintf("域名ID %d: %v", id, err))
				failCount++
				continue
			}
			successCount++
		}

	default:
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不支持的操作类型"))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     fmt.Sprintf("batch_%s_domain", req.Operation),
			Resource:   "domain",
			ResourceId: 0, // 批量操作没有单一资源ID
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回操作结果
	resp := types.BatchOperationResp{
		Total:        len(req.Ids),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
