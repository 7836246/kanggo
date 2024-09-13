package kanggo

import "strings"

// Group 结构定义了一个路由组
type Group struct {
	prefix string  // 路由组的前缀
	router *Router // 引用 Router
}

// NewGroup 创建一个新的路由组
func (r *Router) NewGroup(prefix string) *Group {
	return &Group{
		prefix: strings.TrimSuffix(prefix, "/"), // 去除尾部的 "/"
		router: r,
	}
}

// GET 方法为路由组注册一个 GET 请求处理函数
func (g *Group) GET(pattern string, handler HandlerFunc) {
	g.router.Handle("GET", g.prefix+pattern, handler)
}

// POST 方法为路由组注册一个 POST 请求处理函数
func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.router.Handle("POST", g.prefix+pattern, handler)
}

// PUT 方法为路由组注册一个 PUT 请求处理函数
func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.router.Handle("PUT", g.prefix+pattern, handler)
}

// DELETE 方法为路由组注册一个 DELETE 请求处理函数
func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.router.Handle("DELETE", g.prefix+pattern, handler)
}

// PATCH 方法为路由组注册一个 PATCH 请求处理函数
func (g *Group) PATCH(pattern string, handler HandlerFunc) {
	g.router.Handle("PATCH", g.prefix+pattern, handler)
}

// OPTIONS 方法为路由组注册一个 OPTIONS 请求处理函数
func (g *Group) OPTIONS(pattern string, handler HandlerFunc) {
	g.router.Handle("OPTIONS", g.prefix+pattern, handler)
}

// HEAD 方法为路由组注册一个 HEAD 请求处理函数
func (g *Group) HEAD(pattern string, handler HandlerFunc) {
	g.router.Handle("HEAD", g.prefix+pattern, handler)
}

// TRACE 方法为路由组注册一个 TRACE 请求处理函数
func (g *Group) TRACE(pattern string, handler HandlerFunc) {
	g.router.Handle("TRACE", g.prefix+pattern, handler)
}

// CONNECT 方法为路由组注册一个 CONNECT 请求处理函数
func (g *Group) CONNECT(pattern string, handler HandlerFunc) {
	g.router.Handle("CONNECT", g.prefix+pattern, handler)
}
