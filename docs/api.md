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

## 1. 管理员接口

### 1.1 管理员登录
- **URL**：`/admin/login`
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
      "token": "登录令牌（用于后续管理员接口）",
      "expires_at": "令牌过期时间（ISO 8601格式）"
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

### 1.2 修改管理员密码
- **URL**：`/admin/password`
- **方法**：`PUT`
- **请求头**：`Authorization: Bearer {token}`（登录后获取的令牌）
- **请求体**：
  ```json
  {
    "old_password": "原密码",
    "new_password": "新密码"
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "密码修改成功"
  }
  ```

## 2. 软件管理 API

### 2.1 创建软件
- **URL**：`/software`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **请求体**：
  ```json
  {
    "name": "软件名称",  // 必选，字符串
    "description": "软件描述"  // 可选，字符串
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "akey": "软件唯一标识（自动生成）",
      "name": "软件名称",
      "description": "软件描述",
      "created_at": "创建时间（ISO 8601格式）"
    }
  }
  ```

### 2.2 获取软件列表
- **URL**：`/software`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **查询参数**：`page=1&size=10`（分页参数，可选）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "total": 10,  // 总软件数
      "list": [
        {
          "akey": "soft_xxx",
          "name": "软件名称",
          "description": "软件描述",
          "created_at": "创建时间",
          "version_count": 3
        }
        // ... 更多软件
      ]
    }
  }
  ```

> **说明**：软件列表接口现在始终包含软件描述字段，不再需要单独的软件信息接口。

### 2.3 更新软件信息
- **URL**：`/software/{akey}`
- **方法**：`PUT`
- **路径参数**：`akey`（软件唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **请求体**：
  ```json
  {
    "name": "新软件名称",  // 可选
    "description": "新软件描述"  // 可选
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```

### 2.4 删除软件
- **URL**：`/software/{akey}`
- **方法**：`DELETE`
- **路径参数**：`akey`（软件唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "删除成功"
  }
  ```

## 3. 版本管理 API

### 3.1 发布新版本
- **URL**：`/software/{akey}/versions`
- **方法**：`POST`
- **路径参数**：`akey`（软件唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **请求体**：
  ```json
  {
    "version": "版本号（如1.0.0、v2.1-beta）",  // 必选，字符串
    "description": "版本描述",  // 可选，字符串
    "is_latest": false  // 可选，布尔值，是否标记为最新版本（默认自动更新最新）
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "data": {
      "vkey": "版本唯一标识（自动生成）",
      "akey": "软件唯一标识",
      "version": "版本号",
      "description": "版本描述",
      "is_latest": true,
      "created_at": "发布时间"
    }
  }
  ```

### 3.2 获取版本列表
- **URL**：`/software/{akey}/versions`
- **方法**：`GET`
- **路径参数**：`akey`（软件唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
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
          "created_at": "发布时间"
        },
        // ... 更多版本
      ]
    }
  }
  ```

### 3.3 更新版本信息
- **URL**：`/versions/{vkey}`
- **方法**：`PUT`
- **路径参数**：`vkey`（版本唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **请求体**：
  ```json
  {
    "version": "新版本号",  // 可选，更新版本号
    "description": "新的版本描述",  // 可选，更新版本描述
    "is_latest": true  // 可选，是否标记为最新版本
  }
  ```
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```

### 3.4 删除版本
- **URL**：`/versions/{vkey}`
- **方法**：`DELETE`
- **路径参数**：`vkey`（版本唯一标识）
- **请求头**：`Authorization: Bearer {token}`（管理员令牌）
- **成功响应**（200）：
  ```json
  {
    "code": 200,
    "message": "删除成功"
  }
  ```

## 4. 合法性与更新检测 API

### 4.1 校验 AKey 和 VKey 合法性（POST 方法）
- **URL**：`/check/legality`
- **方法**：`POST`（避免参数在 URL 中泄露）
- **请求体**：
  ```json
  {
    "akey": "软件唯一标识",  // 必选
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
    "akey": "软件唯一标识",  // 必选
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
