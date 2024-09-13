package kanggo

import (
	"fmt"
	"net/http"
)

// KangGo 核心结构
type KangGo struct {
	Router *Router
	Config Config
}

// Default 创建一个带有默认设置的 KangGo 实例
func Default() *KangGo {
	cfg := DefaultConfig() // 使用默认配置
	k := New(cfg)
	return k
}

// New 创建一个新的 KangGo 实例
func New(cfg Config) *KangGo {
	// 创建 KangGo 实例
	k := &KangGo{
		Router: NewRouter(cfg), // 将配置传递给 NewRouter
		Config: cfg,
	}

	// 根据配置决定是否打印横幅
	if cfg.ShowBanner {
		PrintWelcomeBanner()
	}

	return k
}

// GET 注册一个 GET 请求路由
func (k *KangGo) GET(pattern string, handler HandlerFunc) {
	k.Router.Handle("GET", pattern, handler)
}

// POST 注册一个 POST 请求路由
func (k *KangGo) POST(pattern string, handler HandlerFunc) {
	k.Router.Handle("POST", pattern, handler)
}

// PUT 注册一个 PUT 请求路由
func (k *KangGo) PUT(pattern string, handler HandlerFunc) {
	k.Router.Handle("PUT", pattern, handler)
}

// DELETE 注册一个 DELETE 请求路由
func (k *KangGo) DELETE(pattern string, handler HandlerFunc) {
	k.Router.Handle("DELETE", pattern, handler)
}

// PATCH 注册一个 PATCH 请求路由
func (k *KangGo) PATCH(pattern string, handler HandlerFunc) {
	k.Router.Handle("PATCH", pattern, handler)
}

// OPTIONS 注册一个 OPTIONS 请求路由
func (k *KangGo) OPTIONS(pattern string, handler HandlerFunc) {
	k.Router.Handle("OPTIONS", pattern, handler)
}

// HEAD 注册一个 HEAD 请求路由
func (k *KangGo) HEAD(pattern string, handler HandlerFunc) {
	k.Router.Handle("HEAD", pattern, handler)
}

// TRACE 注册一个 TRACE 请求路由
func (k *KangGo) TRACE(pattern string, handler HandlerFunc) {
	k.Router.Handle("TRACE", pattern, handler)
}

// CONNECT 注册一个 CONNECT 请求路由
func (k *KangGo) CONNECT(pattern string, handler HandlerFunc) {
	k.Router.Handle("CONNECT", pattern, handler)
}

// Run 启动 HTTP 服务器
func (k *KangGo) Run(addr string) error {
	// 根据配置决定是否打印路由信息
	if k.Config.PrintRoutes {
		k.Router.PrintRoutes() // 打印所有注册的路由信息
	}
	// 创建一个自定义的 HTTP 服务器配置
	server := &http.Server{
		Addr:         addr,
		Handler:      k.Router,              // 使用 KangGo 的路由器作为请求处理器
		IdleTimeout:  k.Config.IdleTimeout,  // 设置空闲连接超时时间
		ReadTimeout:  k.Config.ReadTimeout,  // 设置读取请求超时时间
		WriteTimeout: k.Config.WriteTimeout, // 设置写入响应超时时间
	}

	fmt.Printf("KangGo 服务器正在运行，地址: %s\n", addr)
	return server.ListenAndServe()
}
