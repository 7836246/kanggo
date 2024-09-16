package etag

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/7836246/kanggo"
)

func TestETagMiddleware(t *testing.T) {
	// 创建 KangGo 实例
	app := kanggo.Default()

	// 使用 ETag 中间件
	app.Use(New())

	// 注册一个简单的 GET 路由
	app.GET("/etag", func(ctx *kanggo.Context) error {
		return ctx.SendString("Hello, KangGo with ETag!")
	})

	// 创建第一个请求
	req, _ := http.NewRequest("GET", "/etag", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp, req)

	// 验证响应状态码和 ETag 头
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	etag := resp.Header().Get("ETag")
	if etag == "" {
		t.Error("未生成 ETag 头")
	}

	// 创建第二个请求，使用 If-None-Match 头
	req2, _ := http.NewRequest("GET", "/etag", nil)
	req2.Header.Set("If-None-Match", etag)
	resp2 := httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp2, req2)

	// 验证响应状态码
	if status := resp2.Code; status != http.StatusNotModified {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusNotModified)
	}
}
