package kanggo

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/7836246/kanggo/core"
)

// HandlerFunc 定义处理函数签名
type HandlerFunc func(ctx *Context) error

// RadixNode 是 Radix Tree 的节点
type RadixNode struct {
	path     string
	handler  HandlerFunc
	children map[string]*RadixNode
	isLeaf   bool
	isParam  bool   // 标识该节点是否为路径参数
	paramKey string // 路径参数的键（如 :id）
}

// RouteInfo 存储动态路由的信息
type RouteInfo struct {
	Method  string
	Pattern string
}

// StaticRouteInfo 存储普通静态路由的信息
type StaticRouteInfo struct {
	Method  string // 新增字段，存储请求方法
	Prefix  string
	Handler HandlerFunc
}

// FileRouteInfo 存储文件路由的信息
type FileRouteInfo struct {
	Method  string // 新增字段，存储请求方法
	Prefix  string
	Root    string
	Handler HandlerFunc
}

// Router 路由结构
type Router struct {
	staticRoutes   []StaticRouteInfo      // 普通静态路由列表（用于打印）
	staticRouteMap map[string]HandlerFunc // 静态路由哈希表，O(1) 查找
	fileRoutes     []FileRouteInfo        // 文件路由列表
	dynamicRoot    *RadixNode             // 动态路由的 Radix Tree 根节点
	routes         []RouteInfo            // 存储所有注册的动态路由信息
	config         Config                 // 添加配置到 Router 中
	middleware     []core.MiddlewareFunc  // 中间件切片
	routeCache     *RouteCache            // 路由缓存（Phase 1 优化）
	regexRouter    *RegexRouter           // 正则表达式路由器（Phase 3）
}

// Use 方法注册中间件到路由器
func (r *Router) Use(mw core.MiddlewareFunc) {
	r.middleware = append(r.middleware, mw)
}

// NewRouter 创建一个新的路由器
func NewRouter(cfg Config) *Router {
	router := &Router{
		staticRoutes:   []StaticRouteInfo{},              // 初始化普通静态路由列表
		staticRouteMap: make(map[string]HandlerFunc, 32), // 预分配哈希表容量
		fileRoutes:     []FileRouteInfo{},                // 初始化文件路由列表
		dynamicRoot:    &RadixNode{children: make(map[string]*RadixNode)},
		config:         cfg,
		routes:         []RouteInfo{},       // 初始化路由信息列表
		routeCache:     NewRouteCache(1000), // 初始化路由缓存，默认 1000 条
	}

	// 根据配置决定是否启用路由缓存
	if !cfg.EnableRouteCache {
		router.routeCache.Disable()
	}

	return router
}

// RegisterStaticRoute 注册普通静态路由信息
func (r *Router) RegisterStaticRoute(method, pattern string, handler HandlerFunc) {
	r.staticRoutes = append(r.staticRoutes, StaticRouteInfo{
		Method:  method,
		Prefix:  pattern,
		Handler: handler,
	})
	// 同时添加到哈希表，使用 method:pattern 作为键
	key := method + ":" + pattern
	r.staticRouteMap[key] = handler
}

// RegisterFileRoute 注册文件路由信息
func (r *Router) RegisterFileRoute(method, pattern, root string, handler HandlerFunc) {
	r.fileRoutes = append(r.fileRoutes, FileRouteInfo{
		Method:  method,
		Prefix:  pattern,
		Root:    root,
		Handler: handler,
	})
}

