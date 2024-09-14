package kanggo

import (
	"github.com/7836246/kanggo/constants" // 引入 constants 包
	"strings"
)

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
	g.Router.Handle(constants.MethodGet, g.Prefix+pattern, handler)
}

// POST 方法为路由组注册一个 POST 请求处理函数
func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodPost, g.Prefix+pattern, handler)
}

// PUT 方法为路由组注册一个 PUT 请求处理函数
func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodPut, g.Prefix+pattern, handler)
}

// DELETE 方法为路由组注册一个 DELETE 请求处理函数
func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodDelete, g.Prefix+pattern, handler)
}

// PATCH 方法为路由组注册一个 PATCH 请求处理函数
func (g *Group) PATCH(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodPatch, g.Prefix+pattern, handler)
}

// OPTIONS 方法为路由组注册一个 OPTIONS 请求处理函数
func (g *Group) OPTIONS(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodOptions, g.Prefix+pattern, handler)
}

// HEAD 方法为路由组注册一个 HEAD 请求处理函数
func (g *Group) HEAD(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodHead, g.Prefix+pattern, handler)
}

// TRACE 方法为路由组注册一个 TRACE 请求处理函数
func (g *Group) TRACE(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodTrace, g.Prefix+pattern, handler)
}

// CONNECT 方法为路由组注册一个 CONNECT 请求处理函数
func (g *Group) CONNECT(pattern string, handler HandlerFunc) {
	g.Router.Handle(constants.MethodConnect, g.Prefix+pattern, handler)
}

// Add 方法允许您指定一个方法作为值来注册一个路由
func (g *Group) Add(method, pattern string, handlers ...HandlerFunc) {
	for _, handler := range handlers {
		g.Router.Handle(method, g.Prefix+pattern, handler)
	}
}

// All 方法将给定路径注册到所有 HTTP 方法
func (g *Group) All(pattern string, handlers ...HandlerFunc) {
	methods := []string{
		constants.MethodGet,
		constants.MethodPost,
		constants.MethodPut,
		constants.MethodDelete,
		constants.MethodPatch,
		constants.MethodOptions,
		constants.MethodHead,
		constants.MethodConnect,
		constants.MethodTrace,
	}

	for _, method := range methods {
		g.Add(method, pattern, handlers...)
	}
}
