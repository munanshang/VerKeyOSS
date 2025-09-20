# 部署文档

## 部署步骤

发布版本为预编译的二进制文件，无需安装Go环境和编译过程，直接按照以下步骤操作即可使用。

### 步骤 1：配置文件设置

需要将示例配置文件复制并重命名为 `config.yaml`，并根据您的环境修改配置内容。

#### Windows 系统

1. **使用文件资源管理器**（推荐）：
   - 右键点击 `config.example.yaml` 文件，选择 "复制" 
   - 在同一文件夹中右键点击空白处，选择 "粘贴" 
   - 将新生成的文件重命名为 `config.yaml`

#### Linux 系统

1. **使用命令行**：
   - 打开终端
   - 切换到程序所在目录
   - 执行以下命令：
     ```bash
     cp config.example.yaml config.yaml
     ```

### 步骤 2：编辑配置文件

使用文本编辑器打开 `config.yaml` 文件，修改以下内容：

```yaml
# 数据库配置
db:
  host: localhost      # 数据库主机地址
  port: 3306             # 数据库端口
  user: verkeyoss      # 数据库用户名
  password: verkeyoss  # 数据库密码
  name: verkeyoss      # 数据库名称

# 服务器配置
server:
  port: 8913             # API 服务端口
  debug: false           # 设置为true启用调试模式（允许所有域名访问，仅用于开发环境）

# JWT 配置
jwt:
  secret: 32字符的随机密钥  # JWT 密钥，正式环境中请使用强密钥
  expire_hours: 24        # token 过期时间（小时）

# 管理员配置（首次运行时系统会自动生成）
admin:
  username: verkeyoss
  password: # 系统自动生成加密后的密码
```

**重要提示**：
- 确保您已经在 MySQL 数据库中创建了名为 `verkeyoss` 的数据库（或配置文件中指定的数据库名称）
- 正式环境中，请务必修改 `jwt.secret` 为一个强密钥
- 管理员配置由系统自动生成，默认用户名和密码都是 `verkeyoss`
- 调试模式下会允许所有域名访问，正式环境中必须设置 `debug: false`

**配置自动生成说明**：
如果配置文件不存在，系统在首次启动时会自动创建一个包含默认值的配置文件。

### 步骤 3：启动服务

完成配置后，可以启动 VerKeyOSS 核心服务：

#### Windows 系统

1. **双击运行**（推荐）：
   - 直接双击下载的 `VerKeyOSS_0.0.1-beta_win_amd64.exe` 文件 （根据实际版本）
   - 服务将在配置的端口（默认为 8913）上启动

#### Linux 系统

1. **使用终端**：
   - 打开终端
   - 切换到程序所在目录
   - 设置可执行权限（如果尚未设置）：
     ```bash
     chmod +x VerKeyOSS_0.0.1-beta_linux_amd64 （根据实际版本替换）
     ```
   - 执行以下命令：
     ```bash
     ./VerKeyOSS_0.0.1-beta_linux_amd64 （根据实际版本替换）
     ```

### 步骤 4：访问管理界面

服务启动成功后：

1. **打开浏览器**，访问：`http://localhost:8913`
2. **使用默认管理员账号登录**：
   - 用户名：`verkeyoss`
   - 密码：`verkeyoss`
3. **首次登录后请立即修改密码**（通过Web界面或API接口）

**注意：** 集成版本将前端管理界面嵌入到Go应用中，单一可执行文件即包含了完整的前后端功能。

### API 访问地址

- **Web 管理界面**：`http://localhost:8913`
- **API 接口**：`http://localhost:8913/api`
- **健康检查**：`http://localhost:8913/api/check/health`