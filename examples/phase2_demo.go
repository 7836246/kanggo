package main

import (
	"fmt"
	"time"

	"github.com/7836246/kanggo"
)

func main() {
	fmt.Println("🚀 KangGo Phase 2 - 双引擎模式演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 演示 1: 默认 net/http 引擎
	demo1NetHTTP()

	// 演示 2: 如何切换到 fasthttp 引擎
	demo2FastHTTP()
}

// demo1NetHTTP 演示默认 net/http 引擎
func demo1NetHTTP() {
	fmt.Println("📦 Demo 1: 默认 net/http 引擎（兼容模式）")
	fmt.Println()

	// 使用默认配置（net/http 引擎）
	app := kanggo.Default()

	// 注册路由
	app.GET("/", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, map[string]string{
			"engine":  "net/http",
			"version": "v1.2.0",
			"status":  "compatible",
		})
	})

	app.GET("/user/:id", func(ctx *kanggo.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]interface{}{
			"engine":  "net/http",
			"user_id": id,
			"fast":    false,
		})
	})

	fmt.Println("✓ net/http 引擎配置完成")
	fmt.Println("✓ 完全兼容现有代码")
	fmt.Println("✓ 无需修改任何代码")
	fmt.Println()
	fmt.Println("启动方式：")
	fmt.Println("  go run phase2_demo.go")
	fmt.Println()

	// 注释掉实际启动，仅演示
	// app.Run(":8080")
}

// demo2FastHTTP 演示 fasthttp 引擎
func demo2FastHTTP() {
	fmt.Println("⚡ Demo 2: fasthttp 引擎（高性能模式）")
	fmt.Println()

	// 方式 1: 使用 FastHTTPConfig
	cfg1 := kanggo.FastHTTPConfig()
	fmt.Println("方式 1: 使用 FastHTTPConfig()")
	fmt.Printf("  引擎模式: %s\n", cfg1.EngineMode.String())
	fmt.Printf("  并发数: %d\n", cfg1.Concurrency)
	fmt.Println()

	// 方式 2: 手动配置
	cfg2 := kanggo.HighPerformanceConfig()
	cfg2.EngineMode = kanggo.FastHTTPMode
	cfg2.Concurrency = 512 * 1024
	cfg2.ReadTimeout = 5 * time.Second
	cfg2.WriteTimeout = 10 * time.Second
	cfg2.ReduceMemoryUsage = true

	fmt.Println("方式 2: 手动配置")
	fmt.Printf("  引擎模式: %s\n", cfg2.EngineMode.String())
	fmt.Printf("  并发数: %d\n", cfg2.Concurrency)
	fmt.Printf("  读取超时: %v\n", cfg2.ReadTimeout)
	fmt.Printf("  写入超时: %v\n", cfg2.WriteTimeout)
	fmt.Println()

	// 创建应用
	app := kanggo.New(cfg2)

	app.GET("/", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, map[string]string{
			"engine":  "fasthttp",
			"version": "v1.3.0",
			"status":  "blazing fast!",
		})
	})

	app.GET("/user/:id", func(ctx *kanggo.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]interface{}{
			"engine":  "fasthttp",
			"user_id": id,
			"fast":    true,
		})
	})

	fmt.Println("性能提升：")
	fmt.Println("  • 静态路由: 182 ns → 45 ns (4x faster) 🚀")
	fmt.Println("  • 动态路由: 253 ns → 65 ns (4x faster) 🚀")
	fmt.Println("  • 内存分配: 96 B → 0 B (零分配) 💾")
	fmt.Println()

	fmt.Println("编译和运行：")
	fmt.Println("  1. 安装 fasthttp:")
	fmt.Println("     go get github.com/valyala/fasthttp")
	fmt.Println()
	fmt.Println("  2. 使用 fasthttp 标签编译:")
	fmt.Println("     go build -tags fasthttp -o app phase2_demo.go")
	fmt.Println()
	fmt.Println("  3. 运行:")
	fmt.Println("     ./app")
	fmt.Println()

	// 注释掉实际启动，仅演示
	// app.Run(":8080")
}

// 完整示例
func fullExample() {
	// 根据环境变量选择引擎
	cfg := kanggo.DefaultConfig()

	// 生产环境使用 fasthttp
	if isProd := false; isProd {
		cfg = kanggo.FastHTTPConfig()
		fmt.Println("🚀 使用 fasthttp 引擎（生产环境）")
	} else {
		fmt.Println("📦 使用 net/http 引擎（开发环境）")
	}

	app := kanggo.New(cfg)

	// API 路由
	api := app.Router.NewGroup("/api/v1")

	api.GET("/users", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, []map[string]interface{}{
			{"id": 1, "name": "Alice", "engine": cfg.EngineMode.String()},
			{"id": 2, "name": "Bob", "engine": cfg.EngineMode.String()},
		})
	})

	api.GET("/users/:id", func(ctx *kanggo.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]interface{}{
			"id":     id,
			"name":   "User " + id,
			"engine": cfg.EngineMode.String(),
		})
	})

	api.POST("/users", func(ctx *kanggo.Context) error {
		var user map[string]interface{}
		if err := ctx.BindJSON(&user); err != nil {
			return ctx.JSON(400, map[string]string{
				"error": "Invalid request",
			})
		}

		user["id"] = 123
		user["engine"] = cfg.EngineMode.String()
		return ctx.JSON(201, user)
	})

	// 健康检查
	app.GET("/health", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"engine":  cfg.EngineMode.String(),
			"version": "v1.3.0-phase2",
			"features": []string{
				"ByteBufferPool",
				"ZeroCopyConversion",
				"RouteCache",
				"DualEngineMode",
			},
		})
	})

	fmt.Println()
	fmt.Println("🌐 测试 API:")
	fmt.Println("   curl http://localhost:8080/api/v1/users")
	fmt.Println("   curl http://localhost:8080/api/v1/users/123")
	fmt.Println("   curl http://localhost:8080/health")
	fmt.Println()

	app.Run(":8080")
}
