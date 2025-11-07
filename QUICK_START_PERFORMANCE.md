# KangGo 高性能快速开始指南

## 🚀 5分钟上手高性能 KangGo

### 1. 安装

```bash
go get -u github.com/7836246/kanggo@latest
```

### 2. 基础高性能配置

```go
package main

import (
    "github.com/7836246/kanggo"
    "time"
)

func main() {
    // 使用高性能配置
    app := kanggo.New(kanggo.HighPerformanceConfig())
    
    // 配置超时
    app.Config.ReadTimeout = 5 * time.Second
    app.Config.WriteTimeout = 10 * time.Second
    
    // 注册路由
    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, High Performance!")
    })
    
    app.Run(":8080")
}
```

### 3. 性能优化特性

#### 自动启用的优化
以下优化无需配置，框架自动启用：

✅ **Context 对象池** - 自动复用 Context 对象
✅ **静态路由哈希表** - O(1) 路由查找
✅ **零内存分配路径解析** - 减少内存分配
✅ **预分配容量** - 优化内存使用

#### 可选的优化配置

##### 静态文件缓存

```go
cfg := kanggo.StaticConfig{
    EnableCache:   true,              // 启用缓存
    MaxCacheSize:  1024 * 1024,       // 1MB
    CacheDuration: 10 * time.Second,  // 缓存10秒
    MaxAge:        3600,               // 浏览器缓存1小时
}
app.Static("/static", "./public", cfg)
```

##### 高性能 JSON（可选）

```go
import "github.com/bytedance/sonic"

cfg := kanggo.DefaultConfig()
cfg.JSONEncoder = sonic.Marshal
cfg.JSONDecoder = sonic.Unmarshal
app := kanggo.New(cfg)
```

### 4. 性能最佳实践

#### ✅ DO - 推荐做法

```go
// 1. 优先使用静态路由（最快）
app.GET("/api/users", handleUsers)

// 2. 合理设计动态路由（不超过5层）
app.GET("/api/users/:id", handleUser)

// 3. 使用高性能配置
app := kanggo.New(kanggo.HighPerformanceConfig())

// 4. 预分配已知大小的切片
users := make([]User, 0, 100)

// 5. 复用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}
```

#### ❌ DON'T - 避免做法

```go
// 1. 避免过深的动态路由层级
app.GET("/a/:b/c/:d/e/:f/g/:h", handler) // 太深了

// 2. 避免在热路径创建临时对象
app.GET("/hot", func(ctx *kanggo.Context) error {
    // 不好：每次请求都创建新对象
    data := make([]byte, 1024)
    return ctx.SendString(string(data))
})

// 3. 避免注册过多不必要的中间件
app.Use(middleware1)
app.Use(middleware2)
app.Use(middleware3) // 只注册必要的

// 4. 避免在请求处理中进行耗时操作
app.GET("/slow", func(ctx *kanggo.Context) error {
    time.Sleep(5 * time.Second) // 不好！使用异步处理
    return ctx.SendString("Done")
})
```

### 5. 性能监控

#### 启用 pprof

```go
import _ "net/http/pprof"

func main() {
    // 启动 pprof 服务器
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    app := kanggo.New(kanggo.HighPerformanceConfig())
    // ... 路由配置
    app.Run(":8080")
}
```

#### 查看性能数据

```bash
# CPU Profile
go tool pprof http://localhost:6060/debug/pprof/profile

# Memory Profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### 6. 基准测试

#### 运行框架基准测试

```bash
cd $GOPATH/src/github.com/7836246/kanggo
go test -bench=. -benchmem -benchtime=3s
```

#### 编写自己的基准测试

```go
func BenchmarkMyHandler(b *testing.B) {
    app := kanggo.New(kanggo.HighPerformanceConfig())
    app.GET("/test", myHandler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        app.Router.ServeHTTP(w, req)
        w.Body.Reset()
    }
}
```

### 7. 完整示例

```go
package main

import (
    "github.com/7836246/kanggo"
    "time"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // 高性能配置
    cfg := kanggo.HighPerformanceConfig()
    cfg.ReadTimeout = 5 * time.Second
    cfg.WriteTimeout = 10 * time.Second
    cfg.IdleTimeout = 120 * time.Second
    cfg.MaxRequestBodySize = 10 * 1024 * 1024 // 10MB
    
    app := kanggo.New(cfg)
    
    // API 路由组
    api := app.Router.NewGroup("/api")
    
    // 静态路由（最快）
    api.GET("/health", func(ctx *kanggo.Context) error {
        return ctx.JSON(200, map[string]string{
            "status": "ok",
        })
    })
    
    // 动态路由
    api.GET("/users/:id", func(ctx *kanggo.Context) error {
        id := ctx.Param("id")
        return ctx.JSON(200, User{
            ID:   1,
            Name: "User " + id,
        })
    })
    
    // POST 请求
    api.POST("/users", func(ctx *kanggo.Context) error {
        var user User
        if err := ctx.BindJSON(&user); err != nil {
            return ctx.JSON(400, map[string]string{
                "error": "Invalid request",
            })
        }
        return ctx.JSON(201, user)
    })
    
    // 静态文件（带缓存）
    staticCfg := kanggo.StaticConfig{
        EnableCache:   true,
        MaxCacheSize:  1024 * 1024,
        CacheDuration: 10 * time.Second,
        MaxAge:        3600,
    }
    app.Static("/static", "./public", staticCfg)
    
    // 启动服务器
    app.Run(":8080")
}
```

### 8. 性能指标参考

基于 Intel Core i5-13500H, Go 1.23.1：

| 场景 | 性能 | 内存 | 说明 |
|------|------|------|------|
| 静态路由 | ~180 ns/op | 96 B/op | 最快 |
| 动态路由 | ~250 ns/op | 96 B/op | 很快 |
| JSON 响应 | ~1000 ns/op | 785 B/op | 正常 |
| 并发请求 | ~55 ns/op | 88 B/op | 极快 |

### 9. 常见问题

#### Q: 如何进一步提升性能？

A: 
1. 使用 sonic JSON 库
2. 启用静态文件缓存
3. 减少中间件数量
4. 使用连接池
5. 启用 HTTP/2

#### Q: 内存占用如何优化？

A:
1. 框架已自动使用 Context 对象池
2. 预分配已知大小的切片和 map
3. 复用临时对象
4. 避免在热路径创建大对象

#### Q: 如何测试我的应用性能？

A:
```bash
# 使用 wrk 压测
wrk -t12 -c400 -d30s http://localhost:8080/

# 使用 ab 压测
ab -n 10000 -c 100 http://localhost:8080/

# 使用 hey 压测
hey -n 10000 -c 100 http://localhost:8080/
```

### 10. 下一步

- 📖 阅读 [PERFORMANCE.md](PERFORMANCE.md) 了解详细优化
- 📊 查看 [OPTIMIZATION_REPORT.md](OPTIMIZATION_REPORT.md) 了解优化细节
- 🔧 运行 `examples/performance_demo.go` 查看示例
- 🧪 运行基准测试了解性能表现

---

**需要帮助？**
- GitHub Issues: https://github.com/7836246/kanggo/issues
- 文档: https://github.com/7836246/kanggo

**Happy Coding! 🚀**

