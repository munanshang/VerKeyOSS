# 完整 API 文档

本文档详细描述 VerKeyOSS 的所有 API 接口，包括接口路径、参数、响应格式及示例。

## 基础信息
- **基础路径**：`/api`
- **数据格式**：JSON
- **服务端口**：8913（默认）
- **默认管理员账号**：
  - 用户名：`verkeyoss`
  - 密码：`verkeyoss`
  - 首次登录后请立即修改密码
- **状态码说明**：
  - 200：成功
  - 400：请求参数错误
  - 401：未授权（如管理员接口未登录）
  - 404：资源不存在（如 AKey 或 VKey 无效）
  - 500：服务器内部错误

## API 分类说明

本API文档按照使用场景分为以下两个大类：
1. **管理接口**：需要管理员认证才能访问的系统管理功能，包括认证、应用管理、版本管理等
2. **应用调用接口**：第三方应用调用的合法性验证和更新检测功能

## 部署模式说明

VerKeyOSS 支持两种部署模式：

### 集成模式（Embed模式，推荐）
- 前端管理界面嵌入到 Go 应用中
- 单一可执行文件包含完整的前后端功能
- 访问地址：`http://localhost:8913`（Web 管理界面）
- API 访问地址：`http://localhost:8913/api`

### 分离模式
- 前后端分开部署
- 适用于开发环境或需要独立部署的场景

## 1. 管理接口

系统采用单管理员模式，使用config.yaml配置文件进行管理。所有管理接口都需要管理员认证才能访问。

> JWT令牌中包含admin:true声明，用于验证管理员身份，但此信息不会在API响应中返回给客户端。

### 1.1 登录
- **URL**：`/api/auth/login`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "username": "管理员用户名",
    "password": "管理员密码"
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "token": "登录令牌（用于后续管理接口）",
      "expires_at": "令牌过期时间（ISO 8601格式）",
      "user_info": {
        "username": "verkeyoss"
      }
    }
  }
  ```
- **失败响应**（401）：
  ```json
  {
    "code": 401,
    "message": "用户名或密码错误"
  }
  ```

### 1.2 修改密码
- **URL**：`/api/auth/password`
- **方法**：`PUT`
- **请求头**：`Authorization: Bearer {token}`（登录后获取的令牌）
- **请求体**：
   ```json
   {
     "old_password": "原密码",  // 必选
     "new_password": "新密码"  // 必选
   }
   ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "密码修改成功"
  }
  ```
- **失败响应**（401）：
  ```json
  {
    "code": 401,
    "message": "原密码错误或用户不存在"
  }
  ```

### 1.3 通过token获取用户信息
- **URL**：`/api/auth/user-info`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录后获取的令牌）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "username": "verkeyoss"
    }
  }
  ```
- **失败响应**（401）：
  ```json
  {
    "code": 401,
    "message": "未授权访问"
  }
  ```

### 1.5 应用管理接口

#### 1.5.1 创建应用

- **URL**: `/api/app`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: 
```json
{
  "name": "应用名称",  // 必选
  "description": "应用描述",  // 可选
  "is_paid": false  // 是否收费应用
}
```
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "akey": "应用唯一标识",
    "name": "应用名称",
    "description": "应用描述",
    "is_paid": false,  // 是否收费应用
    "created_at": "创建时间（ISO 8601格式）"
  }
}
```

#### 1.5.2 获取应用列表

- **URL**: `/api/app`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**: 
  - `page`: 页码，默认1
  - `size`: 每页数量，默认10
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "total": 总记录数,
    "list": [
      {
        "akey": "应用唯一标识",
        "name": "应用名称",
        "description": "应用描述",
        "is_paid": false,  // 是否收费应用
        "version_count": 版本数量,
        "created_at": "创建时间（ISO 8601格式）"
      }
      // 更多应用
    ]
  }
}
```

#### 1.5.3 更新应用

- **URL**: `/api/app/:akey`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `akey`: 应用唯一标识
- **请求体**: 
```json
{
  "name": "新应用名称",
  "description": "新应用描述",
  "is_paid": false  // 是否收费应用
}
```
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

#### 1.5.4 删除应用

- **URL**: `/api/app/:akey`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `akey`: 应用唯一标识
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

### 1.6 版本管理接口

#### 1.6.1 创建新版本

- **URL**: `/api/app/:akey/versions`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `akey`: 应用唯一标识
- **请求体**: 
```json
{
  "version": "版本号",  // 必选
  "description": "版本描述",  // 可选
  "is_latest": true,  // 是否为最新版本
  "is_forced_update": false  // 是否强制更新
}
```
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "vkey": "版本唯一标识",
    "akey": "应用唯一标识",
    "version": "版本号",
    "description": "版本描述",
    "is_latest": true,
    "is_forced_update": false,  // 是否强制更新
    "created_at": "创建时间（ISO 8601格式）"
  }
}
```

#### 1.6.2 获取版本列表

- **URL**: `/api/app/:akey/versions`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `akey`: 应用唯一标识
- **请求参数**: 
  - `page`: 页码，默认1
  - `size`: 每页数量，默认10
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "total": 总记录数,
    "list": [
      {
        "vkey": "版本唯一标识",
        "version": "版本号",
        "description": "版本描述",
        "is_latest": true,
        "is_forced_update": false,  // 是否强制更新
        "created_at": "创建时间（ISO 8601格式）"
      }
      // 更多版本
    ]
  }
}
```

