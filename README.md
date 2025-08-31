<div align="center">
  <img src="docs/images/Logo.svg" alt="Logo" width="180" height="180">
</div>

<h3 align="center">VerKeyOSS</h3>

一款开源的软件版本管理系统，用于管理软件标识（AKey）和版本标识（VKey），支持合法性校验和新版本检测，适用于各类软件的版本管控场景。

## 核心功能

- 创建软件并生成唯一标识 `AKey`（软件唯一标识，保密）
- 为软件发布版本并生成唯一标识 `VKey`（版本唯一标识，保密）
- 管理软件信息（名称、描述等）和版本信息（版本号、发布时间等）
- 通过 API 校验 `AKey` 和 `VKey` 的合法性（基于 POST 方法，避免参数泄露）
- 检测当前版本是否存在更新（仅返回公开的版本号和发布时间）

## 开源协议

本项目采用 **GNU Affero General Public License v3.0 (AGPLv3)** 协议开源：
- 允许二次开发、修改和分发，但必须保留原作者署名
- 任何基于本项目的修改或衍生作品必须以 AGPLv3 协议开源
- 若将本项目或其衍生作品作为网络服务提供（如部署为在线 API），必须公开完整源代码
- 协议全文见 [LICENSE](LICENSE)（[官方原文](https://www.gnu.org/licenses/agpl-3.0.txt)）

## 管理员账号

系统默认内置管理员账号：
- 用户名：`verkey`
- 密码：`verkey`
- 首次登录后请立即修改密码（通过管理员接口或前端界面）

## 安装部署

### 环境要求
- Go 1.19+
- MySQL 8.0+（或兼容的关系型数据库）

### 部署步骤

1. 克隆代码库
   ```bash
   git clone https://github.com/munanshang/verkeyoss.git
   cd verkeyoss
   ```

2. 配置数据库
   - 新建 MySQL 数据库（如 `verkeyoss`）
   - 复制配置文件模板并修改数据库信息
     ```bash
     cp config.example.yaml config.yaml
     ```
   - 配置文件内容示例：
     ```yaml
     db:
       host: "localhost"
       port: 3306
       user: "verkeyoss"
       password: "verkeyoss"
       name: "verkeyoss"
     server:
       port: 8080  # API 服务端口
     ```

3. 初始化数据库
   ```bash
   go run cmd/migrate/main.go  # 执行数据库迁移，创建表结构
   ```

4. 启动服务
   ```bash
   go run cmd/server/main.go
   # 或编译为二进制文件
   go build -o verkeyoss cmd/server/main.go
   ./verkeyoss
   ```

## 快速使用示例

以下示例基于本地部署的服务（`http://localhost:8080`），所有接口路径均以 `/api` 开头。

### 1. 管理员登录（获取令牌）
- **URL**：`/api/admin/login`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "username": "管理员用户名",
    "password": "管理员密码"
  }
  ```
- **响应**：
  
  ```json
  {
    "code": 200,
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2023-10-08T12:00:00Z"
    }
  }
  ```

### 2. 创建软件（获取 AKey）
- **URL**：`/api/software`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer {admin_token}`  // 替换为登录获取的token
- **请求体**：
  ```json
  {
    "name": "我的软件",
    "description": "这是一个测试软件"
  }
  ```
- **响应**：
  ```json
  {
    "code": 200,
    "data": {
      "akey": "soft_5f8d888a-7d8f-4e3b-9c1a-1b7d6e3c8d1b",
      "name": "我的软件",
      "description": "这是一个测试软件",
      "created_at": "2023-10-01T12:00:00Z"
    }
  }
  ```

### 3. 发布版本（获取 VKey）
- **URL**：`/api/software/{akey}/versions`  // 替换{akey}为实际软件标识
- **方法**：`POST`
- **请求头**：`Authorization: Bearer {admin_token}`
- **请求体**：
  
  ```json
  {
    "version": "1.0.0",
    "description": "初始版本"
  }
  ```
- **响应**：
  ```json
  {
    "code": 200,
    "data": {
      "vkey": "ver_7e6c5d4b-8c7d-4f2e-7b6a-2c5d4e3f2a1b",
      "akey": "soft_5f8d888a-7d8f-4e3b-9c1a-1b7d6e3c8d1b",
      "version": "1.0.0",
      "description": "初始版本",
      "is_latest": true,
      "created_at": "2023-10-01T12:30:00Z"
    }
  }
  ```

### 4. 校验合法性
- **URL**：`/api/check/legality`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "akey": "soft_5f8d888a-7d8f-4e3b-9c1a-1b7d6e3c8d1b",
    "vkey": "ver_7e6c5d4b-8c7d-4f2e-7b6a-2c5d4e3f2a1b"
  }
  ```
- **响应**：
  ```json
  {
    "code": 200,
    "data": {
      "legal": true,
      "message": "AKey和VKey合法"
    }
  }
  ```

### 5. 检测是否有新版本
- **URL**：`/api/check/update`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "akey": "soft_5f8d888a-7d8f-4e3b-9c1a-1b7d6e3c8d1b",
    "vkey": "ver_7e6c5d4b-8c7d-4f2e-7b6a-2c5d4e3f2a1b"
  }
  ```
- **响应（若存在新版本）**：
  
  ```json
  {
    "code": 200,
    "data": {
      "has_update": true,
      "latest_version": "2.0.0",
      "release_time": "2023-10-02T10:00:00Z"
    }
  }
  ```

## 完整 API 文档

查看完整接口说明：[API 文档](docs/api.md)

## 项目结构
```
verkeyoss/
├── cmd/
│   ├── server/        # 服务相关逻辑
│   └── migrate/       # 数据库迁移工具
├── internal/
│   ├── api/           # API 处理器（路由和请求处理）
│   ├── model/         # 数据模型（结构体定义）
│   ├── service/       # 业务逻辑层
│   └── store/         # 数据库操作层（与数据库交互）
├── templates/         # 前端模板及静态资源
│   ├── public/        # 静态文件（CSS、JS、图片等）
│   ├── admin/         # 管理员界面模板
│   │   └── index.html # 管理员首页
│   └── index.html     # 前端入口页面
├── docs/
│   └── api.md         # 完整API文档
├── config.example.yaml # 配置模板（需复制为config.yaml使用）
├── LICENSE            # AGPLv3 协议文本
└── README.md          # 项目说明文档
```

## 贡献指南

1. Fork 本仓库
2. 创建特性分支（`git checkout -b feature/xxx`）
3. 提交代码（`git commit -m "add: 新增xxx功能"`）
4. 推送分支（`git push origin feature/xxx`）
5. 提交 Pull Request

## 联系方式

- 项目地址：https://github.com/munanshang/verkeyoss
- 问题反馈：提交 Issue 至本仓库    