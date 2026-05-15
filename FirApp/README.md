# FirApp

移动应用分发平台，支持 iOS (.ipa) 和 Android (.apk) 上传、下载，iOS 支持 OTA 在线安装。

## 快速部署

提供两种部署方式：

- **基础模式**（`docker-compose.yaml`）：Go 后端 + MySQL，仅 HTTP，适合已有nginx做为反向代理的场景使用
- **nginx模式**（`docker-compose-nginx.yaml`）：Nginx + Go 后端 + MySQL，支持 HTTPS，适合新手，无需自己部署nginx服务

### 基础模式

```bash
docker compose -f docker/docker-compose.yaml up -d
```

访问 `http://你的IP:80`。

### nginx模式

#### 1. 准备 SSL 证书

iOS OTA 安装强制 HTTPS。推荐[阿里云免费 SSL 证书](https://yundun.console.aliyun.com/)，申请后下载 Nginx 格式，放到 `docker/ssl/`：

```
docker/ssl/
├── nginx.pem    # 证书文件
└── nginx.key    # 私钥文件
```

#### 2. 修改域名

编辑 `docker/nginx.conf`，把 `server_name xxx.xxx.com` 换成你的实际域名。

#### 3. 启动

```bash
docker compose -f docker/docker-compose-nginx.yaml up -d
```

首次启动会自动构建前端（Vue）和后端（Go）镜像，MySQL 也会自动初始化数据库表。

## 常用命令

```bash
# 改代码后重新构建（基础模式）
docker compose -f docker/docker-compose.yaml up -d --build

# 改代码后重新构建（生产模式）
docker compose -f docker/docker-compose-nginx.yaml up -d --build

# 修改 template.plist 后只需重启 app 服务
docker compose -f docker/docker-compose-nginx.yaml restart app

# 停止（基础模式）
docker compose -f docker/docker-compose.yaml down

# 停止（生产模式）
docker compose -f docker/docker-compose-nginx.yaml down

# 停止并清理数据（重置数据库）
docker compose -f docker/docker-compose-nginx.yaml down -v
```

## 修改 template.plist

`docker/template.plist` 是 iOS OTA 安装的描述文件模板，上传 IPA 时后端会基于此模板生成对应 App 的 plist。需要提前修改以下字段：

| 字段 | 位置 | 说明 |
|------|------|------|
| `ipaurl` | assets → software-package → url | **不用改**，后端上传 IPA 后会自动替换为实际的下载地址 |
| `bundle-identifier` | metadata → bundle-identifier | 填写 App 的 Bundle ID，如 `com.example.myapp` |
| `bundle-version` | metadata → bundle-version | 填写 App 的版本号，如 `1.0.0` |
| `title` | metadata → title | 安装时显示的 App 名称 |
| 大图标 URL | assets → full-size-image → url | 安装时显示的大图标，尺寸要求 **1024×1024**，需要提前上传到可访问的 URL |
| 小图标 URL | assets → display-image → url | 安装时显示的小图标，尺寸要求 **180×180**，需要提前上传到可访问的 URL |

修改完 `template.plist` 后，重启 app 服务使其生效：

```bash
docker compose -f docker/docker-compose-nginx.yaml restart app
```

### 注意事项

1. **`ipaurl` 不要修改** — 这是占位符，上传 IPA 后后端会自动替换为 `{ipauri}{文件名}`，其中 `ipauri` 来自 `docker/config.yaml` 的配置
2. iOS OTA 安装强制 HTTPS，确保 `docker/config.yaml` 中 `ipauri` 配置为 `https://` 协议
3. 图标文件必须是 PNG 格式，URL 必须使用 HTTPS，否则 iOS 设备可能加载不出来

## 项目结构

| 目录 | 说明 |
|------|------|
| `docker/` | Dockerfile、compose、配置、模板 |
| `fir_go/` | Go 后端，详见 [fir_go/readme.md](fir_go/readme.md) |
| `fir_vue/` | Vue 前端，详见 [fir_vue/README.md](fir_vue/README.md) |
