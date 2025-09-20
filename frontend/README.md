# VerKeyOSS Frontend

VerKeyOSS 应用版本管理系统的前端项目，基于 Vue 3 + TypeScript + Element Plus 构建。

## 项目特性

- 🚀 基于 Vue 3 Composition API
- 📱 响应式设计，支持移动端
- 🎨 现代化 UI 设计，使用 Element Plus 组件库
- 🔐 完整的认证和权限管理
- 📊 数据可视化仪表盘
- 🛠️ 应用和版本管理
- ✅ AKey/VKey 校验工具
- 🔄 版本更新检测
- 💾 状态管理（Pinia）
- 🌐 路由守卫和导航管理

## 技术栈

- **框架**: Vue 3
- **语言**: TypeScript
- **构建工具**: Vite
- **UI 库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **日期处理**: Day.js
- **图标**: @element-plus/icons-vue

## 项目结构

```
src/
├── api/                    # API 接口封装
├── assets/                 # 静态资源
│   ├── icons/             # 图标文件
│   └── styles/            # 全局样式
├── components/            # 通用组件
│   └── common/           # 公共组件
├── layouts/               # 布局组件
├── router/                # 路由配置
├── stores/                # 状态管理
├── types/                 # TypeScript 类型定义
├── utils/                 # 工具函数
├── views/                 # 页面组件
│   ├── auth/             # 认证相关页面
│   ├── dashboard/        # 仪表盘
│   ├── apps/             # 应用管理
│   ├── versions/         # 版本管理
│   └── tools/            # 校验工具
├── App.vue               # 根组件
└── main.ts              # 入口文件
```

## 功能模块

### 1. 用户认证
- 管理员登录
- Token 管理
- 路由权限控制
- 密码修改

### 2. 仪表盘
- 统计数据展示
- 最近应用和版本
- 系统公告
- 快速操作入口

### 3. 应用管理
- 应用列表查看
- 创建新应用
- 编辑应用信息
- 删除应用
- 应用详情查看

### 4. 版本管理
- 版本列表管理
- 创建新版本
- 编辑版本信息
- 删除版本
- 设置最新版本
- 强制更新标记

### 5. 校验工具
- AKey/VKey 合法性校验
- 版本更新检测
- API 使用说明
- 集成示例

## 开发命令

### 安装依赖
```bash
npm install
```

### 开发模式
```bash
npm run dev
```

### 类型检查
```bash
npm run type-check
```

### 构建正式版本
```bash
npm run build
```

### 预览构建结果
```bash
npm run preview
```

## 环境配置

### 开发环境
项目默认连接到 `http://localhost:8913` 的后端 API。

如需修改 API 地址，可以设置环境变量 `VITE_API_BASE_URL` 或 `VITE_BACKEND_URL`，或编辑 `src/utils/request.ts` 文件中的配置。

### 正式环境
在正式环境下（embed模式），前端会被嵌入到Go应用中，单一可执行文件即包含了完整的前后端功能。

构建时请使用：
```bash
npm run build
```

构建产物会生成在 `dist/` 目录中，然后被Go应用通过embed功能嵌入。

## API 接口

项目与后端的接口交互包括：

- **认证接口** (`/api/auth/*`)
  - POST `/api/auth/login` - 用户登录
  - GET `/api/auth/user-info` - 获取用户信息
  - PUT `/api/auth/password` - 修改密码

- **应用管理** (`/api/app/*`)
  - GET `/api/app` - 获取应用列表
  - POST `/api/app` - 创建应用
  - PUT `/api/app/:akey` - 更新应用
  - DELETE `/api/app/:akey` - 删除应用

- **版本管理** (`/api/app/:akey/versions/*`, `/api/versions/*`)
  - GET `/api/app/:akey/versions` - 获取版本列表
  - POST `/api/app/:akey/versions` - 创建版本
  - PUT `/api/versions/:vkey` - 更新版本
  - DELETE `/api/versions/:vkey` - 删除版本

- **校验接口** (`/api/check/*`)
  - POST `/api/check/validate` - 校验 AKey/VKey
  - POST `/api/check/update` - 检查更新

- **仪表盘** (`/api/dashboard/*`)
  - GET `/api/dashboard/stats` - 获取统计数据
  - GET `/api/dashboard/announcements` - 获取公告

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 开发说明

### 代码规范
- 使用 TypeScript 进行类型约束
- 遵循 Vue 3 Composition API 最佳实践
- 组件采用 `<script setup>` 语法
- CSS 使用 scoped 样式

### 状态管理
使用 Pinia 进行状态管理，主要包括：
- `useAuthStore` - 认证状态
- `useAppStore` - 应用和版本数据
- `useDashboardStore` - 仪表盘数据

### 路由设计
- 公开路由：`/login`
- 受保护路由：`/dashboard`, `/apps`, `/tools` 等
- 路由守卫自动处理认证检查

## 部署说明

1. 构建项目：
   ```bash
   npm run build
   ```

2. 将 `dist` 目录中的文件部署到 Web 服务器

3. 配置 Web 服务器支持 SPA 路由（History 模式）

4. 确保后端 API 的 CORS 配置正确

## License

本项目采用 AGPLv3 开源协议。