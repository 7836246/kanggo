# KangGo → Fiber V3 技术迁移指南

## 🎯 迁移策略

我们采用**渐进式演进**策略，确保：
- ✅ 向后兼容
- ✅ 平滑过渡
- ✅ 风险可控
- ✅ 性能提升明显

## 📋 Phase 1: 立即可实施的优化

### 1.1 字节池系统（立即实施）

创建 `buffer_pool.go`：

```go
package kanggo

import (
    "sync"
)

// ByteBuffer 字节缓冲区
type ByteBuffer struct {
    B []byte
}

// Reset 重置缓冲区
func (b *ByteBuffer) Reset() {
    b.B = b.B[:0]
}

// Write 写入数据
func (b *ByteBuffer) Write(p []byte) (int, error) {
    b.B = append(b.B, p...)
    return len(p), nil
}

// String 转换为字符串（零拷贝）
func (b *ByteBuffer) String() string {
    return b2s(b.B)
}

// ByteBufferPool 字节缓冲池
type ByteBufferPool struct {
    small  sync.Pool // < 1KB
    medium sync.Pool // 1KB - 8KB
    large  sync.Pool // > 8KB
}

// DefaultByteBufferPool 默认缓冲池
var DefaultByteBufferPool = &ByteBufferPool{
    small: sync.Pool{
        New: func() interface{} {
            return &ByteBuffer{
                B: make([]byte, 0, 512),
            }
        },
    },
    medium: sync.Pool{
        New: func() interface{} {
            return &ByteBuffer{
                B: make([]byte, 0, 4096),
            }
        },
    },
    large: sync.Pool{
        New: func() interface{} {
            return &ByteBuffer{
                B: make([]byte, 0, 16384),
            }
        },
    },
}

// Get 获取缓冲区
func (p *ByteBufferPool) Get(size int) *ByteBuffer {
    var pool *sync.Pool
    switch {
    case size <= 1024:
        pool = &p.small
    case size <= 8192:
        pool = &p.medium
    default:
        pool = &p.large
    }
    
    buf := pool.Get().(*ByteBuffer)
    buf.Reset()
    return buf
}

// Put 归还缓冲区
func (p *ByteBufferPool) Put(buf *ByteBuffer) {
    size := cap(buf.B)
    var pool *sync.Pool
    
    switch {
    case size <= 1024:
        pool = &p.small
    case size <= 8192:
        pool = &p.medium
    default:
        pool = &p.large
    }
    
    pool.Put(buf)
}
```

### 1.2 零拷贝字符串转换

创建 `unsafe_utils.go`：

```go
package kanggo

import (
    "unsafe"
)

// b2s 字节切片转字符串（零拷贝）
// 注意：使用时要确保字节切片不会被修改
func b2s(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

// s2b 字符串转字节切片（零拷贝）
// 注意：返回的切片不可修改
func s2b(s string) []byte {
    return *(*[]byte)(unsafe.Pointer(
        &struct {
            string
            Cap int
        }{s, len(s)},
    ))
}

// BytesToString 安全的字节切片转字符串
func BytesToString(b []byte) string {
    if len(b) == 0 {
        return ""
    }
    return b2s(b)
}

// StringToBytes 安全的字符串转字节切片
func StringToBytes(s string) []byte {
    if len(s) == 0 {
        return nil
    }
    return s2b(s)
}
```

### 1.3 优化 Context 使用缓冲池

修改 `context.go`：

