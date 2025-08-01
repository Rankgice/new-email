# 邮件管理系统测试文档

## 测试概述

本目录包含了邮件管理系统的完整测试套件，涵盖了用户认证、邮箱管理、邮件收发等核心功能的测试。

## 测试文件结构

```
test/
├── auth/                           # 用户认证相关测试
│   ├── user_register_test.go      # Go单元测试
│   ├── register_api_test.http     # 用户注册API测试
│   ├── register_with_code_test.http # 带验证码的注册测试
│   ├── login_test.http            # 用户登录测试
│   ├── profile_test.http          # 用户资料管理测试
│   └── change_password_test.http  # 密码修改测试
├── mailbox/                       # 邮箱管理相关测试
│   ├── mailbox_create_test.http   # 邮箱添加测试
│   ├── mailbox_test_connection.http # 邮箱连接测试
│   └── mailbox_management_test.http # 邮箱管理测试
├── email/                         # 邮件收发相关测试
│   ├── email_send_test.http       # 邮件发送测试
│   └── email_receive_test.http    # 邮件接收测试
├── integration_test.http          # 集成测试
├── run_tests.sh                   # Linux/Mac自动化测试脚本
├── run_tests.bat                  # Windows自动化测试脚本
├── test_report_template.md        # 测试报告模板
└── README.md                      # 本文件
```

## 测试环境要求

### 系统要求
- 操作系统：Windows 10+, macOS 10.15+, 或 Linux
- Go 1.19+ (用于Go单元测试)
- curl (用于HTTP API测试)

### 服务要求
- 邮件管理系统服务运行在 `http://localhost:8080`
- 数据库服务正常运行
- 邮件服务配置正确

## 测试运行方法

### 1. 自动化测试

#### Linux/Mac系统
```bash
# 给脚本执行权限
chmod +x run_tests.sh

# 运行测试
./run_tests.sh
```

#### Windows系统
```cmd
# 直接运行批处理文件
run_tests.bat
```

### 2. 手动API测试

#### 使用VS Code REST Client插件
1. 安装VS Code的REST Client插件
2. 打开对应的`.http`文件
3. 点击"Send Request"按钮执行测试

#### 使用Postman
1. 导入`.http`文件到Postman
2. 配置环境变量
3. 逐个执行测试用例

#### 使用curl命令
```bash
# 示例：测试用户注册
curl -X POST http://localhost:8080/api/public/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "测试用户"
  }'
```

### 3. Go单元测试

```bash
# 进入测试目录
cd test/auth

# 运行Go测试
go test -v
```

## 测试用例说明

### 用户认证测试 (auth/)

#### 用户注册测试
- **文件**: `register_api_test.http`
- **测试内容**:
  - 正常用户注册
  - 用户名为空（应该失败）
  - 邮箱格式错误（应该失败）
  - 密码为空（应该失败）
  - 重复用户名注册（应该失败）
  - 重复邮箱注册（应该失败）

#### 带验证码注册测试
- **文件**: `register_with_code_test.http`
- **测试内容**:
  - 发送邮箱验证码
  - 使用验证码注册
  - 使用错误验证码注册（应该失败）
  - 使用过期验证码注册（应该失败）

#### 用户登录测试
- **文件**: `login_test.http`
- **测试内容**:
  - 用户名登录
  - 邮箱登录
  - 错误密码登录（应该失败）
  - 不存在用户登录（应该失败）
  - Token验证测试

#### 用户资料管理测试
- **文件**: `profile_test.http`
- **测试内容**:
  - 获取用户资料
  - 更新用户昵称
  - 更新用户头像
  - 权限验证测试

#### 密码修改测试
- **文件**: `change_password_test.http`
- **测试内容**:
  - 正常密码修改
  - 错误旧密码（应该失败）
  - 新旧密码相同（应该失败）
  - 新密码太短（应该失败）

### 邮箱管理测试 (mailbox/)

#### 邮箱添加测试
- **文件**: `mailbox_create_test.http`
- **测试内容**:
  - 添加Gmail邮箱
  - 添加Outlook邮箱
  - 添加QQ邮箱
  - 添加163邮箱
  - 添加自建邮箱
  - 重复邮箱添加（应该失败）
  - 无效邮箱格式（应该失败）

#### 邮箱连接测试
- **文件**: `mailbox_test_connection.http`
- **测试内容**:
  - 测试各种邮箱服务商连接
  - 错误密码连接测试（应该失败）
  - 错误服务器配置测试（应该失败）

