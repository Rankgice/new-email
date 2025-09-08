# Docker 邮件系统部署指南

## 概述

这是一个完整的Docker化邮件管理系统，包含邮件服务器功能（SMTP/IMAP）、Web管理界面、SQLite数据库和MinIO对象存储。

## 邮件系统架构

- **email-system**: 完整邮件服务器（SMTP + IMAP + Web管理 + SQLite）
- **minio**: MinIO 对象存储服务

### 邮件服务器功能

- **SMTP服务器**: 
  - 端口25 (MTA) - 接收来自外部邮件服务器的邮件
  - 端口587 (MSA) - 用户认证提交邮件
- **IMAP服务器**: 端口993 (SSL) - 邮件客户端访问
- **Web界面**: 端口8080 - 管理界面和API

## 快速启动

### 1. 构建并启动所有服务

```bash
# 构建并启动
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f email-system
```

### 2. 访问应用

- **Web管理界面**: http://localhost:8080
- **邮件客户端配置**:
  - SMTP发送服务器: localhost:587 (需要认证)
  - IMAP接收服务器: localhost:993 (SSL)
- **MinIO控制台**: http://localhost:9001 (用户名/密码: minioadmin/minioadmin)

### 3. 默认账户

- **管理员**: admin / admin123456
- **测试邮箱**: test@email.host / test123

## 邮件客户端配置

### 常用邮件客户端设置

**发送邮件 (SMTP)**:
- 服务器: localhost (或服务器IP)
- 端口: 587
- 加密: STARTTLS或无
- 认证: 是

**接收邮件 (IMAP)**:
- 服务器: localhost (或服务器IP)  
- 端口: 993
- 加密: SSL/TLS
- 认证: 是

### Thunderbird配置示例

```
账户设置:
  邮箱地址: test@email.host
  密码: test123

服务器设置:
  接收邮件服务器:
    - 协议: IMAP
    - 服务器: localhost
    - 端口: 993
    - 连接安全性: SSL/TLS
    
  发送邮件服务器:
    - 协议: SMTP
    - 服务器: localhost
    - 端口: 587
    - 连接安全性: STARTTLS
```

## 端口映射详情

| 服务 | 内部端口 | 外部端口 | 协议 | 说明 |
|------|----------|----------|------|------|
| email-system | 8080 | 8080 | HTTP | Web管理界面和API |
| email-system | 25 | 25 | SMTP | 邮件接收 (MTA) |
| email-system | 587 | 587 | SMTP | 邮件提交 (MSA) |
| email-system | 993 | 993 | IMAP | 邮件访问 (SSL) |
| minio | 9000,9001 | 9000,9001 | HTTP | MinIO API和控制台 |

## 数据持久化

所有重要数据都通过Docker卷进行持久化：

- `minio_data`: MinIO对象存储数据
- `./data`: 应用数据目录（包含SQLite数据库文件）
  - `email.db`: SQLite数据库
  - `attachments/`: 邮件附件
  - `logs/`: 日志文件

## 健康检查

- **email-system**: `curl -f http://localhost:8080/api/health`
- **minio**: `curl -f http://localhost:9000/minio/health/live`

## 邮件系统特点

### 完整的邮件服务器

1. **SMTP服务器** (发送邮件):
   - MTA功能: 端口25接收外部邮件
   - MSA功能: 端口587用户提交邮件
   - 支持邮件转发和中继

2. **IMAP服务器** (接收邮件):
   - 端口993提供SSL加密访问
   - 支持文件夹管理
   - 支持邮件搜索和标记

3. **Web管理界面**:
   - 用户邮件管理
   - 管理员后台
   - API接口

### 域名配置

系统默认域名: `email.host`

**DNS配置建议**:
```dns
# MX记录
email.host.    IN  MX  10  your-server.com.

# A记录
your-server.com.    IN  A   YOUR_SERVER_IP
```

## 开发模式