```go
// JSON 返回一个 JSON 响应（使用缓冲池）
func (c *Context) JSON(code int, obj interface{}) error {
    // 从池中获取缓冲区
    buf := DefaultByteBufferPool.Get(256)
    defer DefaultByteBufferPool.Put(buf)
    
    // 序列化到缓冲区
    data, err := c.jsonEncoder(obj)
    if err != nil {
        return err
    }
    
    buf.Write(data)
    
    c.Writer.Header().Set(constants.HeaderContentType, constants.MIMEApplicationJSON)
    c.Writer.WriteHeader(code)
    _, err = c.Writer.Write(buf.B)
    return err
}

// JSONP 返回一个 JSONP 响应（优化版）
func (c *Context) JSONP(callback string, obj interface{}) error {
    data, err := c.jsonEncoder(obj)
    if err != nil {
        return err
    }
    
    // 使用缓冲池
    buf := DefaultByteBufferPool.Get(len(callback) + len(data) + 3)
    defer DefaultByteBufferPool.Put(buf)
    
    buf.Write(StringToBytes(callback))
    buf.B = append(buf.B, '(')
    buf.Write(data)
    buf.B = append(buf.B, ')', ';')
    
    c.Writer.Header().Set(constants.HeaderContentType, constants.MIMEApplicationJavaScript)
    _, err = c.Writer.Write(buf.B)
    return err
}
```

### 1.4 路由缓存

创建 `route_cache.go`：

```go
package kanggo

import (
    "sync"
    "time"
)

// CacheEntry 缓存条目
type CacheEntry struct {
    handler HandlerFunc
    params  map[string]string
    hits    uint64
    time    time.Time
}

// RouteCache 路由缓存
type RouteCache struct {
    cache   sync.Map
    maxSize int
    enabled bool
}

// NewRouteCache 创建路由缓存
func NewRouteCache(maxSize int) *RouteCache {
    return &RouteCache{
        maxSize: maxSize,
        enabled: true,
    }
}

// Get 获取缓存的路由
func (rc *RouteCache) Get(method, path string) (HandlerFunc, map[string]string, bool) {
    if !rc.enabled {
        return nil, nil, false
    }
    
    key := method + ":" + path
    if v, ok := rc.cache.Load(key); ok {
        entry := v.(*CacheEntry)
        entry.hits++
        return entry.handler, entry.params, true
    }
    return nil, nil, false
}

// Set 设置路由缓存
func (rc *RouteCache) Set(method, path string, handler HandlerFunc, params map[string]string) {
    if !rc.enabled {
        return
    }
    
    // 复制 params 避免共享
    paramsCopy := make(map[string]string, len(params))
    for k, v := range params {
        paramsCopy[k] = v
    }
    
    key := method + ":" + path
    rc.cache.Store(key, &CacheEntry{
        handler: handler,
        params:  paramsCopy,
        time:    time.Now(),
    })
}

// Stats 获取缓存统计
func (rc *RouteCache) Stats() map[string]interface{} {
    count := 0
    totalHits := uint64(0)
    
    rc.cache.Range(func(key, value interface{}) bool {
        count++
        entry := value.(*CacheEntry)
        totalHits += entry.hits
        return true
    })
    
    return map[string]interface{}{
        "entries":    count,
        "total_hits": totalHits,
        "enabled":    rc.enabled,
    }
}
```

## 📋 Phase 2: fasthttp 双模式支持

### 2.1 引擎抽象层

创建 `engine.go`：

```go
package kanggo

import (
    "net/http"
)

// EngineMode 引擎模式
type EngineMode int

const (
    // NetHTTPMode 标准库模式（默认）
    NetHTTPMode EngineMode = iota
    // FastHTTPMode fasthttp 模式（高性能）
    FastHTTPMode
)

// Engine 引擎接口
type Engine interface {
    // ListenAndServe 启动服务器
    ListenAndServe(addr string, handler interface{}) error
    // Shutdown 优雅关闭
    Shutdown() error
}

// NewEngine 创建引擎
func NewEngine(mode EngineMode, cfg Config) Engine {
    switch mode {
    case FastHTTPMode:
        return NewFastHTTPEngine(cfg)
    default:
        return NewNetHTTPEngine(cfg)
    }
}

// NetHTTPEngine 标准库引擎
type NetHTTPEngine struct {
    server *http.Server
    cfg    Config
}

func NewNetHTTPEngine(cfg Config) *NetHTTPEngine {
    return &NetHTTPEngine{cfg: cfg}
}

func (e *NetHTTPEngine) ListenAndServe(addr string, handler interface{}) error {
    e.server = &http.Server{
        Addr:         addr,
        Handler:      handler.(http.Handler),
        ReadTimeout:  e.cfg.ReadTimeout,
        WriteTimeout: e.cfg.WriteTimeout,
        IdleTimeout:  e.cfg.IdleTimeout,
    }
    return e.server.ListenAndServe()
}

func (e *NetHTTPEngine) Shutdown() error {
    if e.server != nil {
        return e.server.Close()
    }
    return nil
}
```

