# Monitor Webhook 通知转发服务

将 HTTP 请求转发到钉钉机器人和/或邮件，支持按路由组合不同的通知渠道。典型场景：接收 Prometheus AlertManager 告警，解析后通过钉钉 + 邮件双通道发送。

## 快速开始

```bash
# 本地编译运行
CGO_ENABLED=0 go build -o webhook main.go
./webhook

# Docker Compose 部署
docker compose up -d
```

服务默认监听 `:9099`。

---

## 配置文件 `webhook.yaml`

### 完整示例

```yaml
smtp:
  host: smtp.example.com       # SMTP 服务器地址
  port: 465                 # 端口，465 为 SMTPS（直接 TLS）
  user: noreply@example.com # 发件邮箱
  password: your_password   # 邮箱密码或授权码

routes:
  alert:                    # Prometheus 告警专用路由，钉钉 + 邮件同时发
    use_ding: true
    dingtalk: https://oapi.dingtalk.com/robot/send?access_token=xxx
    use_email: true
    email:
      - admin@example.com
      - oncall@example.com

  newschange:               # 仅发钉钉
    use_ding: true
    dingtalk: https://oapi.dingtalk.com/robot/send?access_token=xxx

  devops:                   # 仅发邮件
    use_email: true
    email:
      - devops@example.com
```

### 配置字段说明

#### `smtp` - SMTP 邮件配置

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `host` | string | 是 | SMTP 服务器地址，如 `smtp.sina.com` |
| `port` | int | 是 | SMTP 端口，465（SMTPS 直连 TLS）或 587（STARTTLS 不支持） |
| `user` | string | 是 | 发件邮箱账号 |
| `password` | string | 是 | 邮箱密码或 SMTP 授权码 |

> **注意**：当前实现仅支持 465 端口的直接 TLS 连接（SMTPS），不支持 587 端口的 STARTTLS 升级。

#### `routes` - 路由配置

每个路由是一个键值对，key 为路由名称（对应 API 路径中的 `:group` 参数），value 包含以下字段：

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `use_ding` | bool | 否 | `true` | 是否启用钉钉通知，显式设为 `false` 关闭 |
| `dingtalk` | string | 否 | - | 钉钉机器人 Webhook URL |
| `use_email` | bool | 否 | `true` | 是否启用邮件通知，显式设为 `false` 关闭 |
| `email` | []string | 否 | - | 收件邮箱列表 |

> **开关逻辑**：`use_ding` / `use_email` 不写默认启用；写 `false` 强制关闭，无需删除 URL 或邮箱配置。如果 `dingtalk` 为空或 `email` 列表为空，即使开关为 `true` 也不会创建对应的通知器。

---

## API 接口

### 接口一览

| 方法 | 路径 | Content-Type | 说明 |
|------|------|-------------|------|
| `GET` | `/health` | - | 健康检查 |
| `POST` | `/send/:group` | `application/json` | 通用转发，原始 body 直接转发 |
| `POST` | `/recieve` | `application/json` | Prometheus AlertManager 告警接收，解析后格式化发送 |

所有接口响应均为 `application/json`。

---

### `GET /health` - 健康检查

**响应示例：**

```json
{"status": "ok"}
```

用于负载均衡器或容器编排的健康探测。

---

### `POST /send/:group` - 通用转发

将请求体**原样**转发到 `group` 对应路由配置的所有通知渠道。

**路径参数：**

| 参数 | 说明 |
|------|------|
| `group` | 路由名称，对应 `webhook.yaml` 中 `routes` 下的 key |

**请求示例：**

```bash
# 发送到 newschange 路由（钉钉机器人）
curl -X POST http://localhost:9099/send/newschange \
  -H "Content-Type: application/json" \
  -d '{"title":"服务上线通知","content":"v2.3.1 已发布"}'

# 发送到 devops 路由（邮件）
curl -X POST http://localhost:9099/send/devops \
  -H "Content-Type: application/json" \
  -d '{"msg":"磁盘使用率超过 80%"}'
```

**成功响应：**

```json
{"msg": "ok"}
```

**失败响应：**

```json
{"error": "unknown group: testgroup"}
```

> 路由不存在时返回 HTTP 404。

---

### `POST /recieve` - Prometheus 告警接收

接收 Prometheus AlertManager 的 Webhook JSON，解析后通过 `alert` 路由发送。

**要求**：`webhook.yaml` 中必须配置 key 为 `alert` 的路由。

**请求格式（Prometheus AlertManager）：**

```json
{
  "status": "firing",
  "alerts": [
    {
      "annotations": {
        "summary": "CPU 使用率过高",
        "value": "92%"
      },
      "labels": {
        "desc": "node_cpu_usage"
      }
    }
  ]
}
```

