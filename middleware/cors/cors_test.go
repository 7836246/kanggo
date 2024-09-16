package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/7836246/kanggo"
)

// 测试 CORS 中间件
func TestCORSMiddleware(t *testing.T) {
	// 创建 KangGo 实例
	app := kanggo.Default()

	// 使用 CORS 中间件
	app.Use(New())

	// 注册一个简单的 GET 路由
	app.GET("/test", func(ctx *kanggo.Context) error {
		return ctx.SendString("Hello, KangGo!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp, req)

	// 验证 CORS 头
	if origin := resp.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("CORS 头错误: 得到 %v, 期待 %v", origin, "*")
	}

	// 验证状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 创建一个实际请求来验证 CORS
	req, _ = http.NewRequest("GET", "/test", nil)
	resp = httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp, req)

	// 验证状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "Hello, KangGo!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}
