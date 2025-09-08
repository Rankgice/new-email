# 邮件系统后端 Dockerfile（支持CGO和SQLite）

# 构建阶段
FROM golang:1.24-bullseye AS builder

# 安装构建依赖
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

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
    -o email-system main.go

# 运行阶段
FROM debian:bullseye-slim

# 安装运行时依赖
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    curl \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

# 创建非root用户
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -m appuser

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