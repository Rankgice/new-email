package handler

import (
	"crypto/rand"
	"fmt"
	"github.com/rankgice/new-email/internal/result"
	"github.com/rankgice/new-email/internal/svc"
	"github.com/rankgice/new-email/internal/types"
	"math/big"
	"net/http"
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

	// 生成验证码
	code, err := h.generateVerificationCode(req.Length)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成验证码失败"))
		return
	}

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(req.ExpireMinutes) * time.Minute)

	// 发送验证码
	if req.Type == "email" {
		// 发送邮件验证码
		if err := h.svcCtx.ServiceManager.SendVerificationEmail(req.Target, code); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("邮件发送失败: "+err.Error()))
			return
		}
	} else if req.Type == "sms" {
		// 发送短信验证码
		if _, err := h.svcCtx.ServiceManager.SendVerificationSMS(req.Target, code); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("短信发送失败: "+err.Error()))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不支持的验证码类型"))
		return
	}

	// 返回响应（不包含实际验证码）
	resp := types.SendCodeResp{
		Success:   true,
		Message:   "验证码发送成功",
		ExpiresAt: expiresAt,
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
	// 生成验证码ID
	captchaId := fmt.Sprintf("captcha_%d", time.Now().UnixNano())

	// 生成4位数字验证码
	_, err := h.generateVerificationCode(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成验证码失败"))
		return
	}

	// 这里应该生成图形验证码图片，暂时返回简单的文本验证码
	// 在实际项目中，可以使用 github.com/mojocn/base64Captcha 等库
	// 生成base64编码的图片验证码

	// 模拟base64图片数据
	imageData := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

	// 将验证码存储到缓存中（这里简化处理，实际应该存储到Redis等缓存中）
	// 暂时存储到内存中，过期时间5分钟
	expiresAt := time.Now().Add(5 * time.Minute)

	// 返回验证码信息
	resp := types.CaptchaResp{
		CaptchaId: captchaId,
		ImageData: imageData,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// VerifyCaptcha 验证图形验证码
func (h *CommonHandler) VerifyCaptcha(c *gin.Context) {
	var req types.CaptchaVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorBindingParam.AddError(err))
		return
	}

	// 这里应该从缓存中获取验证码进行验证
	// 暂时简化处理，只要验证码不为空就认为验证成功
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码不能为空"))
		return
	}

	if req.CaptchaId == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码ID不能为空"))
		return
	}

	// 模拟验证逻辑
	// 实际应该从缓存中获取对应的验证码进行比较
	success := len(req.Code) == 4 // 简单验证：4位数字

	if !success {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码错误"))
		return
	}

	// 验证成功后应该删除缓存中的验证码
	resp := types.CaptchaVerifyResp{
		Success: true,
		Message: "验证码验证成功",
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}

// Upload 文件上传
func (h *CommonHandler) Upload(c *gin.Context) {
	// 获取上传类型
	uploadType := c.PostForm("type")
	if uploadType == "" {
		uploadType = "attachment" // 默认为附件
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("获取上传文件失败"))
		return
	}
	defer file.Close()

	// 检查文件大小（10MB限制）
	maxSize := int64(10 * 1024 * 1024)
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("文件大小超过限制"))
		return
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		"image/jpeg":         true,
		"image/png":          true,
		"image/gif":          true,
		"image/webp":         true,
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不支持的文件类型"))
		return
	}

	// 使用存储服务上传文件
	if h.svcCtx.ServiceManager.Storage == nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("存储服务未配置"))
		return
	}

	fileInfo, err := h.svcCtx.ServiceManager.Storage.UploadFile(file, header, uploadType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("文件上传失败: "+err.Error()))
		return
	}

	resp := types.UploadResp{
		Url:      fileInfo.URL,
		Filename: fileInfo.Filename,
		Size:     fileInfo.Size,
		Type:     uploadType,
	}

	c.JSON(http.StatusOK, result.SuccessResult(resp))
}
