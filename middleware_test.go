package kanggo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 Logger 中间件
func TestLoggerMiddleware(t *testing.T) {
	// 创建 KangGo 实例
	app := Default()

	// 使用 Logger 中间件
	app.Use(Logger())

	// 注册一个简单的 GET 路由
	app.GET("/test", func(ctx *Context) error {
		return ctx.SendString("Hello, KangGo!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "Hello, KangGo!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}
