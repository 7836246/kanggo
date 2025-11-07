package kanggo

import (
	"net/http"

	"github.com/7836246/kanggo/core"
)

// CompiledMiddleware 预编译的中间件链
type CompiledMiddleware struct {
	handler http.HandlerFunc
}

// CompileMiddleware 预编译中间件链，减少运行时开销
func CompileMiddleware(middlewares []core.MiddlewareFunc, finalHandler http.HandlerFunc) *CompiledMiddleware {
	// 从后向前应用中间件
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &CompiledMiddleware{
		handler: handler,
	}
}

// ServeHTTP 实现 http.Handler 接口
func (cm *CompiledMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.handler(w, r)
}

// MiddlewareChain 中间件链优化器
type MiddlewareChain struct {
	middlewares []core.MiddlewareFunc
	compiled    *CompiledMiddleware
	dirty       bool // 标记是否需要重新编译
}

// NewMiddlewareChain 创建新的中间件链
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]core.MiddlewareFunc, 0, 8), // 预分配容量
		dirty:       true,
	}
}

// Use 添加中间件
func (mc *MiddlewareChain) Use(middleware core.MiddlewareFunc) {
	mc.middlewares = append(mc.middlewares, middleware)
	mc.dirty = true
}

// Compile 编译中间件链
func (mc *MiddlewareChain) Compile(finalHandler http.HandlerFunc) *CompiledMiddleware {
	if !mc.dirty && mc.compiled != nil {
		return mc.compiled
	}

	mc.compiled = CompileMiddleware(mc.middlewares, finalHandler)
	mc.dirty = false
	return mc.compiled
}

// Apply 应用中间件链到处理函数
func (mc *MiddlewareChain) Apply(finalHandler http.HandlerFunc) http.HandlerFunc {
	if len(mc.middlewares) == 0 {
		return finalHandler
	}

	// 直接应用，避免额外的函数调用
	handler := finalHandler
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		handler = mc.middlewares[i](handler)
	}
	return handler
}

// Len 返回中间件数量
func (mc *MiddlewareChain) Len() int {
	return len(mc.middlewares)
}
