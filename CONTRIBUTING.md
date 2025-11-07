# 贡献指南

感谢你对 KangGo 的关注！我们欢迎各种形式的贡献。

## 📋 贡献方式

### 1. 报告 Bug 🐛
- 在 GitHub Issues 中创建新 issue
- 描述问题和复现步骤
- 提供系统环境信息

### 2. 提出功能建议 💡
- 在 GitHub Issues 中标记为 `enhancement`
- 详细描述功能需求和使用场景
- 说明为什么这个功能有价值

### 3. 改进文档 📝
- 修正文档中的错误
- 添加示例和说明
- 翻译文档

### 4. 提交代码 🔧
- Fork 项目
- 创建功能分支
- 提交 Pull Request

## 🔧 开发流程

### 1. 设置开发环境

```bash
# 克隆项目
git clone https://github.com/7836246/kanggo.git
cd kanggo

# 安装依赖
make install-deps

# 运行测试
make test
```

### 2. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 3. 编写代码

**代码规范：**
- 遵循 Go 代码规范
- 添加必要的注释
- 编写单元测试
- 运行 `make fmt` 格式化代码

### 4. 运行测试

```bash
# 运行所有测试
make test

# 运行基准测试
make bench

# 代码检查
make lint
```

### 5. 提交代码

```bash
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name
```

**提交信息格式：**
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `perf:` 性能优化
- `test:` 测试相关
- `chore:` 构建/工具相关

### 6. 创建 Pull Request

- 在 GitHub 上创建 PR
- 填写 PR 模板
- 等待代码审查

## 📝 代码规范

### Go 代码风格

```go
// ✅ 好的例子
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctx := AcquireContext(w, req, r.config)
    defer ReleaseContext(ctx)
    
    // ... 处理逻辑
}

// ❌ 不好的例子
func (r *Router) ServeHTTP(w http.ResponseWriter,req *http.Request){
    ctx:=AcquireContext(w,req,r.config)
    // 没有 defer 释放
    // ... 处理逻辑
}
```

### 文档注释

```go
// ✅ 好的注释
// AcquireContext 从对象池中获取一个 Context 实例
// 使用完毕后必须调用 ReleaseContext 归还
func AcquireContext(w http.ResponseWriter, req *http.Request, cfg Config) *Context {
    // ...
}

// ❌ 不好的注释
// 获取 ctx
func AcquireContext(w http.ResponseWriter, req *http.Request, cfg Config) *Context {
    // ...
}
```

## 🧪 测试要求

### 单元测试

```go
func TestRouterStaticRoute(t *testing.T) {
    // 准备
    router := NewRouter(DefaultConfig())
    router.RegisterStaticRoute("GET", "/hello", func(ctx *Context) error {
        return ctx.SendString("Hello")
    })
    
    // 执行
    req := httptest.NewRequest("GET", "/hello", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // 断言
    if w.Code != 200 {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

### 基准测试

```go
func BenchmarkStaticRoute(b *testing.B) {
    app := Default()
    app.GET("/hello", func(ctx *Context) error {
        return ctx.SendString("Hello")
    })
    
    req := httptest.NewRequest("GET", "/hello", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        app.Router.ServeHTTP(w, req)
    }
}
```

## 📚 文档规范

### 文档结构

```markdown
# 标题

## 概述

简短描述功能

## 使用方法

### 基础用法

​```go
// 示例代码
​```

### 高级用法

​```go
// 高级示例
​```

## 注意事项

- 注意点 1
- 注意点 2
```

## 🎯 贡献重点

### 当前优先级

**P0 - 最高优先级：**
- Bug 修复
- 性能优化
- 安全问题

**P1 - 高优先级：**
- Phase 3 功能实现
- 文档完善
- 示例程序

**P2 - 中优先级：**
- 新功能开发
- 工具优化
- 测试覆盖

**P3 - 低优先级：**
- 代码重构
- 样式调整
- 注释改进

## 🏆 贡献者

感谢所有为 KangGo 做出贡献的开发者！

**核心贡献者：**
- @7836246 - 项目创建者和维护者

**贡献者列表：**
（在此列出所有贡献者）

## 📞 联系方式

**GitHub：** https://github.com/7836246/kanggo
**Issues：** https://github.com/7836246/kanggo/issues
**Discussions：** https://github.com/7836246/kanggo/discussions

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

**感谢你的贡献！让我们一起把 KangGo 打造成最好的 Go Web 框架！** 🚀

