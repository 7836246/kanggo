# CORS 中间件

## 概述

`CORS`（跨域资源共享）中间件用于解决浏览器的跨域请求限制问题。通过配置 `CORS` 中间件，你可以指定允许哪些域、方法和请求头访问你的
API。

## 功能

- **允许的跨域请求**：通过设置允许的域名、HTTP 方法和请求头，控制哪些请求可以跨域访问。
- **处理复杂请求**：处理预检请求（`OPTIONS` 请求），确保客户端能够正常发起跨域请求。
- **增强安全性**：可以通过配置限制跨域请求的来源和类型，增强 API 的安全性。

## 使用方法

### 1. 安装 `kanggo` 框架

确保项目中已经包含了 `kanggo` 框架。

### 2. 创建 `CORS` 中间件

在 `middleware/cors` 目录下创建一个 `cors.go` 文件，并实现以下代码：

```go
package cors

import (
    "net/http"

    "github.com/7836246/kanggo/core" // 根据你的项目实际情况调整路径
)

// CORS 中间件，用于处理跨域请求
func CORS(allowedOrigins string, allowedMethods string, allowedHeaders string) core.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
            w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
            w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
            
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### 3. 使用 `CORS` 中间件

在 `main.go` 文件中集成 `CORS` 中间件：

```go
package main

import (
    "github.com/7836246/kanggo"
    "github.com/7836246/kanggo/middleware/cors" // 引入 CORS 中间件
)

func main() {
    app := bakenggo.Default()

    // 使用 CORS 中间件
    app.Use(cors.CORS("*", "GET,POST,OPTIONS", "Content-Type"))

    app.GET("/", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, KangGo with CORS!")
    })

    app.Run(":8080")
}
```

### 4. 运行项目

在终端中运行以下命令，启动服务器：

```bash
go run main.go
```

访问 `http://localhost:8080`，你将能够从不同的域名进行跨域访问。

## 说明

可以根据需求调整允许的域名、方法和请求头，以适应特定的安全要求。
