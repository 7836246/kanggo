# KangGo 向 Fiber V3 技术路线演进规划

## 📋 Fiber V3 核心技术特点分析

### 1. 零内存分配架构
- fasthttp 底层（而非 net/http）
- 零拷贝技术
- 字节池复用
- 预分配缓冲区

### 2. 极致的性能优化
- Context 池化（我们已实现 ✅）
- 字符串零拷贝
- 路由树优化
- 并发安全的内存管理

### 3. 先进的路由算法
- Radix Tree + 哈希表混合（我们已实现 ✅）
- 路由缓存
- 通配符优化
- 正则表达式支持

### 4. 现代化 API 设计
- 链式调用
- 中间件系统
- 错误处理机制
- 类型安全

### 5. 高级特性
- WebSocket 支持
- Server-Sent Events (SSE)
- HTTP/2 推送
- 自定义协议

## 🎯 技术演进路线图

### Phase 1: 底层优化（当前完成度：70%）

#### ✅ 已完成
- [x] Context 对象池
- [x] 静态路由哈希表
- [x] 零内存分配路径解析
- [x] 预分配切片容量
- [x] 静态文件缓存

#### 🔄 待实现
```go
// 1. 字节池系统
type ByteBufferPool struct {
    pools []*sync.Pool // 不同大小的池
}

func (p *ByteBufferPool) Get(size int) *[]byte {
    // 根据大小选择合适的池
    idx := size / 1024
    if idx >= len(p.pools) {
        idx = len(p.pools) - 1
    }
    return p.pools[idx].Get().(*[]byte)
}

// 2. 零拷贝字符串转换
func b2s(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) []byte {
    return *(*[]byte)(unsafe.Pointer(
        &struct {
            string
            Cap int
        }{s, len(s)},
    ))
}

// 3. 请求/响应池
var requestPool = sync.Pool{
    New: func() interface{} {
        return &Request{
            header: make(map[string]string, 16),
            query:  make(map[string]string, 8),
        }
    },
}
```

### Phase 2: fasthttp 集成（可选）

#### 为什么考虑 fasthttp？

**性能对比：**
```
net/http:  ~200 ns/op, 96 B/op
fasthttp:  ~50 ns/op,  0 B/op  (4x faster, 零分配)
```

#### 实现方案

**方案 A：完全迁移到 fasthttp（激进）**
```go
package kanggo

import (
    "github.com/valyala/fasthttp"
)

type KangGo struct {
    server   *fasthttp.Server
    router   *Router
    config   Config
}

func (k *KangGo) Run(addr string) error {
    return fasthttp.ListenAndServe(addr, k.router.Handler)
}
```

**优点：**
- 性能提升 3-4 倍
- 零内存分配
- 完全兼容 Fiber 生态

**缺点：**
- 破坏向后兼容
- 需要重写大量代码
- 学习成本高

**方案 B：双模式支持（稳健）✅ 推荐**
```go
package kanggo

import (
    "net/http"
    "github.com/valyala/fasthttp"
)

type EngineMode int

const (
    NetHTTPMode EngineMode = iota  // 默认，兼容性好
    FastHTTPMode                    // 高性能，零分配
)

type Config struct {
    // ... 现有配置
    EngineMode EngineMode  // 新增：引擎模式
}

type KangGo struct {
    mode       EngineMode
    netRouter  http.Handler
    fastRouter fasthttp.RequestHandler
}

// 统一的 Context 接口
type Context interface {
    Param(key string) string
    Query(key string) string
    JSON(code int, obj interface{}) error
    // ...
}

// net/http 实现
type NetHTTPContext struct {
    Writer  http.ResponseWriter
    Request *http.Request
    // ...
}

// fasthttp 实现
type FastHTTPContext struct {
    RequestCtx *fasthttp.RequestCtx
    // ...
}
```

**优势：**
- ✅ 向后兼容
- ✅ 渐进式升级
- ✅ 用户可选择
- ✅ 最佳实践

### Phase 3: 高级路由特性

#### 3.1 路由缓存系统
```go
type RouteCache struct {
    cache sync.Map  // 缓存已匹配的路由
    size  int       // 缓存大小限制
}

func (rc *RouteCache) Get(path string) (HandlerFunc, map[string]string, bool) {
    if v, ok := rc.cache.Load(path); ok {
        entry := v.(*CacheEntry)
        return entry.handler, entry.params, true
    }
    return nil, nil, false
}

func (rc *RouteCache) Set(path string, handler HandlerFunc, params map[string]string) {
    // LRU 策略
    rc.cache.Store(path, &CacheEntry{
        handler: handler,
        params:  params,
        time:    time.Now(),
    })
}
```

#### 3.2 正则表达式路由
```go
// 支持正则表达式路由
app.GET("/user/:id<\\d+>", handler)  // 只匹配数字
app.GET("/file/:name<.*\\.pdf>", handler)  // 只匹配 PDF 文件
```

#### 3.3 路由约束
```go
type Constraint interface {
    Match(value string) bool
}

app.GET("/user/:id", handler).
    Where("id", &NumericConstraint{Min: 1, Max: 9999})
```

