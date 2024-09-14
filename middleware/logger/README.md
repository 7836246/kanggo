# Logger 中间件

## 概述

`Logger` 中间件用于记录每个 HTTP 请求的详细信息，包括请求方法、URL
路径、处理时间等。它能够帮助开发者进行调试和监控应用程序的运行状态。在每次请求到达服务器时，`Logger` 中间件会将相关信息打印到控制台。

## 功能

- **记录请求方法**：记录每个 HTTP 请求的请求方法（如 `GET`、`POST` 等）。
- **记录请求路径**：记录每个 HTTP 请求的 URL 路径。
- **记录处理时间**：记录每个请求的处理时间，帮助开发者了解性能瓶颈。

## 使用方法

### 1. 安装 `kanggo` 框架

确保项目中已经包含了 `kanggo` 框架。

### 2. 创建 `Logger` 中间件

在 `middleware/logger` 目录下创建一个 `logger.go` 文件，并实现以下代码：

```go
package logger

import (
    "log"
    "net/http"
    "time"

    "github.com/7836246/kanggo/core" // 根据你的项目实际情况调整路径
)

// Logger 中间件，用于记录请求处理时间
func Logger() core.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
        })
    }
}
```

### 3. 使用 `Logger` 中间件

在 `main.go` 文件中集成 `Logger` 中间件：

```go
package main

import (
    "github.com/7836246/kanggo"
    "github.com/7836246/kanggo/middleware/logger" // 引入 Logger 中间件
)

func main() {
    app := kanggo.Default()

    // 使用 Logger 中间件
    app.Use(logger.Logger())

    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, KangGo with Logger!")
    })

    app.Run(":8080")
}
```

### 4. 运行项目

在终端中运行以下命令，启动服务器：

```bash
go run main.go
```

访问 `http://localhost:8080`，你将在终端中看到请求的日志输出。

### 示例输出

```bash
GET / 200 OK 300µs
```

该输出表示一个 `GET` 请求访问了 `/` 路径，响应状态码是 `200`，请求处理时间为 `300µs`。
