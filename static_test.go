package kanggo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// 测试静态文件服务的功能
func TestStatic(t *testing.T) {
	// 创建一个测试目录 ./public 并在其中创建一个测试文件 test.txt
	err := os.MkdirAll("./public", os.ModePerm)
	if err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	defer os.RemoveAll("./public") // 测试完成后删除该目录

	// 创建测试文件
	fileContent := []byte("Hello, Static File!")
	err = os.WriteFile("./public/test.txt", fileContent, 0644)
	if err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	// 确认文件是否存在
	if _, err := os.Stat("./public/test.txt"); os.IsNotExist(err) {
		t.Fatalf("测试文件不存在: %v", err)
	}

	// 创建 KangGo 应用实例
	app := New(Config{})

	// 注册静态文件服务，路由前缀为 /static，根目录为 ./public
	app.Static("/static", "./public", StaticConfig{
		Browse:        false,
		Download:      false,
		CacheDuration: 15 * time.Second,
		MaxAge:        3600,
		ModifyResponse: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "KangGo-Test")
		},
	})

	// 创建一个 HTTP 请求，访问 /static/test.txt
	req := httptest.NewRequest(http.MethodGet, "/static/test.txt", nil)
	// 创建一个 HTTP 响应记录器
	resp := httptest.NewRecorder()

	// 处理请求
	app.router.ServeHTTP(resp, req)

	// 打印调试信息
	t.Logf("请求路径: %s\n", req.URL.Path)
	t.Logf("响应状态码: %d\n", resp.Code)
	t.Logf("响应内容: %s\n", resp.Body.String())

	// 验证响应状态码
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v 期待 %v", status, http.StatusOK)
	}

	// 验证响应内容
	expectedBody := "Hello, Static File!"
	if resp.Body.String() != expectedBody {
		t.Errorf("响应内容错误: 得到 %v 期待 %v", resp.Body.String(), expectedBody)
	}

	// 验证自定义响应头
	if header := resp.Header().Get("X-Custom-Header"); header != "KangGo-Test" {
		t.Errorf("自定义响应头错误: 得到 %v 期待 %v", header, "KangGo-Test")
	}

	// 验证 Cache-Control 头
	if cacheControl := resp.Header().Get("Cache-Control"); cacheControl != "public, max-age=3600" {
		t.Errorf("Cache-Control 头错误: 得到 %v 期待 %v", cacheControl, "public, max-age=3600")
	}
}