### Phase 4: 中间件增强

#### 4.1 中间件链优化（已部分实现 ✅）
```go
// 预编译中间件链 - 已实现
chain := NewMiddlewareChain()
chain.Use(logger)
chain.Use(recovery)
compiled := chain.Compile(handler)

// 增强：中间件条件执行
app.Use(logger).If(func(ctx *Context) bool {
    return !strings.HasPrefix(ctx.Path(), "/health")
})

// 增强：中间件分组
app.Group("/api").Use(auth, rateLimit)
```

#### 4.2 内置高级中间件
```go
// 1. 自适应限流
import "github.com/7836246/kanggo/middleware/ratelimit"

app.Use(ratelimit.New(ratelimit.Config{
    Max:        100,              // 每秒 100 请求
    Duration:   time.Second,
    KeyFunc:    ratelimit.KeyByIP,  // 按 IP 限流
    LimitReached: func(ctx *Context) error {
        return ctx.JSON(429, map[string]string{
            "error": "Too many requests",
        })
    },
}))

// 2. 断路器
import "github.com/7836246/kanggo/middleware/circuit"

app.Use(circuit.New(circuit.Config{
    Threshold:   5,     // 5 次失败后打开
    Timeout:     30,    // 30 秒后尝试恢复
    MaxRequests: 3,     // 半开状态最多 3 个请求
}))

// 3. 请求 ID
import "github.com/7836246/kanggo/middleware/requestid"

app.Use(requestid.New())

// 4. 压缩
import "github.com/7836246/kanggo/middleware/compress"

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))
```

### Phase 5: WebSocket 和实时通信

#### 5.1 WebSocket 支持
```go
package websocket

import "github.com/fasthttp/websocket"

type Config struct {
    HandshakeTimeout time.Duration
    ReadBufferSize   int
    WriteBufferSize  int
}

// 使用示例
app.GET("/ws", websocket.New(func(conn *websocket.Conn) {
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        
        // 回显消息
        if err := conn.WriteMessage(msgType, msg); err != nil {
            break
        }
    }
}))
```

#### 5.2 Server-Sent Events (SSE)
```go
package sse

// 使用示例
app.GET("/events", func(ctx *Context) error {
    ctx.Set("Content-Type", "text/event-stream")
    ctx.Set("Cache-Control", "no-cache")
    ctx.Set("Connection", "keep-alive")
    
    for i := 0; i < 10; i++ {
        ctx.WriteString(fmt.Sprintf("data: Message %d\n\n", i))
        ctx.Flush()
        time.Sleep(time.Second)
    }
    return nil
})
```

### Phase 6: 性能监控和调优

#### 6.1 内置性能监控
```go
package monitor

type Metrics struct {
    RequestCount  uint64
    ErrorCount    uint64
    TotalDuration time.Duration
    AvgDuration   time.Duration
}

app.Use(monitor.New(monitor.Config{
    Next: func(ctx *Context) bool {
        // 跳过健康检查
        return ctx.Path() == "/health"
    },
    OnStats: func(stats *Metrics) {
        log.Printf("QPS: %d, Avg: %v", 
            stats.RequestCount, stats.AvgDuration)
    },
}))

// 访问指标
app.GET("/metrics", monitor.Handler())
```

#### 6.2 自适应性能调优
```go
package adaptive

// 根据负载自动调整池大小
type AdaptivePool struct {
    pool     *sync.Pool
    size     int
    maxSize  int
    loadAvg  float64
}

func (p *AdaptivePool) Tune() {
    if p.loadAvg > 0.8 {
        // 高负载，增加池大小
        p.size = min(p.size*2, p.maxSize)
    } else if p.loadAvg < 0.3 {
        // 低负载，减少池大小
        p.size = max(p.size/2, 10)
    }
}
```

### Phase 7: 高级特性

#### 7.1 请求验证
```go
package validator

type UserRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=130"`
}

