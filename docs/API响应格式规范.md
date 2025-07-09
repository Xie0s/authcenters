# API响应格式规范

## 概述

本文档定义了RBAC企业级认证授权系统的API响应格式标准，确保前端和后端的响应格式一致性。

## 响应格式标准

### 1. 后端响应格式

后端使用统一的响应结构，定义在 `pkg/response/response.go` 中：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

#### 成功响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

#### 错误响应示例
```json
{
  "code": 401,
  "message": "Token刷新失败",
  "error": "无效的Refresh Token"
}
```

### 2. 前端响应格式

前端的 `apiRequest` 函数将后端响应包装成统一格式：

```javascript
{
  status: httpStatusCode,  // HTTP状态码
  data: backendResponse    // 后端的完整响应
}
```

#### 成功响应示例
```json
{
  "status": 200,
  "data": {
    "code": 200,
    "message": "success",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 3600
    }
  }
}
```

#### 错误响应示例
```json
{
  "status": 401,
  "data": {
    "code": 401,
    "message": "Token刷新失败",
    "error": "无效的Refresh Token"
  }
}
```

## 状态码规范

### HTTP状态码
- `200` - 请求成功
- `400` - 请求参数错误
- `401` - 未授权/认证失败
- `403` - 权限不足
- `404` - 资源不存在
- `429` - 请求频率限制
- `500` - 服务器内部错误

### 业务状态码
业务状态码与HTTP状态码保持一致，便于前端处理。

## API端点响应格式

### 1. 认证相关API

#### 用户登录 `POST /api/v1/auth/login`

**成功响应 (200)**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "string",
    "refresh_token": "string",
    "expires_in": 3600,
    "user": {
      "id": "string",
      "username": "string",
      "roles": ["string"]
    }
  }
}
```

**失败响应 (401)**
```json
{
  "code": 401,
  "message": "登录失败",
  "error": "用户名或密码错误"
}
```

#### Token刷新 `POST /api/v1/auth/refresh`

**成功响应 (200)**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "string",
    "refresh_token": "string",
    "expires_in": 3600
  }
}
```

**失败响应 (401)**
```json
{
  "code": 401,
  "message": "Token刷新失败",
  "error": "无效的Refresh Token"
}
```

#### Token验证 `POST /api/v1/auth/verify`

**成功响应 (200)**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "valid": true,
    "user_id": "string",
    "username": "string",
    "roles": ["string"],
    "permissions": ["string"],
    "has_access": true
  }
}
```

**失败响应 (401)**
```json
{
  "code": 401,
  "message": "Token验证失败",
  "error": "Token已过期"
}
```

### 2. 用户管理API

#### 获取用户列表 `GET /api/v1/users`

**成功响应 (200)**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": [
      {
        "id": "string",
        "username": "string",
        "email": "string",
        "status": "active",
        "roles": ["string"],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

## 前端处理规范

### 1. 响应处理
```javascript
// 正确的响应处理方式
const response = await authAPI.login(username, password);

if (response.status === 200) {
    // 成功处理
    const tokenData = response.data.data;
    console.log('登录成功:', tokenData);
} else {
    // 错误处理
    const errorInfo = response.data;
    console.error('登录失败:', errorInfo.message, errorInfo.error);
}
```

### 2. 错误处理
```javascript
// 统一的错误处理函数
function handleApiError(response) {
    const errorData = response.data;
    
    switch (response.status) {
        case 400:
            showMessage('参数错误: ' + errorData.message, 'warning');
            break;
        case 401:
            showMessage('认证失败: ' + errorData.message, 'error');
            // 可能需要重新登录
            redirectToLogin();
            break;
        case 403:
            showMessage('权限不足: ' + errorData.message, 'error');
            break;
        case 500:
            showMessage('服务器错误: ' + errorData.message, 'error');
            break;
        default:
            showMessage('未知错误: ' + errorData.message, 'error');
    }
}
```

## 测试验证

使用 `test/response_format_test.html` 页面可以测试响应格式的一致性：

1. 打开测试页面
2. 点击各个测试按钮
3. 检查响应格式是否符合规范
4. 验证错误处理是否正确

## 注意事项

1. **格式一致性**: 所有API响应必须遵循统一格式
2. **错误信息**: 错误信息应该清晰明确，便于调试
3. **状态码**: HTTP状态码和业务状态码保持一致
4. **向后兼容**: 新增字段时保持向后兼容性
5. **文档更新**: API变更时及时更新文档
