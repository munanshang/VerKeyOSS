# 完整 API 文档

本文档详细描述 VerKeyOSS 的所有 API 接口，包括接口路径、参数、响应格式及示例。

## 基础信息
- 基础路径：`/api`
- 数据格式：JSON
- 状态码说明：
  - 200：成功
  - 400：请求参数错误
  - 401：未授权（如管理员接口未登录）
  - 404：资源不存在（如 AKey 或 VKey 无效）
  - 500：服务器内部错误

## API 分类说明

本API文档按照使用场景分为以下四个大类：
1. **验证接口**：用于用户认证、登录、密码修改等基础功能
2. **管理员接口**：需要管理员权限才能访问的系统管理功能
3. **用户接口**：前端页面使用的应用和版本管理功能
4. **应用调用接口**：第三方应用调用的合法性验证和更新检测功能

## 1. 验证接口

> **注意**：系统已从简单的管理员标识 (`is_admin`) 迁移到基于权限组的访问控制系统。用户通过 `group_id` 字段关联到相应的权限组。当前系统中，`group_id` 为 1 的用户具有管理员权限。

### 1.1 登录
- **URL**：`/auth/login`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "username": "用户名",
    "password": "密码"
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
        "user_id": 1,
        "username": "admin",
        "group_id": 1
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
- **URL**：`/auth/password`
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
- **URL**：`/auth/user-info`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录后获取的令牌）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "user_id": 1,
      "username": "admin",
      "group_id": 1
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

### 1.4 获取所有权限组信息
- **URL**：`/permission-groups`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录后获取的令牌）
- **访问权限**：所有登录用户
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": [
      {
        "id": 1,
        "group_name": "管理员",
        "permission": "{\"read\": true, \"write\": true, \"delete\": true}",
        "description": "系统管理员权限组",
        "created_at": "2023-05-15T10:00:00Z",
        "updated_at": "2023-05-15T10:00:00Z"
      },
      {
        "id": 2,
        "group_name": "普通用户",
        "permission": "{\"read\": true, \"write\": false, \"delete\": false}",
        "description": "普通用户权限组",
        "created_at": "2023-05-15T10:00:00Z",
        "updated_at": "2023-05-15T10:00:00Z"
      }
    ]
  }
  ```
- **失败响应**（401）：
  ```json
  {
    "code": 401,
    "message": "未授权访问"
  }
  ```

## 2. 管理员接口

以下接口需要管理员权限才能访问，用于系统管理和监控。

### 2.1 获取所有用户列表

- **URL**: `/api/admin/users`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**: 
  - `page`: 页码，默认1
  - `size`: 每页数量，默认10
- **访问权限**: 管理员

**成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "is_banned": false,
        "ban_reason": null,
        "group_id": 1,
        "permission_group": {
          "id": 1,
          "group_name": "管理员",
          "permission": "{\"read\": true, \"write\": true, \"delete\": true}",
          "description": "系统管理员权限组"
        },
        "created_at": "2023-05-15T10:00:00Z",
        "updated_at": "2023-05-15T10:00:00Z"
      }
      // 更多用户
    ],
    "page": 1,
    "size": 10
  }
}
```

### 2.2 封禁/解封用户

- **URL**: `/api/admin/users/:id/ban`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `id`: 用户ID
- **请求体**: 
```json
{
  "banned": true, // true为封禁，false为解封
  "ban_reason": "违规使用系统"
}
```
- **访问权限**: 管理员

**成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

### 2.3 获取所有应用列表

- **URL**: `/api/admin/apps`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**: 
  - `page`: 页码，默认1
  - `size`: 每页数量，默认10
- **访问权限**: 管理员

**成功响应示例**: 
```json
{
  "code": 200,
  "data": {
    "total": 50,
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "akey": "test",
        "name": "测试应用",
        "description": "这是一个用于测试的默认应用",
        "created_at": "2023-05-15T10:00:00Z",
        "version_count": 3,
        "is_banned": false,
        "ban_reason": null,
        "is_paid": false
      }
      // 更多应用
    ],
    "page": 1,
    "size": 10
  }
}
```

