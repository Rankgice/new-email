package handler

import (
	"fmt"
	"net"
	"net/http"
	"new-email/internal/constant"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"strconv"
	"strings"

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

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// Delete 删除域名
func (h *DomainHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	domainId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

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

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// Verify 验证域名DNS配置
func (h *DomainHandler) Verify(c *gin.Context) {
	idStr := c.Param("id")
	domainId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

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

	// 执行DNS验证
	verified, verifyResults := h.verifyDomainDNS(domain.Name)

	// 更新验证状态和DNS记录
	updateData := map[string]interface{}{
		"dns_verified": verified,
	}

	// 如果验证成功，生成并保存DNS记录
	if verified {
		// 生成SPF记录
		spfRecord := "v=spf1 mx a -all"
		updateData["spf_record"] = spfRecord

		// 生成DMARC记录
		dmarcRecord := "v=DMARC1; p=quarantine; rua=mailto:dmarc@" + domain.Name
		updateData["dmarc_record"] = dmarcRecord

		// 生成DKIM记录（简化版）
		dkimRecord := "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC..." // 示例公钥
		updateData["dkim_record"] = dkimRecord
	}

	if err := h.svcCtx.DomainModel.MapUpdate(nil, domainId, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	message := "DNS验证成功"
	if !verified {
		message = "DNS验证失败: " + strings.Join(verifyResults.Errors, "; ")
	}

	c.JSON(http.StatusOK, result.SuccessResult(map[string]interface{}{
		"verified": verified,
		"message":  message,
		"details":  verifyResults,
	}))
}

// DNSVerifyResult DNS验证结果
type DNSVerifyResult struct {
	MXRecords   []string `json:"mx_records"`
	SPFRecord   string   `json:"spf_record"`
	DKIMRecord  string   `json:"dkim_record"`
	DMARCRecord string   `json:"dmarc_record"`
	HasMX       bool     `json:"has_mx"`
	HasSPF      bool     `json:"has_spf"`
	HasDKIM     bool     `json:"has_dkim"`
	HasDMARC    bool     `json:"has_dmarc"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
}

// verifyDomainDNS 验证域名DNS配置
func (h *DomainHandler) verifyDomainDNS(domainName string) (bool, *DNSVerifyResult) {
	result := &DNSVerifyResult{
		MXRecords: make([]string, 0),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
	}

	// 1. 检查MX记录
	mxRecords, err := net.LookupMX(domainName)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("MX记录查询失败: %v", err))
	} else if len(mxRecords) == 0 {
		result.Errors = append(result.Errors, "未找到MX记录")
	} else {
		result.HasMX = true
		for _, mx := range mxRecords {
			result.MXRecords = append(result.MXRecords, fmt.Sprintf("%s (优先级:%d)", mx.Host, mx.Pref))
		}
	}

	// 2. 检查SPF记录
	txtRecords, err := net.LookupTXT(domainName)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("TXT记录查询失败: %v", err))
	} else {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				result.HasSPF = true
				result.SPFRecord = txt
				break
			}
		}
		if !result.HasSPF {
			result.Warnings = append(result.Warnings, "未找到SPF记录，建议添加")
		}
	}

	// 3. 检查DKIM记录（检查default选择器）
	dkimDomain := "default._domainkey." + domainName
	dkimRecords, err := net.LookupTXT(dkimDomain)
	if err != nil {
		result.Warnings = append(result.Warnings, "未找到DKIM记录，建议配置")
	} else {
		for _, txt := range dkimRecords {
			if strings.HasPrefix(txt, "v=DKIM1") {
				result.HasDKIM = true
				result.DKIMRecord = txt
				break
			}
		}
	}

	// 4. 检查DMARC记录
	dmarcDomain := "_dmarc." + domainName
	dmarcRecords, err := net.LookupTXT(dmarcDomain)
	if err != nil {
		result.Warnings = append(result.Warnings, "未找到DMARC记录，建议配置")
	} else {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(txt, "v=DMARC1") {
				result.HasDMARC = true
				result.DMARCRecord = txt
				break
			}
		}
	}

	// 验证通过条件：必须有MX记录
	verified := result.HasMX

	return verified, result
}

// GetById 根据ID获取域名信息
func (h *DomainHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的域名ID"))
		return
	}

	domain, err := h.svcCtx.DomainModel.GetById(id)
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

	// 返回操作结果
	resp := types.BatchOperationResp{
		Total:        len(req.Ids),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
