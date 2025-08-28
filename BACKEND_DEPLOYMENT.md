# 后端API部署指南

## 概述

这是邮件系统的后端API部署配置，不包含前端界面。适用于前后端分离的部署场景。

## CI/CD 流水线

### 自动构建和部署

当代码推送到 `main` 分支时，GitHub Actions 会自动：

1. 构建 Go 应用程序
2. 创建 Docker 镜像（仅包含后端）
3. 推送到 GitHub Container Registry (GHCR)

### 流水线配置

```yaml
# .github/workflows/deploy.yml
- 构建 Go 应用: CGO_ENABLED=1 GOOS=linux
- 创建 Docker 镜像: 使用 Dockerfile.backend
- 推送到: ghcr.io/your-repo/app:latest
```

## 部署方式

### 1. 使用预构建镜像部署

```bash
# 拉取最新镜像
docker pull ghcr.io/your-repo/app:latest

# 运行容器
docker run -d \
  --name email-system \
  -p 8080:8080 \
  -p 25:25 \
  -p 587:587 \
  -p 993:993 \
  -v ./data:/app/data \
  -v ./etc:/app/etc \
  ghcr.io/your-repo/app:latest
```

### 2. 使用 docker-compose 部署

```yaml
# docker-compose.backend.yml
version: '3.8'

services:
  email-system:
    image: ghcr.io/your-repo/app:latest
    container_name: email-system
    restart: unless-stopped
    ports:
      - "8080:8080"    # API接口
      - "25:25"        # SMTP接收
      - "587:587"      # SMTP提交
      - "993:993"      # IMAP访问
    volumes:
      - ./data:/app/data
      - ./etc:/app/etc
    environment:
      - TZ=Asia/Shanghai
    depends_on:
      minio:
        condition: service_healthy

  minio:
    image: minio/minio:latest
    container_name: email-minio
    restart: unless-stopped
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"

volumes:
  minio_data:
```

## API 接口

### 健康检查

```bash
curl http://localhost:8080/api/health
```

响应：
```json
{
  "status": "ok",
  "message": "邮件系统运行正常",
  "version": "1.0.0"
}
```

### 主要API端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/` | GET | API状态页面 |
| `/api/health` | GET | 健康检查 |
| `/api/public/user/login` | POST | 用户登录 |
| `/api/public/admin/login` | POST | 管理员登录 |
| `/api/user/emails` | GET | 获取邮件列表 |
| `/api/user/emails/send` | POST | 发送邮件 |

完整API文档请参考 `API_DOCUMENTATION.md`

## 邮件服务

### 端口配置

| 端口 | 协议 | 用途 |
|------|------|------|
| 8080 | HTTP | API接口 |
| 25 | SMTP | 邮件接收 (MTA) |
| 587 | SMTP | 邮件提交 (MSA) |
| 993 | IMAP | 邮件访问 (SSL) |

### 邮件客户端配置

**接收邮件 (IMAP)**:
- 服务器: your-server.com
- 端口: 993
- 加密: SSL/TLS
- 用户名: 邮箱地址
- 密码: 邮箱密码

**发送邮件 (SMTP)**:
- 服务器: your-server.com
- 端口: 587
- 加密: STARTTLS
- 认证: 是
- 用户名: 邮箱地址
- 密码: 邮箱密码

## 配置文件

### 必要的配置文件

```
etc/
└── config.yaml    # 主配置文件
```

确保配置文件中的以下设置正确：

```yaml
# etc/config.yaml
app:
  debug: false    # 生产环境设为false

web:
  port: 8080     # API端口

database:
  type: "sqlite"
  sqlite:
    path: "./data/email.db"

minio:
  endpoint: "minio:9000"    # 或外部MinIO地址
  access_key_id: "minioadmin"
  secret_access_key: "minioadmin"
```

## 前端部署

由于这是纯后端部署，前端需要单独部署：

### 选项1: 静态文件服务器

```bash
# 构建前端
cd web
npm install
npm run build

# 部署到nginx、Apache等
cp -r dist/* /var/www/html/
```

### 选项2: CDN部署

将前端构建产物上传到CDN或对象存储：

```bash
# 上传到云存储
aws s3 sync dist/ s3://your-bucket/
# 或
gsutil -m cp -r dist/* gs://your-bucket/
```

### 选项3: 容器化前端

```dockerfile
# Dockerfile.frontend
FROM nginx:alpine
COPY web/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
```

## 监控和维护

### 日志查看

```bash
# 容器日志
docker logs email-system

# 应用日志
docker exec email-system tail -f /app/data/logs/app.log
```

### 数据备份

```bash
# 备份SQLite数据库
docker cp email-system:/app/data/email.db ./backup/

# 备份MinIO数据
docker exec email-minio mc cp --recursive /data ./backup/
```

### 更新部署

```bash
# 拉取最新镜像
docker pull ghcr.io/your-repo/app:latest

# 重启容器
docker-compose down
docker-compose up -d
```

## 安全注意事项

1. **API安全**:
   - 启用HTTPS（建议使用反向代理）
   - 配置CORS策略
   - 使用强JWT密钥

2. **邮件安全**:
   - 配置防火墙规则
   - 启用邮件认证
   - 监控异常流量

3. **数据安全**:
   - 定期备份数据库
   - 加密敏感配置
   - 限制容器权限

## 故障排除

### 常见问题

1. **容器启动失败**
   ```bash
   docker logs email-system
   ```

2. **API无法访问**
   ```bash
   # 检查端口映射
   docker port email-system
   
   # 测试健康检查
   curl http://localhost:8080/api/health
   ```

3. **邮件服务无法连接**
   ```bash
   # 检查端口监听
   netstat -tlnp | grep -E "25|587|993"
   
   # 测试SMTP连接
   telnet localhost 587
   ```

## 技术栈

- **后端**: Go 1.24 + Gin
- **邮件服务**: emersion/go-smtp + emersion/go-imap
- **数据库**: SQLite
- **对象存储**: MinIO
- **容器化**: Docker
- **CI/CD**: GitHub Actions 