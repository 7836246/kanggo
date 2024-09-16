package recovery

import (
	"github.com/7836246/kanggo"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 Recovery 中间件
func TestRecoveryMiddleware(t *testing.T) {
	// 创建 KangGo 实例
	app := kanggo.Default()

	// 使用 Recovery 中间件
	app.Use(New())

	// 注册一个会产生 panic 的路由
	app.GET("/panic", func(ctx *kanggo.Context) error {
		panic("这是一个测试 panic!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/panic", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	app.Router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusInternalServerError {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusInternalServerError)
	}

	// 验证响应内容
	expected := "500 Internal Server Error"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}
