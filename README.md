# NextTrace Exporter

[English](#english) | [中文](#中文)

---

<a name="english"></a>

## English

A Prometheus exporter for [NextTrace](https://github.com/nxtrace/NTrace-core), providing continuous route tracing metrics with RTT, packet loss, ASN, and geolocation information.

### ✨ Features

- 🔄 **Continuous Background Execution** - Runs nexttrace periodically, not on-demand
- 📊 **Rich Prometheus Metrics** - Hop-by-hop RTT, packet loss, ASN info
- 🎯 **Multi-Target Support** - Configure multiple targets with individual intervals
- 🔄 **Hot Reload** - Update config without restart (SIGHUP or HTTP POST)
- 🐳 **Docker Ready** - Full Docker and docker-compose support
- 📝 **Structured Logging** - Detailed operational logs using Go slog

### 📦 Quick Start

#### Prerequisites

Install [NextTrace](https://github.com/nxtrace/NTrace-core):
```bash
curl -sL nxtrace.org/nt | sudo bash
```

#### Installation

**From Source:**
```bash
git clone https://github.com/vinsec/nexttrace_exporter.git
cd nexttrace_exporter
make build
sudo make install
```

**Using Docker:**
```bash
docker-compose up -d
```

#### Configuration

Create `config.yml`:
```yaml
# HTTP Server Configuration (optional)
server:
  listen_address: localhost:9101  # Default: localhost:9101 (local only)
  metrics_path: /metrics          # Default: /metrics

targets:
  - host: 8.8.8.8
    name: google_dns
    interval: 5m
    max_hops: 30
```

**Server Configuration:**
| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `listen_address` | string | No | localhost:9101 | HTTP server listen address |
| `metrics_path` | string | No | /metrics | Path where metrics are exposed |

> **Security Note**: Default is `localhost:9101` (local only). Use `0.0.0.0:9101` or `:9101` to listen on all interfaces.

**Target Configuration:**
| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `host` | string | Yes | - | Target hostname or IP |
| `name` | string | No | host | Friendly name (used in labels) |
| `interval` | duration | No | 5m | Execution interval (e.g., 30s, 5m, 1h) |
| `max_hops` | int | No | 30 | Maximum hops (1-64) |

> **Note**: The exporter automatically uses `nexttrace -j` for JSON output.

#### Running

**Standalone:**
```bash
sudo nexttrace_exporter --config.file=config.yml
```

**With systemd:**
```bash
sudo cp examples/systemd/nexttrace_exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now nexttrace_exporter
```

**With Docker:**
```bash
docker run -d \
  -p 9101:9101 \
  -v $(pwd)/config.yml:/etc/nexttrace_exporter/config.yml:ro \
  --cap-add=NET_RAW \
  --name nexttrace_exporter \
  nexttrace_exporter:latest
```

### 📊 Prometheus Metrics

The exporter provides the following metrics:

- `nexttrace_hop_rtt_milliseconds` - RTT per hop (with IP, hostname, ASN labels)
- `nexttrace_hop_loss_ratio` - Packet loss ratio per hop (0.0-1.0)
- `nexttrace_total_hops` - Total number of hops to target
- `nexttrace_execution_duration_seconds` - Execution time
- `nexttrace_executions_total` - Total executions counter (with status label)
- `nexttrace_last_execution_timestamp` - Last successful execution timestamp

### 🔧 Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config.file` | `config.yml` | Configuration file path |
| `--web.listen-address` | `localhost:9101` | HTTP listen address (overrides config file) |
| `--web.telemetry-path` | `/metrics` | Metrics endpoint path (overrides config file) |
| `--nexttrace.binary` | `nexttrace` | Path to nexttrace binary |
| `--nexttrace.timeout` | `2m` | Execution timeout |
| `--log.level` | `info` | Log level (debug/info/warn/error) |

> **Note**: Command-line flags take precedence over configuration file values.

### 🔄 Hot Reload

Reload configuration without restart:
```bash
# Send SIGHUP signal
sudo kill -HUP $(pgrep nexttrace_exporter)

# Or use HTTP endpoint
curl -X POST http://localhost:9101/-/reload
```

### 🌐 HTTP Endpoints

- `/metrics` - Prometheus metrics
- `/` - Web interface showing configured targets
- `/-/healthy` - Health check endpoint
- `/-/reload` - Configuration reload (POST)

### 📈 Prometheus Configuration

Add to your `prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'nexttrace'
    static_configs:
      - targets: ['localhost:9101']
    scrape_interval: 30s
```

See `examples/prometheus.yml` for a complete example with alerting rules.

### 🐛 Troubleshooting

**Permission Issues:**
```bash
# Grant required capabilities (recommended)
sudo setcap cap_net_raw+ep /usr/local/bin/nexttrace_exporter
sudo setcap cap_net_raw+ep $(which nexttrace)

# Or run as root (not recommended)
sudo nexttrace_exporter --config.file=config.yml
```

**Test nexttrace works:**
```bash
sudo nexttrace -j 8.8.8.8
```

**Debug mode:**
```bash
nexttrace_exporter --config.file=config.yml --log.level=debug
```

### 🤝 Contributing

Contributions are welcome! Please see [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

### 🔒 Security

- **Network Binding**: Default binding is `localhost:9101` (local only) for security
  - For remote access, use `0.0.0.0:9101` or specific IP
  - Consider firewall rules when binding to public interfaces
- **Run with minimal privileges**: Use `CAP_NET_RAW` capability instead of root
- **Authentication**: Use reverse proxy (nginx/caddy) with authentication for public access
- **TLS**: Enable HTTPS through reverse proxy for encrypted communication
- See [docs/SECURITY.md](docs/SECURITY.md) for more information

### 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

### 🙏 Acknowledgments

- [NextTrace](https://github.com/nxtrace/NTrace-core) - The underlying route tracing tool
- [Prometheus](https://prometheus.io/) - Metrics and monitoring system

---

<a name="中文"></a>

## 中文

一个用于 [NextTrace](https://github.com/nxtrace/NTrace-core) 的 Prometheus Exporter，提供持续的路由追踪指标，包括 RTT、丢包率、ASN 和地理位置信息。

### ✨ 特性

- 🔄 **后台持续执行** - 周期性运行 nexttrace，而非按需执行
- 📊 **丰富的 Prometheus 指标** - 逐跳 RTT、丢包率、ASN 信息
- 🎯 **多目标支持** - 为每个目标配置独立的执行间隔
- 🔄 **热重载** - 无需重启即可更新配置（SIGHUP 或 HTTP POST）
- 🐳 **Docker 就绪** - 完整的 Docker 和 docker-compose 支持
- 📝 **结构化日志** - 使用 Go slog 的详细操作日志

### 📦 快速开始

#### 前置要求

安装 [NextTrace](https://github.com/nxtrace/NTrace-core)：
```bash
curl -sL nxtrace.org/nt | sudo bash
```

#### 安装

**从源码构建：**
```bash
git clone https://github.com/vinsec/nexttrace_exporter.git
cd nexttrace_exporter
make build
sudo make install
```

**使用 Docker：**
```bash
docker-compose up -d
```

#### 配置

创建 `config.yml`：
```yaml
# HTTP 服务器配置（可选）
server:
  listen_address: localhost:9101  # 默认：localhost:9101（仅本地访问）
  metrics_path: /metrics          # 默认：/metrics

targets:
  - host: 8.8.8.8
    name: google_dns
    interval: 5m
    max_hops: 30
```

**服务器配置：**
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `listen_address` | string | 否 | localhost:9101 | HTTP 服务器监听地址 |
| `metrics_path` | string | 否 | /metrics | 指标暴露路径 |

> **安全提示**：默认值为 `localhost:9101`（仅本地访问）。使用 `0.0.0.0:9101` 或 `:9101` 可监听所有网络接口。

**目标配置：**
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `host` | string | 是 | - | 目标主机名或 IP 地址 |
| `name` | string | 否 | host | 友好名称（用于标签） |
| `interval` | duration | 否 | 5m | 执行间隔（如：30s, 5m, 1h） |
| `max_hops` | int | 否 | 30 | 最大跳数（1-64） |

> **注意**：Exporter 会自动使用 `nexttrace -j` 获取 JSON 输出。

#### 运行

**独立运行：**
```bash
sudo nexttrace_exporter --config.file=config.yml
```

**使用 systemd：**
```bash
sudo cp examples/systemd/nexttrace_exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now nexttrace_exporter
```

**使用 Docker：**
```bash
docker run -d \
  -p 9101:9101 \
  -v $(pwd)/config.yml:/etc/nexttrace_exporter/config.yml:ro \
  --cap-add=NET_RAW \
  --name nexttrace_exporter \
  nexttrace_exporter:latest
```

### 📊 Prometheus 指标

Exporter 提供以下指标：

- `nexttrace_hop_rtt_milliseconds` - 每跳的 RTT（带 IP、主机名、ASN 标签）
- `nexttrace_hop_loss_ratio` - 每跳的丢包率（0.0-1.0）
- `nexttrace_total_hops` - 到达目标的总跳数
- `nexttrace_execution_duration_seconds` - 执行耗时
- `nexttrace_executions_total` - 总执行次数（带状态标签）
- `nexttrace_last_execution_timestamp` - 最后一次成功执行的时间戳

### 🔧 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config.file` | `config.yml` | 配置文件路径 |
| `--web.listen-address` | `localhost:9101` | HTTP 监听地址（覆盖配置文件） |
| `--web.telemetry-path` | `/metrics` | 指标端点路径（覆盖配置文件） |
| `--nexttrace.binary` | `nexttrace` | nexttrace 二进制文件路径 |
| `--nexttrace.timeout` | `2m` | 执行超时时间 |
| `--log.level` | `info` | 日志级别（debug/info/warn/error） |

> **注意**：命令行参数的优先级高于配置文件。

### 🔄 热重载

无需重启即可重载配置：
```bash
# 发送 SIGHUP 信号
sudo kill -HUP $(pgrep nexttrace_exporter)

# 或使用 HTTP 端点
curl -X POST http://localhost:9101/-/reload
```

### 🌐 HTTP 端点

- `/metrics` - Prometheus 指标
- `/` - Web 界面，显示已配置的目标
- `/-/healthy` - 健康检查端点
- `/-/reload` - 配置重载（POST）

### 📈 Prometheus 配置

添加到你的 `prometheus.yml`：
```yaml
scrape_configs:
  - job_name: 'nexttrace'
    static_configs:
      - targets: ['localhost:9101']
    scrape_interval: 30s
```

查看 `examples/prometheus.yml` 获取包含告警规则的完整示例。

### 🐛 故障排除

**权限问题：**
```bash
# 授予所需权限（推荐）
sudo setcap cap_net_raw+ep /usr/local/bin/nexttrace_exporter
sudo setcap cap_net_raw+ep $(which nexttrace)

# 或以 root 运行（不推荐）
sudo nexttrace_exporter --config.file=config.yml
```

**测试 nexttrace 是否工作：**
```bash
sudo nexttrace -j 8.8.8.8
```

**调试模式：**
```bash
nexttrace_exporter --config.file=config.yml --log.level=debug
```

### 🤝 贡献

欢迎贡献！请查看 [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) 了解详情。

### 🔒 安全

- **网络绑定**：默认绑定到 `localhost:9101`（仅本地访问）以确保安全
  - 如需远程访问，使用 `0.0.0.0:9101` 或指定 IP
  - 绑定到公网接口时请配置防火墙规则
- **最小权限运行**：使用 `CAP_NET_RAW` 权限而非 root 运行
- **身份认证**：公网访问时使用带认证的反向代理（nginx/caddy）
- **TLS 加密**：通过反向代理启用 HTTPS 进行加密通信
- 查看 [docs/SECURITY.md](docs/SECURITY.md) 了解更多信息

### 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件。

### 🙏 致谢

- [NextTrace](https://github.com/nxtrace/NTrace-core) - 底层路由追踪工具
- [Prometheus](https://prometheus.io/) - 指标和监控系统

---

## 📁 Project Structure

```
nexttrace_exporter/
├── main.go                    # Entry point
├── config/                    # Configuration handling
├── executor/                  # NextTrace execution logic
├── collector/                 # Prometheus metrics collection
├── parser/                    # JSON parsing
├── examples/                  # Example configs
│   ├── config.yml            # Configuration example
│   ├── prometheus.yml        # Prometheus config
│   ├── alert_rules.yml       # Alert rules
│   ├── grafana_dashboard.json # Grafana dashboard
│   └── systemd/              # Systemd service file
├── Dockerfile                # Container image
├── docker-compose.yml        # Docker Compose setup
└── Makefile                  # Build automation
```
