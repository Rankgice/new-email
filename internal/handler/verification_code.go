package handler

import (
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// VerificationCodeHandler 验证码处理器
type VerificationCodeHandler struct {
	svcCtx *svc.ServiceContext
}

// NewVerificationCodeHandler 创建验证码处理器
func NewVerificationCodeHandler(svcCtx *svc.ServiceContext) *VerificationCodeHandler {
	return &VerificationCodeHandler{
		svcCtx: svcCtx,
	}
}

// List 验证码列表
func (h *VerificationCodeHandler) List(c *gin.Context) {
	var req types.VerificationCodeListReq
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
	params := model.VerificationCodeListParams{
		BaseListParams: model.BaseListParams{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		// BaseTimeRangeParams: model.BaseTimeRangeParams{ // VerificationCodeListParams中没有BaseTimeRangeParams
		//     CreatedAtStart: req.CreatedAtStart,
		//     CreatedAtEnd:   req.CreatedAtEnd,
		// },
		EmailId: 0, // 默认值
		RuleId:  0, // 默认值
		Code:    req.Code,
		Source:  req.Source,
		IsUsed:  req.IsUsed,
	}

	// 处理可选参数
	if req.EmailId != nil {
		params.EmailId = *req.EmailId
	}

	// 查询验证码列表
	codes, total, err := h.svcCtx.VerificationCodeModel.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	// 转换为响应格式
	var codeList []types.VerificationCodeResp
	for _, code := range codes {
		codeList = append(codeList, types.VerificationCodeResp{
			Id:        code.Id,
			UserId:    0, // VerificationCode模型中没有UserId字段
			EmailId:   code.EmailId,
			Code:      code.Code,
			Source:    code.Source,
			IsUsed:    code.IsUsed,
			IsExpired: false, // VerificationCode模型中没有IsExpired字段
			UsedAt:    code.UsedAt,
			ExpiresAt: time.Time{}, // VerificationCode模型中没有ExpiresAt字段
			CreatedAt: code.CreatedAt,
		})
	}

	// 返回分页结果
	resp := types.PageResp{
		List:     codeList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetById 获取验证码详情
func (h *VerificationCodeHandler) GetById(c *gin.Context) {
	// 获取验证码ID
	idStr := c.Param("id")
	codeId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询验证码详情
	code, err := h.svcCtx.VerificationCodeModel.GetById(codeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if code == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	// TODO: 检查权限（只能查看自己的验证码）
	// VerificationCode模型中没有UserId字段，需要通过EmailId关联查询
	// if code.UserId != currentUserId {
	//     c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此验证码"))
	//     return
	// }

	// 返回验证码详情
	resp := types.VerificationCodeResp{
		Id:        code.Id,
		UserId:    0, // VerificationCode模型中没有UserId字段
		EmailId:   code.EmailId,
		Code:      code.Code,
		Source:    code.Source,
		IsUsed:    code.IsUsed,
		IsExpired: false, // VerificationCode模型中没有IsExpired字段
		UsedAt:    code.UsedAt,
		ExpiresAt: time.Time{}, // VerificationCode模型中没有ExpiresAt字段
		CreatedAt: code.CreatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// MarkAsUsed 标记验证码已使用
func (h *VerificationCodeHandler) MarkAsUsed(c *gin.Context) {
	// 获取验证码ID
	idStr := c.Param("id")
	codeId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询验证码
	code, err := h.svcCtx.VerificationCodeModel.GetById(codeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if code == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	// TODO: 检查权限（只能操作自己的验证码）
	// VerificationCode模型中没有UserId字段，需要通过EmailId关联查询
	// if code.UserId != currentUserId {
	//     c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此验证码"))
	//     return
	// }

	// 检查验证码是否已使用
	if code.IsUsed {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码已使用"))
		return
	}

	// 标记为已使用
	if err := h.svcCtx.VerificationCodeModel.MarkAsUsed(codeId); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "mark_verification_code_used",
		Resource:   "verification_code",
		ResourceId: codeId,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("标记成功"))
}

// GetLatest 获取最新验证码
func (h *VerificationCodeHandler) GetLatest(c *gin.Context) {
	var req types.VerificationCodeLatestReq
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

	// 查询最新验证码
	code, err := h.svcCtx.VerificationCodeModel.GetLatestBySource(req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if code == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("未找到验证码"))
		return
	}

	// 返回验证码详情
	resp := types.VerificationCodeResp{
		Id:        code.Id,
		UserId:    0, // VerificationCode模型中没有UserId字段
		EmailId:   code.EmailId,
		Code:      code.Code,
		Source:    code.Source,
		IsUsed:    code.IsUsed,
		IsExpired: false, // VerificationCode模型中没有IsExpired字段
		UsedAt:    code.UsedAt,
		ExpiresAt: time.Time{}, // VerificationCode模型中没有ExpiresAt字段
		CreatedAt: code.CreatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// GetStatistics 获取验证码统计
func (h *VerificationCodeHandler) GetStatistics(c *gin.Context) {
	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 获取统计数据
	stats, err := h.svcCtx.VerificationCodeModel.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(stats))
}

// MarkUsed 标记验证码已使用
func (h *VerificationCodeHandler) MarkUsed(c *gin.Context) {
	// 获取验证码ID
	idStr := c.Param("id")
	codeId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	// 获取当前用户ID
	currentUserId := middleware.GetCurrentUserId(c)
	if currentUserId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return
	}

	// 查询验证码
	code, err := h.svcCtx.VerificationCodeModel.GetById(codeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if code == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	// TODO: 检查权限（只能操作自己的验证码）
	// VerificationCode模型中没有UserId字段，需要通过EmailId关联查询
	// if code.UserId != currentUserId {
	//     c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此验证码"))
	//     return
	// }

	// 检查验证码是否已使用
	if code.IsUsed {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码已使用"))
		return
	}

	// 标记为已使用
	if err := h.svcCtx.VerificationCodeModel.MarkAsUsed(codeId); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	log := &model.OperationLog{
		UserId:     currentUserId,
		Action:     "mark_verification_code_used",
		Resource:   "verification_code",
		ResourceId: codeId,
		Method:     "PUT",
		Path:       c.Request.URL.Path,
		Ip:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Status:     http.StatusOK,
	}
	h.svcCtx.OperationLogModel.Create(log)

	c.JSON(http.StatusOK, result.SimpleResult("标记已使用"))
}