#### 邮箱管理测试
- **文件**: `mailbox_management_test.http`
- **测试内容**:
  - 获取邮箱列表
  - 分页和筛选
  - 更新邮箱信息
  - 启用/禁用邮箱
  - 删除邮箱
  - 邮箱同步

### 邮件收发测试 (email/)

#### 邮件发送测试
- **文件**: `email_send_test.http`
- **测试内容**:
  - 发送文本邮件
  - 发送HTML邮件
  - 发送带抄送/密送邮件
  - 发送给多个收件人
  - 各种错误情况测试

#### 邮件接收测试
- **文件**: `email_receive_test.http`
- **测试内容**:
  - 获取收件箱邮件
  - 获取已发送邮件
  - 邮件搜索和筛选
  - 邮件操作（标记已读、星标、删除）
  - 批量操作

### 集成测试

#### 完整流程测试
- **文件**: `integration_test.http`
- **测试内容**:
  - 用户注册 → 登录 → 添加邮箱 → 发送邮件 → 邮件管理 → 密码修改
  - 模拟真实用户使用场景

## 测试数据配置

### 环境变量
在运行测试前，请确保以下配置正确：

```bash
# 基础配置
BASE_URL=http://localhost:8080
CONTENT_TYPE=application/json

# 测试用户信息
TEST_EMAIL=autotest@example.com
TEST_USERNAME=autotest_user
TEST_PASSWORD=AutoTest123!
```

### 测试邮箱配置
为了测试邮箱功能，需要准备以下测试邮箱：

```json
{
  "gmail": {
    "email": "test@gmail.com",
    "password": "app_password_here"
  },
  "outlook": {
    "email": "test@outlook.com", 
    "password": "password123"
  },
  "qq": {
    "email": "test@qq.com",
    "password": "authorization_code"
  }
}
```

**注意**: 
- Gmail需要使用应用专用密码
- QQ邮箱需要使用授权码
- 请不要在测试中使用真实的生产邮箱

## 测试结果分析

### 自动化测试结果
运行自动化测试脚本后，会生成详细的测试报告：

```
=== 测试结果统计 ===
总测试数: 25
通过测试: 23
失败测试: 2
通过率: 92%
```

### 日志文件
测试日志会保存在以下文件中：
- `test_results_YYYYMMDD_HHMMSS.log`

日志包含：
- 每个测试的请求和响应详情
- HTTP状态码
- 响应内容
- 错误信息

### 常见问题排查

#### 1. 连接失败
```
[FAIL] 系统健康检查 (Expected: 200, Got: 000)
```
**解决方案**: 检查服务是否在localhost:8080运行

#### 2. 认证失败
```
[FAIL] 获取用户资料 (Expected: 200, Got: 401)
```
**解决方案**: 检查Token是否正确获取和传递

#### 3. 邮箱连接失败
```
[FAIL] 测试Gmail连接 (Expected: 200, Got: 500)
```
**解决方案**: 检查邮箱密码和服务器配置

## 测试最佳实践

### 1. 测试隔离
- 每次测试使用独立的测试数据
- 测试完成后清理测试数据
- 避免测试之间的相互影响

### 2. 数据准备
- 使用专门的测试邮箱账户
- 准备各种边界条件的测试数据
- 模拟真实的使用场景

### 3. 错误处理
- 测试各种错误情况
- 验证错误信息的准确性
- 确保系统的健壮性

### 4. 性能测试
- 监控API响应时间
- 测试并发访问情况
- 验证系统在负载下的表现

## 持续集成

### GitHub Actions配置示例
```yaml
name: API Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Start services
        run: docker-compose up -d
      - name: Run tests
        run: ./test/run_tests.sh
      - name: Upload test results
        uses: actions/upload-artifact@v2
        with:
          name: test-results
          path: test/test_results_*.log
```

## 贡献指南

### 添加新测试
1. 在对应目录下创建新的测试文件
2. 遵循现有的命名规范
3. 添加详细的测试说明
4. 更新本README文档

### 测试文件格式
```http
### 测试名称
POST {{baseUrl}}/api/endpoint
Authorization: Bearer {{authToken}}
Content-Type: {{contentType}}

{
  "param1": "value1",
  "param2": "value2"
}
```

### 提交测试结果
1. 运行完整的测试套件
2. 记录测试结果和发现的问题
3. 提交测试报告和改进建议

## 联系方式

如有测试相关问题，请联系：
- 开发团队：dev@example.com
- 测试团队：qa@example.com
- 项目经理：pm@example.com