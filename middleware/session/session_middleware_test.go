package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionMiddleware(t *testing.T) {
	store := NewMemoryStore() // 使用内存存储
	middleware := Middleware(store)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// 创建一个处理程序，模拟使用会话中间件的请求处理
	handler := middleware(func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{ResponseWriter: w, Request: r, SessionData: map[string]interface{}{}}

		// 设置并获取会话数据
		ctx.SetSessionValue("username", "kanggo")
		username := ctx.GetSessionValue("username").(string)

		if username != "kanggo" {
			t.Errorf("期望的用户名 'kanggo'，但得到 '%s'", username)
		}

		w.WriteHeader(http.StatusOK)
	})

	// 执行请求
	handler.ServeHTTP(rr, req)

	// 验证响应状态码
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("状态码错误: 得到 %v, 期待 %v", status, http.StatusOK)
	}
}
