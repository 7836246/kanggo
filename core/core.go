package core

import "net/http"

// MiddlewareFunc 定义中间件函数的类型
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
