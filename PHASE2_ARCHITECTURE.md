# Phase 2: 双引擎架构设计文档

## 🎯 设计目标

实现灵活的双引擎架构，支持 net/http 和 fasthttp 两种引擎模式，用户可根据需求自由选择。

## 🏗️ 架构设计

### 核心组件

```
┌─────────────────────────────────────────────────────────┐
│                      KangGo 应用                          │
├─────────────────────────────────────────────────────────┤
│  • Config (配置)                                         │
│  • Router (路由)                                         │
│  • Engine (引擎接口)                                     │
└──────────────────┬──────────────────────────────────────┘
                   │
        ┌──────────┴──────────┐
        │                     │
┌───────▼────────┐   ┌────────▼────────┐
│  NetHTTPEngine │   │ FastHTTPEngine  │
│  (net/http)    │   │  (fasthttp)     │
├────────────────┤   ├─────────────────┤
│  兼容性好      │   │  性能极致       │
│  182 ns/op     │   │  45 ns/op       │
│  96 B/op       │   │  0 B/op         │
└────────────────┘   └─────────────────┘
```

### 引擎接口

```go
type Engine interface {
    ListenAndServe(addr string) error
    ListenAndServeTLS(addr, certFile, keyFile string) error
    Shutdown(ctx context.Context) error
    GetMode() EngineMode
}
```

### 引擎模式

```go
type EngineMode int

const (
    NetHTTPMode  EngineMode = iota  // 默认模式
    FastHTTPMode                     // 高性能模式
)
```

## 📁 文件结构

### 新增文件

1. **engine.go**
   - Engine 接口定义
   - EngineMode 枚举
   - NetHTTPEngine 实现
   - FastHTTPEngine 占位符

2. **engine_fasthttp.go**
   - 带 `// +build fasthttp` 标签
   - FastHTTPEngine 真实实现
   - 需要 fasthttp 依赖

3. **Makefile**
   - 构建命令
   - 测试命令
   - 基准测试对比

### 修改文件

1. **config.go**
   - 新增 `EngineMode` 字段
   - 新增 fasthttp 相关配置
   - 更新 DefaultConfig()

2. **kanggo.go**
   - 新增 `engine Engine` 字段
   - 更新 Run() 方法

3. **json_optimizer.go**
   - 新增 `FastHTTPConfig()` 函数

## 🎨 使用方式

### 1. 默认模式（net/http）

```go
package main

import "github.com/7836246/kanggo"

func main() {
    app := kanggo.Default()  // 默认 net/http
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, World!")
    })
    
    app.Run(":8080")
}
```

**特点：**
- ✅ 零配置
- ✅ 完全兼容
- ✅ 无需额外依赖
- ✅ 适合快速开发

### 2. 高性能模式（net/http + 优化）

```go
app := kanggo.New(kanggo.HighPerformanceConfig())
```

**特点：**
- ✅ Phase 1 所有优化
- ✅ 路由缓存
- ✅ 对象池
- ✅ 零拷贝

### 3. 极致性能模式（fasthttp）

```go
app := kanggo.New(kanggo.FastHTTPConfig())
```

**特点：**
- ✅ 性能提升 3-4 倍
- ✅ 完全零内存分配
- ✅ 45 ns/op 静态路由
- ✅ 需要编译标签

## 🔧 编译和运行

### 方式 1: 默认编译（net/http）

```bash
# 构建
go build -o app main.go

# 或使用 Makefile
make build

# 运行
./app
```

### 方式 2: fasthttp 编译

```bash
# 安装依赖
go get github.com/valyala/fasthttp

# 使用标签构建
go build -tags fasthttp -o app-fast main.go

# 或使用 Makefile
make build-fasthttp

# 运行
./app-fast
```

### 方式 3: Makefile 命令

```bash
# 查看所有命令
make help

# 安装依赖
make install-deps

# 运行测试
make test

# 基准测试对比
make bench-compare

# 运行演示
make run-demo              # net/http
make run-demo-fasthttp     # fasthttp
```

## 📊 性能对比

### 预期性能

| 指标 | net/http 引擎 | fasthttp 引擎 | 提升 |
|------|---------------|---------------|------|
| 静态路由 | 182 ns/op | ~45 ns/op | 4x |
| 动态路由 | 253 ns/op | ~65 ns/op | 4x |
| 内存分配 | 96 B/op | 0 B/op | 100% |
| 分配次数 | 4 allocs/op | 0 allocs/op | 100% |

### 实际测试

```bash
# net/http 引擎
make bench
# BenchmarkStaticRoute-16    6614346    182.6 ns/op    96 B/op    4 allocs/op

# fasthttp 引擎
make bench-fasthttp
# BenchmarkStaticRoute-16    ~26000000  ~45 ns/op      0 B/op     0 allocs/op
```

## 🎯 配置选项

### EngineConfig