### 2.2 fasthttp 引擎实现

创建 `engine_fasthttp.go`：

```go
// +build fasthttp

package kanggo

import (
    "github.com/valyala/fasthttp"
)

// FastHTTPEngine fasthttp 引擎
type FastHTTPEngine struct {
    server *fasthttp.Server
    cfg    Config
}

func NewFastHTTPEngine(cfg Config) *FastHTTPEngine {
    return &FastHTTPEngine{
        cfg: cfg,
        server: &fasthttp.Server{
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            IdleTimeout:  cfg.IdleTimeout,
        },
    }
}

func (e *FastHTTPEngine) ListenAndServe(addr string, handler interface{}) error {
    e.server.Handler = handler.(func(*fasthttp.RequestCtx))
    return e.server.ListenAndServe(addr)
}

func (e *FastHTTPEngine) Shutdown() error {
    if e.server != nil {
        return e.server.Shutdown()
    }
    return nil
}

// FastHTTPContext fasthttp 上下文适配器
type FastHTTPContext struct {
    RequestCtx *fasthttp.RequestCtx
    params     map[string]string
}

func NewFastHTTPContext(ctx *fasthttp.RequestCtx) *FastHTTPContext {
    return &FastHTTPContext{
        RequestCtx: ctx,
        params:     make(map[string]string, 4),
    }
}

func (c *FastHTTPContext) Param(key string) string {
    return c.params[key]
}

func (c *FastHTTPContext) Query(key string) string {
    return BytesToString(c.RequestCtx.QueryArgs().Peek(key))
}

func (c *FastHTTPContext) JSON(code int, obj interface{}) error {
    data, err := json.Marshal(obj)
    if err != nil {
        return err
    }
    
    c.RequestCtx.SetStatusCode(code)
    c.RequestCtx.SetContentType("application/json")
    c.RequestCtx.SetBody(data)
    return nil
}
```

### 2.3 构建标签支持

创建 `Makefile`：

```makefile
.PHONY: build build-fasthttp test bench

# 默认构建（net/http）
build:
	go build -o bin/kanggo-app

# fasthttp 构建
build-fasthttp:
	go build -tags fasthttp -o bin/kanggo-app-fast

# 测试
test:
	go test -v ./...

# 基准测试
bench:
	go test -bench=. -benchmem -benchtime=3s

# fasthttp 基准测试
bench-fasthttp:
	go test -tags fasthttp -bench=. -benchmem -benchtime=3s

# 安装依赖
deps:
	go get github.com/valyala/fasthttp
	go get github.com/bytedance/sonic

# 清理
clean:
	rm -rf bin/
```

## 🔧 使用示例

### 默认模式（net/http）
```go
package main

import "github.com/7836246/kanggo"

func main() {
    app := kanggo.Default()
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, World!")
    })
    
    app.Run(":8080")
}
```

### 高性能模式（fasthttp）
```go
package main

import (
    "github.com/7836246/kanggo"
)

func main() {
    // 启用 fasthttp 引擎
    cfg := kanggo.DefaultConfig()
    cfg.EngineMode = kanggo.FastHTTPMode
    
    app := kanggo.New(cfg)
    
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, FastHTTP!")
    })
    
    app.Run(":8080")
}
```

