package logger

import (
	"github.com/7836246/kanggo/core"
	"log"
	"net/http"
	"time"
)

// New Logger 中间件，用于记录请求处理时间
func New() core.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r) // 调用下一个处理器
			log.Printf("请求 %s %s 处理时间: %v\n", r.Method, r.URL.Path, time.Since(start))
		}
	}
}
