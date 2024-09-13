package kanggo

import "strings"

// Group 结构定义了一个路由组
type Group struct {
	Prefix string  // 路由组的前缀
	Router *Router // 引用 Router
}

// NewGroup 创建一个新的路由组
func (r *Router) NewGroup(prefix string) *Group {
	return &Group{
		Prefix: strings.TrimSuffix(prefix, "/"), // 去除尾部的 "/"
		Router: r,
	}
}

// GET 方法为路由组注册一个 GET 请求处理函数
func (g *Group) GET(pattern string, handler HandlerFunc) {
	g.Router.Handle("GET", g.Prefix+pattern, handler)
}

// POST 方法为路由组注册一个 POST 请求处理函数
func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.Router.Handle("POST", g.Prefix+pattern, handler)
}

// PUT 方法为路由组注册一个 PUT 请求处理函数
func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.Router.Handle("PUT", g.Prefix+pattern, handler)
}

// DELETE 方法为路由组注册一个 DELETE 请求处理函数
func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.Router.Handle("DELETE", g.Prefix+pattern, handler)
}

// PATCH 方法为路由组注册一个 PATCH 请求处理函数
func (g *Group) PATCH(pattern string, handler HandlerFunc) {
	g.Router.Handle("PATCH", g.Prefix+pattern, handler)
}

// OPTIONS 方法为路由组注册一个 OPTIONS 请求处理函数
func (g *Group) OPTIONS(pattern string, handler HandlerFunc) {
	g.Router.Handle("OPTIONS", g.Prefix+pattern, handler)
}

// HEAD 方法为路由组注册一个 HEAD 请求处理函数
func (g *Group) HEAD(pattern string, handler HandlerFunc) {
	g.Router.Handle("HEAD", g.Prefix+pattern, handler)
}

// TRACE 方法为路由组注册一个 TRACE 请求处理函数
func (g *Group) TRACE(pattern string, handler HandlerFunc) {
	g.Router.Handle("TRACE", g.Prefix+pattern, handler)
}

// CONNECT 方法为路由组注册一个 CONNECT 请求处理函数
func (g *Group) CONNECT(pattern string, handler HandlerFunc) {
	g.Router.Handle("CONNECT", g.Prefix+pattern, handler)
}
