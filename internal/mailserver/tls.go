package mailserver

import (
	"crypto/tls"
	"log"
)

func loadOptionalTLSConfig(serverName string, useTLS bool, certPath, keyPath string) (*tls.Config, bool) {
	if !useTLS {
		return nil, false
	}

	if certPath == "" || keyPath == "" {
		log.Printf("⚠️ %s已启用TLS，但证书路径未完整配置，将以非TLS模式启动", serverName)
		return nil, false
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Printf("❌ 无法加载%s的TLS证书: %v。%s将以非TLS模式启动", serverName, err, serverName)
		return nil, false
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, true
}
