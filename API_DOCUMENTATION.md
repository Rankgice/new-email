# 📧 邮件系统 API 文档

## 📋 目录

- [基础信息](#基础信息)
- [认证方式](#认证方式)
- [公共接口](#公共接口)
- [用户接口](#用户接口)
- [管理员接口](#管理员接口)
- [API密钥接口](#api密钥接口)
- [错误码说明](#错误码说明)

## 🔧 基础信息

### 服务器地址
- **开发环境**: `http://localhost:8081`
- **生产环境**: `https://your-domain.com`

### 请求格式
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`

### 响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🔐 认证方式

### 1. JWT Token 认证 (用户/管理员)
```http
Authorization: Bearer <token>
```

### 2. API Key 认证
```http
X-API-Key: <api_key>
```

## 🌐 公共接口

### 健康检查
```http
GET /api/health
```

### 用户注册
```http
POST /api/public/user/register
Content-Type: application/json

{
  "username": "string",
  "email": "string", 
  "password": "string",
  "nickname": "string"
}
```

### 用户登录
```http
POST /api/public/user/login
Content-Type: application/json

{
  "email": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": "string",
      "username": "string",
      "email": "string",
      "nickname": "string",
      "role": "user"
    }
  }
}
```

### 管理员登录
```http
POST /api/public/admin/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}
```

### 发送验证码
```http
POST /api/public/send-code
Content-Type: application/json

{
  "email": "string",
  "type": "register|reset_password|verify_email"
}
```

### 验证验证码
```http
POST /api/public/verify-code
Content-Type: application/json

{
  "email": "string",
  "code": "string",
  "type": "register|reset_password|verify_email"
}
```

## 👤 用户接口

> 所有用户接口需要 JWT Token 认证

### 用户信息管理

#### 获取用户资料
```http
GET /api/user/profile
Authorization: Bearer <token>
```

#### 更新用户资料
```http
PUT /api/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "string",
  "avatar": "string"
}
```

#### 修改密码
```http
POST /api/user/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "string",
  "new_password": "string"
}
```

### 邮箱管理

#### 获取邮箱列表
```http
GET /api/user/mailboxes
Authorization: Bearer <token>
```

#### 创建邮箱
```http
POST /api/user/mailboxes
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "string",
  "password": "string",
  "imap_server": "string",
  "imap_port": 993,
  "smtp_server": "string", 
  "smtp_port": 587,
  "use_ssl": true
}
```

#### 更新邮箱
```http
PUT /api/user/mailboxes/:id
Authorization: Bearer <token>
```

#### 删除邮箱
```http
DELETE /api/user/mailboxes/:id
Authorization: Bearer <token>
```

#### 同步邮箱
```http
POST /api/user/mailboxes/:id/sync
Authorization: Bearer <token>
```

#### 测试邮箱连接
```http
POST /api/user/mailboxes/:id/test
Authorization: Bearer <token>
```

### 邮件管理

#### 获取邮件列表
```http
GET /api/user/emails?page=1&limit=20&folder=inbox&search=keyword
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `limit`: 每页数量 (默认: 20)
- `folder`: 文件夹 (inbox/sent/drafts/trash)
- `search`: 搜索关键词
- `is_read`: 是否已读 (true/false)
- `is_starred`: 是否加星 (true/false)

#### 获取邮件详情
```http
GET /api/user/emails/:id
Authorization: Bearer <token>
```

#### 发送邮件
```http
POST /api/user/emails/send
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": ["email1@example.com"],
  "cc": ["email2@example.com"],
  "bcc": ["email3@example.com"],
  "subject": "string",
  "content": "string",
  "html_content": "string",
  "attachments": ["file_id1", "file_id2"]
}
```

#### 标记已读/未读
```http
PUT /api/user/emails/:id/read
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_read": true
}
```

#### 标记星标
```http
PUT /api/user/emails/:id/star
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_starred": true
}
```

#### 删除邮件
```http
DELETE /api/user/emails/:id
Authorization: Bearer <token>
```

#### 批量操作
```http
POST /api/user/emails/batch
Authorization: Bearer <token>
Content-Type: application/json

{
  "email_ids": ["id1", "id2"],
  "action": "read|unread|star|unstar|delete|move",
  "folder": "trash"
}
```

### 草稿管理

#### 获取草稿列表
```http
GET /api/user/drafts
Authorization: Bearer <token>
```

#### 创建草稿
```http
POST /api/user/drafts
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": ["email@example.com"],
  "subject": "string",
  "content": "string"
}
```

#### 更新草稿
```http
PUT /api/user/drafts/:id
Authorization: Bearer <token>
```

#### 删除草稿
```http
DELETE /api/user/drafts/:id
Authorization: Bearer <token>
```

#### 发送草稿
```http
POST /api/user/drafts/:id/send
Authorization: Bearer <token>
```

### 邮件模板

#### 获取模板列表
```http
GET /api/user/templates
Authorization: Bearer <token>
```

#### 创建模板
```http
POST /api/user/templates
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "subject": "string",
  "content": "string"
}
```

#### 更新模板
```http
PUT /api/user/templates/:id
Authorization: Bearer <token>
```

#### 删除模板
```http
DELETE /api/user/templates/:id
Authorization: Bearer <token>
```

### 邮件签名

#### 获取签名列表
```http
GET /api/user/signatures
Authorization: Bearer <token>
```

#### 创建签名
```http
POST /api/user/signatures
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "content": "string",
  "is_default": false
}
```

### 规则管理

#### 验证码规则

##### 获取验证码规则列表
```http
GET /api/user/rules/verification
Authorization: Bearer <token>
```

##### 创建验证码规则
```http
POST /api/user/rules/verification
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "pattern": "string",
  "description": "string"
}
```

#### 转发规则

##### 获取转发规则列表
```http
GET /api/user/rules/forward
Authorization: Bearer <token>
```

##### 创建转发规则
```http
POST /api/user/rules/forward
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "from_pattern": "string",
  "to_email": "string",
  "conditions": {}
}
```

### 验证码记录

#### 获取验证码记录
```http
GET /api/user/verification-codes
Authorization: Bearer <token>
```

#### 获取验证码详情
```http
GET /api/user/verification-codes/:id
Authorization: Bearer <token>
```

#### 标记验证码已使用
```http
PUT /api/user/verification-codes/:id/used
Authorization: Bearer <token>
```

### API密钥管理

#### 获取API密钥列表
```http
GET /api/user/api-keys
Authorization: Bearer <token>
```

#### 创建API密钥
```http
POST /api/user/api-keys
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "description": "string",
  "permissions": ["read", "write"]
}
```

### 日志查询

#### 获取操作日志
```http
GET /api/user/logs/operation
Authorization: Bearer <token>
```

#### 获取邮件日志
```http
GET /api/user/logs/email
Authorization: Bearer <token>
```

## 👑 管理员接口

> 所有管理员接口需要管理员 JWT Token 认证

### 用户管理

#### 获取用户列表
```http
GET /api/admin/users?page=1&limit=20&search=keyword
Authorization: Bearer <admin_token>
```

#### 获取用户详情
```http
GET /api/admin/users/:id
Authorization: Bearer <admin_token>
```

#### 创建用户
```http
POST /api/admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string",
  "role": "user|admin"
}
```

#### 更新用户
```http
PUT /api/admin/users/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "role": "user|admin",
  "status": "active|inactive"
}
```

#### 删除用户
```http
DELETE /api/admin/users/:id
Authorization: Bearer <admin_token>
```

#### 重置用户密码
```http
POST /api/admin/users/:id/reset-password
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "new_password": "string"
}
```

### 邮箱管理

#### 获取所有邮箱
```http
GET /api/admin/mailboxes
Authorization: Bearer <admin_token>
```

#### 获取邮箱详情
```http
GET /api/admin/mailboxes/:id
Authorization: Bearer <admin_token>
```

#### 强制同步邮箱
```http
POST /api/admin/mailboxes/:id/force-sync
Authorization: Bearer <admin_token>
```

#### 禁用/启用邮箱
```http
PUT /api/admin/mailboxes/:id/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": "active|inactive"
}
```

### 邮件监控

#### 获取所有邮件
```http
GET /api/admin/emails?page=1&limit=20
Authorization: Bearer <admin_token>
```

#### 获取邮件统计
```http
GET /api/admin/emails/stats
Authorization: Bearer <admin_token>
```

#### 删除邮件
```http
DELETE /api/admin/emails/:id
Authorization: Bearer <admin_token>
```

### 系统管理

#### 获取系统统计
```http
GET /api/admin/stats
Authorization: Bearer <admin_token>
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "total_users": 100,
    "active_users": 85,
    "total_emails": 5000,
    "today_emails": 50,
    "system_status": "healthy",
    "storage_used": "2.5GB",
    "storage_limit": "10GB"
  }
}
```

#### 获取系统配置
```http
GET /api/admin/config
Authorization: Bearer <admin_token>
```

#### 更新系统配置
```http
PUT /api/admin/config
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "smtp_settings": {},
  "storage_settings": {},
  "security_settings": {}
}
```

### 日志管理

#### 获取系统日志
```http
GET /api/admin/logs/system
Authorization: Bearer <admin_token>
```

#### 获取操作日志
```http
GET /api/admin/logs/operation
Authorization: Bearer <admin_token>
```

#### 获取错误日志
```http
GET /api/admin/logs/error
Authorization: Bearer <admin_token>
```

## 🔑 API密钥接口

> 使用 X-API-Key 头部认证

### 邮件操作

#### 发送邮件
```http
POST /api/v1/emails/send
X-API-Key: <api_key>
Content-Type: application/json

{
  "to": ["email@example.com"],
  "subject": "string",
  "content": "string"
}
```

#### 获取邮件列表
```http
GET /api/v1/emails
X-API-Key: <api_key>
```

#### 获取验证码
```http
GET /api/v1/verification-codes
X-API-Key: <api_key>
```

## ❌ 错误码说明

| 错误码 | 说明 | 描述 |
|--------|------|------|
| 200 | 成功 | 请求成功 |
| 400 | 请求错误 | 请求参数错误 |
| 401 | 未授权 | 需要登录或token无效 |
| 403 | 禁止访问 | 权限不足 |
| 404 | 未找到 | 资源不存在 |
| 409 | 冲突 | 资源已存在 |
| 422 | 验证失败 | 数据验证失败 |
| 429 | 请求过多 | 超出速率限制 |
| 500 | 服务器错误 | 内部服务器错误 |

### 错误响应格式
```json
{
  "code": 400,
  "message": "请求参数错误",
  "error": "validation_failed",
  "details": {
    "field": "email",
    "message": "邮箱格式不正确"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 📝 请求示例

### 用户登录示例
```bash
curl -X POST http://localhost:8081/api/public/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 获取邮件列表示例
```bash
curl -X GET "http://localhost:8081/api/user/emails?page=1&limit=10" \
  -H "Authorization: Bearer your_jwt_token"
```

### 发送邮件示例
```bash
curl -X POST http://localhost:8081/api/user/emails/send \
  -H "Authorization: Bearer your_jwt_token" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "测试邮件",
    "content": "这是一封测试邮件"
  }'
```

---

## 📚 更新日志

- **v1.0.0** (2024-01-15): 初始版本
- 支持用户注册、登录、邮件管理
- 支持管理员功能
- 支持API密钥认证
