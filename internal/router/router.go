package router

import (
	"new-email/internal/handler"
	"new-email/internal/middleware"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine, svcCtx *svc.ServiceContext) {
	// 全局中间件
	r.Use(middleware.CorsMiddleware())

	// 静态文件服务
	r.Static("/static", "./web/dist/static")
	r.StaticFile("/", "./web/dist/index.html")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

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

	// 前端路由（SPA）
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}
