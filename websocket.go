package kanggo

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHandler WebSocket 处理函数类型
type WebSocketHandler func(*WebSocketConn) error

// WebSocketConfig WebSocket 配置
type WebSocketConfig struct {
	// 读取缓冲区大小
	ReadBufferSize int
	// 写入缓冲区大小
	WriteBufferSize int
	// 握手超时时间
	HandshakeTimeout time.Duration
	// 允许的源（CORS）
	CheckOrigin func(r *http.Request) bool
	// 子协议
	Subprotocols []string
	// 压缩支持
	EnableCompression bool
}

// DefaultWebSocketConfig 返回默认的 WebSocket 配置
func DefaultWebSocketConfig() WebSocketConfig {
	return WebSocketConfig{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		HandshakeTimeout:  10 * time.Second,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: false,
	}
}

// WebSocketConn WebSocket 连接封装
type WebSocketConn struct {
	conn      *websocket.Conn
	mu        sync.Mutex
	writeMu   sync.Mutex
	closeChan chan struct{}
	closed    bool
	ID        string                 // 连接唯一标识
	Data      map[string]interface{} // 自定义数据存储
}

// newWebSocketConn 创建新的 WebSocket 连接
func newWebSocketConn(conn *websocket.Conn, id string) *WebSocketConn {
	return &WebSocketConn{
		conn:      conn,
		closeChan: make(chan struct{}),
		closed:    false,
		ID:        id,
		Data:      make(map[string]interface{}),
	}
}

// ReadMessage 读取消息
func (c *WebSocketConn) ReadMessage() (messageType int, message []byte, err error) {
	return c.conn.ReadMessage()
}

// WriteMessage 写入消息（线程安全）
func (c *WebSocketConn) WriteMessage(messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.closed {
		return fmt.Errorf("connection closed")
	}

	return c.conn.WriteMessage(messageType, data)
}

// WriteJSON 写入 JSON 消息
func (c *WebSocketConn) WriteJSON(v interface{}) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.closed {
		return fmt.Errorf("connection closed")
	}

	return c.conn.WriteJSON(v)
}

// ReadJSON 读取 JSON 消息
func (c *WebSocketConn) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}

// Close 关闭连接
func (c *WebSocketConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	close(c.closeChan)
	return c.conn.Close()
}

// IsClosed 检查连接是否已关闭
func (c *WebSocketConn) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// SetReadDeadline 设置读取超时
func (c *WebSocketConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时
func (c *WebSocketConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// SetPingHandler 设置 Ping 处理器
func (c *WebSocketConn) SetPingHandler(h func(string) error) {
	c.conn.SetPingHandler(h)
}

// SetPongHandler 设置 Pong 处理器
func (c *WebSocketConn) SetPongHandler(h func(string) error) {
	c.conn.SetPongHandler(h)
}

// SetCloseHandler 设置关闭处理器
func (c *WebSocketConn) SetCloseHandler(h func(code int, text string) error) {
	c.conn.SetCloseHandler(h)
}

// WebSocket 注册 WebSocket 路由
func (k *KangGo) WebSocket(pattern string, handler WebSocketHandler, config ...WebSocketConfig) *KangGo {
	cfg := DefaultWebSocketConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// 配置 upgrader
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:    cfg.ReadBufferSize,
		WriteBufferSize:   cfg.WriteBufferSize,
		HandshakeTimeout:  cfg.HandshakeTimeout,
		CheckOrigin:       cfg.CheckOrigin,
		Subprotocols:      cfg.Subprotocols,
		EnableCompression: cfg.EnableCompression,
	}

	// 注册 HTTP 处理器
	httpHandler := func(ctx *Context) error {
		// 升级 HTTP 连接为 WebSocket
		conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return fmt.Errorf("websocket upgrade failed: %w", err)
		}

		// 创建 WebSocket 连接封装
		wsConn := newWebSocketConn(conn, generateConnID())

		// 调用用户处理器
		if err := handler(wsConn); err != nil {
			wsConn.Close()
			return err
		}

		return nil
	}

	// 注册到路由
	k.GET(pattern, httpHandler)

	return k
}

// generateConnID 生成连接 ID
func generateConnID() string {
	return fmt.Sprintf("ws-%d", time.Now().UnixNano())
}
