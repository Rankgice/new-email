package middleware

import (
	"errors"
	"net/http"
	"new-email/internal/constant"
	"new-email/internal/result"
	"new-email/internal/svc"
	"new-email/pkg/auth"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 用户认证中间件
func AuthMiddleware(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		newToken, claims, statusCode := validateAndRefreshToken(c, svcCtx, false)
		if statusCode != http.StatusOK {
			c.Abort()
			return
		}

		// 检查用户是否存在且状态正常
		user, err := svcCtx.UserModel.GetById(claims.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			c.Abort()
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, result.ErrorUserNotFound)
			c.Abort()
			return
		}

		if user.Status != constant.StatusEnabled {
			c.JSON(http.StatusUnauthorized, result.ErrorUserDisabled)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("userType", "user")
		c.Set("user", user)

		// 如果有新token，设置到响应头
		if newToken != "" {
			c.Header("New-Token", newToken)
		}

		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		newToken, claims, statusCode := validateAndRefreshToken(c, svcCtx, true)
		if statusCode != http.StatusOK {
			c.Abort()
			return
		}

		// 检查管理员是否存在且状态正常
		admin, err := svcCtx.AdminModel.GetById(claims.UserId) // 使用UserId，因为IsAdmin=true时代表adminId
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			c.Abort()
			return
		}

		if admin == nil {
			c.JSON(http.StatusUnauthorized, result.ErrorUserNotFound)
			c.Abort()
			return
		}

		if admin.Status != constant.StatusEnabled {
			c.JSON(http.StatusUnauthorized, result.ErrorUserDisabled)
			c.Abort()
			return
		}

		// 将管理员信息存储到上下文
		c.Set("adminId", claims.UserId) // 使用UserId，因为IsAdmin=true时代表adminId
		c.Set("userId", claims.UserId)  // 为了兼容性也设置userId
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("userType", "admin")
		c.Set("admin", admin)

		// 如果有新token，设置到响应头
		if newToken != "" {
			c.Header("New-Token", newToken)
		}

		c.Next()
	}
}

// validateAndRefreshToken 统一的token验证和刷新逻辑
// 返回值：newToken(如果刷新了), claims, statusCode
func validateAndRefreshToken(c *gin.Context, svcCtx *svc.ServiceContext, requireAdmin bool) (string, *auth.Claims, int) {
	// 获取Authorization头
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
		return "", nil, http.StatusUnauthorized
	}

	// 检查Bearer前缀
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid)
		return "", nil, http.StatusUnauthorized
	}

	// 提取token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid)
		return "", nil, http.StatusUnauthorized
	}

	// 解析token
	claims, err := auth.ParseToken(token, svcCtx.Config.JWT.Secret)
	if err != nil {
		// 检查是否是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) && claims != nil {
			// Token过期但有效，检查是否在可刷新时间内（比如过期后24小时内）
			refreshWindow := time.Duration(svcCtx.Config.JWT.RefreshExpireHours) * time.Hour
			if time.Since(claims.ExpiresAt.Time) <= refreshWindow {
				// 生成新token
				newToken, genErr := auth.GenerateToken(
					claims.UserId,
					claims.Username,
					claims.IsAdmin,
					claims.Role,
					svcCtx.Config.JWT.Secret,
					svcCtx.Config.JWT.ExpireHours,
				)
				if genErr != nil {
					c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid.AddError(genErr))
					return "", nil, http.StatusUnauthorized
				}

				// 检查用户类型是否匹配
				if claims.IsAdmin != requireAdmin {
					c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid)
					return "", nil, http.StatusUnauthorized
				}

				// 返回状态码401表示token已刷新
				c.JSON(401, result.Result{
					Code: 401,
					Msg:  "Token refreshed",
					Data: map[string]string{"newToken": newToken},
				})
				return newToken, claims, 401
			}
		}
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid.AddError(err))
		return "", nil, http.StatusUnauthorized
	}

	// 检查用户类型是否匹配
	if claims.IsAdmin != requireAdmin {
		c.JSON(http.StatusUnauthorized, result.ErrorTokenInvalid)
		return "", nil, http.StatusUnauthorized
	}

	return "", claims, http.StatusOK
}

// SuperAdminMiddleware 超级管理员权限中间件
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, result.ErrorUnauthorized)
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, result.ErrorPermissionDeny)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ApiKeyMiddleware API密钥认证中间件
func ApiKeyMiddleware(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取API密钥
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, result.ErrorAPIKeyInvalid)
			c.Abort()
			return
		}

		// 验证API密钥
		key, err := svcCtx.ApiKeyModel.GetByKey(apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSelect.AddError(err))
			c.Abort()
			return
		}

		if key == nil {
			c.JSON(http.StatusUnauthorized, result.ErrorAPIKeyInvalid)
			c.Abort()
			return
		}

		if key.Status != 1 {
			c.JSON(http.StatusUnauthorized, result.ErrorAPIKeyDisabled)
			c.Abort()
			return
		}

		// 检查是否过期
		if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, result.ErrorAPIKeyExpired)
			c.Abort()
			return
		}

		// 更新最后使用时间
		go func() {
			svcCtx.ApiKeyModel.UpdateLastUsed(key.Id)
		}()

		// 将API密钥信息存储到上下文
		c.Set("apiKeyId", key.Id)
		c.Set("userId", key.UserId)
		c.Set("permissions", key.Permissions)
		c.Set("userType", "api")

		c.Next()
	}
}

// GetCurrentUserId 获取当前用户ID
func GetCurrentUserId(c *gin.Context) int64 {
	if userId, exists := c.Get("userId"); exists {
		if id, ok := userId.(int64); ok {
			return id
		}
		// 尝试从字符串转换
		if idStr, ok := userId.(string); ok {
			if id, err := strconv.ParseInt(idStr, 10, 32); err == nil {
				return id
			}
		}
	}
	return 0
}

// GetCurrentAdminId 获取当前管理员ID
func GetCurrentAdminId(c *gin.Context) int64 {
	if adminId, exists := c.Get("adminId"); exists {
		if id, ok := adminId.(int64); ok {
			return id
		}
		// 尝试从字符串转换
		if idStr, ok := adminId.(string); ok {
			if id, err := strconv.ParseInt(idStr, 10, 32); err == nil {
				return id
			}
		}
	}
	return 0
}

// GetCurrentUserType 获取当前用户类型
func GetCurrentUserType(c *gin.Context) string {
	if userType, exists := c.Get("userType"); exists {
		if t, ok := userType.(string); ok {
			return t
		}
	}
	return ""
}

// IsAdmin 检查是否为管理员
func IsAdmin(c *gin.Context) bool {
	return GetCurrentUserType(c) == "admin"
}

// IsSuperAdmin 检查是否为超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	if !IsAdmin(c) {
		return false
	}
	if role, exists := c.Get("role"); exists {
		return role == "admin"
	}
	return false
}