// PrintRoutes 打印所有注册的路由信息，区分目录文件路由、单文件路由、普通静态路由和动态路由
func (r *Router) PrintRoutes() {
	// 打印表头
	fmt.Println("\n📋 已注册的路由信息:")
	fmt.Println(strings.Repeat("=", 66))
	fmt.Printf("| %-10s | %-10s | %-20s | %-20s |\n", "类型", "请求方式", "路由前缀", "映射路径")
	fmt.Println(strings.Repeat("=", 66))

	// 打印文件路由和目录路由
	for _, fileRoute := range r.fileRoutes {
		// 判断是文件还是目录
		routeType := "目录"
		if isFile(fileRoute.Root) {
			routeType = "文件"
		}
		fmt.Printf("| %-10s | %-10s | %-20s | %-20s |\n", routeType, fileRoute.Method, fileRoute.Prefix, fileRoute.Root)
	}

	// 打印普通静态路由
	for _, staticRoute := range r.staticRoutes {
		fmt.Printf("| %-10s | %-10s | %-20s | %-20s |\n", "静态", staticRoute.Method, staticRoute.Prefix, "-")
	}

	// 打印动态路由
	for _, route := range r.routes {
		fmt.Printf("| %-10s | %-10s | %-20s | %-20s |\n", "动态", route.Method, route.Pattern, "-")
	}

	// 打印表格结束线
	fmt.Println(strings.Repeat("=", 66))
}

// isFile 检查给定的路径是否是文件
func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Handle 注册路由
func (r *Router) Handle(method, pattern string, handler HandlerFunc) {
	// 根据配置决定是否对路由进行大小写转换
	if !r.config.CaseSensitiveRouting {
		pattern = strings.ToLower(pattern)
	}

	// 根据配置决定是否启用严格路由模式
	if !r.config.StrictRouting {
		pattern = strings.TrimSuffix(pattern, "/")
	}

	// 根据配置决定是否对路径进行解码
	if r.config.UnescapePath {
		unescapedPattern, err := url.PathUnescape(pattern)
		if err == nil {
			pattern = unescapedPattern
		}
	}

	// 判断是否为文件路由
	if strings.HasSuffix(pattern, "/*") {
		// 注册文件路由
		filePattern := strings.TrimSuffix(pattern, "/*")
		r.RegisterFileRoute(method, filePattern, "", handler)
	} else if isStaticRoute(pattern) {
		// 判断是否为普通静态路由
		r.RegisterStaticRoute(method, pattern, handler)
	} else {
		// 动态路由，存入 Radix Tree
		r.insertDynamicRoute(method, pattern, handler)

		// 记录动态路由信息
		r.routes = append(r.routes, RouteInfo{Method: method, Pattern: pattern})
	}
}

// isStaticRoute 判断是否为普通静态路由（不包含 ":" 或 "*"）
func isStaticRoute(pattern string) bool {
	return !strings.Contains(pattern, ":") && !strings.Contains(pattern, "*")
}

// insertDynamicRoute 向 Radix Tree 中插入动态路由
func (r *Router) insertDynamicRoute(_ /* method */, pattern string, handler HandlerFunc) {
	// 预分配切片容量，减少扩容
	parts := make([]string, 0, 8)
	start := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '/' {
			if i > start {
				parts = append(parts, pattern[start:i])
			}
			start = i + 1
		}
	}
	if start < len(pattern) {
		parts = append(parts, pattern[start:])
	}

	node := r.dynamicRoot
	for _, part := range parts {
		if part == "" {
			continue
		}
		isParam := len(part) > 0 && part[0] == ':'
		childKey := part
		if isParam {
			childKey = ":param" // 所有参数化路径的标识符
		}

		child, ok := node.children[childKey]
		if !ok {
			child = &RadixNode{
				path:     part,
				children: make(map[string]*RadixNode, 2), // 预分配容量
				isParam:  isParam,
			}
			if isParam {
				child.paramKey = part[1:] // 存储参数的键名，例如 "id"
			}
			node.children[childKey] = child
		}
		node = child
	}
	node.handler = handler
	node.isLeaf = true
}

