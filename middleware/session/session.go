package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/7836246/kanggo/core"
)

// Store 定义会话存储的接口
type Store interface {
	Get(sessionID string) (map[string]interface{}, bool)
	Set(sessionID string, data map[string]interface{})
	Delete(sessionID string)
}

// MemoryStore 使用内存存储会话数据
type MemoryStore struct {
	sessions map[string]map[string]interface{}
	mu       sync.RWMutex
}

// NewMemoryStore 创建一个新的 MemoryStore 实例
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]map[string]interface{}),
	}
}

// Get 获取会话数据
func (store *MemoryStore) Get(sessionID string) (map[string]interface{}, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	data, exists := store.sessions[sessionID]
	return data, exists
}

// Set 设置会话数据
func (store *MemoryStore) Set(sessionID string, data map[string]interface{}) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions[sessionID] = data
}

// Delete 删除会话数据
func (store *MemoryStore) Delete(sessionID string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.sessions, sessionID)
}

// generateSessionID 生成一个新的随机 Session ID
func generateSessionID() string {
	bytes := make([]byte, 16) // 16字节的随机数，等同于128位
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // 如果读取随机字节失败，则抛出错误
	}
	return hex.EncodeToString(bytes) // 返回十六进制编码的字符串
}

// New Middleware 中间件函数，用于管理会话
func New(store Store) core.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查请求中是否存在 session ID
			sessionID, err := r.Cookie("session_id")
			if err != nil || sessionID.Value == "" {
				// 如果没有 session ID，则创建一个新的
				newSessionID := generateSessionID()
				http.SetCookie(w, &http.Cookie{
					Name:  "session_id",
					Value: newSessionID,
					Path:  "/",
				})
				sessionID = &http.Cookie{Name: "session_id", Value: newSessionID}
			}

			// 获取或初始化会话数据
			sessionData, exists := store.Get(sessionID.Value)
			if !exists {
				sessionData = make(map[string]interface{})
			}

			// 创建一个自定义的上下文，以便处理会话数据
			ctx := &Context{
				ResponseWriter: w,
				Request:        r,
				SessionID:      sessionID.Value,
				SessionData:    sessionData,
				Store:          store,
			}

			// 调用下一个处理程序
			next(ctx.ResponseWriter, r)

			// 更新会话数据
			store.Set(sessionID.Value, ctx.SessionData)
		}
	}
}

// Context 用于处理会话数据的自定义上下文
type Context struct {
	http.ResponseWriter
	Request     *http.Request
	SessionID   string
	SessionData map[string]interface{}
	Store       Store
}

// GetSessionValue 获取会话中的值
func (ctx *Context) GetSessionValue(key string) interface{} {
	return ctx.SessionData[key]
}

// SetSessionValue 设置会话中的值
func (ctx *Context) SetSessionValue(key string, value interface{}) {
	ctx.SessionData[key] = value
}