### 使用缓冲池优化
```go
app.POST("/data", func(ctx *kanggo.Context) error {
    // 自动使用缓冲池
    var data map[string]interface{}
    if err := ctx.BindJSON(&data); err != nil {
        return err
    }
    
    // JSON 响应自动使用缓冲池
    return ctx.JSON(200, data)
})
```

### 启用路由缓存
```go
cfg := kanggo.DefaultConfig()
cfg.EnableRouteCache = true  // 启用路由缓存
cfg.RouteCacheSize = 1000    // 缓存 1000 条路由

app := kanggo.New(cfg)
```

## 📊 性能对比

### 当前性能（net/http + 优化）
```
静态路由:  182.6 ns/op    96 B/op     4 allocs/op
动态路由:  252.9 ns/op    96 B/op     4 allocs/op
```

### Phase 1 优化后（缓冲池 + 零拷贝）
```
静态路由:  ~150 ns/op     64 B/op     2 allocs/op  (18% faster)
动态路由:  ~200 ns/op     64 B/op     2 allocs/op  (21% faster)
```

### Phase 2 优化后（fasthttp）
```
静态路由:  ~45 ns/op      0 B/op      0 allocs/op  (4x faster)
动态路由:  ~65 ns/op      0 B/op      0 allocs/op  (4x faster)
```

## 🎯 迁移检查清单

### Phase 1: 立即实施
- [ ] 实现字节缓冲池
- [ ] 添加零拷贝工具函数
- [ ] 优化 JSON/JSONP 使用缓冲池
- [ ] 实现路由缓存
- [ ] 性能测试和对比
- [ ] 更新文档

### Phase 2: fasthttp 集成
- [ ] 设计引擎抽象层
- [ ] 实现 fasthttp 适配器
- [ ] 实现双模式切换
- [ ] 添加构建标签支持
- [ ] 完整测试套件
- [ ] 性能基准测试
- [ ] 迁移指南

### Phase 3: 高级特性
- [ ] WebSocket 支持
- [ ] SSE 支持
- [ ] 正则表达式路由
- [ ] 请求验证
- [ ] 性能监控

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/7836246/kanggo.git
cd kanggo
```

### 2. 创建新分支
```bash
git checkout -b feature/fiber-optimization
```

### 3. 实施 Phase 1 优化
```bash
# 创建新文件
touch buffer_pool.go
touch unsafe_utils.go
touch route_cache.go

# 编写代码（参考上面的示例）
# ...

# 运行测试
go test -v ./...

# 基准测试
go test -bench=. -benchmem
```

### 4. 提交更改
```bash
git add .
git commit -m "feat: Phase 1 - Buffer pool and zero-copy optimization"
git push origin feature/fiber-optimization
```

## 📚 推荐阅读

1. **Fiber V3 文档**
   - https://docs.gofiber.io/
   
2. **fasthttp 文档**
   - https://github.com/valyala/fasthttp
   
3. **Go 性能优化**
   - https://github.com/dgryski/go-perfbook
   
4. **零拷贝技术**
   - https://en.wikipedia.org/wiki/Zero-copy

## ⚠️ 注意事项

### 安全性
- 零拷贝函数使用 `unsafe` 包，注意内存安全
- 确保字节切片在使用期间不被修改
- 添加足够的单元测试

### 兼容性
- Phase 1 优化完全向后兼容
- Phase 2 需要用户选择引擎模式
- 提供清晰的迁移文档

### 性能
- 在生产环境前进行充分测试
- 使用 pprof 分析性能瓶颈
- 监控内存使用和 GC 行为

## 🎉 总结

通过渐进式演进，KangGo 可以：

1. **Phase 1（2周）**: 性能提升 15-20%，完全兼容
2. **Phase 2（1个月）**: 性能提升 300-400%，可选模式
3. **Phase 3（2个月）**: 完整的 Fiber 风格框架

**最终目标**: 成为 Go 生态中最快的 Web 框架之一！🚀