// ServeHTTP 实现 http.Handler 接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 设置响应头中的 Server 字段
	if r.config.ServerHeader != "" {
		w.Header().Set("Server", r.config.ServerHeader)
	}

	// 检查请求体大小是否超过配置的最大限制
	if req.ContentLength > int64(r.config.MaxRequestBodySize) {
		http.Error(w, "请求体过大", http.StatusRequestEntityTooLarge)
		return
	}

	// 根据配置决定是否对路径进行解码
	path := req.URL.Path
	if r.config.UnescapePath {
		unescapedPath, err := url.PathUnescape(path)
		if err != nil {
			http.Error(w, "路径解码错误", http.StatusBadRequest)
			return
		}
		path = unescapedPath
	}

	// 从对象池获取 Context 实例
	ctx := AcquireContext(w, req, r.config)
	defer ReleaseContext(ctx) // 请求结束后归还到对象池

	// 最终的处理函数，实际处理请求逻辑
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 查找文件路由
		for _, fileRoute := range r.fileRoutes {
			if path == fileRoute.Prefix || strings.HasPrefix(path, fileRoute.Prefix+"/") {
				req.URL.Path = strings.TrimPrefix(path, fileRoute.Prefix)
				if err := fileRoute.Handler(ctx); err != nil {
					r.handleError(w, err)
				}
				return
			}
		}

		// 使用哈希表查找静态路由，O(1) 时间复杂度
		staticKey := req.Method + ":" + path
		if handler, found := r.staticRouteMap[staticKey]; found {
			if err := handler(ctx); err != nil {
				r.handleError(w, err)
			}
			return
		}

		// 先检查路由缓存
		if handler, params, found := r.routeCache.Get(req.Method, path); found {
			// 从缓存获取路由
			for k, v := range params {
				ctx.Params[k] = v
			}
			if err := handler(ctx); err != nil {
				r.handleError(w, err)
			}
			return
		}

		// 查找动态路由
		if handler, found := r.searchDynamicRoute(req.Method, path, ctx); found {
			// 将路由添加到缓存
			r.routeCache.Set(req.Method, path, handler, ctx.Params)

			if err := handler(ctx); err != nil {
				r.handleError(w, err)
			}
			return
		}

		// 尝试正则表达式路由匹配（Phase 3）
		if handler, found := matchRegexRoute(r, req.Method, path, ctx); found {
			if err := handler(ctx); err != nil {
				r.handleError(w, err)
			}
			return
		}

		// 如果没有匹配的路由，返回 404
		http.NotFound(w, req)
	})

	// 应用所有中间件，注意类型转换
	wrappedHandler := finalHandler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		wrappedHandler = r.middleware[i](wrappedHandler)
	}

	// 处理最终的请求
	wrappedHandler.ServeHTTP(w, req)
}

// searchDynamicRoute 在 Radix Tree 中查找动态路由
func (r *Router) searchDynamicRoute(_ /* method */, path string, ctx *Context) (HandlerFunc, bool) {
	// 手动分割路径，避免 strings.Split 的内存分配
	node := r.dynamicRoot
	var child *RadixNode // 在循环外声明 child 变量
	var ok bool          // 在循环外声明 ok 变量

	start := 0
	for i := 0; i <= len(path); i++ {
		if i == len(path) || path[i] == '/' {
			if i > start {
				part := path[start:i]

				// 先尝试静态部分匹配
				child, ok = node.children[part]
				if ok {
					node = child
				} else {
					// 再尝试参数化部分匹配
					child, ok = node.children[":param"]
					if ok {
						ctx.Params[child.paramKey] = part
						node = child
					} else {
						return nil, false
					}
				}
			}
			start = i + 1
		}
	}

	if node.isLeaf {
		return node.handler, true
	}
	return nil, false
}

// handleError 统一的错误处理
func (r *Router) handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// CacheStats 获取路由缓存统计信息
func (r *Router) CacheStats() map[string]interface{} {
	if r.routeCache == nil {
		return map[string]interface{}{
			"enabled": false,
		}
	}
	return r.routeCache.Stats()
}
