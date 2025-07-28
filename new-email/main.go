package main

import (
	"flag"
	"fmt"
	"log"
	"new-email/internal/config"
	"new-email/internal/router"
	"new-email/internal/svc"

	"github.com/gin-gonic/gin"
)

var configFile = flag.String("f", "etc/config.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置文件
	c := config.NewConfig(*configFile)

	// 设置Gin模式
	if c.Web.Mode != "" {
		gin.SetMode(c.Web.Mode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(c)
	log.Println("✅ 服务上下文初始化成功")

	// 设置路由
	router.SetupRouter(r, svcCtx)
	log.Println("✅ 路由设置完成")

	// 确定端口
	port := c.Web.Port
	if port == 0 {
		port = 8081
	}

	log.Printf("🚀 邮件管理系统启动成功")
	log.Printf("📱 管理端: http://localhost:%d/admin", port)
	log.Printf("👤 用户端: http://localhost:%d/user", port)
	log.Printf("🔗 API文档: http://localhost:%d/api/docs", port)

	// 启动服务器
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("服务器启动失败：", err)
	}
}
