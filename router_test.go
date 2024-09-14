package kanggo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// 测试静态路由的注册和处理
func TestStaticRoute(t *testing.T) {
	router := NewRouter(DefaultConfig())

	// 注册静态路由
	router.Handle("GET", "/home", func(ctx *Context) error {
		return ctx.SendString("Welcome to the home page!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/home", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "Welcome to the home page!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}

// 测试文件路由的注册和处理
func TestFileRoute(t *testing.T) {
	app := Default()

	testRoot := "./testdata"
	testFile := "test.txt"
	os.MkdirAll(testRoot, 0755)
	defer os.RemoveAll(testRoot)

	filePath := filepath.Join(testRoot, testFile)
	os.WriteFile(filePath, []byte("Hello, KangGo!"), 0644)

	app.Static("/files", testRoot)

	req, err := http.NewRequest("GET", "/files/"+testFile, nil)
	if err != nil {
		t.Fatalf("无法创建请求: %v", err)
	}
	resp := httptest.NewRecorder()

	app.Router.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	expected := "Hello, KangGo!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}

// 测试动态路由的注册和处理
func TestDynamicRoute(t *testing.T) {
	router := NewRouter(DefaultConfig())

	// 注册动态路由
	router.Handle("GET", "/user/:id", func(ctx *Context) error {
		userID := ctx.Param("id")
		return ctx.SendString("User ID: " + userID)
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/user/123", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "User ID: 123"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}

// 测试 Add 方法的注册和处理
func TestAddMethod(t *testing.T) {
	router := NewRouter(DefaultConfig())

	// 使用 Add 方法注册自定义方法
	router.Handle("CUSTOM", "/custom", func(ctx *Context) error {
		return ctx.SendString("Custom method route!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("CUSTOM", "/custom", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "Custom method route!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}

// 测试 All 方法的注册和处理
func TestAllMethod(t *testing.T) {
	router := NewRouter(DefaultConfig())

	// 使用 All 方法注册处理函数
	router.Handle("GET", "/all", func(ctx *Context) error {
		return ctx.SendString("All methods route!")
	})

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/all", nil)
	resp := httptest.NewRecorder()

	// 处理请求
	router.ServeHTTP(resp, req)

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expected := "All methods route!"
	if resp.Body.String() != expected {
		t.Errorf("响应内容错误: 得到 %v, 期待 %v", resp.Body.String(), expected)
	}
}
