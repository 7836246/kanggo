package encryptcookie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/7836246/kanggo"
)

// 测试 EncryptCookie 中间件
func TestEncryptCookieMiddleware(t *testing.T) {
	// 创建 KangGo 实例
	app := kanggo.Default()

	// 使用 EncryptCookie 中间件
	app.Use(New(Config{
		Key: GenerateKey(),
	}))

	// 注册一个设置 Cookie 的路由
	app.GET("/set", func(ctx *kanggo.Context) error {
		http.SetCookie(ctx.Writer, &http.Cookie{Name: "test_cookie", Value: "test_value"})
		return ctx.SendString("Cookie 已设置")
	})

	// 注册一个获取 Cookie 的路由
	app.GET("/get", func(ctx *kanggo.Context) error {
		cookie, err := ctx.Request.Cookie("test_cookie")
		if err != nil {
			return ctx.SendString("未找到 Cookie")
		}
		return ctx.SendString("Cookie 值: " + cookie.Value)
	})

	// 测试设置 Cookie
	req := httptest.NewRequest("GET", "/set", nil)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	// 验证响应状态码
	if resp.Code != http.StatusOK {
		t.Fatalf("期望状态码 %d，得到 %d", http.StatusOK, resp.Code)
	}

	// 获取加密后的 Cookie 值
	cookies := resp.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("期望找到一个 Cookie，但未找到")
	}

	// 测试获取 Cookie
	req = httptest.NewRequest("GET", "/get", nil)
	req.AddCookie(cookies[0]) // 使用设置的 Cookie
	resp = httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	// 验证响应状态码
	if resp.Code != http.StatusOK {
		t.Fatalf("期望状态码 %d，得到 %d", http.StatusOK, resp.Code)
	}

	// 验证响应内容
	expectedBody := "Cookie 值: test_value"
	if resp.Body.String() != expectedBody {
		t.Fatalf("期望响应内容 '%s'，但得到 '%s'", expectedBody, resp.Body.String())
	}
}
