package kanggo

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SSEEvent SSE 事件
type SSEEvent struct {
	ID    string // 事件 ID（可选）
	Event string // 事件类型（可选）
	Data  string // 事件数据
	Retry int    // 重连时间（毫秒，可选）
}

// Format 格式化 SSE 事件为标准格式
func (e *SSEEvent) Format() string {
	result := ""

	if e.ID != "" {
		result += fmt.Sprintf("id: %s\n", e.ID)
	}

	if e.Event != "" {
		result += fmt.Sprintf("event: %s\n", e.Event)
	}

	if e.Retry > 0 {
		result += fmt.Sprintf("retry: %d\n", e.Retry)
	}

	result += fmt.Sprintf("data: %s\n\n", e.Data)

	return result
}

// SSEConn SSE 连接封装
type SSEConn struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	mu      sync.Mutex
	closed  bool
	ID      string                 // 连接唯一标识
	Data    map[string]interface{} // 自定义数据存储
}

// newSSEConn 创建新的 SSE 连接
func newSSEConn(w http.ResponseWriter, id string) (*SSEConn, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	return &SSEConn{
		writer:  w,
		flusher: flusher,
		closed:  false,
		ID:      id,
		Data:    make(map[string]interface{}),
	}, nil
}

// Send 发送 SSE 事件
func (c *SSEConn) Send(event *SSEEvent) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("connection closed")
	}

	// 写入事件数据
	if _, err := fmt.Fprint(c.writer, event.Format()); err != nil {
		return err
	}

	// 立即刷新
	c.flusher.Flush()

	return nil
}

// SendData 发送简单数据（无事件类型和 ID）
func (c *SSEConn) SendData(data string) error {
	return c.Send(&SSEEvent{Data: data})
}

// SendEvent 发送带类型的事件
func (c *SSEConn) SendEvent(eventType, data string) error {
	return c.Send(&SSEEvent{
		Event: eventType,
		Data:  data,
	})
}

// Close 关闭连接
func (c *SSEConn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
}

// IsClosed 检查连接是否已关闭
func (c *SSEConn) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// SSEHandler SSE 处理函数类型
type SSEHandler func(*SSEConn) error

// SSE 注册 SSE 路由
func (k *KangGo) SSE(pattern string, handler SSEHandler) *KangGo {
	httpHandler := func(ctx *Context) error {
		// 设置 SSE 响应头
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 创建 SSE 连接
		conn, err := newSSEConn(ctx.Writer, generateSSEConnID())
		if err != nil {
			return err
		}

		// 调用用户处理器
		return handler(conn)
	}

	k.GET(pattern, httpHandler)

	return k
}

// generateSSEConnID 生成 SSE 连接 ID
func generateSSEConnID() string {
	return fmt.Sprintf("sse-%d", time.Now().UnixNano())
}

// SSEBroadcaster SSE 广播器
type SSEBroadcaster struct {
	clients    map[string]*SSEConn
	mu         sync.RWMutex
	register   chan *SSEConn
	unregister chan *SSEConn
	broadcast  chan *SSEEvent
	closed     bool
}

// NewSSEBroadcaster 创建新的 SSE 广播器
func NewSSEBroadcaster() *SSEBroadcaster {
	b := &SSEBroadcaster{
		clients:    make(map[string]*SSEConn),
		register:   make(chan *SSEConn, 16),
		unregister: make(chan *SSEConn, 16),
		broadcast:  make(chan *SSEEvent, 256),
		closed:     false,
	}

	// 启动广播器
	go b.run()

	return b
}

// run 广播器主循环
func (b *SSEBroadcaster) run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			if !b.closed {
				b.clients[client.ID] = client
			}
			b.mu.Unlock()

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client.ID]; ok {
				delete(b.clients, client.ID)
				client.Close()
			}
			b.mu.Unlock()

		case event := <-b.broadcast:
			b.mu.RLock()
			for _, client := range b.clients {
				if err := client.Send(event); err != nil {
					// 发送失败，标记为需要移除
					go func(c *SSEConn) {
						b.unregister <- c
					}(client)
				}
			}
			b.mu.RUnlock()
		}
	}
}

// Register 注册客户端
func (b *SSEBroadcaster) Register(client *SSEConn) {
	b.register <- client
}

// Unregister 注销客户端
func (b *SSEBroadcaster) Unregister(client *SSEConn) {
	b.unregister <- client
}

// Broadcast 广播事件
func (b *SSEBroadcaster) Broadcast(event *SSEEvent) {
	b.broadcast <- event
}

// BroadcastData 广播简单数据
func (b *SSEBroadcaster) BroadcastData(data string) {
	b.Broadcast(&SSEEvent{Data: data})
}

// BroadcastEvent 广播带类型的事件
func (b *SSEBroadcaster) BroadcastEvent(eventType, data string) {
	b.Broadcast(&SSEEvent{
		Event: eventType,
		Data:  data,
	})
}

// ClientCount 获取客户端数量
func (b *SSEBroadcaster) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// Close 关闭广播器
func (b *SSEBroadcaster) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true

	// 关闭所有客户端
	for _, client := range b.clients {
		client.Close()
	}

	b.clients = make(map[string]*SSEConn)
}
