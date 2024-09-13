package kanggo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 测试路由组的 GET 方法
func TestGroupGET(t *testing.T) {
	app := New(Config{})
	group := app.Router.NewGroup("/api") // 使用大写 Router

	// 注册 GET 路由
	group.GET("/users", func(ctx *Context) error {
		ctx.SendString("Get Users")
		return nil
	})

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/api/users", nil)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	// 检查响应状态码
	if resp.Code != http.StatusOK {
		t.Errorf("状态码错误: 得到 %d, 期待 %d", resp.Code, http.StatusOK)
	}

	// 检查响应内容
	if body := resp.Body.String(); body != "Get Users" {
		t.Errorf("响应内容错误: 得到 %s, 期待 %s", body, "Get Users")
	}
}

// 测试路由组的 POST 方法
func TestGroupPOST(t *testing.T) {
	app := New(Config{})
	group := app.Router.NewGroup("/api")

	// 注册 POST 路由
	group.POST("/users", func(ctx *Context) error {
		ctx.SendString("Create User")
		return nil
	})

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/api/users", nil)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	// 检查响应状态码
	if resp.Code != http.StatusOK {
		t.Errorf("状态码错误: 得到 %d, 期待 %d", resp.Code, http.StatusOK)
	}

	// 检查响应内容
	if body := resp.Body.String(); body != "Create User" {
		t.Errorf("响应内容错误: 得到 %s, 期待 %s", body, "Create User")
	}
}

// 测试路由组的其他方法 (PUT, DELETE, PATCH, OPTIONS, HEAD, TRACE, CONNECT)
func TestGroupOtherMethods(t *testing.T) {
	app := New(Config{})
	group := app.Router.NewGroup("/api")

	methods := []struct {
		method  string
		pattern string
		body    string
	}{
		{"PUT", "/users/:id", "Update User"},
		{"DELETE", "/users/:id", "Delete User"},
		{"PATCH", "/users/:id", "Patch User"},
		{"OPTIONS", "/users", "Options"},
		{"HEAD", "/users", ""},
		{"TRACE", "/trace", "Trace route"},
		{"CONNECT", "/connect", "Connect route"},
	}

	// 注册路由和测试
	for _, m := range methods {
		group.router.Handle(m.method, "/api"+m.pattern, func(ctx *Context) error {
			ctx.SendString(m.body)
			return nil
		})

		req, _ := http.NewRequest(m.method, "/api"+m.pattern, nil)
		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		// 检查响应状态码
		if resp.Code != http.StatusOK {
			t.Errorf("状态码错误: 方法 %s, 路径 %s, 得到 %d, 期待 %d", m.method, m.pattern, resp.Code, http.StatusOK)
		}

		// 检查响应内容
		if m.method != "HEAD" { // HEAD 请求没有响应体
			if body := strings.TrimSpace(resp.Body.String()); body != m.body {
				t.Errorf("响应内容错误: 方法 %s, 路径 %s, 得到 %s, 期待 %s", m.method, m.pattern, body, m.body)
			}
		}
	}
}
