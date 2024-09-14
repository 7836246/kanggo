package kanggo

import (
	"log"
	"net/http"
	"time"
)

// MiddlewareFunc 定义中间件函数的类型
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// Logger 中间件，用于记录请求处理时间
func Logger() MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r) // 调用下一个处理器
			log.Printf("请求 %s %s 处理时间: %v\n", r.Method, r.URL.Path, time.Since(start))
		}
	}
}
