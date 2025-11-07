package kanggo

import (
	"net/http/httptest"
	"testing"
)

// BenchmarkByteBufferPool 测试字节缓冲池性能
func BenchmarkByteBufferPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := AcquireByteBuffer(256)
		buf.WriteString("Hello, World!")
		_ = buf.String()
		ReleaseByteBuffer(buf)
	}
}

// BenchmarkByteBufferPoolLarge 测试大缓冲区
func BenchmarkByteBufferPoolLarge(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := AcquireByteBuffer(8192)
		for j := 0; j < 100; j++ {
			buf.WriteString("Test data ")
		}
		_ = buf.String()
		ReleaseByteBuffer(buf)
	}
}

// BenchmarkZeroCopyConversion 测试零拷贝转换
func BenchmarkZeroCopyConversion(b *testing.B) {
	data := []byte("Hello, World! This is a test string for zero-copy conversion")

	b.Run("ZeroCopy", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = BytesToString(data)
		}
	})

	b.Run("Standard", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = string(data)
		}
	})
}

// BenchmarkRouteCacheHit 测试路由缓存命中
func BenchmarkRouteCacheHit(b *testing.B) {
	cache := NewRouteCache(1000)
	handler := func(ctx *Context) error {
		return ctx.SendString("OK")
	}

	// 预热缓存
	cache.Set("GET", "/user/123", handler, map[string]string{"id": "123"})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = cache.Get("GET", "/user/123")
	}
}

// BenchmarkRouteCacheMiss 测试路由缓存未命中
func BenchmarkRouteCacheMiss(b *testing.B) {
	cache := NewRouteCache(1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = cache.Get("GET", "/user/999")
	}
}

// BenchmarkDynamicRouteWithCache 测试带缓存的动态路由
func BenchmarkDynamicRouteWithCache(b *testing.B) {
	cfg := HighPerformanceConfig()
	cfg.EnableRouteCache = true
	app := New(cfg)

	app.GET("/user/:id", func(ctx *Context) error {
		id := ctx.Param("id")
		return ctx.SendString("User ID: " + id)
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		app.Router.ServeHTTP(w, req)
	}
}

// BenchmarkDynamicRouteWithoutCache 测试不带缓存的动态路由
func BenchmarkDynamicRouteWithoutCache(b *testing.B) {
	cfg := HighPerformanceConfig()
	cfg.EnableRouteCache = false
	app := New(cfg)

	app.GET("/user/:id", func(ctx *Context) error {
		id := ctx.Param("id")
		return ctx.SendString("User ID: " + id)
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		app.Router.ServeHTTP(w, req)
	}
}

// BenchmarkMultipleRoutesWithCache 测试多路由缓存效果
func BenchmarkMultipleRoutesWithCache(b *testing.B) {
	cfg := HighPerformanceConfig()
	cfg.EnableRouteCache = true
	app := New(cfg)

	// 注册多个路由
	app.GET("/api/users/:id", func(ctx *Context) error {
		return ctx.SendString("User")
	})
	app.GET("/api/posts/:id", func(ctx *Context) error {
		return ctx.SendString("Post")
	})
	app.GET("/api/comments/:id", func(ctx *Context) error {
		return ctx.SendString("Comment")
	})

	requests := []string{
		"/api/users/1",
		"/api/posts/2",
		"/api/comments/3",
		"/api/users/1", // 重复，应该命中缓存
		"/api/posts/2", // 重复，应该命中缓存
	}

	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", requests[i%len(requests)], nil)
		w.Body.Reset()
		app.Router.ServeHTTP(w, req)
	}
}

// BenchmarkContextPooling 测试 Context 对象池
func BenchmarkContextPooling(b *testing.B) {
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

// BenchmarkJSONResponseWithBuffer 测试使用缓冲的 JSON 响应
func BenchmarkJSONResponseWithBuffer(b *testing.B) {
	app := New(HighPerformanceConfig())

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
		w.Body.Reset()
		app.Router.ServeHTTP(w, req)
	}
}
