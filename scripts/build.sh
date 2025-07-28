#!/bin/bash

# 邮件管理系统构建脚本

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
VERSION="1.0.0"
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

print_message "开始构建 $PROJECT_NAME v$VERSION"
print_message "构建时间: $BUILD_TIME"
print_message "Git提交: $GIT_COMMIT"

# 检查必要的命令
print_step "检查构建环境..."
check_command "go"

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_message "Go版本: $GO_VERSION"

# 进入项目根目录
cd "$(dirname "$0")/.."
PROJECT_ROOT=$(pwd)
print_message "项目根目录: $PROJECT_ROOT"

# 创建构建目录
BUILD_DIR="$PROJECT_ROOT/build"
mkdir -p "$BUILD_DIR"

# 清理旧的构建文件
print_step "清理旧的构建文件..."
rm -rf "$BUILD_DIR"/*

# 下载依赖
print_step "下载Go依赖..."
go mod tidy
go mod download

# 运行测试（如果存在）
if [ -f "$PROJECT_ROOT/go.mod" ] && grep -q "testing" "$PROJECT_ROOT"/**/*.go 2>/dev/null; then
    print_step "运行测试..."
    go test ./... -v
else
    print_warning "未找到测试文件，跳过测试"
fi

# 构建标志
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime='$BUILD_TIME' -X main.GitCommit=$GIT_COMMIT"

# 构建不同平台的二进制文件
print_step "构建二进制文件..."

# Linux amd64
print_message "构建 Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${PROJECT_NAME}-linux-amd64" main.go

# Linux arm64
print_message "构建 Linux arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${PROJECT_NAME}-linux-arm64" main.go

# Windows amd64
print_message "构建 Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${PROJECT_NAME}-windows-amd64.exe" main.go

# macOS amd64
print_message "构建 macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${PROJECT_NAME}-darwin-amd64" main.go

# macOS arm64 (Apple Silicon)
print_message "构建 macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/${PROJECT_NAME}-darwin-arm64" main.go

# 本地平台
print_message "构建本地平台..."
go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$PROJECT_NAME" main.go

# 复制配置文件和其他资源
print_step "复制资源文件..."
cp -r "$PROJECT_ROOT/etc" "$BUILD_DIR/"
cp -r "$PROJECT_ROOT/web/dist" "$BUILD_DIR/web/"
cp "$PROJECT_ROOT/README.md" "$BUILD_DIR/"

# 创建数据目录
mkdir -p "$BUILD_DIR/data/attachments"
mkdir -p "$BUILD_DIR/data/logs"

# 创建启动脚本
print_step "创建启动脚本..."

# Linux/macOS 启动脚本
cat > "$BUILD_DIR/start.sh" << 'EOF'
#!/bin/bash

# 邮件管理系统启动脚本

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 检查配置文件
if [ ! -f "etc/config.yaml" ]; then
    echo "错误: 配置文件 etc/config.yaml 不存在"
    echo "请复制 etc/config.yaml.example 并修改配置"
    exit 1
fi

# 创建必要的目录
mkdir -p data/attachments
mkdir -p data/logs

# 启动服务
echo "启动邮件管理系统..."
./email-system -f etc/config.yaml
EOF

# Windows 启动脚本
cat > "$BUILD_DIR/start.bat" << 'EOF'
@echo off
chcp 65001 >nul

cd /d "%~dp0"

if not exist "etc\config.yaml" (
    echo 错误: 配置文件 etc\config.yaml 不存在
    echo 请复制 etc\config.yaml.example 并修改配置
    pause
    exit /b 1
)

if not exist "data" mkdir data
if not exist "data\attachments" mkdir data\attachments
if not exist "data\logs" mkdir data\logs

echo 启动邮件管理系统...
email-system-windows-amd64.exe -f etc\config.yaml
pause
EOF

# 设置执行权限
chmod +x "$BUILD_DIR/start.sh"
chmod +x "$BUILD_DIR"/${PROJECT_NAME}*

# 创建配置文件示例
print_step "创建配置文件示例..."
if [ -f "$PROJECT_ROOT/etc/config.yaml" ]; then
    cp "$PROJECT_ROOT/etc/config.yaml" "$BUILD_DIR/etc/config.yaml.example"
fi

# 显示构建结果
print_step "构建完成！"
print_message "构建目录: $BUILD_DIR"
print_message "构建的文件:"
ls -la "$BUILD_DIR"

# 计算文件大小
print_message "文件大小:"
du -sh "$BUILD_DIR"/*

print_message "构建成功完成！"
print_message "使用方法:"
print_message "  1. 进入构建目录: cd $BUILD_DIR"
print_message "  2. 复制配置文件: cp etc/config.yaml.example etc/config.yaml"
print_message "  3. 修改配置文件: vim etc/config.yaml"
print_message "  4. 启动服务: ./start.sh (Linux/macOS) 或 start.bat (Windows)"

print_message "🎉 构建完成！"
