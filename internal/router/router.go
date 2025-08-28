package router

import (
	"github.com/rankgice/new-email/internal/handler"
	"github.com/rankgice/new-email/internal/middleware"
	"github.com/rankgice/new-email/internal/svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine, svcCtx *svc.ServiceContext) {
	// 全局中间件
	r.Use(middleware.CorsMiddleware())

	// 健康检查端点
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "邮件系统运行正常",
			"version": "1.0.0",
		})
	})

	// API状态页面（简单的HTML页面）
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>邮件管理系统</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .container { max-width: 600px; margin: 0 auto; }
        .status { color: #28a745; font-size: 24px; }
        .info { margin: 20px 0; }
        .endpoint { background: #f8f9fa; padding: 10px; margin: 10px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>邮件管理系统</h1>
        <div class="status">✅ 系统运行正常</div>
        
        <div class="info">
            <h3>邮件服务端口</h3>
            <div class="endpoint">SMTP接收: 25端口 (MTA)</div>
            <div class="endpoint">SMTP提交: 587端口 (MSA)</div>
            <div class="endpoint">IMAP访问: 993端口 (SSL)</div>
        </div>
        
        <div class="info">
            <h3>API接口</h3>
            <div class="endpoint">健康检查: <a href="/api/health">/api/health</a></div>
            <div class="endpoint">用户登录: POST /api/public/user/login</div>
            <div class="endpoint">管理员登录: POST /api/public/admin/login</div>
        </div>
        
        <div class="info">
            <h3>使用说明</h3>
            <p>这是一个纯API模式的邮件系统，请使用邮件客户端或API调用来操作。</p>
            <p>前端界面需要单独部署。</p>
        </div>
    </div>
</body>
</html>`)
	})

	// 创建handler实例
	userHandler := handler.NewUserHandler(svcCtx)
	adminHandler := handler.NewAdminHandler(svcCtx)
	commonHandler := handler.NewCommonHandler(svcCtx)
	mailboxHandler := handler.NewMailboxHandler(svcCtx)
	emailHandler := handler.NewEmailHandler(svcCtx)
	apiKeyHandler := handler.NewApiKeyHandler(svcCtx)
	domainHandler := handler.NewDomainHandler(svcCtx)
	apiHandler := handler.NewApiHandler(svcCtx)

	// API路由组
	api := r.Group("/api")
	{
		// 公共API（无需认证）
		public := api.Group("/public")
		{
			// 用户注册登录
			public.POST("/user/register", userHandler.Register)
			public.POST("/user/login", userHandler.Login)

			// 管理员登录
			public.POST("/admin/login", adminHandler.Login)

			// 验证码相关
			public.POST("/send-code", commonHandler.SendCode)
		}

		// 用户API（需要用户认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(svcCtx))
		{
			// 用户信息
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.POST("/change-password", userHandler.ChangePassword)

			// 邮箱管理
			mailbox := user.Group("/mailboxes")
			{
				mailbox.GET("", mailboxHandler.List)
				mailbox.POST("", mailboxHandler.Create)
				mailbox.PUT("/:id", mailboxHandler.Update)
				mailbox.DELETE("/:id", mailboxHandler.Delete)
				mailbox.GET("/stats", mailboxHandler.GetStats)
				mailbox.GET("/:id", mailboxHandler.GetById)
				mailbox.POST("/:id/sync", mailboxHandler.Sync)
			}

			// 邮件管理
			email := user.Group("/emails")
			{
				email.GET("", emailHandler.List)
				email.GET("/:id", emailHandler.GetById)
				email.POST("/send", emailHandler.Send)
				email.PUT("/:id/read", emailHandler.MarkRead)
				email.PUT("/:id/star", emailHandler.MarkStar)
				email.DELETE("/:id", emailHandler.Delete)
				email.POST("/batch", emailHandler.BatchOperation)
				email.GET("/export", emailHandler.Export)
				email.GET("/download/:filename", emailHandler.Download)
			}

			// API密钥管理
			apiKeys := user.Group("/api-keys")
			{
				apiKeys.GET("", apiKeyHandler.List)
				apiKeys.POST("", apiKeyHandler.Create)
				apiKeys.PUT("/:id", apiKeyHandler.Update)
				apiKeys.DELETE("/:id", apiKeyHandler.Delete)
			}
		}

		// 管理员API（需要管理员认证）
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware(svcCtx))
		{
			// 管理员信息
			admin.GET("/profile", adminHandler.GetProfile)
			admin.PUT("/profile", adminHandler.UpdateProfile)
			admin.POST("/change-password", adminHandler.ChangePassword)

			// 系统管理
			system := admin.Group("/system")
			{
				system.GET("/settings", adminHandler.GetSystemSettings)
				system.PUT("/settings", adminHandler.UpdateSystemSettings)
			}

			// 用户管理
			users := admin.Group("/users")
			{
				users.GET("", adminHandler.List)
				users.GET("/:id", adminHandler.GetById)
				users.POST("", adminHandler.Create)
				users.PUT("/:id", adminHandler.Update)
				users.DELETE("/:id", adminHandler.Delete)
				users.POST("/batch", adminHandler.BatchOperationUsers)
				users.POST("/import", adminHandler.ImportUsers)
				users.GET("/export", adminHandler.ExportUsers)
			}

			// 管理员管理（仅超级管理员）
			admins := admin.Group("/admins")
			admins.Use(middleware.SuperAdminMiddleware())
			{
				admins.GET("", adminHandler.ListAdmins)
				admins.POST("", adminHandler.CreateAdmin)
				admins.PUT("/:id", adminHandler.UpdateAdmin)
				admins.DELETE("/:id", adminHandler.DeleteAdmin)
				admins.POST("/batch", adminHandler.BatchOperationAdmins)
			}

			// 域名管理
			domains := admin.Group("/domains")
			{
				domains.GET("", domainHandler.List)
				domains.POST("", domainHandler.Create)
				domains.PUT("/:id", domainHandler.Update)
				domains.DELETE("/:id", domainHandler.Delete)
				domains.POST("/:id/verify", domainHandler.Verify)
				domains.POST("/batch", domainHandler.BatchOperation)
			}

			// 系统设置
			settings := admin.Group("/settings")
			{
				settings.GET("", adminHandler.GetSystemSettings)
				settings.PUT("", adminHandler.UpdateSystemSettings)
			}
		}

		// API密钥访问接口
		apiAccess := api.Group("/v1")
		apiAccess.Use(middleware.ApiKeyMiddleware(svcCtx))
		{
			// 邮件API
			apiAccess.GET("/emails/:id", apiHandler.GetEmail)
			apiAccess.POST("/emails/send", apiHandler.SendEmail)
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "请求的资源不存在",
			"path":    c.Request.URL.Path,
		})
	})
}