### 2.4 封禁/解封应用

- **URL**: `/api/admin/apps/:akey/ban`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**: 
  - `akey`: 应用唯一标识
- **请求体**: 
```json
{
  "banned": true, // true为封禁，false为解封
  "ban_reason": "存在安全隐患"
}
```
- **访问权限**: 管理员

**成功响应示例**: 
```json
{
  "code": 200,
  "data": "操作成功"
}
```

## 3. 用户接口

以下接口主要用于前端页面使用，提供应用和版本的管理功能。

### 3.1 应用管理

#### 3.1.1 创建应用
- **URL**：`/app`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：所有登录用户
- **请求体**：
  ```json
  {
    "name": "应用名称",  // 必选，字符串
    "description": "应用描述",  // 可选，字符串
    "is_paid": false  // 可选，布尔值，是否为付费应用（默认false）
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "akey": "应用唯一标识（自动生成）",
      "name": "应用名称",
      "description": "应用描述",
      "created_at": "创建时间（ISO 8601格式）",
      "user_id": 1,  // 创建者用户ID
      "is_banned": false,  // 是否被封禁
      "is_paid": false  // 是否为付费应用
    }
  }
  ```

#### 3.1.2 获取应用列表
- **URL**：`/app`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **查询参数**：`page=1&size=10`（分页参数，可选）
- **访问权限**：所有登录用户
- **成功响应**（200）：
  
  管理员用户响应示例：
  ```json
  {
    "code": 200,
    "data": {
      "total": 10,  // 总应用数
      "list": [
        {
          "akey": "soft_xxx",
          "name": "应用名称",
          "description": "应用描述",
          "created_at": "创建时间",
          "version_count": 3,
          "user_id": 1,  // 创建者用户ID
          "username": "admin",  // 创建者用户名
          "is_banned": false,  // 是否被封禁
          "is_paid": false  // 是否为付费应用
        }
        // ... 更多应用
      ]
    }
  }
  ```
  
  普通用户响应示例：
  ```json
  {
    "code": 200,
    "data": {
      "total": 5,  // 总应用数
      "list": [
        {
          "akey": "app_yyy",
          "name": "应用名称",
          "description": "应用描述",
          "created_at": "创建时间",
          "version_count": 2,
          "is_banned": false,  // 是否被封禁
          "is_paid": false  // 是否为付费应用
        }
        // ... 更多应用
      ]
    }
  }
  ```

> **说明**：应用列表接口根据用户角色返回不同内容。管理员用户可以看到应用的完整信息，包括创建者的用户ID和用户名；普通用户只能看到应用的基本信息，不包含任何用户相关信息。

#### 3.1.3 更新应用信息
- **URL**：`/app/{akey}`
- **方法**：`PUT`
- **路径参数**：`akey`（应用唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能更新自己创建的应用
  - 管理员用户：可以更新所有用户的应用
- **请求体**：
  ```json
  {
    "name": "新应用名称",  // 可选
    "description": "新应用描述",  // 可选
    "is_paid": false,  // 可选，布尔值，是否为付费应用
    "is_banned": false,  // 可选，布尔值，是否封禁应用
    "ban_reason": "封禁理由"  // 可选，字符串，封禁理由（仅在is_banned为true时有效）
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```

