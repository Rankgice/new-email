# 邮件管理系统

一个现代化的邮件管理系统，支持多邮箱管理、验证码自动提取、邮件规则处理等功能。

## 🚀 功能特性

### 核心功能
- **多邮箱管理** - 支持自建邮箱和第三方邮箱（Gmail、Outlook、QQ等）
- **邮件收发** - 自动收取邮件，支持发送邮件和草稿管理
- **验证码提取** - 智能识别和提取邮件中的验证码
- **规则引擎** - 支持邮件转发规则和反垃圾邮件规则
- **域名管理** - DNS验证、DKIM、SPF、DMARC配置
- **API接口** - 完整的RESTful API，支持第三方集成

### 管理功能
- **用户管理** - 用户注册、登录、权限控制
- **管理员面板** - 系统监控、用户管理、全局设置
- **日志审计** - 完整的操作日志和邮件日志
- **数据导入导出** - 支持批量操作和数据迁移

### 已实现的Handler方法
- **日志管理** - 操作日志和邮件日志的查询和管理
- **签名管理** - 邮件签名的创建、更新、删除和设置默认
- **验证码管理** - 验证码列表查询、详情查看、标记使用状态
- **模板管理** - 邮件模板的完整CRUD操作和分类管理
- **API密钥管理** - API密钥的创建、更新、删除和权限管理
- **规则管理** - 验证码规则和转发规则的完整管理
- **通用功能** - 验证码发送验证、系统信息查询
- **草稿管理** - 邮件草稿的创建、编辑、删除和发送

## 🛠️ 技术栈

### 后端
- **Go 1.21+** - 主要编程语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **SQLite/MySQL** - 数据库
- **JWT** - 身份认证
- **Zap** - 日志框架

### 前端
- **Vue 3** - 前端框架
- **Element Plus** - UI组件库
- **Vite** - 构建工具
- **Axios** - HTTP客户端

### 邮件处理
- **go-mail** - SMTP发送
- **go-imap** - IMAP接收
- **go-message** - 邮件解析

## 📦 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- SQLite 3 或 MySQL 5.7+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd new-email
```

2. **安装依赖**
```bash
# 后端依赖
go mod tidy

# 前端依赖（如果需要开发前端）
cd web && npm install
```

3. **配置文件**
```bash
# 复制配置文件
cp etc/config.yaml.example etc/config.yaml

# 编辑配置文件
vim etc/config.yaml
```

4. **启动服务**
```bash
# 开发模式
go run main.go -f etc/config.yaml

# 生产模式
go build -o email-system main.go
./email-system -f etc/config.yaml
```

5. **访问系统**
- 系统首页: http://localhost:8081
- 管理员面板: http://localhost:8081/admin
- 用户面板: http://localhost:8081/user
- API文档: http://localhost:8081/api/health

## ⚙️ 配置说明

### 基本配置
```yaml
# 应用配置
app:
  name: "邮件管理系统"
  version: "1.0.0"
  debug: true

# Web服务配置
web:
  port: 8081
  mode: "debug"

# 数据库配置
database:
  type: "sqlite"  # sqlite 或 mysql
  sqlite:
    path: "./data/email.db"
```

### 邮件配置
```yaml
# 邮件配置
email:
  default_smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "your-email@gmail.com"
    password: "your-app-password"
    use_tls: true
