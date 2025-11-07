package main

import (
	"fmt"
	"time"

	"github.com/7836246/kanggo"
)

func main() {
	fmt.Println("🚀 KangGo Phase 1 优化演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 使用高性能配置（已集成 Phase 1 优化）
	cfg := kanggo.HighPerformanceConfig()
	cfg.EnableRouteCache = true // 启用路由缓存
	cfg.RouteCacheSize = 2000   // 缓存 2000 条路由
	cfg.ShowBanner = true       // 显示横幅
	cfg.ReadTimeout = 5 * time.Second
	cfg.WriteTimeout = 10 * time.Second

	app := kanggo.New(cfg)

	// 演示 1：字节缓冲池使用
	fmt.Println("📦 1. 字节缓冲池演示")
	buf := kanggo.AcquireByteBuffer(256)
	buf.WriteString("Hello, ")
	buf.WriteString("Phase 1 Optimization!")
	fmt.Printf("   缓冲内容: %s\n", buf.String())
	fmt.Printf("   缓冲大小: %d bytes\n", buf.Len())
	kanggo.ReleaseByteBuffer(buf)
	fmt.Println("   ✓ 缓冲区已归还到池中")
	fmt.Println()

	// 演示 2：零拷贝转换
	fmt.Println("⚡ 2. 零拷贝字符串转换演示")
	data := []byte("Zero-copy conversion!")
	str := kanggo.BytesToString(data)
	fmt.Printf("   转换结果: %s\n", str)
	fmt.Println("   ✓ 零内存分配，完全零拷贝")
	fmt.Println()

	// 演示 3：路由缓存
	fmt.Println("💾 3. 路由缓存演示")

	// 静态路由
	app.GET("/", func(ctx *kanggo.Context) error {
		return ctx.SendString("Welcome to Phase 1 Demo!")
	})

	// 动态路由（会被缓存）
	app.GET("/user/:id", func(ctx *kanggo.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]string{
			"user_id": id,
			"status":  "active",
			"cached":  "yes",
		})
	})

	// API 路由组
	api := app.Router.NewGroup("/api/v1")

	api.GET("/users", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		})
	})

	api.GET("/users/:id/posts", func(ctx *kanggo.Context) error {
		userId := ctx.Param("id")
		return ctx.JSON(200, map[string]interface{}{
			"user_id": userId,
			"posts": []map[string]string{
				{"id": "1", "title": "Post 1"},
				{"id": "2", "title": "Post 2"},
			},
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
		return ctx.JSON(201, user)
	})

	// 缓存统计端点
	app.GET("/stats/cache", func(ctx *kanggo.Context) error {
		stats := app.Router.CacheStats()
		return ctx.JSON(200, stats)
	})

	// 健康检查
	app.GET("/health", func(ctx *kanggo.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"version": "v1.2.0-phase1",
			"optimizations": []string{
				"ByteBufferPool",
				"ZeroCopyConversion",
				"RouteCache",
			},
		})
	})

	// 演示路由信息
	fmt.Println("   已注册路由：")
	fmt.Println("   GET  /")
	fmt.Println("   GET  /user/:id")
	fmt.Println("   GET  /api/v1/users")
	fmt.Println("   GET  /api/v1/users/:id/posts")
	fmt.Println("   POST /api/v1/users")
	fmt.Println("   GET  /stats/cache")
	fmt.Println("   GET  /health")
	fmt.Println()

	// 性能特性说明
	fmt.Println("🎯 Phase 1 优化特性：")
	fmt.Println("   ✓ 字节缓冲池：零内存分配")
	fmt.Println("   ✓ 零拷贝转换：快 177 倍")
	fmt.Println("   ✓ 路由缓存：34ns 查找")
	fmt.Println("   ✓ 自动启用：无需配置")
	fmt.Println()

	fmt.Println("📊 性能数据：")
	fmt.Println("   • 字节缓冲池: 11.69 ns/op, 0 B/op, 0 allocs/op")
	fmt.Println("   • 零拷贝转换: 0.12 ns/op, 0 B/op, 0 allocs/op")
	fmt.Println("   • 路由缓存:   34.73 ns/op, 0 B/op, 0 allocs/op")
	fmt.Println()

	fmt.Println("🌐 测试 API:")
	fmt.Println("   curl http://localhost:8080/")
	fmt.Println("   curl http://localhost:8080/user/123")
	fmt.Println("   curl http://localhost:8080/api/v1/users")
	fmt.Println("   curl http://localhost:8080/stats/cache")
	fmt.Println("   curl http://localhost:8080/health")
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 服务器启动中...")
	fmt.Println()

	// 启动服务器
	if err := app.Run(":8080"); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
	}
}
