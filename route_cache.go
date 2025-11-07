package kanggo

import (
	"sync"
	"sync/atomic"
	"time"
)

// CacheEntry 缓存条目
type CacheEntry struct {
	handler HandlerFunc
	params  map[string]string
	hits    uint64
	time    time.Time
}

// RouteCache 路由缓存
type RouteCache struct {
	cache     sync.Map
	maxSize   int
	size      int32
	enabled   bool
	hitCount  uint64
	missCount uint64
}

// NewRouteCache 创建路由缓存
func NewRouteCache(maxSize int) *RouteCache {
	if maxSize <= 0 {
		maxSize = 1000 // 默认缓存 1000 条
	}

	return &RouteCache{
		maxSize: maxSize,
		enabled: true,
	}
}

// Get 获取缓存的路由
func (rc *RouteCache) Get(method, path string) (HandlerFunc, map[string]string, bool) {
	if !rc.enabled {
		return nil, nil, false
	}

	key := method + ":" + path
	if v, ok := rc.cache.Load(key); ok {
		entry := v.(*CacheEntry)
		atomic.AddUint64(&entry.hits, 1)
		atomic.AddUint64(&rc.hitCount, 1)
		return entry.handler, entry.params, true
	}

	atomic.AddUint64(&rc.missCount, 1)
	return nil, nil, false
}

// Set 设置路由缓存
func (rc *RouteCache) Set(method, path string, handler HandlerFunc, params map[string]string) {
	if !rc.enabled {
		return
	}

	// 检查缓存大小限制
	currentSize := atomic.LoadInt32(&rc.size)
	if int(currentSize) >= rc.maxSize {
		// 缓存已满，使用简单的策略：不添加新条目
		// 在生产环境中可以实现 LRU 策略
		return
	}

	// 复制 params 避免共享
	paramsCopy := make(map[string]string, len(params))
	for k, v := range params {
		paramsCopy[k] = v
	}

	key := method + ":" + path
	entry := &CacheEntry{
		handler: handler,
		params:  paramsCopy,
		hits:    0,
		time:    time.Now(),
	}

	// 只在新增时增加计数
	if _, loaded := rc.cache.LoadOrStore(key, entry); !loaded {
		atomic.AddInt32(&rc.size, 1)
	}
}

// Clear 清空缓存
func (rc *RouteCache) Clear() {
	rc.cache.Range(func(key, value interface{}) bool {
		rc.cache.Delete(key)
		return true
	})
	atomic.StoreInt32(&rc.size, 0)
	atomic.StoreUint64(&rc.hitCount, 0)
	atomic.StoreUint64(&rc.missCount, 0)
}

// Enable 启用缓存
func (rc *RouteCache) Enable() {
	rc.enabled = true
}

// Disable 禁用缓存
func (rc *RouteCache) Disable() {
	rc.enabled = false
}

// Stats 获取缓存统计
func (rc *RouteCache) Stats() map[string]interface{} {
	count := atomic.LoadInt32(&rc.size)
	hits := atomic.LoadUint64(&rc.hitCount)
	misses := atomic.LoadUint64(&rc.missCount)

	totalRequests := hits + misses
	hitRate := float64(0)
	if totalRequests > 0 {
		hitRate = float64(hits) / float64(totalRequests) * 100
	}

	// 收集详细条目信息
	topEntries := make([]map[string]interface{}, 0, 10)
	rc.cache.Range(func(key, value interface{}) bool {
		entry := value.(*CacheEntry)
		if len(topEntries) < 10 {
			topEntries = append(topEntries, map[string]interface{}{
				"route": key,
				"hits":  atomic.LoadUint64(&entry.hits),
				"time":  entry.time,
			})
		}
		return true
	})

	return map[string]interface{}{
		"enabled":        rc.enabled,
		"entries":        count,
		"max_size":       rc.maxSize,
		"hit_count":      hits,
		"miss_count":     misses,
		"hit_rate":       hitRate,
		"total_requests": totalRequests,
		"top_routes":     topEntries,
	}
}

// Size 返回当前缓存条目数
func (rc *RouteCache) Size() int {
	return int(atomic.LoadInt32(&rc.size))
}

// HitRate 返回缓存命中率
func (rc *RouteCache) HitRate() float64 {
	hits := atomic.LoadUint64(&rc.hitCount)
	misses := atomic.LoadUint64(&rc.missCount)
	total := hits + misses

	if total == 0 {
		return 0
	}

	return float64(hits) / float64(total) * 100
}
