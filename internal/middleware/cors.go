package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	// 允许所有来源（生产环境应该限制具体域名）
	config.AllowAllOrigins = true

	// 允许的请求头
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Connection",
		"Host",
		"Referer",
		"User-Agent",
		"X-API-Key",
	}

	// 允许的请求方法
	config.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
		"OPTIONS",
	}

	// 允许携带凭证
	config.AllowCredentials = true

	// 暴露的响应头
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"Authorization",
	}

	return cors.New(config)
}