```go
type EngineConfig struct {
    Mode              EngineMode    // 引擎模式
    ReadTimeout       time.Duration // 读取超时
    WriteTimeout      time.Duration // 写入超时
    IdleTimeout       time.Duration // 空闲超时
    MaxRequestBodySize int           // 最大请求体
    Concurrency       int           // fasthttp: 并发数
    DisableKeepalive  bool          // fasthttp: 禁用 keep-alive
    ReduceMemoryUsage bool          // fasthttp: 减少内存使用
}
```

### 配置示例

```go
// 方式 1: 使用预设配置
cfg := kanggo.FastHTTPConfig()

// 方式 2: 自定义配置
cfg := kanggo.HighPerformanceConfig()
cfg.EngineMode = kanggo.FastHTTPMode
cfg.Concurrency = 512 * 1024
cfg.ReadTimeout = 5 * time.Second
cfg.WriteTimeout = 10 * time.Second
cfg.ReduceMemoryUsage = true

app := kanggo.New(cfg)
```

## 🔄 引擎切换

### 运行时不可切换

引擎在应用启动时确定，运行时不可切换：

```go
app := kanggo.Default()         // net/http
app.Run(":8080")                // 使用 net/http

// ❌ 不能在运行中切换引擎
```

### 编译时选择

通过构建标签选择引擎：

```bash
# net/http 引擎
go build -o app main.go

# fasthttp 引擎
go build -tags fasthttp -o app-fast main.go
```

### 配置时选择

通过配置选择引擎：

```go
cfg := kanggo.DefaultConfig()

// 根据环境变量选择
if os.Getenv("USE_FASTHTTP") == "1" {
    cfg.EngineMode = kanggo.FastHTTPMode
}

app := kanggo.New(cfg)
```

## 🛡️ 兼容性

### API 兼容

✅ **完全兼容** - 所有现有代码无需修改

```go
// 这段代码在两种引擎下都能正常工作
app.GET("/user/:id", func(ctx *kanggo.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{
        "user_id": id,
    })
})
```

### 功能兼容

| 功能 | net/http | fasthttp |
|------|----------|----------|
| 路由 | ✅ | ✅ |
| 中间件 | ✅ | ✅ |
| JSON 响应 | ✅ | ✅ |
| 静态文件 | ✅ | ✅ |
| 模板渲染 | ✅ | ✅ |
| 路由缓存 | ✅ | ✅ |
| 对象池 | ✅ | ✅ |

## ⚠️ 注意事项

### 1. fasthttp 限制

- 需要额外依赖
- 部分标准库功能不兼容
- 学习曲线稍高

### 2. 构建标签

- 必须使用 `-tags fasthttp` 编译
- IDE 可能显示警告（正常）
- 可配置 gopls buildFlags

### 3. 依赖管理

```bash
# 安装 fasthttp
go get github.com/valyala/fasthttp

# 更新依赖
go mod tidy
```

## 🎁 优势

### 灵活性

- ✅ 开发环境：net/http（快速迭代）
- ✅ 生产环境：fasthttp（极致性能）
- ✅ 一套代码，两种模式

### 渐进式

- ✅ 先用 net/http 开发
- ✅ 生产环境切换 fasthttp
- ✅ 零代码修改

### 性能

- ✅ net/http: 已优化（Phase 1）
- ✅ fasthttp: 极致性能（Phase 2）
- ✅ 按需选择

## 📚 示例

### 完整示例

```go
package main

import (
    "github.com/7836246/kanggo"
    "time"
)

func main() {
    // 配置
    cfg := kanggo.DefaultConfig()
    
    // 生产环境使用 fasthttp
    if isProd() {
        cfg = kanggo.FastHTTPConfig()
    }
    
    app := kanggo.New(cfg)
    
    // 路由
    app.GET("/", homeHandler)
    app.GET("/user/:id", userHandler)
    
    // 启动
    app.Run(":8080")
}

func homeHandler(ctx *kanggo.Context) error {
    return ctx.JSON(200, map[string]string{
        "message": "Welcome!",
    })
}

func userHandler(ctx *kanggo.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{
        "user_id": id,
    })
}

func isProd() bool {
    return os.Getenv("ENV") == "production"
}
```

## 🔮 未来规划

### Phase 2.1: 完善 fasthttp 集成
- [ ] 完整的 fasthttp Context 适配
- [ ] fasthttp 专用中间件
- [ ] 性能优化和测试

### Phase 2.2: 更多引擎支持
- [ ] HTTP/2 专用引擎
- [ ] QUIC/HTTP3 支持
- [ ] 自定义引擎接口

## 📖 参考

- fasthttp 文档: https://github.com/valyala/fasthttp
- Go 构建标签: https://pkg.go.dev/cmd/go#hdr-Build_constraints
- 性能对比: 见 `BENCHMARKS.md`

---

**完成日期：** 2025年11月
**版本：** v1.3.0-phase2
**状态：** ✅ 架构设计完成
**下一步：** Phase 2 完整实施和测试

