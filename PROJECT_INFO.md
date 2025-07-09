# AuthCenter 项目下载指南

## 📦 下载方式

### 方式一：Git Clone（推荐）
```bash
git clone https://github.com/Xie0s/authcenters.git
cd authcenters
```

### 方式二：下载压缩包
项目已打包为 `authcenters-project.tar.gz`（90KB），包含所有源代码。

解压方式：
```bash
tar -xzf authcenters-project.tar.gz
cd authcenters
```

## 🛠️ 系统要求

- **Go**: 1.19 或更高版本
- **Docker**: 用于运行 MongoDB
- **Node.js**: 用于数据库初始化脚本

## 🚀 部署步骤

### 1. 检查系统要求
```bash
./deploy-local.sh check
```

### 2. 一键部署
```bash
./deploy-local.sh
```

### 3. 分步部署
```bash
# 启动 MongoDB
./deploy-local.sh mongodb

# 初始化数据库
./deploy-local.sh init

# 构建项目
./deploy-local.sh build

# 启动服务
./deploy-local.sh start
```

## 🌐 访问地址

部署成功后，可以通过以下地址访问：

- **API 服务**: http://localhost:8080
- **测试页面**: http://localhost:8080/test/
- **健康检查**: http://localhost:8080/health

## 🔧 新功能说明

### 注册功能
- 用户名现在是必填字段
- 系统会验证用户名唯一性
- 支持邮箱和手机号（可选）

### 登录功能
支持4种登录方式：

1. **用户名登录**
   ```json
   {
     "username": "testuser",
     "password": "password123",
     "type": "username"
   }
   ```

2. **邮箱登录**
   ```json
   {
     "email": "test@example.com",
     "password": "password123",
     "type": "email"
   }
   ```

3. **自动识别登录**
   ```json
   {
     "username": "testuser",  // 或 "email": "test@example.com"
     "password": "password123",
     "type": "auto"
   }
   ```

4. **手机号登录**（原有功能）
   ```json
   {
     "phone": "13800138000",
     "code": "123456",
     "type": "phone"
   }
   ```

## 📁 项目结构

```
authcenters/
├── cmd/                    # 应用程序入口
├── configs/               # 配置文件
├── docs/                  # 文档
├── internal/              # 内部代码
│   ├── auth/             # 认证模块
│   ├── user/             # 用户管理
│   ├── role/             # 角色管理
│   ├── permission/       # 权限管理
│   └── ...
├── pkg/                   # 公共包
├── scripts/              # 数据库初始化脚本
├── test/                 # 测试页面
├── deploy-local.sh       # 一键部署脚本
├── go.mod               # Go 模块文件
└── README.md            # 项目说明
```

## 🧪 测试功能

### 使用测试页面
1. 访问 http://localhost:8080/test/
2. 使用界面进行注册和登录测试
3. 点击"一键测试所有功能"进行完整测试

### 使用 API 测试
```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# 用户名登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"testuser","password":"password123","type":"username"}'
```

## 🛠️ 故障排除

### MongoDB 连接问题
```bash
# 检查 MongoDB 容器状态
docker ps | grep mongodb

# 重启 MongoDB
docker restart mongodb

# 查看 MongoDB 日志
docker logs mongodb
```

### 端口占用问题
```bash
# 检查端口占用
ss -tlnp | grep :8080

# 修改配置文件中的端口
vim configs/config.yaml
```

### 数据库初始化问题
```bash
# 重新初始化数据库
cd scripts
node run_init_new.js
```

## 📞 技术支持

如果遇到问题，请检查：
1. 系统要求是否满足
2. MongoDB 是否正常运行
3. 端口 8080 是否被占用
4. 防火墙设置是否正确

## 🎯 功能特点

- ✅ 完整的 RBAC 权限管理
- ✅ JWT Token 认证
- ✅ 多种登录方式支持
- ✅ 企业级安全特性
- ✅ RESTful API 设计
- ✅ 完整的测试界面
- ✅ 一键部署脚本
- ✅ 详细的文档说明

现在您可以轻松地将项目部署到本地环境并开始使用新的用户名登录功能！
