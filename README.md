# KangGo

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-v1.3.0--phase2-green.svg)](https://github.com/7836246/kanggo/releases)

**KangGo** 是一个极简且高性能的 Go Web 框架，专为 2025 年现代化 Web 应用设计。采用双引擎架构，支持 `net/http` 和 `fasthttp` 两种模式，让你在兼容性和极致性能之间自由选择。

```
🚀 性能卓越    ⚡ 双引擎架构    🔧 零配置开箱即用    📦 完全模块化
```

## ✨ 核心特性

### 🏎️ 双引擎架构（Phase 2 新特性）

- **net/http 引擎**：完全兼容，开箱即用，已优化至超越 Gin
- **fasthttp 引擎**：极致性能，零内存分配，接近 Fiber V3
- **灵活切换**：配置驱动，一套代码，两种模式

```go
// 默认模式（net/http）
app := kanggo.Default()

// 极致模式（fasthttp）
app := kanggo.New(kanggo.FastHTTPConfig())
```

### ⚡ Phase 1 性能优化

- 🎯 **智能路由缓存**：34 ns/op，O(1) 查找
- 🔄 **对象池复用**：零 GC 压力的 Context 池
- 🚀 **零拷贝转换**：unsafe 优化，快 177 倍
- 💾 **字节缓冲池**：11.69 ns/op，零分配

### 🎨 开发者友好

- ⚡️ **超高性能路由**：O(1) 静态路由，近 O(1) 动态路由
- 🔌 **模块化中间件**：灵活的全局和路由级中间件
- 🧩 **插件式架构**：核心功能模块化，易于扩展
- 📦 **最小依赖**：标准库优先，可选 fasthttp
- 🔧 **自定义配置**：灵活的配置选项，支持多种场景

## 📊 性能对比

| 框架 | 引擎 | 静态路由 | 动态路由 | 内存 | 分配 |
|------|------|---------|---------|------|------|
| Gin | net/http | ~250 ns | ~400 ns | 600 B | 8x |
| **KangGo** | **net/http** | **182 ns** | **253 ns** | **96 B** | **4x** |
| **KangGo** | **fasthttp** | **~45 ns** | **~65 ns** | **0 B** | **0x** |
| Fiber V3 | fasthttp | ~40 ns | ~60 ns | 0 B | 0x |

**结论：** 
- ✅ net/http 模式：已超越 Gin（快 27%）
- ✅ fasthttp 模式：接近 Fiber V3（性能相当）

## 🚀 快速开始

### 安装

```bash
go get -u github.com/7836246/kanggo@latest
```

### 使用示例

```go
package main

import (
    "fmt"
    "github.com/7836246/kanggo"
    "encoding/json"
)

func main() {
    // 使用自定义的 JSON 编解码器配置
    cfg := kanggo.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
        ShowBanner:  true,
    }

    // 使用 Default() 初始化框架
    app := kanggo.New(cfg)

    // 注册静态路由
    app.GET("/home", func(ctx *kanggo.Context) error {
        return ctx.SendString("Welcome to the home page!")
    })

    // 注册带命名参数的动态路由
    app.GET("/user/:id", func(ctx *kanggo.Context) error {
        id := ctx.Param("id")
        if id == "" {
            return fmt.Errorf("id 参数缺失")
        }
        msg := fmt.Sprintf("用户ID：%s", id)
        return ctx.SendString(msg)
    })

    // 启动服务器
    app.Run(":8080")
}
```

## 高级特性

### 高性能配置

```go
// 使用高性能配置（2025 优化版）
app := kanggo.New(kanggo.HighPerformanceConfig())

// 自定义高性能 JSON 库（如 sonic）
cfg := kanggo.DefaultConfig()
cfg.JSONEncoder = sonic.Marshal
cfg.JSONDecoder = sonic.Unmarshal
app := kanggo.New(cfg)
```

### 性能优化特性

- **Context 对象池**：使用 sync.Pool 自动管理 Context 对象，减少 60-80% 的内存分配
- **静态路由哈希表**：O(1) 时间复杂度的路由查找，性能提升 10-100 倍
- **零内存分配路径解析**：手动实现路径分割，避免不必要的内存分配
- **静态文件缓存**：内置文件缓存，支持配置缓存大小和过期时间

### 其他特性

- **路由组**：支持路由分组，方便 API 管理。
- **内置中间件**：日志、恢复、跨域、ETag 等常用中间件开箱即用。
- **内存池**：高效的内存管理，减少 GC 开销，提高并发处理能力。

## 性能基准测试

运行基准测试查看性能表现：

```bash
go test -bench=. -benchmem -benchtime=3s
```

## 📚 完整文档

KangGo 提供了详尽的文档系统，包含 25+ 份完整文档：

### 快速开始
- **[文档中心](docs/README.md)** - 所有文档的导航中心
- **[项目状态](docs/PROJECT_STATUS.md)** - 当前项目状态
- **[快速开始 - Phase 2](docs/QUICK_START_PHASE2.md)** - 5分钟快速上手

### 性能文档
- **[性能优化指南](docs/PERFORMANCE.md)** - 完整的性能优化指南
- **[基准测试报告](docs/BENCHMARKS.md)** - 详细的性能测试结果
- **[优化报告](docs/OPTIMIZATION_REPORT.md)** - 优化详细报告

### Phase 报告
- **[Phase 1 完成报告](docs/PHASE1_COMPLETED.md)** - 底层优化详解
- **[Phase 2 架构设计](docs/PHASE2_ARCHITECTURE.md)** - 双引擎架构
- **[Phase 3 完成报告](docs/PHASE3_COMPLETED.md)** - 高级特性详解

### 对比与路线图
- **[与 Gin 对比](docs/COMPARISON_WITH_GIN.md)** - KangGo vs Gin
- **[Fiber V3 路线图](docs/FIBER_V3_ROADMAP.md)** - 技术演进路线
- **[终极总结](docs/ULTIMATE_SUMMARY.md)** - 完整项目总结

> 📖 **查看所有文档**: [docs/README.md](docs/README.md)

## 🎯 示例程序

KangGo 提供了丰富的示例程序，展示各种功能的使用：

```bash
# Phase 1: 性能优化演示
go run examples/phase1_demo.go

# Phase 2: 双引擎演示
go run examples/phase2_demo.go

# WebSocket 实时通信
go run examples/websocket_demo.go

# Server-Sent Events
go run examples/sse_demo.go

# 请求验证系统
go run examples/validator_demo.go
```

## 🔮 未来路线图

### Phase 4: 生态扩展（规划中）
- 🔧 CLI 工具和项目生成器
- 🛡️ 更多内置中间件（限流、认证等）
- 💾 数据库集成和 ORM 支持
- 📊 性能监控和链路追踪
- 🌍 国际化和本地化
- 📜 Swagger 自动文档生成

详细路线图请查看: [docs/FIBER_V3_ROADMAP.md](docs/FIBER_V3_ROADMAP.md)

## 🤝 贡献

欢迎参与 KangGo 的开发！我们欢迎各种形式的贡献：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码

请阅读 [贡献指南](CONTRIBUTING.md) 了解详情。