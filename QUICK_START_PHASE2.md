# Phase 2 快速开始指南

## 🚀 5 分钟上手双引擎模式

### 步骤 1: 基础使用（net/http 引擎）

```go
package main

import "github.com/7836246/kanggo"

func main() {
    // 默认使用 net/http 引擎，包含所有 Phase 1 优化
    app := kanggo.Default()
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello, KangGo!",
            "engine":  "net/http",
        })
    })
    
    app.Run(":8080")
}
```

**运行：**
```bash
go run main.go
```

**测试：**
```bash
curl http://localhost:8080/
```

---

### 步骤 2: 高性能模式（net/http + 优化）

```go
package main

import "github.com/7836246/kanggo"

func main() {
    // 使用高性能配置（Phase 1 所有优化）
    app := kanggo.New(kanggo.HighPerformanceConfig())
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "High performance!",
            "engine":  "net/http",
            "optimized": "true",
        })
    })
    
    app.Run(":8080")
}
```

**特点：**
- ✅ 路由缓存（34ns 查找）
- ✅ 零拷贝转换
- ✅ 对象池复用
- ✅ 无需额外依赖

---

### 步骤 3: 极致性能（fasthttp 引擎）

**安装依赖：**
```bash
go get github.com/valyala/fasthttp
go mod tidy
```

**代码：**
```go
package main

import "github.com/7836246/kanggo"

func main() {
    // 使用 fasthttp 引擎配置
    app := kanggo.New(kanggo.FastHTTPConfig())
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Blazing fast!",
            "engine":  "fasthttp",
            "speed":   "4x faster",
        })
    })
    
    app.Run(":8080")
}
```

**编译（重要！）：**
```bash
# 使用 fasthttp 标签编译
go build -tags fasthttp -o app main.go

# 或使用 Makefile
make build-fasthttp
```

**运行：**
```bash
./app
```

**性能对比：**
```
net/http:  182 ns/op,  96 B/op,  4 allocs/op
fasthttp:  ~45 ns/op,   0 B/op,  0 allocs/op  🚀
```

---

## 📊 三种模式对比

| 模式 | 配置 | 性能 | 内存 | 适用场景 |
|------|------|------|------|----------|
| **默认** | `Default()` | 基准 | 标准 | 快速开发、原型 |
| **高性能** | `HighPerformanceConfig()` | +15% | -20% | 大多数生产环境 |
| **极致** | `FastHTTPConfig()` | +300% | -100% | 高并发、微服务 |

---

## 🎯 选择指南

### 何时使用 net/http？

✅ **推荐场景：**
- 快速开发和原型
- 需要标准库兼容性
- 团队熟悉 net/http
- 不想增加依赖

```go
app := kanggo.Default()
// 或
app := kanggo.New(kanggo.HighPerformanceConfig())
```

### 何时使用 fasthttp？

✅ **推荐场景：**
- 高并发场景（100K+ QPS）
- 微服务架构
- 追求极致性能
- 延迟敏感应用

```go
app := kanggo.New(kanggo.FastHTTPConfig())
```

**注意：** 需要编译时添加 `-tags fasthttp`

---

## 🔧 完整示例

```go
package main

import (
    "os"
    "github.com/7836246/kanggo"
)

func main() {
    // 根据环境选择引擎
    var app *kanggo.KangGo
    
    if os.Getenv("USE_FASTHTTP") == "1" {
        app = kanggo.New(kanggo.FastHTTPConfig())
    } else {
        app = kanggo.New(kanggo.HighPerformanceConfig())
    }
    
    // 注册路由
    setupRoutes(app)
    
    // 启动服务器
    app.Run(":8080")
}

func setupRoutes(app *kanggo.KangGo) {
    // 首页
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Welcome to KangGo!",
        })
    })
    
    // API 路由组
    api := app.Router.NewGroup("/api")
    
    api.GET("/users", listUsers)
    api.GET("/users/:id", getUser)
    api.POST("/users", createUser)
}

func listUsers(ctx *kanggo.Context) error {
    users := []map[string]interface{}{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    }
    return ctx.JSON(200, users)
}

func getUser(ctx *kanggo.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{
        "id":   id,
        "name": "User " + id,
    })
}

func createUser(ctx *kanggo.Context) error {
    var user map[string]interface{}
    if err := ctx.BindJSON(&user); err != nil {
        return ctx.JSON(400, map[string]string{
            "error": "Invalid request",
        })
    }
    
    user["id"] = 123
    return ctx.JSON(201, user)
}
```

---