```

## ✅ 已完成的功能模块

### 1. 日志管理 (LogHandler)
- `ListOperationLogs` - 操作日志列表查询，支持分页和筛选
- `ListEmailLogs` - 邮件日志列表查询，支持按邮件ID、邮箱ID等筛选

### 2. 签名管理 (SignatureHandler)
- `List` - 签名列表查询，支持分页和按名称筛选
- `Create` - 创建新的邮件签名
- `Update` - 更新现有签名信息
- `Delete` - 删除签名（软删除）
- `GetById` - 获取签名详情
- `SetDefault` - 设置默认签名

### 3. 验证码管理 (VerificationCodeHandler)
- `List` - 验证码列表查询，支持按来源、使用状态等筛选
- `GetById` - 获取验证码详情
- `MarkAsUsed` - 标记验证码为已使用
- `GetLatest` - 获取指定来源的最新验证码
- `GetStatistics` - 获取验证码统计信息

### 4. 模板管理 (TemplateHandler)
- `List` - 模板列表查询，支持按分类、名称筛选
- `Create` - 创建新的邮件模板
- `Update` - 更新模板内容和配置
- `Delete` - 删除模板（软删除）
- `GetById` - 获取模板详情
- `GetCategories` - 获取模板分类列表
- `Copy` - 复制现有模板

### 5. API密钥管理 (ApiKeyHandler)
- `List` - API密钥列表查询（密钥已脱敏）
- `Create` - 创建新的API密钥（仅创建时返回完整密钥）
- `Update` - 更新API密钥权限和配置
- `Delete` - 删除API密钥（软删除）
- `GetById` - 获取API密钥详情

### 6. 规则管理 (RuleHandler)
- **验证码规则**
  - `ListVerificationRules` - 验证码规则列表
  - `CreateVerificationRule` - 创建验证码提取规则
  - `UpdateVerificationRule` - 更新验证码规则
  - `DeleteVerificationRule` - 删除验证码规则
- **转发规则**
  - `ListForwardRules` - 转发规则列表
  - `CreateForwardRule` - 创建邮件转发规则
  - `UpdateForwardRule` - 更新转发规则
  - `DeleteForwardRule` - 删除转发规则

### 7. 通用功能 (CommonHandler)
- `SendCode` - 发送验证码（邮件/短信）
- `VerifyCode` - 验证验证码有效性
- `GetSystemInfo` - 获取系统基本信息
- `GetCaptcha` - 获取图形验证码（待实现）
- `VerifyCaptcha` - 验证图形验证码（待实现）
- `Upload` - 文件上传（待实现）

### 8. 草稿管理 (DraftHandler)
- `List` - 草稿列表查询，支持按邮箱、主题筛选
- `Create` - 创建新的邮件草稿
- `Update` - 更新草稿内容
- `Delete` - 删除草稿（软删除）
- `GetById` - 获取草稿详情
- `Send` - 发送草稿邮件（待完善邮件发送逻辑）

### 9. 管理员日志 (AdminLogHandler)
- `ListOperationLogs` - 管理员操作日志查询
- `ListEmailLogs` - 管理员邮件日志查询

## 🔌 API 接口

### 认证接口
```bash
# 用户登录
POST /api/public/user/login
{
  "username": "user@example.com",
  "password": "password"
}

# 管理员登录
POST /api/public/admin/login
{
  "username": "admin",
  "password": "admin123"
}
```

### 邮件接口
```bash
# 获取邮件列表
GET /api/user/emails?page=1&pageSize=20

# 发送邮件
POST /api/user/emails/send
{
  "to": ["recipient@example.com"],
  "subject": "测试邮件",
  "content": "邮件内容"
}
```

### API密钥访问
```bash
# 使用API密钥获取邮件
GET /api/v1/emails
Headers: X-API-Key: your-api-key

# 获取验证码
GET /api/v1/verification-codes
Headers: X-API-Key: your-api-key
```

## 📁 项目结构

```
new-email/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块依赖
├── etc/                    # 配置文件
│   └── config.yaml
├── internal/               # 内部代码
│   ├── config/            # 配置处理
│   ├── handler/           # HTTP处理器
│   ├── service/           # 业务逻辑层
│   ├── model/             # 数据模型层
│   ├── result/            # 统一响应结构
│   ├── constant/          # 系统常量
│   ├── types/             # 接口类型定义
│   ├── router/            # 路由配置
│   ├── middleware/        # 中间件
│   └── svc/               # 服务上下文
├── pkg/                   # 通用工具包
│   ├── auth/              # 认证工具
│   └── utils/             # 通用工具
├── data/                  # 数据存储
│   ├── email.db           # SQLite数据库
│   ├── attachments/       # 邮件附件
│   └── logs/              # 日志文件
└── web/                   # 前端项目
    └── dist/              # 前端构建产物
```

## 🔧 开发指南

### 添加新功能
1. 在 `internal/model/` 中定义数据模型
2. 在 `internal/types/` 中定义请求/响应类型
3. 在 `internal/handler/` 中实现HTTP处理器
4. 在 `internal/router/` 中添加路由配置

### 数据库迁移
系统使用GORM的AutoMigrate功能自动迁移数据库结构。

### 编程规范
- 遵循Go语言编程规范
- 使用GORM进行数据库操作
- 统一错误处理和响应格式
- 完整的日志记录

## 🐳 Docker 部署

```dockerfile
# 构建镜像
docker build -t email-system .

# 运行容器
docker run -d \
  --name email-system \
  -p 8081:8081 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/etc:/app/etc \
  email-system
```

## 📝 默认账户

### 管理员账户
- 用户名: `admin`
- 密码: `admin123`
- 邮箱: `admin@example.com`

### 注意事项
- 首次启动后请立即修改默认管理员密码
- 生产环境请使用强密码和HTTPS
- 定期备份数据库和配置文件

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

如果您遇到问题或有建议，请：
1. 查看文档和FAQ
2. 搜索已有的Issues
3. 创建新的Issue描述问题
4. 联系开发团队

---

**邮件管理系统** - 让邮件管理更简单、更智能！
