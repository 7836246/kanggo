package etag

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
)

// ETag 中间件函数，根据响应内容计算 ETag 并将其添加到响应头
func ETag() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 创建一个 ResponseRecorder 来捕获响应
			recorder := NewResponseRecorder(w)

			// 检查请求头中的 If-None-Match
			body := recorder.Body.Bytes()
			hash := md5.Sum(body)
			etag := `"` + hex.EncodeToString(hash[:]) + `"`

			if match := r.Header.Get("If-None-Match"); match != "" {
				if strings.Contains(match, etag) {
					// 设置状态码为 304 并返回
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}

			// 调用下一个中间件/处理程序
			next(recorder, r)

			// 设置 ETag 响应头
			w.Header().Set("ETag", etag)

			// 将捕获的响应写回客户端
			recorder.WriteToResponse(w)
		}
	}
}
