package encryptcookie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/7836246/kanggo"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 200, resp.Code)

	// 获取加密后的 Cookie 值
	cookies := resp.Result().Cookies()
	assert.NotEmpty(t, cookies)

	// 测试获取 Cookie
	req = httptest.NewRequest("GET", "/get", nil)
	req.AddCookie(cookies[0])
	resp = httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "Cookie 值: test_value", resp.Body.String())
}
