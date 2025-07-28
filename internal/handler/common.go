package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"new-email/internal/middleware"
	"new-email/internal/model"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/internal/types"
	"time"

	"github.com/gin-gonic/gin"
)

// CommonHandler 通用处理器
type CommonHandler struct {
	svcCtx *svc.ServiceContext
}

// NewCommonHandler 创建通用处理器
func NewCommonHandler(svcCtx *svc.ServiceContext) *CommonHandler {
	return &CommonHandler{
		svcCtx: svcCtx,
	}
}

// SendCode 发送验证码
func (h *CommonHandler) SendCode(c *gin.Context) {
	var req types.SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID（可选，某些场景下可能不需要登录）
	currentUserId := middleware.GetCurrentUserId(c)

	// 生成验证码
	code, err := h.generateVerificationCode(req.Length)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成验证码失败"))
		return
	}

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(req.ExpireMinutes) * time.Minute)

	// 保存验证码记录
	verificationCode := &model.VerificationCode{
		UserId:    currentUserId,
		Code:      code,
		Type:      req.Type,
		Target:    req.Target, // 邮箱或手机号
		IsUsed:    false,
		IsExpired: false,
		ExpiresAt: expiresAt,
	}

	if err := h.svcCtx.VerificationCodeModel.Create(verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorAdd.AddError(err))
		return
	}

	// TODO: 实际发送验证码（邮件或短信）
	// 这里应该根据req.Type调用相应的发送服务
	// 例如：
	// if req.Type == "email" {
	//     err = h.sendEmailCode(req.Target, code)
	// } else if req.Type == "sms" {
	//     err = h.sendSmsCode(req.Target, code)
	// }

	// 记录操作日志
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "send_verification_code",
			Resource:   "verification_code",
			ResourceId: verificationCode.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回响应（不包含实际验证码）
	resp := types.SendCodeResp{
		Success:   true,
		Message:   "验证码发送成功",
		ExpiresAt: expiresAt,
		CodeId:    verificationCode.Id,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// VerifyCode 验证验证码
func (h *CommonHandler) VerifyCode(c *gin.Context) {
	var req types.VerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 获取当前用户ID（可选）
	currentUserId := middleware.GetCurrentUserId(c)

	// 查询验证码记录
	var verificationCode *model.VerificationCode
	var err error

	if req.CodeId > 0 {
		// 通过ID查询
		verificationCode, err = h.svcCtx.VerificationCodeModel.GetById(req.CodeId)
	} else {
		// 通过目标和类型查询最新的验证码
		verificationCode, err = h.svcCtx.VerificationCodeModel.GetLatestByTargetAndType(req.Target, req.Type)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
		return
	}

	if verificationCode == nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	// 检查验证码是否已使用
	if verificationCode.IsUsed {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码已使用"))
		return
	}

	// 检查验证码是否已过期
	if time.Now().After(verificationCode.ExpiresAt) {
		// 标记为过期
		h.svcCtx.VerificationCodeModel.MarkAsExpired(verificationCode.Id)
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码已过期"))
		return
	}

	// 验证验证码
	if verificationCode.Code != req.Code {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码错误"))
		return
	}

	// 标记验证码为已使用
	if err := h.svcCtx.VerificationCodeModel.MarkAsUsed(verificationCode.Id); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorUpdate.AddError(err))
		return
	}

	// 记录操作日志
	if currentUserId > 0 {
		log := &model.OperationLog{
			UserId:     currentUserId,
			Action:     "verify_code",
			Resource:   "verification_code",
			ResourceId: verificationCode.Id,
			Method:     "POST",
			Path:       c.Request.URL.Path,
			Ip:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     http.StatusOK,
		}
		h.svcCtx.OperationLogModel.Create(log)
	}

	// 返回验证成功响应
	resp := types.VerifyCodeResp{
		Success: true,
		Message: "验证码验证成功",
		CodeId:  verificationCode.Id,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// generateVerificationCode 生成验证码
func (h *CommonHandler) generateVerificationCode(length int) (string, error) {
	if length <= 0 {
		length = 6 // 默认6位
	}

	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

// GetCaptcha 获取图形验证码
func (h *CommonHandler) GetCaptcha(c *gin.Context) {
	// TODO: 实现图形验证码生成
	// 这里可以使用第三方库如 github.com/mojocn/base64Captcha
	c.JSON(http.StatusOK, result.SimpleResult("获取图形验证码接口"))
}

// VerifyCaptcha 验证图形验证码
func (h *CommonHandler) VerifyCaptcha(c *gin.Context) {
	// TODO: 实现图形验证码验证
	c.JSON(http.StatusOK, result.SimpleResult("验证图形验证码接口"))
}

// Upload 文件上传
func (h *CommonHandler) Upload(c *gin.Context) {
	// TODO: 实现文件上传功能
	c.JSON(http.StatusOK, result.SimpleResult("文件上传接口"))
}

// GetSystemInfo 获取系统信息
func (h *CommonHandler) GetSystemInfo(c *gin.Context) {
	info := types.SystemInfoResp{
		Version:     h.svcCtx.Config.App.Version,
		Environment: h.svcCtx.Config.Web.Mode,
		ServerTime:  time.Now(),
		Timezone:    "Asia/Shanghai",
	}

	c.JSON(http.StatusOK, result.SuccessResult(info))
}
