package kanggo

import (
	"context"
	"net/http"
	"time"
)

// EngineMode 引擎模式
type EngineMode int

const (
	// NetHTTPMode 标准库模式（默认，兼容性好）
	NetHTTPMode EngineMode = iota
	// FastHTTPMode fasthttp 模式（高性能，可选）
	FastHTTPMode
)

// String 返回引擎模式的字符串表示
func (m EngineMode) String() string {
	switch m {
	case FastHTTPMode:
		return "fasthttp"
	default:
		return "net/http"
	}
}

// Engine 引擎接口 - 统一 net/http 和 fasthttp
type Engine interface {
	// ListenAndServe 启动服务器
	ListenAndServe(addr string) error

	// ListenAndServeTLS 启动 HTTPS 服务器
	ListenAndServeTLS(addr, certFile, keyFile string) error

	// Shutdown 优雅关闭服务器
	Shutdown(ctx context.Context) error

	// GetMode 获取引擎模式
	GetMode() EngineMode
}

// EngineConfig 引擎配置
type EngineConfig struct {
	Mode               EngineMode
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	MaxRequestBodySize int
	Concurrency        int  // fasthttp 专用：最大并发连接数
	DisableKeepalive   bool // fasthttp 专用：禁用 keep-alive
	ReduceMemoryUsage  bool // fasthttp 专用：减少内存使用
}

// NewEngine 创建引擎
func NewEngine(mode EngineMode, cfg EngineConfig, handler interface{}) Engine {
	switch mode {
	case FastHTTPMode:
		return newFastHTTPEngine(cfg, handler)
	default:
		return newNetHTTPEngine(cfg, handler)
	}
}

// NetHTTPEngine 标准库引擎实现
type NetHTTPEngine struct {
	server  *http.Server
	config  EngineConfig
	handler http.Handler
}

// newNetHTTPEngine 创建 net/http 引擎
func newNetHTTPEngine(cfg EngineConfig, handler interface{}) *NetHTTPEngine {
	httpHandler, ok := handler.(http.Handler)
	if !ok {
		panic("net/http engine requires http.Handler")
	}

	return &NetHTTPEngine{
		config:  cfg,
		handler: httpHandler,
		server: &http.Server{
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// ListenAndServe 启动 HTTP 服务器
func (e *NetHTTPEngine) ListenAndServe(addr string) error {
	e.server.Addr = addr
	e.server.Handler = e.handler
	return e.server.ListenAndServe()
}

// ListenAndServeTLS 启动 HTTPS 服务器
func (e *NetHTTPEngine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	e.server.Addr = addr
	e.server.Handler = e.handler
	return e.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown 优雅关闭服务器
func (e *NetHTTPEngine) Shutdown(ctx context.Context) error {
	if e.server != nil {
		return e.server.Shutdown(ctx)
	}
	return nil
}

// GetMode 获取引擎模式
func (e *NetHTTPEngine) GetMode() EngineMode {
	return NetHTTPMode
}

// FastHTTPEngine fasthttp 引擎实现（占位符）
// 注意：需要构建标签 -tags fasthttp 才能使用
type FastHTTPEngine struct {
	config EngineConfig
	mode   EngineMode
}

// newFastHTTPEngine 创建 fasthttp 引擎（占位符实现）
func newFastHTTPEngine(cfg EngineConfig, _ /* handler */ interface{}) *FastHTTPEngine {
	// 如果没有编译 fasthttp 支持，返回占位符
	return &FastHTTPEngine{
		config: cfg,
		mode:   FastHTTPMode,
	}
}

// ListenAndServe 启动服务器（占位符）
func (e *FastHTTPEngine) ListenAndServe(addr string) error {
	panic("fasthttp engine not available. Build with -tags fasthttp")
}

// ListenAndServeTLS 启动 HTTPS 服务器（占位符）
func (e *FastHTTPEngine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	panic("fasthttp engine not available. Build with -tags fasthttp")
}

// Shutdown 优雅关闭服务器（占位符）
func (e *FastHTTPEngine) Shutdown(ctx context.Context) error {
	return nil
}

// GetMode 获取引擎模式
func (e *FastHTTPEngine) GetMode() EngineMode {
	return e.mode
}
