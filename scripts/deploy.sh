#!/bin/bash

# 邮件管理系统部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 命令未找到，请先安装"
        exit 1
    fi
}

# 项目信息
PROJECT_NAME="email-system"
DEPLOY_MODE=${1:-"standalone"}  # standalone, docker, docker-compose

print_message "开始部署 $PROJECT_NAME"
print_message "部署模式: $DEPLOY_MODE"

# 进入项目根目录
cd "$(dirname "$0")/.."
PROJECT_ROOT=$(pwd)

case $DEPLOY_MODE in
    "standalone")
        print_step "独立部署模式"
        
        # 检查必要的命令
        check_command "go"
        
        # 构建项目
        print_step "构建项目..."
        ./scripts/build.sh
        
        # 创建部署目录
        DEPLOY_DIR="/opt/$PROJECT_NAME"
        print_step "创建部署目录: $DEPLOY_DIR"
        sudo mkdir -p "$DEPLOY_DIR"
        
        # 复制文件
        print_step "复制文件到部署目录..."
        sudo cp -r build/* "$DEPLOY_DIR/"
        sudo chown -R $(whoami):$(whoami) "$DEPLOY_DIR"
        
        # 创建systemd服务文件
        print_step "创建systemd服务..."
        sudo tee /etc/systemd/system/$PROJECT_NAME.service > /dev/null << EOF
[Unit]
Description=Email Management System
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=$DEPLOY_DIR
ExecStart=$DEPLOY_DIR/$PROJECT_NAME -f $DEPLOY_DIR/etc/config.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
        
        # 重新加载systemd并启动服务
        sudo systemctl daemon-reload
        sudo systemctl enable $PROJECT_NAME
        
        print_message "部署完成！"
        print_message "配置文件位置: $DEPLOY_DIR/etc/config.yaml"
        print_message "启动服务: sudo systemctl start $PROJECT_NAME"
        print_message "查看状态: sudo systemctl status $PROJECT_NAME"
        print_message "查看日志: sudo journalctl -u $PROJECT_NAME -f"
        ;;
        
    "docker")
        print_step "Docker部署模式"
        
        # 检查Docker
        check_command "docker"
        
        # 构建Docker镜像
        print_step "构建Docker镜像..."
        docker build -t $PROJECT_NAME:latest .
        
        # 停止并删除旧容器
        print_step "停止旧容器..."
        docker stop $PROJECT_NAME 2>/dev/null || true
        docker rm $PROJECT_NAME 2>/dev/null || true
        
        # 创建数据目录
        mkdir -p data/attachments data/logs
        
        # 运行新容器
        print_step "启动新容器..."
        docker run -d \
            --name $PROJECT_NAME \
            --restart unless-stopped \
            -p 8081:8081 \
            -v $(pwd)/data:/app/data \
            -v $(pwd)/etc:/app/etc \
            $PROJECT_NAME:latest
        
        print_message "Docker部署完成！"
        print_message "容器状态: docker ps | grep $PROJECT_NAME"
        print_message "查看日志: docker logs -f $PROJECT_NAME"
        print_message "进入容器: docker exec -it $PROJECT_NAME sh"
        ;;
        
    "docker-compose")
        print_step "Docker Compose部署模式"
        
        # 检查Docker Compose
        check_command "docker-compose"
        
        # 停止旧服务
        print_step "停止旧服务..."
        docker-compose down 2>/dev/null || true
        
        # 构建并启动服务
        print_step "构建并启动服务..."
        docker-compose up -d --build
        
        print_message "Docker Compose部署完成！"
        print_message "查看服务状态: docker-compose ps"
        print_message "查看日志: docker-compose logs -f"
        print_message "停止服务: docker-compose down"
        ;;
        
    *)
        print_error "不支持的部署模式: $DEPLOY_MODE"
        print_message "支持的部署模式:"
        print_message "  standalone    - 独立部署（systemd服务）"
        print_message "  docker        - Docker容器部署"
        print_message "  docker-compose - Docker Compose部署"
        exit 1
        ;;
esac

# 等待服务启动
print_step "等待服务启动..."
sleep 5

# 健康检查
print_step "执行健康检查..."
for i in {1..10}; do
    if curl -f http://localhost:8081/api/health >/dev/null 2>&1; then
        print_message "✅ 服务启动成功！"
        print_message "🌐 访问地址: http://localhost:8081"
        print_message "👤 默认管理员: admin / admin123"
        break
    else
        print_warning "等待服务启动... ($i/10)"
        sleep 3
    fi
    
    if [ $i -eq 10 ]; then
        print_error "❌ 服务启动失败，请检查日志"
        exit 1
    fi
done

print_message "🎉 部署成功完成！"
