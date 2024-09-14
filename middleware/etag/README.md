# ETag Middleware for KangGo

`ETag` 中间件用于 KangGo 框架，通过为响应内容生成 `ETag`（实体标签），让缓存更加高效，减少带宽使用。当请求内容未更改时，Web
服务器无需重新发送完整的响应。

## 功能

- 计算响应内容的 `ETag` 值。
- 检查 `If-None-Match` 请求头，如果匹配则返回 `304 Not Modified` 状态码，而不是完整响应内容。

## 安装

确保你已经安装并配置好 [KangGo](https://github.com/7836246/kanggo) 框架。

## 使用方法

要在 KangGo 应用程序中使用 `ETag` 中间件，请按照以下步骤操作：

### 1. 导入包

```go
import (
"github.com/7836246/kanggo"
"github.com/7836246/kanggo/middleware/etag"
)
```

### 2. 应用中间件

在 KangGo 应用程序中应用 `ETag` 中间件：

```go
app := kanggo.Default()

// 使用 ETag 中间件
app.Use(etag.ETag())
```

### 3. 注册路由

注册一个简单的路由来测试 `ETag` 中间件：

```go
app.GET("/etag", func (ctx *kanggo.Context) error {
return ctx.SendString("Hello, KangGo with ETag!")
})
```

### 4. 运行应用

```go
app.Run(":8080")
```

## 示例

下面是一个完整的示例应用：

```go
package main

import (
    "github.com/7836246/kanggo"
    "github.com/7836246/kanggo/middleware/etag"
)

func main() {
    app := kambe.Default()

    // 使用 ETag 中间件
    app.Use(etag.ETag())

    // 注册路由
    app.GET("/etag", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, KangGo with ETag!")
    })

    // 运行应用
    app.Run(":8080")
}
```

## 测试

确保你已经安装了 `Go` 语言环境，可以通过运行以下命令来测试 `ETag` 中间件：

```bash
go test -v
```

测试文件可以使用以下代码：

```go
package etag

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/7836246/kanggo"
)

func TestETagMiddleware(t *testing.T) {
    app := kambe.Default()
    app.Use(ETag())

    app.GET("/etag", func(ctx *kanggo.Context) error {
        return ctx.SendString("Hello, KangGo with ETag!")
    })

    req, _ := http.NewRequest("GET", "/etag", nil)
    resp := httptest.NewRecorder()
    app.Router.ServeHTTP(resp, req)

    if status := resp.Code; status != http.StatusOK {
        t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
    }

    etag := resp.Header().Get("ETag")
    if etag == "" {
        t.Error("未生成 ETag 头")
    }

    req2, _ := http.NewRequest("GET", "/etag", nil)
    req2.Header.Set("If-None-Match", etag)
    resp2 := httptest.NewRecorder()
    app.Router.ServeHTTP(resp2, req2)

    if status := resp2.Code; status != http.StatusNotModified {
        t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusNotModified)
    }
}
```

## 许可证

本项目遵循 [MIT 许可证](LICENSE)。
