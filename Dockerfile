# 邮件系统后端 Dockerfile（支持CGO和SQLite）

# 构建阶段
FROM golang:1.24-alpine AS builder

# 更新包索引并安装构建依赖
RUN apk update && apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    git \
    ca-certificates

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建Go应用（启用CGO以支持SQLite）
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o email-system main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    sqlite

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/email-system .

# 复制配置文件
COPY etc ./etc

# 创建数据目录
RUN mkdir -p data/attachments data/logs data/uploads && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080 25 587 993

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/health || exit 1

# 启动命令
CMD ["./email-system", "-f", "etc/config.yaml"]