**解析逻辑：**

| 源字段 | 目标 |
|--------|------|
| `status` | 决定消息前缀：`firing` → "故障警告"，`resolved` → "故障已解决" |
| `alerts.0.annotations.summary` | 告警名称 |
| `alerts.0.labels.desc` | 告警单位/描述 |
| `alerts.0.annotations.value` | 告警状态值 |

> **注意**：仅处理 `alerts` 数组的第一个元素（`alerts.0`），批量告警中额外的告警会被忽略。

**格式化的钉钉消息示例：**

```
故障警告：
告警：CPU 使用率过高
单位：node_cpu_usage
状态(报警时的触发状态，恢复后的详细状态定请访问：http://prometheus.example.com:9090/：92%
```

**请求示例：**

```bash
curl -X POST http://localhost:9099/recieve \
  -H "Content-Type: application/json" \
  -d '{
    "status": "firing",
    "alerts": [{
      "annotations": {"summary": "内存不足", "value": "95%"},
      "labels": {"desc": "memory_usage"}
    }]
  }'
```

---

## 接入 Prometheus AlertManager

在 AlertManager 配置文件中添加 webhook receiver：

```yaml
receivers:
  - name: 'webhook'
    webhook_configs:
      - url: 'http://<服务地址>:9099/recieve'
        send_resolved: true   # 恢复时也发送通知
```

然后修改告警路由规则，将告警发送到该 receiver。

---

## 钉钉机器人配置

1. 在钉钉群聊中，点击 **群设置 → 智能群助手 → 添加机器人 → 自定义（通过 Webhook 接入）**
2. 设置机器人名称（如"监控告警"），选择安全设置（建议选择"自定义关键词"或"加签"）
3. 复制生成的 Webhook URL，填入 `webhook.yaml` 对应路由的 `dingtalk` 字段

钉钉 Webhook URL 格式：
```
https://oapi.dingtalk.com/robot/send?access_token=<你的token>
```

---

## Docker 部署

### 使用 docker-compose（推荐）

```bash
# 启动
docker compose up -d

# 查看日志
docker compose logs -f

# 修改配置后重启
docker compose restart

# 停止
docker compose down
```

`webhook.yaml` 通过 volume 挂载进容器，修改配置后只需重启，无需重新构建镜像。

### 手动构建

```bash
docker build -t monitor-webhook .
docker run -d -p 9099:9099 -v $(pwd)/webhook.yaml:/app/webhook.yaml monitor-webhook
```

---

## 项目结构

```
├── main.go              # 入口：加载配置、注册路由、启动服务、优雅退出
├── webhook.yaml         # 配置文件
├── Dockerfile           # 多阶段构建（golang:1.25-alpine → alpine）
├── docker-compose.yml   # Docker Compose 编排
└── method/
    ├── config.go        # 配置结构体定义 + YAML 加载（LoadConfig）
    ├── VarName.go       # 数据模型 + InitRobot 路由初始化
    ├── notify.go        # Notifier 接口 + 钉钉/邮件实现
    ├── notice.go        # POST /send/:group 通用转发 handler
    └── alart.go         # POST /recieve Prometheus 告警 handler
```

---

## 关键依赖

| 依赖 | 用途 |
|------|------|
| `github.com/gin-gonic/gin` | HTTP 框架 |
| `github.com/go-resty/resty/v2` | HTTP 客户端（向钉钉发请求，10s 超时 + 2 次重试） |
| `github.com/sirupsen/logrus` | 结构化日志 |
| `github.com/tidwall/gjson` | JSON 路径解析（Prometheus 告警字段提取） |
| `gopkg.in/yaml.v3` | YAML 配置解析 |

---

## 常见问题

**Q: 邮件发送失败？**
- 确认 SMTP 服务器和端口正确（当前只支持 465 端口的直接 TLS）
- 确认账号密码无误，部分邮箱需使用 SMTP 授权码而非登录密码
- 检查服务器是否允许访问外部 SMTP 端口

**Q: 钉钉消息收不到？**
- 检查 Webhook URL 中的 `access_token` 是否有效
- 确认服务器可以访问 `oapi.dingtalk.com`
- 查看服务日志，钉钉 API 会返回 `errcode` 和 `errmsg`

**Q: 如何添加新的通知渠道？**
- 在 `webhook.yaml` 的 `routes` 下添加新的 key，配置钉钉 URL 和/或邮箱列表
- 调用 `POST /send/<新key>` 即可使用，无需修改代码

**Q: 服务本身有鉴权吗？**
- 没有。建议在前端部署 nginx 反向代理并配置访问控制，或在 Kubernetes 中使用 NetworkPolicy 限制访问来源。
