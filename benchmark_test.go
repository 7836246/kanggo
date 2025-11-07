package kanggo

import (
	"net/http/httptest"
	"testing"
)

// BenchmarkStaticRoute 测试静态路由性能
func BenchmarkStaticRoute(b *testing.B) {
	app := Default()
	app.GET("/hello", func(ctx *Context) error {
		return ctx.SendString("Hello, World!")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkDynamicRoute 测试动态路由性能
func BenchmarkDynamicRoute(b *testing.B) {
	app := Default()
	app.GET("/user/:id", func(ctx *Context) error {
		id := ctx.Param("id")
		return ctx.SendString("User ID: " + id)
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkComplexDynamicRoute 测试复杂动态路由性能
func BenchmarkComplexDynamicRoute(b *testing.B) {
	app := Default()
	app.GET("/api/v1/users/:userId/posts/:postId/comments/:commentId", func(ctx *Context) error {
		return ctx.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/api/v1/users/123/posts/456/comments/789", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkJSONResponse 测试 JSON 响应性能
func BenchmarkJSONResponse(b *testing.B) {
	app := Default()

	type Response struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	app.GET("/json", func(ctx *Context) error {
		return ctx.JSON(200, Response{
			Status:  "success",
			Message: "Hello, World!",
			Data: map[string]interface{}{
				"id":   123,
				"name": "Test User",
				"age":  25,
			},
		})
	})

	req := httptest.NewRequest("GET", "/json", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkContextPool 测试 Context 对象池性能
func BenchmarkContextPool(b *testing.B) {
	cfg := DefaultConfig()
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := AcquireContext(w, req, cfg)
		ReleaseContext(ctx)
	}
}

// BenchmarkMultipleRoutes 测试多路由性能
func BenchmarkMultipleRoutes(b *testing.B) {
	app := Default()

	// 注册多个路由
	routes := []string{
		"/api/users",
		"/api/posts",
		"/api/comments",
		"/api/likes",
		"/api/shares",
	}

	for _, route := range routes {
		app.GET(route, func(ctx *Context) error {
			return ctx.SendString("OK")
		})
	}

	req := httptest.NewRequest("GET", "/api/comments", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkParallelRequests 测试并发请求性能
func BenchmarkParallelRequests(b *testing.B) {
	app := Default()
	app.GET("/parallel", func(ctx *Context) error {
		return ctx.SendString("OK")
	})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/parallel", nil)
		w := httptest.NewRecorder()

		for pb.Next() {
			app.Router.ServeHTTP(w, req)
			w.Body.Reset()
		}
	})
}

// BenchmarkQueryParams 测试查询参数性能
func BenchmarkQueryParams(b *testing.B) {
	app := Default()
	app.GET("/search", func(ctx *Context) error {
		q := ctx.Query("q")
		page := ctx.Query("page")
		limit := ctx.Query("limit")
		return ctx.SendString(q + page + limit)
	})

	req := httptest.NewRequest("GET", "/search?q=test&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

// BenchmarkHighPerformanceConfig 测试高性能配置
func BenchmarkHighPerformanceConfig(b *testing.B) {
	cfg := HighPerformanceConfig()
	app := New(cfg)

	app.GET("/hp", func(ctx *Context) error {
		return ctx.JSON(200, map[string]string{
			"status":  "ok",
			"message": "high performance",
		})
	})

	req := httptest.NewRequest("GET", "/hp", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}
