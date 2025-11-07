# KangGo v1.4.0-phase3 发布指南

## ✅ 已完成的步骤

1. ✅ 创建中文发布说明文档（`RELEASE_NOTES_ZH.md`）
2. ✅ 更新版本号（`version/version.go`）
3. ✅ 提交更改到 Git
4. ✅ 创建 Git 标签 `v1.4.0-phase3`

---

## 🚀 发布步骤

### 1. 推送代码和标签到远程仓库

```bash
# 推送代码
git push origin main

# 推送标签
git push origin v1.4.0-phase3

# 或者推送所有标签
git push origin --tags
```

### 2. 在 GitHub 上创建 Release

#### 方式一：使用 GitHub Web 界面

1. 访问项目仓库：https://github.com/7836246/kanggo
2. 点击右侧的 "Releases" 或访问：https://github.com/7836246/kanggo/releases
3. 点击 "Draft a new release" 或 "Create a new release"
4. 填写以下信息：
   - **Tag version**: 选择 `v1.4.0-phase3`
   - **Release title**: `🎉 KangGo v1.4.0-phase3 - 凤凰版本`
   - **Description**: 复制 `GITHUB_RELEASE_ZH.txt` 的内容
5. 勾选 "Set as the latest release"
6. 点击 "Publish release"

#### 方式二：使用 GitHub CLI

```bash
# 安装 GitHub CLI（如果还没安装）
# Windows: winget install GitHub.cli
# macOS: brew install gh
# Linux: 参考 https://github.com/cli/cli#installation

# 登录 GitHub
gh auth login

# 创建 Release
gh release create v1.4.0-phase3 \
  --title "🎉 KangGo v1.4.0-phase3 - 凤凰版本" \
  --notes-file GITHUB_RELEASE_ZH.txt \
  --latest

# 或者使用交互式创建
gh release create v1.4.0-phase3 --generate-notes
```

---

## 📋 发布清单

### 代码更改
- [x] 更新 `version/version.go` 到 `v1.4.0-phase3`
- [x] 创建 `RELEASE_NOTES_ZH.md`（完整发布说明）
- [x] 创建 `RELEASE_SUMMARY_ZH.md`（发布摘要）
- [x] 创建 `GITHUB_RELEASE_ZH.txt`（GitHub Release 文本）
- [x] 提交所有更改

### Git 操作
- [x] 创建提交：`Release v1.4.0-phase3 with Chinese release notes`
- [x] 创建标签：`v1.4.0-phase3`
- [ ] 推送到远程仓库
- [ ] 推送标签到远程仓库

### GitHub Release
- [ ] 在 GitHub 上创建 Release
- [ ] 使用 `GITHUB_RELEASE_ZH.txt` 作为发布说明
- [ ] 标记为 Latest Release
- [ ] 发布 Release

### 文档更新
- [ ] 更新 `README.md` 的徽章版本号（如果需要）
- [ ] 更新 `docs/PROJECT_STATUS.md`（如果需要）
- [ ] 发布社交媒体公告（如果需要）

---

## 📝 发布说明文件说明

### 1. RELEASE_NOTES_ZH.md（完整版）
- **用途**：项目根目录的完整发布说明
- **内容**：详细的功能介绍、示例代码、性能对比
- **长度**：约 700 行
- **适合**：用户深入了解所有新功能

### 2. RELEASE_SUMMARY_ZH.md（摘要版）
- **用途**：快速浏览的发布摘要
- **内容**：核心亮点、快速开始、重点功能
- **长度**：约 200 行
- **适合**：快速了解本次更新

### 3. GITHUB_RELEASE_ZH.txt（GitHub 版）
- **用途**：GitHub Release 页面的发布说明
- **内容**：格式化的发布信息，适合网页展示
- **长度**：约 250 行
- **适合**：GitHub Release 页面

---

## 🎯 推荐的发布流程

```bash
# 1. 确认当前状态
git status
git log --oneline -3
git tag --list | tail -5

# 2. 推送代码到远程
git push origin main

# 3. 推送标签到远程
git push origin v1.4.0-phase3

# 4. 验证推送成功
git ls-remote --tags origin

# 5. 在 GitHub 上创建 Release
# 访问：https://github.com/7836246/kanggo/releases/new
# 选择标签：v1.4.0-phase3
# 标题：🎉 KangGo v1.4.0-phase3 - 凤凰版本
# 说明：粘贴 GITHUB_RELEASE_ZH.txt 的内容
# 勾选：Set as the latest release
# 点击：Publish release
```

---

## 📢 发布后的宣传

### 1. 更新项目 README
确保 README.md 中的版本徽章显示最新版本：

```markdown
[![Version](https://img.shields.io/badge/version-v1.4.0--phase3-green.svg)](https://github.com/7836246/kanggo/releases)
```

### 2. 社交媒体
如果需要，可以在以下平台宣传：
- Twitter/X
- Reddit (r/golang)
- 掘金/思否（中文社区）
- V2EX

### 3. Go 社区
- 提交到 awesome-go 列表
- 在 Go Forum 发布公告
- 在相关 Discord/Slack 频道分享

---

## ⚡ 快速命令

```bash
# 推送所有内容
git push origin main && git push origin v1.4.0-phase3

# 查看远程标签
git ls-remote --tags origin

# 删除本地标签（如果需要重新创建）
git tag -d v1.4.0-phase3

# 删除远程标签（如果需要重新创建）
git push origin :refs/tags/v1.4.0-phase3

# 创建新标签
git tag -a v1.4.0-phase3 -m "KangGo v1.4.0-phase3 Release"

# 强制推送标签
git push origin v1.4.0-phase3 --force
```

---

## 🎉 完成标志

发布成功后，您应该能看到：

1. ✅ GitHub Releases 页面显示新版本
2. ✅ 标签 `v1.4.0-phase3` 在远程仓库中
3. ✅ Release 页面有完整的中文发布说明
4. ✅ 用户可以通过 `go get` 安装新版本

```bash
# 测试安装
go get -u github.com/7836246/kanggo@v1.4.0-phase3
```

---

## 📞 遇到问题？

如果在发布过程中遇到问题：

1. **标签推送失败**：检查是否有推送权限
2. **Release 创建失败**：确保标签已推送到远程
3. **版本号冲突**：检查是否已存在同名标签
4. **Go get 失败**：等待 1-5 分钟让 Go proxy 更新

---

**祝发布顺利！** 🚀

如有任何问题，请检查：
- GitHub 仓库权限
- 网络连接
- Git 配置
- GitHub CLI 认证

