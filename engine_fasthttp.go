//go:build fasthttp
// +build fasthttp

package kanggo

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
)

// FastHTTPEngine fasthttp 引擎实现
type FastHTTPEngine struct {
	server  *fasthttp.Server
	config  EngineConfig
	handler fasthttp.RequestHandler
}

// newFastHTTPEngine 创建 fasthttp 引擎
func newFastHTTPEngine(cfg EngineConfig, handler interface{}) *FastHTTPEngine {
	fastHandler, ok := handler.(fasthttp.RequestHandler)
	if !ok {
		panic("fasthttp engine requires fasthttp.RequestHandler")
	}

	server := &fasthttp.Server{
		Handler:            fastHandler,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		MaxRequestBodySize: cfg.MaxRequestBodySize,
		Concurrency:        cfg.Concurrency,
		DisableKeepalive:   cfg.DisableKeepalive,
		ReduceMemoryUsage:  cfg.ReduceMemoryUsage,
	}

	// 默认并发数
	if server.Concurrency <= 0 {
		server.Concurrency = 256 * 1024
	}

	return &FastHTTPEngine{
		server:  server,
		config:  cfg,
		handler: fastHandler,
	}
}

// ListenAndServe 启动 HTTP 服务器
func (e *FastHTTPEngine) ListenAndServe(addr string) error {
	fmt.Printf("🚀 FastHTTP 引擎启动在 %s\n", addr)
	return e.server.ListenAndServe(addr)
}

// ListenAndServeTLS 启动 HTTPS 服务器
func (e *FastHTTPEngine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	fmt.Printf("🔒 FastHTTP HTTPS 引擎启动在 %s\n", addr)
	return e.server.ListenAndServeTLS(addr, certFile, keyFile)
}

// Shutdown 优雅关闭服务器
func (e *FastHTTPEngine) Shutdown(ctx context.Context) error {
	if e.server != nil {
		return e.server.Shutdown()
	}
	return nil
}

// GetMode 获取引擎模式
func (e *FastHTTPEngine) GetMode() EngineMode {
	return FastHTTPMode
}
