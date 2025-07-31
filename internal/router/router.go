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
	r.Use(middleware.LogMiddleware())

	// 静态文件服务
	r.Static("/static", "./web/dist/static")
	r.StaticFile("/", "./web/dist/index.html")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// 创建handler实例
	healthHandler := handler.NewHealthHandler(svcCtx)
	userHandler := handler.NewUserHandler(svcCtx)
	adminHandler := handler.NewAdminHandler(svcCtx)
	commonHandler := handler.NewCommonHandler(svcCtx)
	mailboxHandler := handler.NewMailboxHandler(svcCtx)
	emailHandler := handler.NewEmailHandler(svcCtx)
	draftHandler := handler.NewDraftHandler(svcCtx)
	templateHandler := handler.NewTemplateHandler(svcCtx)
	signatureHandler := handler.NewSignatureHandler(svcCtx)
	ruleHandler := handler.NewRuleHandler(svcCtx)
	verificationCodeHandler := handler.NewVerificationCodeHandler(svcCtx)
	apiKeyHandler := handler.NewApiKeyHandler(svcCtx)
	logHandler := handler.NewLogHandler(svcCtx)
	domainHandler := handler.NewDomainHandler(svcCtx)
	adminRuleHandler := handler.NewAdminRuleHandler(svcCtx)
	adminLogHandler := handler.NewAdminLogHandler(svcCtx)
	apiHandler := handler.NewApiHandler(svcCtx)

	// API路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", healthHandler.Health)

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
			public.POST("/verify-code", commonHandler.VerifyCode)
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
			}

			// 草稿管理
			draft := user.Group("/drafts")
			{
				draft.GET("", draftHandler.List)
				draft.POST("", draftHandler.Create)
				draft.PUT("/:id", draftHandler.Update)
				draft.DELETE("/:id", draftHandler.Delete)
				draft.POST("/:id/send", draftHandler.Send)
			}

			// 邮件模板
			template := user.Group("/templates")
			{
				template.GET("", templateHandler.List)
				template.POST("", templateHandler.Create)
				template.PUT("/:id", templateHandler.Update)
				template.DELETE("/:id", templateHandler.Delete)
			}

			// 邮件签名
			signature := user.Group("/signatures")
			{
				signature.GET("", signatureHandler.List)
				signature.POST("", signatureHandler.Create)
				signature.PUT("/:id", signatureHandler.Update)
				signature.DELETE("/:id", signatureHandler.Delete)
			}

			// 规则管理
			rules := user.Group("/rules")
			{
				// 验证码规则
				verification := rules.Group("/verification")
				{
					verification.GET("", ruleHandler.ListVerificationRules)
					verification.POST("", ruleHandler.CreateVerificationRule)
					verification.PUT("/:id", ruleHandler.UpdateVerificationRule)
					verification.DELETE("/:id", ruleHandler.DeleteVerificationRule)
				}

				// 转发规则
				forward := rules.Group("/forward")
				{
					forward.GET("", ruleHandler.ListForwardRules)
					forward.POST("", ruleHandler.CreateForwardRule)
					forward.PUT("/:id", ruleHandler.UpdateForwardRule)
					forward.DELETE("/:id", ruleHandler.DeleteForwardRule)
				}
			}

			// 验证码记录
			codes := user.Group("/verification-codes")
			{
				codes.GET("", verificationCodeHandler.List)
				codes.GET("/:id", verificationCodeHandler.GetById)
				codes.PUT("/:id/used", verificationCodeHandler.MarkUsed)
			}

			// API密钥管理
			apiKeys := user.Group("/api-keys")
			{
				apiKeys.GET("", apiKeyHandler.List)
				apiKeys.POST("", apiKeyHandler.Create)
				apiKeys.PUT("/:id", apiKeyHandler.Update)
				apiKeys.DELETE("/:id", apiKeyHandler.Delete)
			}

			// 日志查询
			logs := user.Group("/logs")
			{
				logs.GET("/operation", logHandler.ListOperationLogs)
				logs.GET("/email", logHandler.ListEmailLogs)
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

			// 仪表板
			admin.GET("/dashboard", adminHandler.Dashboard)

			// 用户管理
			users := admin.Group("/users")
			{
				users.GET("", adminHandler.ListUsers)
				users.POST("", adminHandler.CreateUser)
				users.PUT("/:id", adminHandler.UpdateUser)
				users.DELETE("/:id", adminHandler.DeleteUser)
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

			// 全局规则管理
			globalRules := admin.Group("/global-rules")
			{
				// 全局验证码规则
				verification := globalRules.Group("/verification")
				{
					verification.GET("", adminRuleHandler.ListGlobalVerificationRules)
					verification.POST("", adminRuleHandler.CreateGlobalVerificationRule)
					verification.PUT("/:id", adminRuleHandler.UpdateGlobalVerificationRule)
					verification.DELETE("/:id", adminRuleHandler.DeleteGlobalVerificationRule)
				}

				// 全局反垃圾规则
				antiSpam := globalRules.Group("/anti-spam")
				{
					antiSpam.GET("", adminRuleHandler.ListAntiSpamRules)
					antiSpam.POST("", adminRuleHandler.CreateAntiSpamRule)
					antiSpam.PUT("/:id", adminRuleHandler.UpdateAntiSpamRule)
					antiSpam.DELETE("/:id", adminRuleHandler.DeleteAntiSpamRule)
				}
			}

			// 系统日志
			systemLogs := admin.Group("/logs")
			{
				systemLogs.GET("/operation", adminLogHandler.ListOperationLogs)
				systemLogs.GET("/email", adminLogHandler.ListEmailLogs)
				systemLogs.GET("/system", adminLogHandler.ListSystemLogs)
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
			apiAccess.GET("/emails", apiHandler.ListEmails)
			apiAccess.GET("/emails/:id", apiHandler.GetEmail)
			apiAccess.POST("/emails/send", apiHandler.SendEmail)

			// 验证码API
			apiAccess.GET("/verification-codes", apiHandler.ListVerificationCodes)
			apiAccess.GET("/verification-codes/:id", apiHandler.GetVerificationCode)
		}
	}

	// 前端路由（SPA）
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}
