# 测试目录说明

## 目录结构

```
test/
├── README.md                   # 测试说明文档
├── api/                        # API接口测试
│   ├── auth_test.go           # 认证相关测试
│   ├── user_test.go           # 用户相关测试
│   ├── mailbox_test.go        # 邮箱相关测试
│   ├── email_test.go          # 邮件相关测试
│   └── admin_test.go          # 管理员相关测试
├── unit/                       # 单元测试
│   ├── model_test.go          # 数据模型测试
│   ├── service_test.go        # 服务层测试
│   └── utils_test.go          # 工具函数测试
├── integration/                # 集成测试
│   ├── email_flow_test.go     # 邮件收发流程测试
│   ├── auth_flow_test.go      # 认证流程测试
│   └── admin_flow_test.go     # 管理员流程测试
├── data/                       # 测试数据
│   ├── test_config.yaml       # 测试配置文件
│   ├── test_emails.json       # 测试邮件数据
│   └── test_users.json        # 测试用户数据
├── fixtures/                   # 测试固件
│   ├── database.go            # 数据库测试固件
│   ├── server.go              # 服务器测试固件
│   └── auth.go                # 认证测试固件
└── scripts/                    # 测试脚本
    ├── setup.sh               # 测试环境设置
    ├── cleanup.sh             # 测试清理
    └── run_tests.sh           # 运行所有测试
```

## 测试规范

### 1. 命名规范
- 测试文件以 `_test.go` 结尾
- 测试函数以 `Test` 开头
- 基准测试以 `Benchmark` 开头
- 示例测试以 `Example` 开头

### 2. 测试分类
- **单元测试**: 测试单个函数或方法
- **集成测试**: 测试多个组件的交互
- **API测试**: 测试HTTP接口
- **端到端测试**: 测试完整的业务流程

### 3. 测试数据
- 使用独立的测试数据库
- 每个测试用例使用独立的数据
- 测试完成后清理数据

### 4. 运行测试
```bash
# 运行所有测试
go test ./test/...

# 运行特定包的测试
go test ./test/api/

# 运行特定测试
go test -run TestUserLogin ./test/api/

# 运行测试并显示覆盖率
go test -cover ./test/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

## 测试环境

### 数据库
- 使用SQLite内存数据库进行测试
- 每个测试用例创建独立的数据库实例

### 配置
- 使用专门的测试配置文件
- 禁用外部服务依赖（使用Mock）

### 日志
- 测试时使用简化的日志输出
- 错误日志保存到测试目录
