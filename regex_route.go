package kanggo

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// RegexRoute 正则表达式路由
type RegexRoute struct {
	pattern *regexp.Regexp
	handler HandlerFunc
	params  []string // 捕获组名称
}

// RegexRouter 正则表达式路由器
type RegexRouter struct {
	routes map[string][]*RegexRoute // method -> routes
}

// NewRegexRouter 创建新的正则表达式路由器
func NewRegexRouter() *RegexRouter {
	return &RegexRouter{
		routes: make(map[string][]*RegexRoute),
	}
}

// AddRoute 添加正则表达式路由
func (r *RegexRouter) AddRoute(method, pattern string, handler HandlerFunc) error {
	// 解析命名捕获组
	params := extractParamNames(pattern)

	// 编译正则表达式
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	route := &RegexRoute{
		pattern: regex,
		handler: handler,
		params:  params,
	}

	if r.routes[method] == nil {
		r.routes[method] = make([]*RegexRoute, 0)
	}

	r.routes[method] = append(r.routes[method], route)

	return nil
}

// Match 匹配路由
func (r *RegexRouter) Match(method, path string, ctx *Context) (HandlerFunc, bool) {
	routes, ok := r.routes[method]
	if !ok {
		return nil, false
	}

	for _, route := range routes {
		if matches := route.pattern.FindStringSubmatch(path); matches != nil {
			// 提取参数
			if len(route.params) > 0 && len(matches) > 1 {
				for i, paramName := range route.params {
					if i+1 < len(matches) {
						ctx.Params[paramName] = matches[i+1]
					}
				}
			}

			return route.handler, true
		}
	}

	return nil, false
}

// extractParamNames 从正则表达式中提取命名捕获组
// 例如: `^/user/(?P<id>\d+)$` -> ["id"]
func extractParamNames(pattern string) []string {
	// 匹配命名捕获组: (?P<name>...)
	namedGroupRegex := regexp.MustCompile(`\(\?P<(\w+)>`)
	matches := namedGroupRegex.FindAllStringSubmatch(pattern, -1)

	params := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}

	return params
}

// Regex 注册正则表达式路由
func (k *KangGo) Regex(method, pattern string, handler HandlerFunc) error {
	// 初始化正则路由器（如果需要）
	if k.Router.regexRouter == nil {
		k.Router.regexRouter = NewRegexRouter()
	}

	return k.Router.regexRouter.AddRoute(method, pattern, handler)
}

// RegexGET 注册 GET 正则路由
func (k *KangGo) RegexGET(pattern string, handler HandlerFunc) error {
	return k.Regex("GET", pattern, handler)
}

// RegexPOST 注册 POST 正则路由
func (k *KangGo) RegexPOST(pattern string, handler HandlerFunc) error {
	return k.Regex("POST", pattern, handler)
}

// RegexPUT 注册 PUT 正则路由
func (k *KangGo) RegexPUT(pattern string, handler HandlerFunc) error {
	return k.Regex("PUT", pattern, handler)
}

// RegexDELETE 注册 DELETE 正则路由
func (k *KangGo) RegexDELETE(pattern string, handler HandlerFunc) error {
	return k.Regex("DELETE", pattern, handler)
}

// RegexPATCH 注册 PATCH 正则路由
func (k *KangGo) RegexPATCH(pattern string, handler HandlerFunc) error {
	return k.Regex("PATCH", pattern, handler)
}

// 更新 Router 结构以支持正则路由
// 需要在 router.go 中添加:
// regexRouter *RegexRouter

// 路由匹配辅助函数
func matchRegexRoute(router *Router, method, path string, ctx *Context) (HandlerFunc, bool) {
	if router.regexRouter != nil {
		return router.regexRouter.Match(method, path, ctx)
	}
	return nil, false
}

// RouteConstraint 路由约束
type RouteConstraint func(value string) bool

// 预定义约束
var (
	// IntConstraint 整数约束
	IntConstraint = func(value string) bool {
		return regexp.MustCompile(`^\d+$`).MatchString(value)
	}

	// UUIDConstraint UUID 约束
	UUIDConstraint = func(value string) bool {
		return regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(value)
	}

	// AlphaConstraint 字母约束
	AlphaConstraint = func(value string) bool {
		return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(value)
	}

	// AlphanumConstraint 字母数字约束
	AlphanumConstraint = func(value string) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value)
	}

	// SlugConstraint Slug 约束 (URL 友好)
	SlugConstraint = func(value string) bool {
		return regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(value)
	}
)

