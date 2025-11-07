package kanggo

import (
	"fmt"
	"sync"
)

// WebSocketRoom WebSocket 房间，用于管理一组连接
type WebSocketRoom struct {
	name        string
	connections map[string]*WebSocketConn
	mu          sync.RWMutex
	broadcast   chan *BroadcastMessage
	register    chan *WebSocketConn
	unregister  chan *WebSocketConn
	closed      bool
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	MessageType int
	Data        []byte
	ExcludeID   string // 排除的连接 ID（可选）
}

// NewWebSocketRoom 创建新的 WebSocket 房间
func NewWebSocketRoom(name string) *WebSocketRoom {
	room := &WebSocketRoom{
		name:        name,
		connections: make(map[string]*WebSocketConn),
		broadcast:   make(chan *BroadcastMessage, 256),
		register:    make(chan *WebSocketConn, 16),
		unregister:  make(chan *WebSocketConn, 16),
		closed:      false,
	}

	// 启动房间管理协程
	go room.run()

	return room
}

// run 房间管理主循环
func (r *WebSocketRoom) run() {
	for {
		select {
		case conn := <-r.register:
			r.mu.Lock()
			if !r.closed {
				r.connections[conn.ID] = conn
			}
			r.mu.Unlock()

		case conn := <-r.unregister:
			r.mu.Lock()
			if _, ok := r.connections[conn.ID]; ok {
				delete(r.connections, conn.ID)
				conn.Close()
			}
			r.mu.Unlock()

		case msg := <-r.broadcast:
			r.mu.RLock()
			for id, conn := range r.connections {
				// 跳过排除的连接
				if msg.ExcludeID != "" && id == msg.ExcludeID {
					continue
				}

				// 尝试发送消息
				if err := conn.WriteMessage(msg.MessageType, msg.Data); err != nil {
					// 发送失败，标记为需要移除
					go func(c *WebSocketConn) {
						r.unregister <- c
					}(conn)
				}
			}
			r.mu.RUnlock()
		}
	}
}

// Join 连接加入房间
func (r *WebSocketRoom) Join(conn *WebSocketConn) error {
	r.mu.RLock()
	if r.closed {
		r.mu.RUnlock()
		return fmt.Errorf("room is closed")
	}
	r.mu.RUnlock()

	r.register <- conn
	return nil
}

// Leave 连接离开房间
func (r *WebSocketRoom) Leave(conn *WebSocketConn) {
	r.unregister <- conn
}

// Broadcast 向房间内所有连接广播消息
func (r *WebSocketRoom) Broadcast(messageType int, data []byte) {
	r.broadcast <- &BroadcastMessage{
		MessageType: messageType,
		Data:        data,
	}
}

// BroadcastExclude 向房间内所有连接广播消息，排除指定连接
func (r *WebSocketRoom) BroadcastExclude(messageType int, data []byte, excludeID string) {
	r.broadcast <- &BroadcastMessage{
		MessageType: messageType,
		Data:        data,
		ExcludeID:   excludeID,
	}
}

// Count 获取房间内连接数量
func (r *WebSocketRoom) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.connections)
}

// GetConnections 获取房间内所有连接
func (r *WebSocketRoom) GetConnections() []*WebSocketConn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conns := make([]*WebSocketConn, 0, len(r.connections))
	for _, conn := range r.connections {
		conns = append(conns, conn)
	}
	return conns
}

// Close 关闭房间
func (r *WebSocketRoom) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return
	}

	r.closed = true

	// 关闭所有连接
	for _, conn := range r.connections {
		conn.Close()
	}

	r.connections = make(map[string]*WebSocketConn)
}

// Name 获取房间名称
func (r *WebSocketRoom) Name() string {
	return r.name
}

// WebSocketHub WebSocket 中心，管理多个房间
type WebSocketHub struct {
	rooms map[string]*WebSocketRoom
	mu    sync.RWMutex
}

// NewWebSocketHub 创建新的 WebSocket 中心
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		rooms: make(map[string]*WebSocketRoom),
	}
}

// GetOrCreateRoom 获取或创建房间
func (h *WebSocketHub) GetOrCreateRoom(name string) *WebSocketRoom {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[name]; ok {
		return room
	}

	room := NewWebSocketRoom(name)
	h.rooms[name] = room
	return room
}

// GetRoom 获取房间
func (h *WebSocketHub) GetRoom(name string) (*WebSocketRoom, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[name]
	return room, ok
}

// DeleteRoom 删除房间
func (h *WebSocketHub) DeleteRoom(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[name]; ok {
		room.Close()
		delete(h.rooms, name)
	}
}

// GetAllRooms 获取所有房间
func (h *WebSocketHub) GetAllRooms() []*WebSocketRoom {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rooms := make([]*WebSocketRoom, 0, len(h.rooms))
	for _, room := range h.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// RoomCount 获取房间数量
func (h *WebSocketHub) RoomCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

// Close 关闭所有房间
func (h *WebSocketHub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, room := range h.rooms {
		room.Close()
	}

	h.rooms = make(map[string]*WebSocketRoom)
}
