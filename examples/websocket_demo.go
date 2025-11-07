package main

import (
	"fmt"
	"log"
	"time"

	"github.com/7836246/kanggo"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("🚀 KangGo WebSocket 演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 创建应用
	app := kanggo.Default()

	// 创建 WebSocket Hub
	hub := kanggo.NewWebSocketHub()

	// 演示 1: 简单的 Echo WebSocket
	demo1SimpleEcho(app)

	// 演示 2: JSON 消息处理
	demo2JSONMessages(app)

	// 演示 3: 聊天室（带房间和广播）
	demo3ChatRoom(app, hub)

	// 演示 4: 心跳检测
	demo4HeartBeat(app)

	// 首页
	app.GET("/", func(ctx *kanggo.Context) error {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>KangGo WebSocket Demo</title>
    <style>
        body { font-family: Arial; max-width: 800px; margin: 50px auto; }
        .demo { border: 1px solid #ddd; padding: 20px; margin: 20px 0; }
        button { padding: 10px 20px; margin: 5px; cursor: pointer; }
        #messages { border: 1px solid #ddd; height: 300px; overflow-y: auto; padding: 10px; }
        .message { margin: 5px 0; }
        input { padding: 10px; width: 300px; }
    </style>
</head>
<body>
    <h1>🚀 KangGo WebSocket 演示</h1>
    
    <div class="demo">
        <h2>演示 1: Echo WebSocket</h2>
        <button onclick="testEcho()">连接并测试</button>
        <div id="echo-status"></div>
    </div>
    
    <div class="demo">
        <h2>演示 2: JSON 消息</h2>
        <button onclick="testJSON()">发送 JSON</button>
        <div id="json-status"></div>
    </div>
    
    <div class="demo">
        <h2>演示 3: 聊天室</h2>
        <input type="text" id="username" placeholder="用户名" value="User1">
        <input type="text" id="message" placeholder="消息">
        <button onclick="sendMessage()">发送</button>
        <div id="messages"></div>
    </div>

    <script>
        let echoWs = null;
        let chatWs = null;

        // 演示 1: Echo
        function testEcho() {
            const status = document.getElementById('echo-status');
            echoWs = new WebSocket('ws://localhost:8080/ws/echo');
            
            echoWs.onopen = () => {
                status.innerHTML = '✅ 连接成功';
                echoWs.send('Hello, WebSocket!');
            };
            
            echoWs.onmessage = (e) => {
                status.innerHTML += '<br>📩 收到: ' + e.data;
            };
            
            echoWs.onerror = (e) => {
                status.innerHTML = '❌ 错误: ' + e;
            };
        }

        // 演示 2: JSON
        function testJSON() {
            const status = document.getElementById('json-status');
            const ws = new WebSocket('ws://localhost:8080/ws/json');
            
            ws.onopen = () => {
                status.innerHTML = '✅ 连接成功';
                ws.send(JSON.stringify({
                    type: 'greeting',
                    name: 'Alice'
                }));
            };
            
            ws.onmessage = (e) => {
                const data = JSON.parse(e.data);
                status.innerHTML += '<br>📩 ' + JSON.stringify(data, null, 2);
            };
        }

        // 演示 3: 聊天室
        if (!chatWs) {
            chatWs = new WebSocket('ws://localhost:8080/ws/chat');
            
            chatWs.onmessage = (e) => {
                const messages = document.getElementById('messages');
                const div = document.createElement('div');
                div.className = 'message';
                div.textContent = e.data;
                messages.appendChild(div);
                messages.scrollTop = messages.scrollHeight;
            };
        }

        function sendMessage() {
            const username = document.getElementById('username').value;
            const message = document.getElementById('message').value;
            
            if (chatWs && message) {
                chatWs.send(JSON.stringify({
                    username: username,
                    message: message
                }));
                document.getElementById('message').value = '';
            }
        }
    </script>
</body>
</html>
`
		return ctx.HTML(200, html)
	})

	fmt.Println("🌐 服务器启动在 http://localhost:8080")
	fmt.Println()
	fmt.Println("测试地址：")
	fmt.Println("  • Echo:     ws://localhost:8080/ws/echo")
	fmt.Println("  • JSON:     ws://localhost:8080/ws/json")
	fmt.Println("  • Chat:     ws://localhost:8080/ws/chat")
	fmt.Println("  • 网页演示:  http://localhost:8080")
	fmt.Println()

	app.Run(":8080")
}

// demo1SimpleEcho 演示 1: 简单的 Echo WebSocket
func demo1SimpleEcho(app *kanggo.KangGo) {
	app.WebSocket("/ws/echo", func(conn *kanggo.WebSocketConn) error {
		log.Printf("新连接: %s", conn.ID)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("读取错误: %v", err)
				break
			}

			log.Printf("收到消息: %s", message)

			// Echo 回去
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Printf("写入错误: %v", err)
				break
			}
		}

		log.Printf("连接关闭: %s", conn.ID)
		return nil
	})
}

// demo2JSONMessages 演示 2: JSON 消息处理
func demo2JSONMessages(app *kanggo.KangGo) {
	type Request struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}

	type Response struct {
		Type      string `json:"type"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	}

	app.WebSocket("/ws/json", func(conn *kanggo.WebSocketConn) error {
		for {
			var req Request
			if err := conn.ReadJSON(&req); err != nil {
				log.Printf("JSON 读取错误: %v", err)
				break
			}

			log.Printf("收到 JSON: %+v", req)

			// 响应
			resp := Response{
				Type:      "response",
				Message:   fmt.Sprintf("Hello, %s!", req.Name),
				Timestamp: time.Now().Unix(),
			}

			if err := conn.WriteJSON(resp); err != nil {
				log.Printf("JSON 写入错误: %v", err)
				break
			}
		}

		return nil
	})
}

// demo3ChatRoom 演示 3: 聊天室
func demo3ChatRoom(app *kanggo.KangGo, hub *kanggo.WebSocketHub) {
	type ChatMessage struct {
		Username string `json:"username"`
		Message  string `json:"message"`
	}

	app.WebSocket("/ws/chat", func(conn *kanggo.WebSocketConn) error {
		// 加入聊天室
		room := hub.GetOrCreateRoom("chat")
		room.Join(conn)

		log.Printf("用户 %s 加入聊天室，当前 %d 人", conn.ID, room.Count())

		// 发送欢迎消息
		welcomeMsg := fmt.Sprintf("欢迎！当前在线: %d 人", room.Count())
		conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

		// 读取消息循环
		for {
			var msg ChatMessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("读取错误: %v", err)
				break
			}

			// 广播消息给所有人
			broadcastMsg := fmt.Sprintf("[%s]: %s", msg.Username, msg.Message)
			room.Broadcast(websocket.TextMessage, []byte(broadcastMsg))
		}

		// 离开房间
		room.Leave(conn)
		log.Printf("用户 %s 离开聊天室，剩余 %d 人", conn.ID, room.Count())

		return nil
	})
}

// demo4HeartBeat 演示 4: 心跳检测
func demo4HeartBeat(app *kanggo.KangGo) {
	app.WebSocket("/ws/heartbeat", func(conn *kanggo.WebSocketConn) error {
		// 设置 Pong 处理器
		conn.SetPongHandler(func(appData string) error {
			log.Printf("收到 Pong: %s", conn.ID)
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		// 启动心跳发送协程
		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						log.Printf("Ping 失败: %v", err)
						return
					}
				case <-done:
					return
				}
			}
		}()

		// 读取消息
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			log.Printf("心跳消息: %s", message)
		}

		close(done)
		return nil
	})
}
