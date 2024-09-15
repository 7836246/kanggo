# Recovery Middleware for KangGo

## 介绍

`Recovery` 中间件用于在路由处理函数中捕获 panic，并防止程序崩溃。它可以捕获运行时的 panic，并返回 HTTP 500
错误码，同时打印错误堆栈信息，帮助开发者快速定位问题。

## 使用方法

### 引入中间件

```go
import (
    "github.com/7836246/kanggo"
    "github.com/7836246/kanggo/middleware/recovery"
)

func main() {
    app := kanggo.Default()

    // 使用 Recovery 中间件
    app.Use(recovery.Recovery())

    app.GET("/panic", func(ctx *kanggo.Context) error {
        panic("这是一个测试 panic!")
    })

    app.Run(":8080")
}
```

### 测试

`Recovery` 中间件附带了一个测试用例，确保其功能正常。你可以通过以下命令运行测试：

```bash
go test -v github.com/7836246/kanggo/middleware/recovery
```

测试通过时的输出：

```
=== RUN   TestRecoveryMiddleware
捕获到 panic: 这是一个测试 panic!
... (堆栈信息)
--- PASS: TestRecoveryMiddleware (0.00s)
PASS
```

## 效果

`Recovery` 中间件确保了在发生 panic 时，程序不会崩溃，并且提供了详细的堆栈信息来帮助调试。

## 贡献

欢迎贡献代码和提出建议！
