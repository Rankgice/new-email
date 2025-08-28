# 邮件管理系统 Dockerfile

# 前端构建阶段
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm ci --only=production

# 复制前端源码
COPY web/ ./

# 构建前端
RUN npm run build

# Go应用构建阶段
FROM golang:1.24-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（SQLite需要CGO）
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o email-system main.go

# 运行阶段
FROM alpine:latest

# 安装必要的包（包括sqlite、curl用于健康检查）
RUN apk --no-cache add ca-certificates tzdata curl sqlite

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=backend-builder /app/email-system .

# 复制配置文件
COPY --from=backend-builder /app/etc ./etc

# 从前端构建阶段复制静态资源
COPY --from=frontend-builder /app/web/dist ./web/dist

# 创建数据目录
RUN mkdir -p data/attachments data/logs data/uploads && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口

# Web管理界面和API
EXPOSE 8080
# SMTP接收端口 (MTA) - 接收外部邮件
EXPOSE 25
# SMTP提交端口 (MSA) - 用户发送邮件
EXPOSE 587
# IMAP端口 (SSL) - 邮件客户端访问
EXPOSE 993

# 健康检查（使用curl替代wget）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/health || exit 1

# 启动命令
CMD ["./email-system", "-f", "etc/config.yaml"]
