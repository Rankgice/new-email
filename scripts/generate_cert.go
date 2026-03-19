package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// 设置证书有效期
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	// 生成序列号
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("无法生成序列号: %v", err)
	}

	// 证书模板
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My Test Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 添加IP和DNS名称
	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
	template.DNSNames = append(template.DNSNames,
		"localhost",
		"email.host",
		"smtp.email.host",
		"imap.email.host",
		"mx.email.host",
		"email1.host",
		"smtp.email1.host",
		"imap.email1.host",
		"mx.email1.host",
	)

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("无法生成私钥: %v", err)
	}

	// 创建证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("无法创建证书: %v", err)
	}

	// 将证书写入cert.pem文件
	certOut, err := os.Create("data/tls/cert.pem")
	if err != nil {
		log.Fatalf("无法创建data/tls/cert.pem文件: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("无法写入证书数据: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("无法关闭data/tls/cert.pem文件: %v", err)
	}
	log.Println("✅ 成功生成 data/tls/cert.pem")

	// 将私钥写入key.pem文件
	keyOut, err := os.Create("data/tls/key.pem")
	if err != nil {
		log.Fatalf("无法创建data/tls/key.pem文件: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		log.Fatalf("无法写入私钥数据: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("无法关闭data/tls/key.pem文件: %v", err)
	}
	log.Println("✅ 成功生成 data/tls/key.pem")
}
