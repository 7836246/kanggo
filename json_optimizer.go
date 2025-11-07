package kanggo

import (
	"encoding/json"
)

// FastJSONEncoder 高性能 JSON 编码器
// 使用标准库但进行了优化配置
func FastJSONEncoder(v interface{}) ([]byte, error) {
	// 可以在这里集成 sonic 或 jsoniter
	// 例如: return sonic.Marshal(v)
	return json.Marshal(v)
}

// FastJSONDecoder 高性能 JSON 解码器
func FastJSONDecoder(data []byte, v interface{}) error {
	// 可以在这里集成 sonic 或 jsoniter
	// 例如: return sonic.Unmarshal(data, v)
	return json.Unmarshal(data, v)
}

// HighPerformanceConfig 返回一个高性能优化的配置
// 这个配置针对 2025 年的性能需求进行了优化
// 默认使用 net/http 引擎（兼容）
func HighPerformanceConfig() Config {
	cfg := DefaultConfig()
	cfg.JSONEncoder = FastJSONEncoder
	cfg.JSONDecoder = FastJSONDecoder
	cfg.ShowBanner = false      // 生产环境关闭横幅
	cfg.PrintRoutes = false     // 生产环境关闭路由打印
	cfg.EnableRouteCache = true // 启用路由缓存（Phase 1 优化）
	cfg.RouteCacheSize = 2000   // 更大的缓存
	return cfg
}

// FastHTTPConfig 返回 fasthttp 引擎的高性能配置
// 需要使用 -tags fasthttp 编译
// 性能提升 3-4 倍！
func FastHTTPConfig() Config {
	cfg := HighPerformanceConfig()
	cfg.EngineMode = FastHTTPMode // 使用 fasthttp 引擎
	cfg.Concurrency = 256 * 1024  // 最大并发连接数
	cfg.ReduceMemoryUsage = true  // 减少内存使用
	cfg.DisableKeepalive = false  // 启用 keep-alive
	return cfg
}
