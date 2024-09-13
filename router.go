package KangGo

import (
	"github.com/7836246/kanggo/config"
	"net/http"
	"strings"
)

// HandlerFunc 定义处理函数签名
type HandlerFunc func(Context) error

// RadixNode 是 Radix Tree 的节点
type RadixNode struct {
	path     string
	handler  HandlerFunc
	children map[string]*RadixNode
	isLeaf   bool
	isParam  bool   // 标识该节点是否为路径参数
	paramKey string // 路径参数的键（如 :id）
}

// Router 路由结构
type Router struct {
	staticRoutes map[string]HandlerFunc // 静态路由哈希表
	dynamicRoot  *RadixNode             // 动态路由的 Radix Tree 根节点
	config       config.Config          // 添加配置到 Router 中
}

// NewRouter 创建一个新的路由器
func NewRouter(cfg config.Config) *Router {
	return &Router{
		staticRoutes: make(map[string]HandlerFunc),
		dynamicRoot:  &RadixNode{children: make(map[string]*RadixNode)},
		config:       cfg,
	}
}

// Handle 注册路由，判断是静态还是动态路由
func (r *Router) Handle(method, pattern string, handler HandlerFunc) {
	routeKey := method + "-" + pattern
	if isStaticRoute(pattern) {
		// 静态路由，存入哈希表
		r.staticRoutes[routeKey] = handler
	} else {
		// 动态路由，存入 Radix Tree
		r.insertDynamicRoute(method, pattern, handler)
	}
}

// isStaticRoute 判断是否为静态路由（不包含 ":" 或 "*"）
func isStaticRoute(pattern string) bool {
	return !strings.Contains(pattern, ":") && !strings.Contains(pattern, "*")
}

// insertDynamicRoute 向 Radix Tree 中插入动态路由
func (r *Router) insertDynamicRoute(method, pattern string, handler HandlerFunc) {
	parts := strings.Split(pattern, "/")
	node := r.dynamicRoot
	for _, part := range parts {
		if part == "" {
			continue
		}
		isParam := strings.HasPrefix(part, ":")
		childKey := part
		if isParam {
			childKey = ":param" // 所有参数化路径的标识符
		}

		child, ok := node.children[childKey]
		if !ok {
			child = &RadixNode{
				path:     part,
				children: make(map[string]*RadixNode),
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
	// 创建 Context 时传递配置参数
	ctx := NewContext(w, req, r.config)
	routeKey := req.Method + "-" + req.URL.Path

	// 优先查找静态路由
	if handler, ok := r.staticRoutes[routeKey]; ok {
		if err := handler(ctx); err != nil {
			r.handleError(w, err) // 处理错误
		}
		return
	}

	// 查找动态路由
	if handler, found := r.searchDynamicRoute(req.Method, req.URL.Path, &ctx); found {
		if err := handler(ctx); err != nil {
			r.handleError(w, err) // 处理错误
		}
		return
	}

	http.NotFound(w, req)
}

// searchDynamicRoute 在 Radix Tree 中查找动态路由
func (r *Router) searchDynamicRoute(method, path string, ctx *Context) (HandlerFunc, bool) {
	parts := strings.Split(path, "/")
	node := r.dynamicRoot
	var child *RadixNode // 在循环外声明 child 变量
	var ok bool          // 在循环外声明 ok 变量

	for _, part := range parts {
		if part == "" {
			continue
		}

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

	if node.isLeaf {
		return node.handler, true
	}
	return nil, false
}

// handleError 统一的错误处理
func (r *Router) handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
