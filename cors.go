package kanggo

import (
	"net/http"
	"strings"
)

// CorsConfig 定义跨域中间件的配置选项
type CorsConfig struct {
	AllowOrigins     string // 允许的来源，多个来源用逗号分隔
	AllowMethods     string // 允许的方法，多个方法用逗号分隔
	AllowHeaders     string // 允许的请求头，多个请求头用逗号分隔
	AllowCredentials bool   // 是否允许携带凭证
	ExposeHeaders    string // 允许客户端获取的响应头
	MaxAge           int    // 浏览器缓存预检请求的结果的时间（以秒为单位）
}

// DefaultCorsConfig 提供默认的 CORS 配置
func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		AllowOrigins:     "*", // 默认允许所有来源
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}
}

// NewCors 创建一个新的 CORS 中间件
func NewCors(config ...CorsConfig) MiddlewareFunc {
	// 使用默认配置，如果用户提供了自定义配置，则使用用户的配置
	cfg := DefaultCorsConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// 返回中间件函数
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 处理跨域请求
			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(origin, cfg.AllowOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
				w.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)

				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if cfg.ExposeHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", cfg.ExposeHeaders)
				}
				if cfg.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", string(rune(cfg.MaxAge)))
				}

				// 如果是预检请求，直接返回
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			// 继续处理下一个请求
			next.ServeHTTP(w, r)
		}
	}
}

// isOriginAllowed 检查 Origin 是否被允许
func isOriginAllowed(origin, allowOrigins string) bool {
	if allowOrigins == "*" {
		return true
	}
	origins := strings.Split(allowOrigins, ",")
	for _, o := range origins {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}
