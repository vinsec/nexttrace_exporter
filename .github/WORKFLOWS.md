# GitHub Workflows 使用指南

本项目包含两个 GitHub Actions workflow：

## 1. CI Workflow (`ci.yml`)

### 触发条件
- 推送到 `main` 分支时自动运行
- 创建 Pull Request 时自动运行

### 执行内容
- **测试**: 运行单元测试并生成代码覆盖率报告
- **构建**: 为 Linux/macOS 的 amd64/arm64 架构构建二进制文件
- **代码检查**: 使用 golangci-lint 进行代码质量检查
- **Docker 构建**: 测试 Docker 镜像构建

### 产物
- 测试覆盖率报告上传到 Codecov
- 构建产物作为 GitHub Actions artifacts 保存

---

## 2. Release Workflow (`release.yml`)

### 触发方式

#### 方式一：手动触发（推荐用于测试）

1. 在 GitHub 仓库页面，点击 **Actions** 标签
2. 选择左侧的 **Release** workflow
3. 点击右侧的 **Run workflow** 按钮
4. 输入版本号（例如：`v0.1.0`）
5. 点击 **Run workflow** 开始构建

#### 方式二：打标签自动触发（推荐用于正式发布）

```bash
# 创建并推送标签
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

### 构建平台

构建支持以下平台和架构：

| 操作系统 | 架构 | 文件名示例 |
|---------|------|-----------|
| Linux | amd64 | `nexttrace_exporter-v0.1.0-linux-amd64.tar.gz` |
| Linux | arm64 | `nexttrace_exporter-v0.1.0-linux-arm64.tar.gz` |
| Linux | armv7 | `nexttrace_exporter-v0.1.0-linux-armv7.tar.gz` |
| macOS | amd64 (Intel) | `nexttrace_exporter-v0.1.0-darwin-amd64.tar.gz` |
| macOS | arm64 (Apple Silicon) | `nexttrace_exporter-v0.1.0-darwin-arm64.tar.gz` |
| Windows | amd64 | `nexttrace_exporter-v0.1.0-windows-amd64.zip` |
| FreeBSD | amd64 | `nexttrace_exporter-v0.1.0-freebsd-amd64.tar.gz` |

### 产物内容

每个发布包包含：
- 编译好的二进制文件
- SHA256 校验和文件
- README.md
- LICENSE
- config.yml 示例配置

### Docker 镜像

同时会构建并推送多架构 Docker 镜像到：
- Docker Hub: `<your-username>/nexttrace_exporter`
- GitHub Container Registry: `ghcr.io/<your-username>/nexttrace_exporter`

支持的架构：
- `linux/amd64`
- `linux/arm64`
- `linux/arm/v7`

镜像标签：
- `latest` - 最新版本
- `v0.1.0` - 完整版本号
- `v0.1` - 主版本号.次版本号
- `v0` - 主版本号

---

## 配置 Secrets

要使用完整的 Release workflow（包括 Docker 镜像发布），需要在 GitHub 仓库设置中配置以下 Secrets：

### Docker Hub (可选)

1. 进入仓库 **Settings** → **Secrets and variables** → **Actions**
2. 添加以下 secrets：
   - `DOCKER_USERNAME`: Docker Hub 用户名
   - `DOCKER_PASSWORD`: Docker Hub 访问令牌（推荐）或密码

### GitHub Container Registry

GitHub Container Registry 使用 `GITHUB_TOKEN`，无需额外配置。

---

## 使用示例

### 发布新版本的完整流程

```bash
# 1. 确保代码已提交
git add .
git commit -m "feat: add new feature"
git push

# 2. 创建并推送标签
git tag -a v0.2.0 -m "Release version 0.2.0"
git push origin v0.2.0

# 3. 等待 GitHub Actions 自动构建和发布
# 可以在 Actions 页面查看进度
```

### 下载和验证发布包

```bash
# 下载发布包
wget https://github.com/vinsec/nexttrace_exporter/releases/download/v0.1.0/nexttrace_exporter-v0.1.0-linux-amd64.tar.gz
wget https://github.com/vinsec/nexttrace_exporter/releases/download/v0.1.0/nexttrace_exporter-v0.1.0-linux-amd64.tar.gz.sha256

# 验证校验和
sha256sum -c nexttrace_exporter-v0.1.0-linux-amd64.tar.gz.sha256

# 解压
tar xzf nexttrace_exporter-v0.1.0-linux-amd64.tar.gz

# 运行
./nexttrace_exporter-v0.1.0-linux-amd64 --version
```

### 使用 Docker 镜像

```bash
# 拉取最新版本
docker pull ghcr.io/vinsec/nexttrace_exporter:latest

# 或拉取特定版本
docker pull ghcr.io/vinsec/nexttrace_exporter:v0.1.0

# 运行
docker run -d \
  -p 9101:9101 \
  -v $(pwd)/config.yml:/etc/nexttrace_exporter/config.yml:ro \
  --cap-add=NET_RAW \
  --name nexttrace_exporter \
  ghcr.io/vinsec/nexttrace_exporter:latest
```

---

## 故障排除

### Release workflow 失败

1. **权限错误**: 确保仓库的 Actions 权限设置正确
   - Settings → Actions → General → Workflow permissions
   - 选择 "Read and write permissions"

2. **Docker 推送失败**: 检查 DOCKER_USERNAME 和 DOCKER_PASSWORD secrets 是否正确配置

3. **标签已存在**: 如果需要重新发布，先删除旧标签：
   ```bash
   git tag -d v0.1.0
   git push origin :refs/tags/v0.1.0
   ```

### 查看构建日志

在 GitHub 仓库页面：
1. 点击 **Actions** 标签
2. 选择对应的 workflow 运行
3. 点击任务查看详细日志

---

## 版本号规范

建议遵循 [Semantic Versioning](https://semver.org/)：

- `v1.0.0` - 正式发布
- `v1.1.0` - 新增功能
- `v1.1.1` - Bug 修复
- `v2.0.0` - 重大变更（不向后兼容）
- `v0.1.0-beta.1` - 预发布版本（可选）
