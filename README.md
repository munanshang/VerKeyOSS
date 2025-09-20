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

2. 构建前端和后端（集成版本）
   
   项目支持将前端嵌入到后端可执行文件中，生成单一的部署文件。提供了自动化构建脚本，支持版本化打包：
   
   **使用自动化构建脚本（推荐）：**
   
   构建脚本会提示输入版本号，并生成版本化的输出目录，同时构建Windows和Linux版本：
   
   ```bash
   # Windows系统
   .\build.bat
   
   # Linux/macOS系统
   chmod +x build.sh
   ./build.sh
   ```
   
   构建完成后，文件保存在 `bin/版本号/` 目录下：
   ```
   bin/
   └── 1.0.0/
       ├── verkeyoss-windows-amd64.exe
       └── verkeyoss-linux-amd64
   ```
   
   **手动构建：**
   
   **Windows系统：**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   go build -ldflags="-s -w" -o verkeyoss.exe .
   ```
   
   **Linux/macOS系统：**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   go build -ldflags="-s -w" -o verkeyoss .
   ```

3. 配置数据库
   - 新建 MySQL 数据库（如 `verkeyoss`）
   - 复制配置文件模板并修改数据库信息
     ```bash
     cp config.example.yaml config.yaml
     ```
   - 配置文件内容示例：
     ```yaml
     # 数据库配置
     db:
       host: localhost
       port: 3306
       user: verkeyoss
       password: verkeyoss
       name: verkeyoss

     # 服务器配置
     server:
       port: 8913  # API 服务端口
       debug: false  # 设置为true启用调试模式（允许所有域名访问，仅用于开发环境）

     # JWT配置
     jwt:
       secret: 32字符的随机密钥
       expire_hours: 24

     # 管理员配置（首次运行时系统会自动生成）
     admin:
       username: verkeyoss
       password: # 系统自动生成加密后的密码
     ```

3. 启动服务（数据库初始化会自动执行）
   ```bash
   # 直接运行可执行文件（内置前端）
   ./verkeyoss
   # 或 Windows下
   .\verkeyoss.exe
   
   # 或直接从源码运行
   go run main.go
   ```

**注意：** 集成版本将前端管理界面嵌入到Go应用中，单一可执行文件即包含了完整的前后端功能。访问 `http://localhost:8913` 即可使用Web管理界面。

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
├── frontend/          # 前端项目（Vue3 + TypeScript + Element Plus）
│   ├── src/           # 前端源代码
│   ├── dist/          # 前端构建产物（会嵌入到Go应用中）
│   ├── package.json   # 前端项目配置
│   └── vite.config.ts # Vite构建配置
├── docs/
│   ├── DEPLOY.md      # 部署文档
│   ├── api.md         # 完整API文档
│   └── images/
│       └── Logo.svg   # 项目Logo
├── BUILD.md           # 构建脚本使用说明
├── build.sh           # Linux/macOS构建脚本
├── build.bat          # Windows构建脚本
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