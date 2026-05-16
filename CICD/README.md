# CICD

一个轻量级 CI/CD 工具集，YAML 定义流水线，Go 单二进制分发。

## 快速开始

```bash
go build -o bin/cicd ./cmd/cicd
./bin/cicd examples/pipeline.yaml
```

## 目录结构

```
CICD/
├── cmd/cicd/main.go           # 入口
├── internal/
│   ├── pipeline/pipeline.go   # YAML 定义、加载、校验
│   ├── executor/executor.go   # 调度器（阶段串行，步骤并行）
│   ├── notifier/notifier.go   # 通知（钉钉、邮件）
│   └── runner/runner.go       # 编排：加载→执行→报告→通知
├── examples/pipeline.yaml     # 示例
├── Makefile
└── go.mod
```

## 流水线定义

```yaml
name: "my-pipeline"

notify:
  dingtalk_url: "https://oapi.dingtalk.com/robot/send?access_token=xxx"
  # email_smtp: "smtp.example.com"
  # email_to: "team@example.com"

stages:
  - name: "build"
    steps:
      - name: "compile"
        cmd: "go build ./..."
        timeout: "5m"
        work_dir: "/path/to/project"
        env:
          - "VERSION=v1.0.0"

      - name: "lint"
        cmd: "golangci-lint run"
```

## 字段说明

### Pipeline

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 是 | 流水线名称 |
| `stages` | 是 | 阶段列表，按顺序串行执行 |
| `notify` | 否 | 通知配置 |

### Stage

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 是 | 阶段名称 |
| `steps` | 是 | 步骤列表，同一阶段内并行执行 |

### Step

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 是 | 步骤名称 |
| `cmd` | 是 | shell 命令，底层执行 `sh -c` |
| `timeout` | 否 | 超时时间，如 `"30s"` `"5m"`，默认 30 分钟 |
| `work_dir` | 否 | 工作目录 |
| `env` | 否 | 环境变量，格式 `"KEY=value"` |

## 执行规则

- **Stage 之间串行**：上一个阶段全部成功才进入下一阶段
- **Step 之间并行**：同一阶段内所有步骤同时启动
- **失败处理**：任一步骤失败 → 取消同阶段其他步骤 → 跳过后续所有阶段

## 通知

### 钉钉

在 `notify.dingtalk_url` 填入机器人 webhook 地址即可。

### 邮件

```yaml
notify:
  email_smtp: "smtp.example.com"
  email_to: "team@example.com"
```

通过环境变量传入邮箱账密：

```bash
CICD_EMAIL_USER=bot@example.com CICD_EMAIL_PASS=xxx ./bin/cicd pipeline.yaml
```

## Make 命令

```bash
make build    # 编译到 bin/cicd
make run      # 编译并运行示例流水线
make clean    # 删除编译产物
```
