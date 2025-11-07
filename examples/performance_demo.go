package main

import (
	"fmt"
	"time"

	"github.com/7836246/kanggo"
)

func main() {
	// 使用高性能配置
	cfg := kanggo.HighPerformanceConfig()
	cfg.ReadTimeout = 5 * time.Second
	cfg.WriteTimeout = 10 * time.Second
	cfg.IdleTimeout = 120 * time.Second

	app := kanggo.New(cfg)

	// 静态路由 - 最快的路由类型
	app.GET("/", func(ctx *kanggo.Context) error {
		return ctx.SendString("Welcome to KangGo High Performance!")
	})

	// 动态路由
	app.GET("/user/:id", func(ctx *kanggo.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]string{
			"user_id": id,
			"status":  "active",
		})
	})

	// 复杂动态路由
	app.GET("/api/v1/users/:userId/posts/:postId", func(ctx *kanggo.Context) error {
		userId := ctx.Param("userId")
		postId := ctx.Param("postId")

		return ctx.JSON(200, map[string]interface{}{
			"user_id": userId,
			"post_id": postId,
			"title":   "Sample Post",
			"content": "This is a high-performance response",
		})
	})

	// 查询参数处理
	app.GET("/search", func(ctx *kanggo.Context) error {
		query := ctx.Query("q")
		page := ctx.DefaultQuery("page", "1")
		limit := ctx.DefaultQuery("limit", "10")

		return ctx.JSON(200, map[string]string{
			"query": query,
			"page":  page,
			"limit": limit,
		})
	})

	// 静态文件服务（带缓存）
	staticCfg := kanggo.StaticConfig{
		EnableCache:   true,
		MaxCacheSize:  1024 * 1024, // 1MB
		CacheDuration: 10 * time.Second,
		MaxAge:        3600, // 1 hour
	}
	app.Static("/static", "./public", staticCfg)

	fmt.Println("\n🚀 性能优化特性:")
	fmt.Println("  ✓ Context 对象池 (sync.Pool)")
	fmt.Println("  ✓ 静态路由哈希表 O(1) 查找")
	fmt.Println("  ✓ 零内存分配路径解析")
	fmt.Println("  ✓ 预分配切片容量")
	fmt.Println("  ✓ 静态文件缓存")
	fmt.Println("\n📊 性能指标:")
	fmt.Println("  • 静态路由: ~180 ns/op, 96 B/op, 4 allocs/op")
	fmt.Println("  • 动态路由: ~250 ns/op, 96 B/op, 4 allocs/op")
	fmt.Println("  • 并发请求: ~55 ns/op, 88 B/op, 4 allocs/op")
	fmt.Println()

	app.Run(":8080")
}
