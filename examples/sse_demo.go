package main

import (
	"fmt"
	"log"
	"time"

	"github.com/7836246/kanggo"
)

func main() {
	fmt.Println("🚀 KangGo Server-Sent Events (SSE) 演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	app := kanggo.Default()

	// 创建 SSE 广播器
	broadcaster := kanggo.NewSSEBroadcaster()

	// 演示 1: 简单的 SSE 流
	demo1SimpleStream(app)

	// 演示 2: 时钟更新
	demo2ClockUpdate(app)

	// 演示 3: 广播（多客户端）
	demo3Broadcast(app, broadcaster)

	// 演示 4: 带事件类型的 SSE
	demo4TypedEvents(app)

	// 首页
	app.GET("/", func(ctx *kanggo.Context) error {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>KangGo SSE Demo</title>
    <style>
        body { font-family: Arial; max-width: 1000px; margin: 50px auto; }
        .demo { border: 1px solid #ddd; padding: 20px; margin: 20px 0; }
        button { padding: 10px 20px; margin: 5px; cursor: pointer; }
        .events { border: 1px solid #ddd; height: 200px; overflow-y: auto; padding: 10px; background: #f5f5f5; }
        .event { margin: 5px 0; padding: 5px; background: white; }
        .status { color: green; font-weight: bold; }
        .error { color: red; }
    </style>
</head>
<body>
    <h1>🚀 KangGo Server-Sent Events 演示</h1>
    
    <div class="demo">
        <h2>演示 1: 简单的 SSE 流</h2>
        <button onclick="startSimpleStream()">开始接收</button>
        <button onclick="stopSimpleStream()">停止</button>
        <div id="simple-status"></div>
        <div id="simple-events" class="events"></div>
    </div>
    
    <div class="demo">
        <h2>演示 2: 实时时钟</h2>
        <button onclick="startClock()">开始</button>
        <button onclick="stopClock()">停止</button>
        <div id="clock-time" style="font-size: 24px; font-weight: bold;"></div>
    </div>
    
    <div class="demo">
        <h2>演示 3: 广播（多客户端同步）</h2>
        <button onclick="startBroadcast()">连接广播</button>
        <button onclick="stopBroadcast()">断开</button>
        <div id="broadcast-status"></div>
        <div id="broadcast-events" class="events"></div>
    </div>
    
    <div class="demo">
        <h2>演示 4: 带类型的事件</h2>
        <button onclick="startTypedEvents()">开始接收</button>
        <button onclick="stopTypedEvents()">停止</button>
        <div id="typed-events" class="events"></div>
    </div>

    <script>
        let simpleSource = null;
        let clockSource = null;
        let broadcastSource = null;
        let typedSource = null;

        // 演示 1: 简单流
        function startSimpleStream() {
            if (simpleSource) {
                simpleSource.close();
            }

            simpleSource = new EventSource('/sse/simple');
            const status = document.getElementById('simple-status');
            const events = document.getElementById('simple-events');
            
            simpleSource.onopen = () => {
                status.innerHTML = '<span class="status">✅ 已连接</span>';
            };
            
            simpleSource.onmessage = (e) => {
                const div = document.createElement('div');
                div.className = 'event';
                div.textContent = new Date().toLocaleTimeString() + ': ' + e.data;
                events.appendChild(div);
                events.scrollTop = events.scrollHeight;
            };
            
            simpleSource.onerror = () => {
                status.innerHTML = '<span class="error">❌ 连接错误</span>';
            };
        }

        function stopSimpleStream() {
            if (simpleSource) {
                simpleSource.close();
                document.getElementById('simple-status').innerHTML = '⏸️ 已停止';
            }
        }

        // 演示 2: 时钟
        function startClock() {
            if (clockSource) {
                clockSource.close();
            }

            clockSource = new EventSource('/sse/clock');
            
            clockSource.onmessage = (e) => {
                document.getElementById('clock-time').textContent = e.data;
            };
        }

        function stopClock() {
            if (clockSource) {
                clockSource.close();
                document.getElementById('clock-time').textContent = '⏸️ 已停止';
            }
        }

        // 演示 3: 广播
        function startBroadcast() {
            if (broadcastSource) {
                broadcastSource.close();
            }

            broadcastSource = new EventSource('/sse/broadcast');
            const status = document.getElementById('broadcast-status');
            const events = document.getElementById('broadcast-events');
            
            broadcastSource.onopen = () => {
                status.innerHTML = '<span class="status">✅ 已连接到广播</span>';
            };
            
            broadcastSource.onmessage = (e) => {
                const div = document.createElement('div');
                div.className = 'event';
                div.textContent = e.data;
                events.appendChild(div);
                events.scrollTop = events.scrollHeight;
            };
        }

        function stopBroadcast() {
            if (broadcastSource) {
                broadcastSource.close();
                document.getElementById('broadcast-status').innerHTML = '⏸️ 已断开';
            }
        }

        // 演示 4: 带类型的事件
        function startTypedEvents() {
            if (typedSource) {
                typedSource.close();
            }

            typedSource = new EventSource('/sse/typed');
            const events = document.getElementById('typed-events');
            
            // 监听不同类型的事件
            typedSource.addEventListener('info', (e) => {
                addTypedEvent(events, 'ℹ️ [INFO]', e.data, '#d1ecf1');
            });
            
            typedSource.addEventListener('warning', (e) => {
                addTypedEvent(events, '⚠️ [WARNING]', e.data, '#fff3cd');
            });
            
            typedSource.addEventListener('error', (e) => {
                addTypedEvent(events, '❌ [ERROR]', e.data, '#f8d7da');
            });
        }

        function addTypedEvent(container, prefix, data, color) {
            const div = document.createElement('div');
            div.className = 'event';
            div.style.backgroundColor = color;
            div.textContent = prefix + ' ' + data;
            container.appendChild(div);
            container.scrollTop = container.scrollHeight;
        }

        function stopTypedEvents() {
            if (typedSource) {
                typedSource.close();
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
	fmt.Println("  • 简单流:    /sse/simple")
	fmt.Println("  • 时钟:      /sse/clock")
	fmt.Println("  • 广播:      /sse/broadcast")
	fmt.Println("  • 类型事件:  /sse/typed")
	fmt.Println("  • 网页演示:  http://localhost:8080")
	fmt.Println()

	app.Run(":8080")
}

// demo1SimpleStream 演示 1: 简单的 SSE 流
func demo1SimpleStream(app *kanggo.KangGo) {
	app.SSE("/sse/simple", func(conn *kanggo.SSEConn) error {
		log.Printf("新 SSE 连接: %s", conn.ID)

		// 发送欢迎消息
		conn.SendData("连接成功！")

		// 每秒发送一条消息
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ticker.C:
				count++
				msg := fmt.Sprintf("消息 #%d - %s", count, time.Now().Format("15:04:05"))

				if err := conn.SendData(msg); err != nil {
					log.Printf("发送失败: %v", err)
					return nil
				}

				// 发送 10 条后停止
				if count >= 10 {
					conn.SendData("流结束")
					return nil
				}
			}
		}
	})
}

// demo2ClockUpdate 演示 2: 实时时钟更新
func demo2ClockUpdate(app *kanggo.KangGo) {
	app.SSE("/sse/clock", func(conn *kanggo.SSEConn) error {
		log.Printf("时钟连接: %s", conn.ID)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				timeStr := t.Format("2006-01-02 15:04:05")

				if err := conn.SendData(timeStr); err != nil {
					return nil
				}
			}
		}
	})
}

// demo3Broadcast 演示 3: 广播到多个客户端
func demo3Broadcast(app *kanggo.KangGo, broadcaster *kanggo.SSEBroadcaster) {
	// 启动广播任务
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ticker.C:
				count++
				msg := fmt.Sprintf("广播消息 #%d - 当前 %d 个连接 - %s",
					count,
					broadcaster.ClientCount(),
					time.Now().Format("15:04:05"))

				broadcaster.BroadcastData(msg)
				log.Printf("广播: %s", msg)
			}
		}
	}()

	// SSE 端点
	app.SSE("/sse/broadcast", func(conn *kanggo.SSEConn) error {
		log.Printf("新客户端连接到广播: %s", conn.ID)

		// 注册到广播器
		broadcaster.Register(conn)
		defer broadcaster.Unregister(conn)

		// 发送欢迎消息
		conn.SendData(fmt.Sprintf("欢迎！你是第 %d 个连接", broadcaster.ClientCount()))

		// 保持连接
		for !conn.IsClosed() {
			time.Sleep(1 * time.Second)
		}

		log.Printf("客户端断开: %s", conn.ID)
		return nil
	})
}

// demo4TypedEvents 演示 4: 带事件类型的 SSE
func demo4TypedEvents(app *kanggo.KangGo) {
	app.SSE("/sse/typed", func(conn *kanggo.SSEConn) error {
		log.Printf("类型事件连接: %s", conn.ID)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		events := []struct {
			eventType string
			message   string
		}{
			{"info", "系统启动完成"},
			{"info", "正在加载数据..."},
			{"warning", "磁盘空间不足"},
			{"info", "数据加载完成"},
			{"error", "无法连接到数据库"},
			{"info", "正在重试..."},
			{"info", "连接成功"},
		}

		for i, event := range events {
			<-ticker.C

			if err := conn.SendEvent(event.eventType, event.message); err != nil {
				return nil
			}

			log.Printf("[%s] %s", event.eventType, event.message)

			if i >= len(events)-1 {
				time.Sleep(2 * time.Second)
				conn.SendEvent("info", "演示结束")
				return nil
			}
		}

		return nil
	})
}
