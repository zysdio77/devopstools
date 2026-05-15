# FirApp - 移动应用分发平台

用于内部分发 iOS (.ipa) 和 Android (.apk) 安装包，支持 iOS OTA 在线安装。

## 项目结构

```
FirApp/
├── docker/                      # Docker 部署
│   ├── Dockerfile               # 多阶段构建（Vue + Go → Alpine）
│   ├── docker-compose.yaml      # MySQL + App 一键启动
│   ├── config.yaml              # Docker 环境配置
│   ├── init.sql                 # 数据库建表（自动执行）
│   └── template.plist           # iOS OTA 安装清单模板
├── fir_go/                      # Go 后端
│   ├── config/config.yaml       # 本地开发配置
│   ├── controler/
│   │   ├── router.go            # 路由定义
│   │   ├── api.go               # HTTP 处理器
│   │   ├── model.go             # 数据库操作
│   │   └── middleware.go        # CORS + 认证中间件
│   ├── template/
│   │   ├── template.plist       # iOS plist 模板
│   │   └── dist/                # 前端构建产物（npm run build 后生成）
│   └── main.go
└── fir_vue/                     # Vue 前端
    └── src/components/          # 页面组件
```

## API 接口

所有接口前缀 `/fir`，支持 Bearer Token 鉴权（可选，配置 `AUTH_TOKEN` 开启）。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/fir/info?name={ios\|android}` | 查询包列表 |
| GET | `/fir/page?name={ios\|android}&page=1` | 分页查询 |
| POST | `/fir/upload` | 上传安装包（multipart/form-data，字段：`upload-file`、`system_type`、`note`） |
| POST | `/fir/delete` | 删除安装包（JSON：`name`、`system_type`） |
| POST | `/fir/update` | 更新包信息（JSON：`name`、`system_type`） |

## 本地开发

```bash
# 1. 启动 MySQL，创建数据库
CREATE DATABASE fir_data DEFAULT CHARACTER SET utf8mb4;
CREATE TABLE fir (
    name VARCHAR(255) NOT NULL,
    system_type VARCHAR(255) NOT NULL,
    note VARCHAR(1000),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

# 2. 修改 fir_go/config/config.yaml 的数据库连接信息

# 3. 构建前端
cd fir_vue
npm install
npm run build
cp -r dist ../fir_go/template/dist

# 4. 启动后端
cd ../fir_go
go run main.go
```

## Docker 部署

```bash
# 启动
docker compose -f docker/docker-compose.yaml up -d

# 修改模板后只需重启，不用重建
docker compose -f docker/docker-compose.yaml restart app

# 停止并清理数据
docker compose -f docker/docker-compose.yaml down -v
```

## iOS OTA 安装

1. 上传 IPA 文件时，系统自动在 `plist/` 目录生成 `.plist` 清单文件
2. 在 iOS 设备上访问 `itms-services://?action=download-manifest&url=https://你的域名/plist/应用名.plist`
3. 下载页面点击「下载安装」按钮会自动跳转

**重要：** iOS OTA 安装强制要求 HTTPS，确保 `config.yaml` 中 `ipauri` 使用 `https://` 协议。

## 配置说明

| 配置项 | 说明 |
|--------|------|
| `server.serverRoot` | 静态文件目录（前端 dist + 上传文件） |
| `server.template` | plist 模板和上传文件存储目录 |
| `server.ipauri` | iOS IPA 包的下载地址前缀，**必须 HTTPS** |
| `server.authToken` | API 鉴权 Token，为空则跳过鉴权 |
| `server.allowedExtensions` | 允许上传的文件类型 |

## 修改 plist 模板

Docker 部署时，模板挂载在 `docker/template.plist`，修改后只需重启容器：

```bash
docker compose restart app
```

需要改的信息：
- `bundle-identifier` — App 的 Bundle ID
- `bundle-version` — 版本号
- `title` — 应用名称
- 图标 URL（full-size-image 和 display-image）