app.POST("/user", func(ctx *Context) error {
    var req UserRequest
    if err := ctx.BindJSON(&req); err != nil {
        return err
    }
    
    if err := validator.Validate(&req); err != nil {
        return ctx.JSON(400, map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    return ctx.JSON(201, req)
})
```

#### 7.2 依赖注入
```go
package di

type Container struct {
    services map[string]interface{}
}

// 注册服务
app.Register("db", NewDatabase())
app.Register("cache", NewCache())

// 在 handler 中使用
app.GET("/user/:id", func(ctx *Context) error {
    db := ctx.Service("db").(*Database)
    user := db.Find(ctx.Param("id"))
    return ctx.JSON(200, user)
})
```

#### 7.3 优雅关闭
```go
package kanggo

func (k *KangGo) RunWithGracefulShutdown(addr string) error {
    server := &http.Server{
        Addr:    addr,
        Handler: k.Router,
    }
    
    // 启动服务器
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 5 秒超时
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return server.Shutdown(ctx)
}
```

## 📊 性能目标对比

### 当前性能 (KangGo v1.1.0)
```
静态路由:   182.6 ns/op    96 B/op     4 allocs/op
动态路由:   252.9 ns/op    96 B/op     4 allocs/op
并发请求:   55.54 ns/op    88 B/op     4 allocs/op
```

### Phase 2 目标 (fasthttp 集成)
```
静态路由:   ~45 ns/op      0 B/op      0 allocs/op  (4x faster)
动态路由:   ~65 ns/op      0 B/op      0 allocs/op  (4x faster)
并发请求:   ~15 ns/op      0 B/op      0 allocs/op  (3.7x faster)
```

### 与 Fiber V3 对比目标
```
指标         KangGo v2.0    Fiber V3      目标
------------------------------------------------------
静态路由     45 ns/op       40 ns/op      持平
动态路由     65 ns/op       60 ns/op      持平
内存分配     0 B/op         0 B/op        持平
并发能力     >20M QPS       >25M QPS      接近
```

## 🗓️ 实施时间表

### Q1 2025: Phase 1 增强
- ✅ 字节池系统
- ✅ 零拷贝字符串
- ✅ 请求/响应池
- ✅ 路由缓存

### Q2 2025: Phase 2 fasthttp 集成
- 🔄 双模式架构
- 🔄 fasthttp 适配器
- 🔄 性能测试
- 🔄 文档更新

### Q3 2025: Phase 3-4 高级特性
- 🔄 正则表达式路由
- 🔄 路由约束
- 🔄 高级中间件
- 🔄 WebSocket 支持

### Q4 2025: Phase 5-7 生态完善
- 🔄 SSE 支持
- 🔄 性能监控
- 🔄 请求验证
- 🔄 依赖注入

## 💻 示例：Fiber 风格 API

```go
package main

import (
    "github.com/7836246/kanggo"
    "github.com/7836246/kanggo/middleware/logger"
    "github.com/7836246/kanggo/middleware/recovery"
)

func main() {
    // 使用 fasthttp 引擎
    app := kanggo.New(kanggo.Config{
        EngineMode: kanggo.FastHTTPMode,  // 高性能模式
        Prefork:    true,                  // 多进程模式
    })
    
    // 中间件
    app.Use(logger.New())
    app.Use(recovery.New())
    
    // 静态文件（带缓存）
    app.Static("/", "./public", kanggo.StaticConfig{
        EnableCache:   true,
        MaxCacheSize:  10 * 1024 * 1024,  // 10MB
        Compress:      true,
    })
    
    // API 路由组
    api := app.Group("/api/v1")
    
    // RESTful 路由
    api.GET("/users", getUsers)
    api.GET("/users/:id<\\d+>", getUser)  // 正则约束
    api.POST("/users", createUser)
    api.PUT("/users/:id", updateUser)
    api.DELETE("/users/:id", deleteUser)
    
    // WebSocket
    app.GET("/ws", websocket.New(handleWebSocket))
    
    // SSE
    app.GET("/events", handleSSE)
    
    // 优雅关闭
    app.RunWithGracefulShutdown(":8080")
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required,min=3"`
    Email string `json:"email" validate:"required,email"`
}

func getUsers(ctx *kanggo.Context) error {
    users := []User{
        {ID: 1, Name: "Alice", Email: "alice@example.com"},
        {ID: 2, Name: "Bob", Email: "bob@example.com"},
    }
    return ctx.JSON(200, users)
}

func getUser(ctx *kanggo.Context) error {
    id := ctx.ParamInt("id")
    user := User{ID: id, Name: "Alice", Email: "alice@example.com"}
    return ctx.JSON(200, user)
}

func createUser(ctx *kanggo.Context) error {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": "Invalid request",
        })
    }
    
    // 验证
    if err := ctx.Validate(&user); err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": err.Error(),
        })
    }
    
    // 保存到数据库
    user.ID = 123
    
    return ctx.Status(201).JSON(user)
}
```

## 🎯 与 Fiber V3 兼容性

### API 风格对齐
```go
// Fiber V3 风格
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})

// KangGo v2.0 风格（兼容）
app.Get("/", func(ctx *kanggo.Context) error {
    return ctx.SendString("Hello, World!")
})
```

### 中间件兼容
```go
// 可以轻松移植 Fiber 中间件
func FiberMiddlewareAdapter(fiberMW fiber.Handler) kanggo.MiddlewareFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // 适配逻辑
            next(w, r)
        }
    }
}
```

## 📚 参考资源

- Fiber V3 文档: https://docs.gofiber.io/
- fasthttp 文档: https://github.com/valyala/fasthttp
- 零拷贝技术: https://en.wikipedia.org/wiki/Zero-copy
- Go 性能优化: https://go.dev/doc/diagnostics

## 🤝 贡献

欢迎参与 KangGo 的技术演进！

- 提交 Issue: https://github.com/7836246/kanggo/issues
- 提交 PR: https://github.com/7836246/kanggo/pulls
- 讨论区: https://github.com/7836246/kanggo/discussions

---

**总结：** 借鉴 Fiber V3 的技术路线，KangGo 可以实现 3-4 倍的性能提升，同时保持向后兼容和简洁易用！🚀