// ConstrainedRoute 带约束的路由
type ConstrainedRoute struct {
	pattern     string
	handler     HandlerFunc
	constraints map[string]RouteConstraint
}

// WithConstraints 添加路由约束
func (k *KangGo) WithConstraints(pattern string, handler HandlerFunc, constraints map[string]RouteConstraint) HandlerFunc {
	return func(ctx *Context) error {
		// 验证所有约束
		for paramName, constraint := range constraints {
			value := ctx.Param(paramName)
			if value == "" || !constraint(value) {
				http.Error(ctx.Writer, "Bad Request: Invalid parameter", http.StatusBadRequest)
				return nil
			}
		}

		// 所有约束通过，调用原处理器
		return handler(ctx)
	}
}

// 路由组合器 - 用于复杂路由
type RouteBuilder struct {
	pattern     string
	method      string
	handler     HandlerFunc
	constraints map[string]RouteConstraint
	middlewares []func(HandlerFunc) HandlerFunc
}

// NewRouteBuilder 创建路由构建器
func NewRouteBuilder(method, pattern string) *RouteBuilder {
	return &RouteBuilder{
		method:      method,
		pattern:     pattern,
		constraints: make(map[string]RouteConstraint),
		middlewares: make([]func(HandlerFunc) HandlerFunc, 0),
	}
}

// Constraint 添加约束
func (rb *RouteBuilder) Constraint(paramName string, constraint RouteConstraint) *RouteBuilder {
	rb.constraints[paramName] = constraint
	return rb
}

// Middleware 添加中间件
func (rb *RouteBuilder) Middleware(mw func(HandlerFunc) HandlerFunc) *RouteBuilder {
	rb.middlewares = append(rb.middlewares, mw)
	return rb
}

// Handler 设置处理器
func (rb *RouteBuilder) Handler(handler HandlerFunc) *RouteBuilder {
	rb.handler = handler
	return rb
}

// Build 构建最终处理器
func (rb *RouteBuilder) Build() HandlerFunc {
	handler := rb.handler

	// 应用约束
	if len(rb.constraints) > 0 {
		handler = func(ctx *Context) error {
			for paramName, constraint := range rb.constraints {
				value := ctx.Param(paramName)
				if value == "" || !constraint(value) {
					http.Error(ctx.Writer, "Bad Request: Invalid parameter", http.StatusBadRequest)
					return nil
				}
			}
			return rb.handler(ctx)
		}
	}

	// 应用中间件（逆序）
	for i := len(rb.middlewares) - 1; i >= 0; i-- {
		handler = rb.middlewares[i](handler)
	}

	return handler
}

// 常用正则模式
var (
	// 数字 ID 模式
	IDPattern = `^/[\w-]+/(?P<id>\d+)$`

	// UUID 模式
	UUIDPattern = `^/[\w-]+/(?P<id>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`

	// Slug 模式
	SlugPattern = `^/[\w-]+/(?P<slug>[a-z0-9-]+)$`

	// 日期模式 (YYYY-MM-DD)
	DatePattern = `^/[\w-]+/(?P<date>\d{4}-\d{2}-\d{2})$`

	// 文件扩展名模式
	FileExtPattern = `^/[\w-/]+\.(?P<ext>\w+)$`
)

// MustRegex 编译正则表达式，失败时 panic
func MustRegex(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %v", err))
	}
	return re
}

// MatchesPath 检查路径是否匹配正则模式
func MatchesPath(pattern, path string) bool {
	re := MustRegex(pattern)
	return re.MatchString(path)
}

// ExtractParams 从路径中提取参数
func ExtractParams(pattern, path string) map[string]string {
	re := MustRegex(pattern)
	matches := re.FindStringSubmatch(path)

	if matches == nil {
		return nil
	}

	params := make(map[string]string)
	paramNames := extractParamNames(pattern)

	for i, name := range paramNames {
		if i+1 < len(matches) {
			params[name] = matches[i+1]
		}
	}

	return params
}

// PathBuilder 路径构建器
type PathBuilder struct {
	segments []string
}

// NewPathBuilder 创建路径构建器
func NewPathBuilder() *PathBuilder {
	return &PathBuilder{
		segments: make([]string, 0),
	}
}

// Add 添加路径段
func (pb *PathBuilder) Add(segment string) *PathBuilder {
	pb.segments = append(pb.segments, segment)
	return pb
}

// Build 构建完整路径
func (pb *PathBuilder) Build() string {
	return "/" + strings.Join(pb.segments, "/")
}
