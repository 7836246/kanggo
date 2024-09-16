# Session Middleware for KangGo

本仓库包含 `KangGo` 框架的会话（Session）中间件实现。会话中间件允许开发者在请求之间存储和访问用户数据。

## 功能特性

- 使用内存存储会话数据（MemoryStore）。
- 支持自定义会话存储。
- 自动生成随机 `session_id`。
- 提供设置和获取会话数据的功能。

## 使用方法

1. 创建会话存储实例，例如 `MemoryStore`：

```go
store := NewMemoryStore()
```

2. 创建中间件并注册到路由器：

```go
app.Use(session.Middleware(store))
```

3. 在处理程序中获取和设置会话数据：

```go
func handler(ctx *session.Context) {
ctx.SetSessionValue("key", "value")
value := ctx.GetSessionValue("key")
}
```

## 测试

你可以使用以下命令运行测试：

```bash
go test -v ./middleware/session
```

## 许可证

此项目基于 MIT 许可证，详情请参阅 LICENSE 文件。
