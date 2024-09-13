package kanggo

import (
	"fmt"
	"net/http"
)

// KangGo 核心结构
type KangGo struct {
	router *Router
	config Config
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
		router: NewRouter(cfg), // 将配置传递给 NewRouter
		config: cfg,
	}

	// 根据配置决定是否打印横幅
	if cfg.ShowBanner {
		PrintWelcomeBanner()
	}

	return k
}

// GET 注册一个 GET 请求路由
func (k *KangGo) GET(pattern string, handler HandlerFunc) {
	k.router.Handle("GET", pattern, handler)
}

// POST 注册一个 POST 请求路由
func (k *KangGo) POST(pattern string, handler HandlerFunc) {
	k.router.Handle("POST", pattern, handler)
}

// PUT 注册一个 PUT 请求路由
func (k *KangGo) PUT(pattern string, handler HandlerFunc) {
	k.router.Handle("PUT", pattern, handler)
}

// DELETE 注册一个 DELETE 请求路由
func (k *KangGo) DELETE(pattern string, handler HandlerFunc) {
	k.router.Handle("DELETE", pattern, handler)
}

// Run 启动 HTTP 服务器
func (k *KangGo) Run(addr string) error {
	// 根据配置决定是否打印路由信息
	if k.config.PrintRoutes {
		k.router.PrintRoutes() // 打印所有注册的路由信息
	}
	// 创建一个自定义的 HTTP 服务器配置
	server := &http.Server{
		Addr:         addr,
		Handler:      k.router,              // 使用 KangGo 的路由器作为请求处理器
		IdleTimeout:  k.config.IdleTimeout,  // 设置空闲连接超时时间
		ReadTimeout:  k.config.ReadTimeout,  // 设置读取请求超时时间
		WriteTimeout: k.config.WriteTimeout, // 设置写入响应超时时间
	}

	fmt.Printf("KangGo 服务器正在运行，地址: %s\n", addr)
	return server.ListenAndServe()
}
