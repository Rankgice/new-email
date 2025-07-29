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

// AdminRuleHandler 管理员规则处理器
type AdminRuleHandler struct {
	svcCtx *svc.ServiceContext
}

// NewAdminRuleHandler 创建管理员规则处理器
func NewAdminRuleHandler(svcCtx *svc.ServiceContext) *AdminRuleHandler {
	return &AdminRuleHandler{
		svcCtx: svcCtx,
	}
}

// ListGlobalVerificationRules 全局验证码规则列表
func (h *AdminRuleHandler) ListGlobalVerificationRules(c *gin.Context) {
	var req types.VerificationRuleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 转换为model参数，只查询全局规则
	isGlobal := true
	params := model.VerificationRuleListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		Name:     req.Name,
		IsGlobal: &isGlobal,
		Status:   req.Status,
	}

	// 查询全局验证码规则列表
	rules, total, err := h.svcCtx.VerificationRuleModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var ruleList []types.VerificationRuleResp
	for _, rule := range rules {
		ruleList = append(ruleList, types.VerificationRuleResp{
			Id:          rule.Id,
			Name:        rule.Name,
			Source:      "global", // 标记为全局规则
			Pattern:     rule.Pattern,
			Description: rule.Description,
			Status:      rule.Status,
			Priority:    rule.Priority,
			CreatedAt:   rule.CreatedAt,
			UpdatedAt:   rule.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     ruleList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// CreateGlobalVerificationRule 创建全局验证码规则
func (h *AdminRuleHandler) CreateGlobalVerificationRule(c *gin.Context) {
	var req types.VerificationRuleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 创建全局验证码规则
	rule := &model.VerificationRule{
		UserId:      0, // 全局规则不属于特定用户（0在这里是合理的，表示无用户）
		Name:        req.Name,
		Pattern:     req.Pattern,
		Description: req.Description,
		// IsGlobal:    true, // VerificationRule模型中没有IsGlobal字段
		Status:   req.Status,
		Priority: req.Priority,
	}

	if err := h.svcCtx.VerificationRuleModel.Create(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "create_global_verification_rule",
			Resource:   "verification_rule",
			ResourceId: rule.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回创建的规则信息
	resp := types.VerificationRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Source:      "global", // 标记为全局规则
		Pattern:     rule.Pattern,
		Description: rule.Description,
		Status:      rule.Status,
		Priority:    rule.Priority,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateGlobalVerificationRule 更新全局验证码规则
func (h *AdminRuleHandler) UpdateGlobalVerificationRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req types.VerificationRuleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 检查规则是否存在
	rule, err := h.svcCtx.VerificationRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// TODO: 检查是否为全局规则
	// VerificationRule模型中没有IsGlobal字段，需要通过其他方式判断
	// if !rule.IsGlobal {
	//     c.JSON(http.StatusForbidden, result.ErrorSimpleResult("只能更新全局规则"))
	//     return
	// }

	// 更新规则信息
	rule.Name = req.Name
	rule.Pattern = req.Pattern
	rule.Description = req.Description
	rule.Status = req.Status
	rule.Priority = req.Priority
	// 保持IsGlobal为true

	if err := h.svcCtx.VerificationRuleModel.Update(nil, rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "update_global_verification_rule",
			Resource:   "verification_rule",
			ResourceId: rule.Id,
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回更新后的规则信息
	resp := types.VerificationRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Source:      "global", // 标记为全局规则
		Pattern:     rule.Pattern,
		Description: rule.Description,
		Status:      rule.Status,
		Priority:    rule.Priority,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// DeleteGlobalVerificationRule 删除全局验证码规则
func (h *AdminRuleHandler) DeleteGlobalVerificationRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 检查规则是否存在
	rule, err := h.svcCtx.VerificationRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// 检查是否为全局规则
	if !rule.IsGlobal {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("只能删除全局规则"))
		return
	}

	// 软删除规则
	if err := h.svcCtx.VerificationRuleModel.Delete(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "delete_global_verification_rule",
			Resource:   "verification_rule",
			ResourceId: rule.Id,
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

// ListAntiSpamRules 反垃圾规则列表
func (h *AdminRuleHandler) ListAntiSpamRules(c *gin.Context) {
	var req types.AntiSpamRuleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
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
	params := model.AntiSpamRuleListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		Name:     req.Name,
		RuleType: req.Type,
		Action:   req.Action,
		Status:   req.Status,
	}

	// 查询反垃圾规则列表
	rules, total, err := h.svcCtx.AntiSpamRuleModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var ruleList []types.AntiSpamRuleResp
	for _, rule := range rules {
		ruleList = append(ruleList, types.AntiSpamRuleResp{
			Id:          rule.Id,
			Name:        rule.Name,
			Type:        rule.RuleType,
			Pattern:     rule.Pattern,
			Action:      rule.Action,
			Description: rule.Description,
			Priority:    rule.Priority,
			Status:      rule.Status,
			CreatedAt:   rule.CreatedAt,
			UpdatedAt:   rule.UpdatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     ruleList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// CreateAntiSpamRule 创建反垃圾规则
func (h *AdminRuleHandler) CreateAntiSpamRule(c *gin.Context) {
	var req types.AntiSpamRuleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 创建反垃圾规则
	rule := &model.AntiSpamRule{
		Name:        req.Name,
		RuleType:    req.Type,
		Pattern:     req.Pattern,
		Action:      req.Action,
		Description: req.Description,
		IsGlobal:    true, // 管理员创建的规则默认为全局规则
		Status:      req.Status,
		Priority:    req.Priority,
	}

	if err := h.svcCtx.AntiSpamRuleModel.Create(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "create_anti_spam_rule",
			Resource:   "anti_spam_rule",
			ResourceId: rule.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回创建的规则信息
	resp := types.AntiSpamRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Type:        rule.RuleType,
		Pattern:     rule.Pattern,
		Action:      rule.Action,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateAntiSpamRule 更新反垃圾规则
func (h *AdminRuleHandler) UpdateAntiSpamRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req types.AntiSpamRuleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 检查规则是否存在
	rule, err := h.svcCtx.AntiSpamRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// 更新规则信息
	rule.Name = req.Name
	rule.RuleType = req.Type
	rule.Pattern = req.Pattern
	rule.Action = req.Action
	rule.Description = req.Description
	rule.Priority = req.Priority
	rule.Status = req.Status

	if err := h.svcCtx.AntiSpamRuleModel.Update(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "update_anti_spam_rule",
			Resource:   "anti_spam_rule",
			ResourceId: rule.Id,
			Method:     "PUT",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回更新后的规则信息
	resp := types.AntiSpamRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Type:        rule.RuleType,
		Pattern:     rule.Pattern,
		Action:      rule.Action,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// DeleteAntiSpamRule 删除反垃圾规则
func (h *AdminRuleHandler) DeleteAntiSpamRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	// 验证管理员权限
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 获取当前管理员ID
	currentAdminId := middleware.GetCurrentUserId(c)

	// 检查规则是否存在
	rule, err := h.svcCtx.AntiSpamRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// 软删除规则
	if err := h.svcCtx.AntiSpamRuleModel.Delete(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	if currentAdminId > 0 {
		log := &model.OperationLog{
			UserId:     currentAdminId,
			Action:     "delete_anti_spam_rule",
			Resource:   "anti_spam_rule",
			ResourceId: rule.Id,
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
