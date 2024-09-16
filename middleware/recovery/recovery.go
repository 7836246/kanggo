package recovery

import (
	"fmt"
	"github.com/7836246/kanggo/core"
	"net/http"
	"runtime"
	"runtime/debug"
)

// New Recovery 中间件，用于捕获 panic 并返回 500 错误
func New() core.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 使用 defer 捕获 panic
			defer func() {
				if err := recover(); err != nil {
					// 打印错误信息及堆栈信息
					fmt.Printf("捕获到 panic: %v\n", err)
					debug.PrintStack()

					// 返回 500 错误
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("500 Internal Server Error"))
				}
			}()

			// 调用下一个中间件或最终的处理程序
			next(w, r)
		}
	}
}

// stack 函数用于获取当前 goroutine 的堆栈信息
func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
