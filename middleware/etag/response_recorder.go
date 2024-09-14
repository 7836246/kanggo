package etag

import (
	"bytes"
	"net/http"
)

// ResponseRecorder 是一个用于捕获 HTTP 响应的自定义结构体
type ResponseRecorder struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

// NewResponseRecorder 创建一个新的 ResponseRecorder
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Body:           &bytes.Buffer{},
	}
}

// Write 捕获写入的内容
func (rec *ResponseRecorder) Write(p []byte) (int, error) {
	rec.Body.Write(p)
	return rec.ResponseWriter.Write(p)
}

// WriteToResponse 将捕获的内容写回原始的 ResponseWriter
func (rec *ResponseRecorder) WriteToResponse(w http.ResponseWriter) {
	w.Write(rec.Body.Bytes())
}
