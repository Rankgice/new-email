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

// RuleHandler 规则处理器
type RuleHandler struct {
	svcCtx *svc.ServiceContext
}

// NewRuleHandler 创建规则处理器
func NewRuleHandler(svcCtx *svc.ServiceContext) *RuleHandler {
	return &RuleHandler{
		svcCtx: svcCtx,
	}
}

// ListVerificationRules 验证码规则列表
func (h *RuleHandler) ListVerificationRules(c *gin.Context) {
	var req types.VerificationRuleListReq
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
		Source:   req.Source,
		Status:   req.Status,
		Priority: req.Priority,
	}

	// 查询验证码规则列表
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
			Source:      rule.Source,
			Pattern:     rule.Pattern,
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

// CreateVerificationRule 创建验证码规则
func (h *RuleHandler) CreateVerificationRule(c *gin.Context) {
	var req types.VerificationRuleCreateReq
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

	// 创建验证码规则
	rule := &model.VerificationRule{
		Name:        req.Name,
		Source:      req.Source,
		Pattern:     req.Pattern,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
	}

	if err := h.svcCtx.VerificationRuleModel.Create(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_verification_rule",
		Resource:   "verification_rule",
		ResourceId: rule.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的规则信息
	resp := types.VerificationRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Source:      rule.Source,
		Pattern:     rule.Pattern,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateVerificationRule 更新验证码规则
func (h *RuleHandler) UpdateVerificationRule(c *gin.Context) {
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

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

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

	// 更新规则信息
	rule.Name = req.Name
	rule.Source = req.Source
	rule.Pattern = req.Pattern
	rule.Description = req.Description
	rule.Priority = req.Priority
	rule.Status = req.Status

	if err := h.svcCtx.VerificationRuleModel.Update(nil, rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_verification_rule",
		Resource:   "verification_rule",
		ResourceId: rule.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的规则信息
	resp := types.VerificationRuleResp{
		Id:          rule.Id,
		Name:        rule.Name,
		Source:      rule.Source,
		Pattern:     rule.Pattern,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// DeleteVerificationRule 删除验证码规则
func (h *RuleHandler) DeleteVerificationRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

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

	// 软删除规则
	if err := h.svcCtx.VerificationRuleModel.Delete(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_verification_rule",
		Resource:   "verification_rule",
		ResourceId: rule.Id,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}

// ListForwardRules 转发规则列表
func (h *RuleHandler) ListForwardRules(c *gin.Context) {
	var req types.ForwardRuleListReq
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
	params := model.ForwardRuleListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		BaseTimeRangeParams: model.BaseTimeRangeParams{
			CreatedAtStart: req.CreatedAtStart,
			CreatedAtEnd:   req.CreatedAtEnd,
		},
		UserId:      &currentUserId,
		Name:        req.Name,
		FromPattern: req.FromPattern,
		ToEmail:     req.ToEmail,
		Status:      req.Status,
	}

	// 查询转发规则列表
	rules, total, err := h.svcCtx.ForwardRuleModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var ruleList []types.ForwardRuleResp
	for _, rule := range rules {
		ruleList = append(ruleList, types.ForwardRuleResp{
			Id:          rule.Id,
			UserId:      rule.UserId,
			Name:        rule.Name,
			FromPattern: rule.FromPattern,
			ToEmail:     rule.ToEmail,
			Conditions:  rule.Conditions,
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

// CreateForwardRule 创建转发规则
func (h *RuleHandler) CreateForwardRule(c *gin.Context) {
	var req types.ForwardRuleCreateReq
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

	// 创建转发规则
	rule := &model.ForwardRule{
		UserId:      currentUserId,
		Name:        req.Name,
		FromPattern: req.FromPattern,
		ToEmail:     req.ToEmail,
		Conditions:  req.Conditions,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
	}

	if err := h.svcCtx.ForwardRuleModel.Create(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "create_forward_rule",
		Resource:   "forward_rule",
		ResourceId: rule.Id,
		Method:     "POST",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回创建的规则信息
	resp := types.ForwardRuleResp{
		Id:          rule.Id,
		UserId:      rule.UserId,
		Name:        rule.Name,
		FromPattern: rule.FromPattern,
		ToEmail:     rule.ToEmail,
		Conditions:  rule.Conditions,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// UpdateForwardRule 更新转发规则
func (h *RuleHandler) UpdateForwardRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req types.ForwardRuleUpdateReq
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

	// 检查规则是否存在
	rule, err := h.svcCtx.ForwardRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// 检查权限（只能更新自己的规则）
	if rule.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此规则"))
		return
	}

	// 更新规则信息
	rule.Name = req.Name
	rule.FromPattern = req.FromPattern
	rule.ToEmail = req.ToEmail
	rule.Conditions = req.Conditions
	rule.Description = req.Description
	rule.Priority = req.Priority
	rule.Status = req.Status

	if err := h.svcCtx.ForwardRuleModel.Update(nil, rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "update_forward_rule",
		Resource:   "forward_rule",
		ResourceId: rule.Id,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	// 返回更新后的规则信息
	resp := types.ForwardRuleResp{
		Id:          rule.Id,
		UserId:      rule.UserId,
		Name:        rule.Name,
		FromPattern: rule.FromPattern,
		ToEmail:     rule.ToEmail,
		Conditions:  rule.Conditions,
		Description: rule.Description,
		Priority:    rule.Priority,
		Status:      rule.Status,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// DeleteForwardRule 删除转发规则
func (h *RuleHandler) DeleteForwardRule(c *gin.Context) {
	// 获取规则ID
	idStr := c.Param("id")
	ruleId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 检查规则是否存在
	rule, err := h.svcCtx.ForwardRuleModel.GetById(uint(ruleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	// 检查权限（只能删除自己的规则）
	if rule.UserId != currentUserId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此规则"))
		return
	}

	// 软删除规则
	if err := h.svcCtx.ForwardRuleModel.Delete(rule); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorDelete.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "delete_forward_rule",
		Resource:   "forward_rule",
		ResourceId: rule.Id,
		Method:     "DELETE",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("删除成功"))
}
