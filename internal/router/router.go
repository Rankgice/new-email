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

	// API路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", handler.NewHealthHandler(svcCtx).Health)

		// 公共API（无需认证）
		public := api.Group("/public")
		{
			// 用户注册登录
			public.POST("/user/register", handler.NewUserHandler(svcCtx).Register)
			public.POST("/user/login", handler.NewUserHandler(svcCtx).Login)
			public.POST("/user/refresh", handler.NewUserHandler(svcCtx).RefreshToken)

			// 管理员登录
			public.POST("/admin/login", handler.NewAdminHandler(svcCtx).Login)
			public.POST("/admin/refresh", handler.NewAdminHandler(svcCtx).RefreshToken)

			// 验证码相关
			public.POST("/send-code", handler.NewCommonHandler(svcCtx).SendCode)
			public.POST("/verify-code", handler.NewCommonHandler(svcCtx).VerifyCode)
		}

		// 用户API（需要用户认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(svcCtx))
		{
			// 用户信息
			user.GET("/profile", handler.NewUserHandler(svcCtx).GetProfile)
			user.PUT("/profile", handler.NewUserHandler(svcCtx).UpdateProfile)
			user.POST("/change-password", handler.NewUserHandler(svcCtx).ChangePassword)

			// 邮箱管理
			mailbox := user.Group("/mailboxes")
			{
				mailbox.GET("", handler.NewMailboxHandler(svcCtx).List)
				mailbox.POST("", handler.NewMailboxHandler(svcCtx).Create)
				mailbox.PUT("/:id", handler.NewMailboxHandler(svcCtx).Update)
				mailbox.DELETE("/:id", handler.NewMailboxHandler(svcCtx).Delete)
				mailbox.POST("/:id/sync", handler.NewMailboxHandler(svcCtx).Sync)
				mailbox.POST("/:id/test", handler.NewMailboxHandler(svcCtx).TestConnection)
			}

			// 邮件管理
			email := user.Group("/emails")
			{
				email.GET("", handler.NewEmailHandler(svcCtx).List)
				email.GET("/:id", handler.NewEmailHandler(svcCtx).GetById)
				email.POST("/send", handler.NewEmailHandler(svcCtx).Send)
				email.PUT("/:id/read", handler.NewEmailHandler(svcCtx).MarkRead)
				email.PUT("/:id/star", handler.NewEmailHandler(svcCtx).MarkStar)
				email.DELETE("/:id", handler.NewEmailHandler(svcCtx).Delete)
				email.POST("/batch", handler.NewEmailHandler(svcCtx).BatchOperation)
			}

			// 草稿管理
			draft := user.Group("/drafts")
			{
				draft.GET("", handler.NewDraftHandler(svcCtx).List)
				draft.POST("", handler.NewDraftHandler(svcCtx).Create)
				draft.PUT("/:id", handler.NewDraftHandler(svcCtx).Update)
				draft.DELETE("/:id", handler.NewDraftHandler(svcCtx).Delete)
				draft.POST("/:id/send", handler.NewDraftHandler(svcCtx).Send)
			}

			// 邮件模板
			template := user.Group("/templates")
			{
				template.GET("", handler.NewTemplateHandler(svcCtx).List)
				template.POST("", handler.NewTemplateHandler(svcCtx).Create)
				template.PUT("/:id", handler.NewTemplateHandler(svcCtx).Update)
				template.DELETE("/:id", handler.NewTemplateHandler(svcCtx).Delete)
			}

			// 邮件签名
			signature := user.Group("/signatures")
			{
				signature.GET("", handler.NewSignatureHandler(svcCtx).List)
				signature.POST("", handler.NewSignatureHandler(svcCtx).Create)
				signature.PUT("/:id", handler.NewSignatureHandler(svcCtx).Update)
				signature.DELETE("/:id", handler.NewSignatureHandler(svcCtx).Delete)
			}

			// 规则管理
			rules := user.Group("/rules")
			{
				// 验证码规则
				verification := rules.Group("/verification")
				{
					verification.GET("", handler.NewRuleHandler(svcCtx).ListVerificationRules)
					verification.POST("", handler.NewRuleHandler(svcCtx).CreateVerificationRule)
					verification.PUT("/:id", handler.NewRuleHandler(svcCtx).UpdateVerificationRule)
					verification.DELETE("/:id", handler.NewRuleHandler(svcCtx).DeleteVerificationRule)
				}

				// 转发规则
				forward := rules.Group("/forward")
				{
					forward.GET("", handler.NewRuleHandler(svcCtx).ListForwardRules)
					forward.POST("", handler.NewRuleHandler(svcCtx).CreateForwardRule)
					forward.PUT("/:id", handler.NewRuleHandler(svcCtx).UpdateForwardRule)
					forward.DELETE("/:id", handler.NewRuleHandler(svcCtx).DeleteForwardRule)
				}
			}

			// 验证码记录
			codes := user.Group("/verification-codes")
			{
				codes.GET("", handler.NewVerificationCodeHandler(svcCtx).List)
				codes.GET("/:id", handler.NewVerificationCodeHandler(svcCtx).GetById)
				codes.PUT("/:id/used", handler.NewVerificationCodeHandler(svcCtx).MarkUsed)
			}

			// API密钥管理
			apiKeys := user.Group("/api-keys")
			{
				apiKeys.GET("", handler.NewApiKeyHandler(svcCtx).List)
				apiKeys.POST("", handler.NewApiKeyHandler(svcCtx).Create)
				apiKeys.PUT("/:id", handler.NewApiKeyHandler(svcCtx).Update)
				apiKeys.DELETE("/:id", handler.NewApiKeyHandler(svcCtx).Delete)
			}

			// 日志查询
			logs := user.Group("/logs")
			{
				logs.GET("/operation", handler.NewLogHandler(svcCtx).ListOperationLogs)
				logs.GET("/email", handler.NewLogHandler(svcCtx).ListEmailLogs)
			}
		}

		// 管理员API（需要管理员认证）
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware(svcCtx))
		{
			// 管理员信息
			admin.GET("/profile", handler.NewAdminHandler(svcCtx).GetProfile)
			admin.PUT("/profile", handler.NewAdminHandler(svcCtx).UpdateProfile)
			admin.POST("/change-password", handler.NewAdminHandler(svcCtx).ChangePassword)

			// 仪表板
			admin.GET("/dashboard", handler.NewAdminHandler(svcCtx).Dashboard)

			// 用户管理
			users := admin.Group("/users")
			{
				users.GET("", handler.NewAdminHandler(svcCtx).ListUsers)
				users.POST("", handler.NewAdminHandler(svcCtx).CreateUser)
				users.PUT("/:id", handler.NewAdminHandler(svcCtx).UpdateUser)
				users.DELETE("/:id", handler.NewAdminHandler(svcCtx).DeleteUser)
				users.POST("/batch", handler.NewAdminHandler(svcCtx).BatchOperationUsers)
				users.POST("/import", handler.NewAdminHandler(svcCtx).ImportUsers)
				users.GET("/export", handler.NewAdminHandler(svcCtx).ExportUsers)
			}

			// 管理员管理（仅超级管理员）
			admins := admin.Group("/admins")
			admins.Use(middleware.SuperAdminMiddleware())
			{
				admins.GET("", handler.NewAdminHandler(svcCtx).ListAdmins)
				admins.POST("", handler.NewAdminHandler(svcCtx).CreateAdmin)
				admins.PUT("/:id", handler.NewAdminHandler(svcCtx).UpdateAdmin)
				admins.DELETE("/:id", handler.NewAdminHandler(svcCtx).DeleteAdmin)
				admins.POST("/batch", handler.NewAdminHandler(svcCtx).BatchOperationAdmins)
			}

			// 域名管理
			domains := admin.Group("/domains")
			{
				domains.GET("", handler.NewDomainHandler(svcCtx).List)
				domains.POST("", handler.NewDomainHandler(svcCtx).Create)
				domains.PUT("/:id", handler.NewDomainHandler(svcCtx).Update)
				domains.DELETE("/:id", handler.NewDomainHandler(svcCtx).Delete)
				domains.POST("/:id/verify", handler.NewDomainHandler(svcCtx).Verify)
				domains.POST("/batch", handler.NewDomainHandler(svcCtx).BatchOperation)
			}

			// 全局规则管理
			globalRules := admin.Group("/global-rules")
			{
				// 全局验证码规则
				verification := globalRules.Group("/verification")
				{
					verification.GET("", handler.NewAdminRuleHandler(svcCtx).ListGlobalVerificationRules)
					verification.POST("", handler.NewAdminRuleHandler(svcCtx).CreateGlobalVerificationRule)
					verification.PUT("/:id", handler.NewAdminRuleHandler(svcCtx).UpdateGlobalVerificationRule)
					verification.DELETE("/:id", handler.NewAdminRuleHandler(svcCtx).DeleteGlobalVerificationRule)
				}

				// 全局反垃圾规则
				antiSpam := globalRules.Group("/anti-spam")
				{
					antiSpam.GET("", handler.NewAdminRuleHandler(svcCtx).ListAntiSpamRules)
					antiSpam.POST("", handler.NewAdminRuleHandler(svcCtx).CreateAntiSpamRule)
					antiSpam.PUT("/:id", handler.NewAdminRuleHandler(svcCtx).UpdateAntiSpamRule)
					antiSpam.DELETE("/:id", handler.NewAdminRuleHandler(svcCtx).DeleteAntiSpamRule)
				}
			}

			// 系统日志
			systemLogs := admin.Group("/logs")
			{
				systemLogs.GET("/operation", handler.NewAdminLogHandler(svcCtx).ListOperationLogs)
				systemLogs.GET("/email", handler.NewAdminLogHandler(svcCtx).ListEmailLogs)
				systemLogs.GET("/system", handler.NewAdminLogHandler(svcCtx).ListSystemLogs)
			}

			// 系统设置
			settings := admin.Group("/settings")
			{
				settings.GET("", handler.NewAdminHandler(svcCtx).GetSystemSettings)
				settings.PUT("", handler.NewAdminHandler(svcCtx).UpdateSystemSettings)
			}
		}

		// API密钥访问接口
		apiAccess := api.Group("/v1")
		apiAccess.Use(middleware.ApiKeyMiddleware(svcCtx))
		{
			// 邮件API
			apiAccess.GET("/emails", handler.NewApiHandler(svcCtx).ListEmails)
			apiAccess.GET("/emails/:id", handler.NewApiHandler(svcCtx).GetEmail)
			apiAccess.POST("/emails/send", handler.NewApiHandler(svcCtx).SendEmail)

			// 验证码API
			apiAccess.GET("/verification-codes", handler.NewApiHandler(svcCtx).ListVerificationCodes)
			apiAccess.GET("/verification-codes/:id", handler.NewApiHandler(svcCtx).GetVerificationCode)
		}
	}

	// 前端路由（SPA）
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}