#### 3.1.4 删除应用
- **URL**：`/app/{akey}`
- **方法**：`DELETE`
- **路径参数**：`akey`（应用唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能删除自己创建的应用
  - 管理员用户：可以删除所有用户的应用
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "删除成功"
  }
  ```

### 3.2 版本管理

#### 3.2.1 发布新版本
- **URL**：`/app/{akey}/versions`
- **方法**：`POST`
- **路径参数**：`akey`（应用唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能为自己创建的应用发布新版本
  - 管理员用户：可以为所有用户的应用发布新版本
- **请求体**：
  ```json
  {
    "version": "版本号（如1.0.0、v2.1-beta）",  // 必选，字符串
    "description": "版本描述",  // 可选，字符串
    "is_latest": false,  // 可选，布尔值，是否标记为最新版本（默认自动更新最新）
    "is_forced_update": false  // 可选，布尔值，是否强制更新
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "vkey": "版本唯一标识（自动生成）",
      "akey": "应用唯一标识",
      "version": "版本号",
      "description": "版本描述",
      "is_latest": true,
      "is_forced_update": false,
      "created_at": "发布时间"
    }
  }
  ```

#### 3.2.2 获取版本列表
- **URL**：`/app/{akey}/versions`
- **方法**：`GET`
- **路径参数**：`akey`（应用唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能获取自己创建的应用的版本列表
  - 管理员用户：可以获取所有用户的应用的版本列表
- **查询参数**：`page=1&size=10`（分页参数，可选）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "total": 5,  // 总版本数
      "list": [
        {
          "vkey": "版本唯一标识",
          "version": "2.0.0",
          "description": "版本描述",
          "is_latest": true,
          "is_forced_update": false,
          "created_at": "发布时间"
        },
        // ... 更多版本
      ]
    }
  }
  ```

#### 3.2.3 更新版本信息
- **URL**：`/versions/{vkey}`
- **方法**：`PUT`
- **路径参数**：`vkey`（版本唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能更新自己创建的应用的版本信息
  - 管理员用户：可以更新所有用户的应用的版本信息
- **请求体**：
  ```json
  {
    "version": "新版本号",  // 可选，更新版本号
    "description": "新的版本描述",  // 可选，更新版本描述
    "is_latest": true,  // 可选，是否标记为最新版本
    "is_forced_update": false  // 可选，是否强制更新
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```

#### 3.2.4 删除版本
- **URL**：`/versions/{vkey}`
- **方法**：`DELETE`
- **路径参数**：`vkey`（版本唯一标识）
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：
  - 非管理员用户：只能删除自己创建的应用的版本
  - 管理员用户：可以删除所有用户的应用的版本
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "删除成功"
  }
  ```

### 3.3 仪表盘与公告

#### 3.3.1 获取仪表盘数据
- **URL**：`/dashboard`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：所有登录用户
- **成功响应**（200）：
  
  普通用户响应示例：
  ```json
  {
    "code": 200,
    "data": {
      "user_app_count": 5  // 用户创建的应用数量
    }
  }
  ```
  
  管理员用户响应示例：
  ```json
  {
    "code": 200,
    "data": {
      "user_app_count": 10,  // 用户创建的应用数量
      "total_users": 25,     // 系统总用户数
      "total_apps": 50       // 系统总应用数
    }
  }
  ```

#### 3.3.2 获取公告列表
- **URL**：`/dashboard/announcements`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（登录令牌）
- **访问权限**：所有登录用户
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": [
      {
          "id": 1,
          "title": "欢迎使用 VerKeyOSS",
          "content": "这是一条测试公告...",
          "is_active": true,
          "publish_date": "2023-05-15T10:00:00Z",
          "url": "https://github.com/munanshang/verkeyoss", // 可选字段，公告链接
          "created_at": "2023-05-15T10:00:00Z",
          "updated_at": "2023-05-15T10:00:00Z"
        }
      // 更多公告
    ]
  }
  ```

## 4. 应用调用接口

以下接口主要用于第三方应用调用，提供应用合法性验证和更新检测功能。

### 4.1 校验 AKey 和 VKey 合法性（POST 方法）
- **URL**：`/check/legality`
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
      "legal": true,  // 布尔值，是否合法
      "message": "AKey和VKey合法"  // 说明信息
    }
  }
  ```
- **失败响应**（404）：
  ```json
  {
    "code": 404,
    "data": {
      "legal": false,
      "message": "AKey不存在"  // 或 "VKey不存在"
    }
  }
  ```

### 4.2 检测是否有新版本（POST 方法）
- **URL**：`/check/update`
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
      "release_time": "最新版本发布时间"  // 仅返回公开的发布时间
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
