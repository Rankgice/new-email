package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rankgice/new-email/internal/config"
	"github.com/rankgice/new-email/internal/mailserver"
	"github.com/rankgice/new-email/internal/router"
	"github.com/rankgice/new-email/internal/svc"

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

	// 启动邮件服务器（使用emersion/go-smtp）
	mailServerConfig := mailserver.Config{
		SMTPReceivePort: c.SMTP.ReceivePort, // 从配置中获取SMTP接收端口
		SMTPSubmitPort:  c.SMTP.Port,        // 从配置中获取SMTP提交端口 (这里假设使用同一个端口，如果需要区分，需要修改config.go)
		IMAPPort:        c.IMAP.Port,        // 从配置中获取IMAP端口
		Domain:          "email.host",       // 从配置中获取主域名
		SMTPUseTLS:      c.SMTP.UseTLS,
		SMTPTLSCertPath: c.SMTP.TLSCertPath,
		SMTPTLSKeyPath:  c.SMTP.TLSKeyPath,
		IMAPUseTLS:      c.IMAP.UseTLS,
		IMAPTLSCertPath: c.IMAP.TLSCertPath,
		IMAPTLSKeyPath:  c.IMAP.TLSKeyPath,
	}
	mailServer := mailserver.NewMailServer(mailServerConfig, svcCtx.DB)
	if err := mailServer.Start(); err != nil {
		log.Fatal("邮件服务器启动失败：", err)
	}

	// 测试所有连接
	results := svcCtx.ServiceManager.TestAllConnections()
	for serviceName, err := range results {
		if err != nil {
			log.Printf("服务 %s 连接失败: %v", serviceName, err)
		} else {
			log.Printf("服务 %s 连接成功", serviceName)
		}
	}

	log.Printf("🚀 邮件管理系统启动成功")
	log.Printf("📱 管理端: http://localhost:%d/admin", port)
	log.Printf("👤 用户端: http://localhost:%d/user", port)
	log.Printf("🔗 API文档: http://localhost:%d/api/docs", port)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中启动Web服务器
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatal("Web服务器启动失败：", err)
		}
	}()

	// 等待退出信号
	<-sigChan
	log.Println("🛑 收到退出信号，正在关闭服务器...")

	// 停止邮件服务器
	if err := mailServer.Stop(); err != nil {
		log.Printf("停止邮件服务器失败: %v", err)
	}

	log.Println("✅ 服务器已安全关闭")
}