#### 1.6.3 更新版本

- **URL**: `/api/versions/:vkey`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `vkey`: 版本唯一标识
- **请求体**: 
```json
{
  "description": "新版本描述",
  "is_latest": false,
  "is_forced_update": false  // 是否强制更新
}
```
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

#### 1.6.4 删除版本

- **URL**: `/api/versions/:vkey`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `vkey`: 版本唯一标识
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

### 1.7 仪表盘接口

#### 1.7.1 获取仪表盘数据

- **URL**: `/api/dashboard/stats`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "total_apps": 5,  // 应用总数
    "total_versions": 20,  // 版本总数
    "recent_apps": [  // 最近创建的应用（最多5个）
      {
        "akey": "应用唯一标识",
        "name": "应用名称",
        "description": "应用描述",
        "is_paid": false,  // 是否收费应用
        "version_count": 3,  // 版本数量
        "created_at": "创建时间（ISO 8601格式）"
      }
      // 更多应用...
    ],
    "recent_versions": [  // 最近创建的版本（最多5个）
      {
        "vkey": "版本唯一标识",
        "akey": "应用唯一标识",
        "version": "版本号",
        "description": "版本描述",
        "is_latest": true,
        "is_forced_update": false,  // 是否强制更新
        "created_at": "创建时间（ISO 8601格式）"
      }
      // 更多版本...
    ]
  }
}
```

#### 1.7.2 获取公告列表

- **URL**: `/api/dashboard/announcements`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **成功响应示例**: 
```json
{
  "code": 200,
  "data": [
    {
      "id": 公告ID,
      "title": "公告标题",
      "content": "公告内容",
      "created_at": "发布时间（ISO 8601格式）"
    }
    // 更多公告
  ]
}
```

## 3. 应用调用接口

以下接口主要用于第三方应用调用，提供应用合法性验证和更新检测功能。

### 3.1 校验 AKey 和 VKey 合法性（POST 方法）
- **URL**：`/api/check/validate`
- **方法**：`POST`（避免参数在 URL 中泄露）
- **请求体**：
  ```json
  {
    "akey": "应用唯一标识",  // 必选
    "vkey": "版本唯一标识"   // 必选
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "valid": true,  // 布尔值，是否合法
      "message": "AKey和VKey合法",  // 说明信息
      "app_name": "应用名称",  // 校验成功时返回应用名称
      "version": "版本号"  // 校验成功时返回版本号
    }
  }
  ```
- **失败响应**（404）：
  ```json
  {
    "code": 404,
    "data": {
      "valid": false,
      "message": "AKey不存在"  // 或 "VKey不存在"
    }
  }
  ```

### 3.2 检测是否有新版本（POST 方法）
- **URL**：`/api/check/update`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "akey": "应用唯一标识",  // 必选
    "vkey": "当前版本的VKey"  // 必选
  }
  ```
- **成功响应**（200，存在更新）：
  ```json
  {
    "code": 200,
    "data": {
      "has_update": true,
      "latest_version": "最新版本号",  // 仅返回公开的版本号
      "release_time": "最新版本发布时间"  // 仅返回公开的发布时间（ISO 8601格式）
    }
  }
  ```
- **成功响应**（200，无更新）：
  ```json
  {
    "code": 200,
    "data": {
      "has_update": false,
      "message": "当前已是最新版本"
    }
  }
  ```

### 3.3 健康检查
- **URL**：`/api/check/health`
- **方法**：`GET`
- **描述**：用于检查服务健康状态，不需要认证
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "OK",
    "data": {
      "status": "healthy",
      "service": "VerKeyOSS",
      "version": "1.0.0"
    }
  }
  ```

**响应字段说明：**
- `status`: 服务状态，始终为 "healthy"
- `service`: 服务名称，固定为 "VerKeyOSS"
- `version`: 后端版本号，由构建时注入

## 附录

### A. 错误代码对照表

| 状态码 | 说明 | 常见场景 |
|--------|------|----------|
| 200 | 成功 | 请求正常处理 |
| 400 | 请求参数错误 | JSON格式错误、必填字段缺失 |
| 401 | 未授权 | token无效、未登录、密码错误 |
| 404 | 资源不存在 | AKey/VKey无效、版本不存在 |
| 500 | 服务器内部错误 | 数据库连接失败、系统异常 |

### B. 身份验证

所有管理接口都需要在请求头中携带JWT令牌：
```
Authorization: Bearer {token}
```

令牌通过登录接口获取，有效期默认24小时。

### C. 分页参数

支持分页的接口都接受以下查询参数：
- `page`: 页码，从1开始，默认1
- `size`: 每页数量，默认10，最大100

### D. 时间格式

所有时间字段统一使用ISO 8601格式：`2024-01-01T12:00:00Z`

### E. 集成建议

1. **安全性**：
   - 妥善保管AKey和VKey，避免在日志中暴露
   - 使用HTTPS传输敏感数据
   - 定期轮换管理员密码

2. **性能优化**：
   - 对校验接口进行合理缓存
   - 使用连接池管理HTTP连接
   - 避免频繁调用管理接口

3. **错误处理**：
   - 实现重试机制（非业务错误）
   - 记录详细的错误日志
   - 对用户提供友好的错误提示

### F. 版本信息

- API版本：v1.0.0
- 文档版本：v1.0.0
- 最后更新：2024年
