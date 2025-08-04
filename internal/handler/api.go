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

// ApiHandler API处理器
type ApiHandler struct {
	svcCtx *svc.ServiceContext
}

// NewApiHandler 创建API处理器
func NewApiHandler(svcCtx *svc.ServiceContext) *ApiHandler {
	return &ApiHandler{
		svcCtx: svcCtx,
	}
}

// GetEmail API获取邮件
func (h *ApiHandler) GetEmail(c *gin.Context) {
	// 获取邮件ID
	idStr := c.Param("id")
	emailId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的邮件ID"))
		return
	}

	// 通过API密钥验证获取用户ID
	// TODO: 实现GetApiUserId方法
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
		return
	}

	// 查询邮件详情
	email, err := h.svcCtx.EmailModel.GetById(emailId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 检查权限（只能查看自己的邮件）
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}
	if mailbox == nil || mailbox.UserId != userId {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限查看此邮件"))
		return
	}

	// 返回邮件详情
	resp := types.EmailResp{
		Id:          email.Id,
		MailboxId:   email.MailboxId,
		Subject:     email.Subject,
		FromEmail:   email.FromEmail,
		ToEmail:     email.ToEmails,  // 使用ToEmails字段
		CcEmail:     email.CcEmails,  // 使用CcEmails字段
		BccEmail:    email.BccEmails, // 使用BccEmails字段
		Content:     email.Content,
		ContentType: email.ContentType,
		Attachments: "", // Email模型中没有Attachments字段
		Status:      0,  // Email模型中没有Status字段
		Type:        "", // Email模型中没有Type字段
		CreatedAt:   email.CreatedAt,
		UpdatedAt:   email.UpdatedAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// SendEmail API发送邮件
func (h *ApiHandler) SendEmail(c *gin.Context) {
	var req types.EmailSendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 通过API密钥验证获取用户ID
	// TODO: 实现GetApiUserId方法
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
		return
	}

	// 检查邮箱是否属于当前用户
	if req.MailboxId > 0 {
		mailbox, err := h.svcCtx.MailboxModel.GetById(req.MailboxId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			return
		}
		if mailbox == nil || mailbox.UserId != userId {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限使用此邮箱"))
			return
		}
	}

	// TODO: 实现实际的邮件发送逻辑
	// 这里应该：
	// 1. 验证邮件内容的完整性
	// 2. 调用邮件发送服务
	// 3. 创建邮件记录
	// 4. 返回发送结果

	// 模拟发送成功
	sendResp := types.EmailSendResp{
		Success: true,
		Message: "邮件发送成功",
		EmailId: 0, // 实际发送后应该返回邮件ID
		SentAt:  time.Now(),
	}

	c.JSON(http.StatusOK, result.SuccessResult(sendResp))
}

// ListVerificationCodes API验证码列表
func (h *ApiHandler) ListVerificationCodes(c *gin.Context) {
	var req types.VerificationCodeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 通过API密钥验证获取用户ID
	// TODO: 实现GetApiUserId方法
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
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

// GetVerificationCode API获取验证码
func (h *ApiHandler) GetVerificationCode(c *gin.Context) {
	// 获取验证码ID
	idStr := c.Param("id")
	codeId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	// 通过API密钥验证获取用户ID
	// TODO: 实现GetApiUserId方法
	userId := middleware.GetCurrentUserId(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("无效的API密钥"))
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

	// 检查权限（只能查看自己的验证码）
	// TODO: VerificationCode模型中没有UserId字段，需要通过EmailId关联查询
	// if code.UserId != userId {
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
