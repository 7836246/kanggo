# KangGo

KangGo 是一个极简且高性能的 Go Web 框架，致力于为开发者提供快速、灵活和高效的开发体验。它采用现代化的设计理念，专注于性能优化和可扩展性，适用于构建各类 Web 应用。框架结构简洁，易于上手，同时具备高度的模块化，方便开发者根据需求进行定制和扩展。

## 特性

- ⚡️ **超高性能路由**：采用自适应 Radix Tree（Adaptive Radix Tree, ART）与暴力哈希表（Perfect Hashing）结合的混合路由算法，实现 O(1) 静态路由查找和近乎 O(1) 的动态路由匹配，确保极致的请求处理性能。
- 🔄 **异步非阻塞 I/O**：核心设计支持异步非阻塞请求处理，充分利用 Go 的 Goroutine 并发机制，提升 I/O 密集型任务的吞吐量。
- 🔌 **模块化中间件体系**：支持全局和路由级中间件，开发者可以灵活组合使用，提供请求处理和响应的高度定制化能力。
- 🧩 **插件式架构**：核心功能模块化设计，支持通过插件方式动态加载扩展功能，方便快速迭代和更新。
- 📦 **无外部依赖**：尽可能使用 Go 标准库实现功能，减少外部依赖，保持框架轻量且安全。
- 📈 **智能内存管理**：优化的内存分配策略，减少垃圾回收（GC）压力，确保高并发场景下的稳定性能表现。
- 🔧 **自定义配置**：支持灵活的配置选项，包括自定义 JSON 编解码器和欢迎横幅显示。

## 快速开始

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

- **路由组**：支持路由分组，方便 API 管理。
- **内置中间件**：日志、恢复、跨域等常用中间件开箱即用。
- **内存池**：高效的内存管理，减少 GC 开销，提高并发处理能力。

## 未来路线图

- 🔧 增强路由解析算法，支持更复杂的路由匹配规则。
- 🛡️ 提供更多内置中间件（如限流、身份验证等）。
- 🌍 支持国际化和本地化。
- 📜 集成 Swagger 自动文档生成。

## 贡献

欢迎参与 KangGo 的开发！请提交 Pull Request 或 Issue 来帮助我们改善这个框架。