```bash
# 只启动MinIO依赖服务
docker-compose -f docker-compose.dev.yml up -d

# 本地运行Go应用 (会自动启动邮件服务器)
go run main.go -f etc/config.yaml

# 本地运行前端（另一个终端）
cd web
npm run dev
```

**开发模式端口**:
- Web界面: http://localhost:8080
- SMTP: localhost:25, localhost:587
- IMAP: localhost:993
- MinIO: http://localhost:9001

## 生产部署

### 1. 防火墙配置

确保以下端口开放：

```bash
# Web管理界面
sudo ufw allow 8080/tcp

# 邮件服务端口
sudo ufw allow 25/tcp     # SMTP接收
sudo ufw allow 587/tcp    # SMTP提交  
sudo ufw allow 993/tcp    # IMAP SSL

# MinIO (可选，如需外部访问)
sudo ufw allow 9000/tcp   # MinIO API
sudo ufw allow 9001/tcp   # MinIO控制台
```

### 2. 域名和DNS配置

```bash
# 修改配置文件中的域名
# etc/config.yaml 中设置真实域名
# main.go 中修改 Domain: "yourdomain.com"
```

### 3. SSL/TLS证书

对于生产环境，建议配置SSL证书：

```bash
# 可以使用Let's Encrypt等免费证书
# 或者在前置反向代理(如Nginx)中配置SSL
```

## 故障排除

### 邮件服务相关

1. **端口25被占用**
   ```bash
   # 检查是否有其他邮件服务占用端口
   sudo netstat -tlnp | grep :25
   sudo systemctl stop postfix  # 如果有postfix服务
   ```

2. **权限问题**
   ```bash
   # 端口25需要root权限，或使用capability
   sudo setcap 'cap_net_bind_service=+ep' /path/to/email-system
   ```

3. **邮件客户端连接失败**
   ```bash
   # 检查邮件服务器日志
   docker-compose logs email-system | grep -i smtp
   docker-compose logs email-system | grep -i imap
   ```

### 测试邮件功能

```bash
# 测试SMTP连接
telnet localhost 587

# 测试IMAP连接  
openssl s_client -connect localhost:993

# 通过API发送测试邮件
curl -X POST http://localhost:8080/api/emails/send \
  -H "Content-Type: application/json" \
  -d '{"to":"test@email.host","subject":"测试","body":"Hello"}'
```

## 监控和维护

### 查看邮件服务状态

```bash
# 查看所有端口监听状态
docker-compose exec email-system netstat -tlnp

# 查看邮件队列
# (根据具体实现查看待发送/已发送邮件)

# 查看邮件日志
docker-compose logs email-system | grep -E "(SMTP|IMAP|邮件)"
```

## 安全注意事项

1. **端口安全**: 
   - 生产环境限制25/587/993端口访问
   - 使用防火墙规则限制IP访问

2. **认证安全**:
   - 更改默认密码
   - 启用强密码策略
   - 定期轮换密钥

3. **邮件安全**:
   - 配置SPF/DKIM/DMARC记录
   - 启用反垃圾邮件过滤
   - 监控邮件发送量

4. **数据备份**:
   - 定期备份SQLite数据库
   - 备份邮件附件和MinIO数据
   - 测试恢复流程

## 技术栈

- **邮件服务器**: Go + emersion/go-smtp + emersion/go-imap
- **后端**: Go 1.21 + Gin + SQLite
- **前端**: Vue 3 + TypeScript + Vite
- **存储**: MinIO对象存储
- **容器化**: Docker + Docker Compose

## 文件结构

```
project/
├── docker-compose.yml        # 生产环境编排配置
├── docker-compose.dev.yml    # 开发环境编排配置
├── Dockerfile               # 容器构建文件
├── data/                    # 数据目录（持久化）
│   ├── email.db            # SQLite数据库
│   ├── attachments/        # 邮件附件
│   ├── uploads/            # 上传文件
│   └── logs/               # 日志文件
├── etc/config.yaml          # 应用配置
├── internal/mailserver/     # 邮件服务器实现
└── web/dist/               # 前端构建产物
``` 