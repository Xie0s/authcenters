# RBAC企业级认证授权系统 (AuthCenter)

基于Go语言和MongoDB的企业级认证授权中心，采用RBAC(Role-Based Access Control)模型。

## 项目结构

```
AuthCenter/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口
├── internal/                     # 内部应用代码（不对外暴露）
│   ├── config/                  # 配置管理
│   │   └── config.go           # 配置结构和加载逻辑
│   ├── database/               # 数据库连接和管理
│   │   ├── connection.go       # MongoDB连接
│   │   └── indexes.go          # 数据库索引创建
│   ├── models/                 # 数据模型
│   │   └── models.go           # 所有数据模型定义
│   ├── auth/                   # 认证服务模块
│   │   ├── handler/            # HTTP请求处理器
│   │   │   └── auth_handler.go
│   │   ├── service/            # 业务逻辑层
│   │   │   └── auth_service.go
│   │   └── repository/         # 数据访问层
│   ├── user/                   # 用户管理模块
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── role/                   # 角色管理模块
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── permission/             # 权限管理模块
│   ├── category/               # 分类管理模块（知识库分类）
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── tag/                    # 标签管理模块
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── ai/                     # AI助手服务模块
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── middleware/             # 中间件
│   │   └── auth.go            # 认证中间件
│   └── router/                 # 路由配置
│       └── router.go          # 主路由设置
├── pkg/                        # 可复用的公共包
│   ├── jwt/                   # JWT工具包
│   │   └── jwt.go
│   ├── utils/                 # 工具函数
│   │   └── password.go        # 密码加密工具
│   ├── response/              # HTTP响应工具
│   │   └── response.go
│   └── logger/                # 日志工具
│       └── logger.go
├── configs/                   # 配置文件
│   └── config.yaml           # 主配置文件
├── scripts/                  # 脚本文件
│   ├── init_db.sh           # MongoDB初始化脚本(Linux/Mac)
│   └── init_db.bat          # MongoDB初始化脚本(Windows)
├── test/                     # 测试文件
│   └── main_test.go
├── docs/                     # 文档
│   └── RBAC企业级认证授权系统.md  # 需求文档
├── go.mod                    # Go模块文件
├── go.sum                    # 依赖校验文件
├── config.yaml              # 配置文件
├── start.sh                 # 启动脚本
└── README.md                # 项目说明
```

## 核心功能模块

### 1. 认证服务 (Authentication)
- 用户注册/登录
- JWT Token生成和验证
- 多种登录方式支持（手机验证码、邮箱密码、第三方OAuth）
- Token刷新机制
- 会话管理

### 2. 授权服务 (Authorization)
- 基于RBAC的权限控制
- 角色和权限管理
- 细粒度权限验证
- 中间件级别的权限控制

### 3. 用户管理
- 用户CRUD操作
- 用户角色分配
- 用户状态管理

### 4. 分类管理
- 层级式分类结构（"书架上的格子"）
- 分类树的CRUD操作
- 文档分类关联

### 5. 标签管理
- 灵活的标签系统（"便利贴"）
- 标签的创建、编辑、删除
- 标签使用统计和推荐

### 6. AI助手服务
- AI对话会话管理
- 消息历史记录
- 权限控制的AI功能访问

## 数据库设计

使用MongoDB作为主数据库，包含以下主要集合：

- **users**: 用户信息
- **roles**: 角色定义
- **permissions**: 权限定义
- **categories**: 分类信息（层级结构）
- **tags**: 标签信息
- **knowledge_documents**: 知识库文档
- **sessions**: 用户会话管理
- **ai_sessions**: AI助手会话
- **ai_messages**: AI助手消息记录

## 技术栈

- **语言**: Go 1.19+
- **Web框架**: Gin
- **数据库**: MongoDB 6.0+
- **认证**: JWT
- **配置**: Viper
- **日志**: Logrus
- **加密**: bcrypt

## 快速开始

### 1. 环境要求
- Go 1.19+
- MongoDB 6.0+

### 2. 环境变量设置
```bash
export MONGODB_URI="mongodb://localhost:27017/auth_center"
export JWT_SECRET="your-secret-key"
```

### 3. 安装依赖
```bash
go mod tidy
```

### 4. 初始化数据库
```bash
# Linux/Mac
chmod +x scripts/init_db.sh
./scripts/init_db.sh

# Windows
scripts/init_db.bat
```

### 5. 启动服务
```bash
go run cmd/server/main.go
```

## API文档

启动服务后，访问以下URL获取API文档：
- 健康检查: `GET /health`
- API前缀: `/api/v1`

### 主要API端点

#### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新Token
- `POST /api/v1/auth/verify` - 验证Token
- `POST /api/v1/auth/logout` - 用户登出

#### 用户管理
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/{id}` - 获取用户详情
- `PUT /api/v1/users/{id}` - 更新用户信息
- `DELETE /api/v1/users/{id}` - 删除用户

#### 角色管理
- `GET /api/v1/roles` - 获取角色列表
- `POST /api/v1/roles` - 创建角色
- `PUT /api/v1/roles/{id}` - 更新角色
- `DELETE /api/v1/roles/{id}` - 删除角色

#### AI助手
- `POST /api/v1/ai/chat` - AI对话
- `GET /api/v1/ai/sessions` - 获取会话列表
- `GET /api/v1/ai/sessions/{session_id}` - 获取会话详情

## 权限系统

### 内置角色
- **Admin**: 系统管理员，拥有所有权限
- **Editor**: 内容管理员，负责知识库内容管理
- **Author**: 内容创作者，专注内容创作
- **User**: 普通用户，基础功能使用

### 权限分类
- **知识库内容权限**: READ, CREATE, UPDATE, DELETE, PUBLISH, APPROVE
- **系统管理权限**: USER_MANAGE, ROLE_MANAGE, CATEGORY_MANAGE, SYSTEM_CONFIG
- **内容组织权限**: TAG_CREATE, TAG_MANAGE
- **交互功能权限**: COMMENT, FAVORITE, SEARCH, AI_ASSISTANT

## 安全特性

- bcrypt密码加密（cost factor ≥ 12）
- JWT访问令牌和刷新令牌机制
- 会话管理和Token吊销
- 登录失败次数限制
- 权限中间件保护
- HTTPS强制传输

## 性能优化

- MongoDB索引优化
- 连接池管理
- JWT无状态认证
- 权限信息缓存
- TTL自动清理过期数据

## 部署说明

### Docker部署
```bash
# 构建镜像
docker build -t authcenter .

# 运行容器
docker run -d -p 8080:8080 --name authcenter \
  -e MONGODB_URI="mongodb://mongodb:27017/auth_center" \
  -e JWT_SECRET="your-secret-key" \
  authcenter
```

### 生产环境配置
- 使用环境变量管理敏感配置
- 配置MongoDB副本集
- 启用HTTPS
- 配置负载均衡
- 监控和日志收集

## 扩展开发

该项目采用模块化设计，支持以下扩展：

1. **新增业务模块**: 在`internal/`下创建新的业务模块
2. **自定义中间件**: 在`internal/middleware/`中添加新的中间件
3. **数据库扩展**: 在`internal/database/`中扩展数据库操作
4. **工具包扩展**: 在`pkg/`中添加可复用的工具包

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交变更
4. 创建Pull Request

## 许可证

本项目采用MIT许可证。