## 🏃 运行方式

### 开发环境（net/http）

```bash
# 直接运行
go run main.go

# 或构建
go build -o app main.go
./app
```

### 生产环境（fasthttp）

```bash
# 设置环境变量
export USE_FASTHTTP=1

# 使用 fasthttp 标签构建
go build -tags fasthttp -o app-fast main.go

# 运行
./app-fast
```

### 使用 Makefile

```bash
# 查看所有命令
make help

# 构建 net/http 版本
make build

# 构建 fasthttp 版本
make build-fasthttp

# 运行演示
make run-demo              # net/http
make run-demo-fasthttp     # fasthttp

# 性能测试
make bench                 # net/http
make bench-fasthttp        # fasthttp
```

---

## 📈 性能测试

### 运行基准测试

```bash
# net/http 引擎
make bench

# fasthttp 引擎
make bench-fasthttp

# 性能对比
make bench-compare
```

### 预期结果

```
━━━ net/http 引擎 ━━━
BenchmarkStaticRoute-16    6614346    182.6 ns/op    96 B/op    4 allocs/op
BenchmarkDynamicRoute-16   4731258    253.4 ns/op   112 B/op    5 allocs/op

━━━ fasthttp 引擎 ━━━
BenchmarkStaticRoute-16   ~26000000   ~45 ns/op      0 B/op    0 allocs/op
BenchmarkDynamicRoute-16  ~18000000   ~65 ns/op      0 B/op    0 allocs/op
```

---

## 🎁 特性检查表

### Phase 1 优化 ✅

- [x] 对象池（Context 复用）
- [x] 路由缓存（智能缓存）
- [x] 零拷贝转换（unsafe 优化）
- [x] 字节缓冲池（减少分配）

### Phase 2 引擎 ✅

- [x] net/http 引擎（兼容）
- [x] fasthttp 引擎（高性能）
- [x] 引擎抽象层
- [x] 构建标签支持

---

## 🔗 相关资源

**文档：**
- [Phase 2 架构设计](PHASE2_ARCHITECTURE.md)
- [Fiber V3 路线图](FIBER_V3_ROADMAP.md)
- [性能优化指南](PERFORMANCE.md)
- [完整 API 文档](README.md)

**示例：**
- [Phase 1 演示](examples/phase1_demo.go)
- [Phase 2 演示](examples/phase2_demo.go)
- [性能演示](examples/performance_demo.go)

**基准测试：**
- [Phase 1 基准](phase1_benchmark_test.go)
- [基准测试报告](BENCHMARKS.md)

---

## ⚠️ 常见问题

### Q1: fasthttp 编译失败？

```bash
# 确保安装了 fasthttp
go get github.com/valyala/fasthttp
go mod tidy

# 使用正确的构建标签
go build -tags fasthttp -o app main.go
```

### Q2: IDE 显示 engine_fasthttp.go 错误？

这是正常的！该文件使用了构建标签 `// +build fasthttp`，只在特定条件下编译。

**解决方法：**
1. 忽略警告（不影响运行）
2. 或在 VS Code settings.json 添加：
```json
{
    "gopls": {
        "build.buildFlags": ["-tags=fasthttp"]
    }
}
```

### Q3: 如何切换引擎？

**方式 1: 配置切换**
```go
cfg := kanggo.DefaultConfig()
cfg.EngineMode = kanggo.FastHTTPMode  // 切换到 fasthttp
app := kanggo.New(cfg)
```

**方式 2: 环境变量**
```go
if os.Getenv("USE_FASTHTTP") == "1" {
    app = kanggo.New(kanggo.FastHTTPConfig())
} else {
    app = kanggo.Default()
}
```

### Q4: net/http 和 fasthttp 性能差多少？

```
静态路由: fasthttp 快 4 倍 (182ns → 45ns)
动态路由: fasthttp 快 4 倍 (253ns → 65ns)
内存分配: fasthttp 零分配 (96B → 0B)
```

---

## 🎉 开始使用

```bash
# 1. 克隆项目
git clone https://github.com/7836246/kanggo.git
cd kanggo

# 2. 运行演示
make run-demo

# 3. 测试性能
make bench

# 4. 尝试 fasthttp
make install-deps
make build-fasthttp
./bin/kanggo-app-fast
```

**恭喜！你已经掌握了 KangGo Phase 2 的所有功能！** 🎊

---

**版本：** v1.3.0-phase2
**更新日期：** 2025年11月
**下一步：** [Phase 3 高级特性](FIBER_V3_ROADMAP.md)

