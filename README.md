<div align="center">
  <img src="docs/images/Logo.svg" alt="Logo" width="180" height="180">
</div>

<h3 align="center">VerKeyOSS</h3>

一款开源的应用版本管理系统，用于管理应用标识（AKey）和版本标识（VKey），支持合法性校验和新版本检测，适用于各类应用的版本管控场景。

## 核心功能

- 创建应用并生成唯一标识 `AKey`（应用唯一标识，保密）
- 为应用发布版本并生成唯一标识 `VKey`（版本唯一标识，保密）
- 管理应用信息（名称、描述等）和版本信息（版本号、发布时间等）
- 付费应用支持：区分免费应用和付费应用
- 强制更新功能：版本发布时可设置是否强制用户更新
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
- 用户名：`verkeyoss`
- 登录密码：`verkeyoss`
- 首次登录后请立即修改密码（通过认证接口或前端界面）

## 安装部署

### 环境要求
- Go 1.19+
- MySQL 8.0+（或兼容的关系型数据库）

### 部署步骤

项目提供两种部署方式，您可以根据需要选择适合的方式。

#### 方式一：使用编译好的版本（推荐）

  [详细部署文档](docs/DEPLOY.md)

#### 方式二：从源码编译部署

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

3. 启动服务（数据库初始化会自动执行）
   ```bash
   go run main.go
   # 或编译为二进制文件
   go build -o verkeyoss main.go
   ./verkeyoss
   ```

## 完整 API 文档

查看完整接口说明：[API 文档](docs/api.md)

## 项目结构
```
VerKeyOSS/
├── internal/
│   ├── api/           # API 处理器（路由和请求处理）
│   ├── initializer/   # 数据库初始化程序
│   ├── model/         # 数据模型（结构体定义）
│   ├── router/        # 路由配置
│   ├── service/       # 业务逻辑层
│   └── store/         # 数据库操作层（与数据库交互）
├── docs/
│   ├── DEPLOY.md      # 部署文档
│   ├── api.md         # 完整API文档
│   └── images/
│       └── Logo.svg   # 项目Logo
├── config.example.yaml # 配置模板（需复制为config.yaml使用）
├── LICENSE            # AGPLv3 协议文本
├── README.md          # 项目说明文档
├── go.mod             # Go模块定义
├── go.sum             # 依赖版本锁定
└── main.go            # 程序入口
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