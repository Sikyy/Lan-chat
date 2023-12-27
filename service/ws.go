package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader 用于升级 HTTP 连接到 WebSocket 连接
var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: 10 * time.Second, // 设置握手超时
	}

	// 连接池
	Clients = make(map[*websocket.Conn]struct{})

	mu sync.Mutex // 互斥锁，用于保护 clients 映射

	// hub实例
	h = hub{
		c: make(map[*connection]bool), //表示连接是否存在
		b: make(chan []byte),          //用于接收要广播的消息数据
		r: make(chan *connection),     //用于接收要移除的连接
		u: make(chan *connection),     //用于接收新连接
	}
)

// Message 消息结构体
// Message 消息结构体
type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

// connection结构体 用于保存每个连接的信息
type connection struct {
	ws   *websocket.Conn //表示 WebSocket 连接的实例
	sc   chan []byte     //表示用于向客户端发送消息的通道
	data Message         //表示用于存储与连接相关的消息数据的结构体
}

// hub结构体，管理所有的连接实例：connection
type hub struct {
	c  map[*connection]bool //连接池--当前连接的集合（数）
	b  chan []byte          //用于接收要广播的消息数据
	r  chan *connection     //用于接收要移除的连接
	u  chan *connection     //用于接收新连接
	mu sync.Mutex           // 互斥锁
}

// 运行hub，监听各个通道的数据
func (h *hub) run() {
	// 启动定时器，每隔5秒执行一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case c := <-h.r: // 有新连接加入
			fmt.Println("新连接加入:", c.ws.RemoteAddr().String())
			h.c[c] = true
			c.data.Type = "handshake"
			c.data.Content = fmt.Sprintf("新用户已连接：%s", c.ws.RemoteAddr().String())
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
			h.broadcastConnectionCount()
		case c := <-h.u: // 有连接断开
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
				fmt.Println("连接已关闭:", c.ws.RemoteAddr().String())
			}
			h.broadcastConnectionCount()
		case data := <-h.b: // 有消息广播
			fmt.Println("广播消息:", string(data))
			for c := range h.c {
				select {
				case c.sc <- data:
				default:
					delete(h.c, c)
					close(c.sc)
					fmt.Println("无法向客户端发送消息。连接已关闭:", c.ws.RemoteAddr().String())
					h.broadcastConnectionCount()
				}
			}
		case <-ticker.C:
			// 定时器触发时打印当前连接数
			h.printConnectionCount()
		}
	}
}

// 打印当前连接数
func (h *hub) printConnectionCount() {
	h.mu.Lock()
	defer h.mu.Unlock()

	count := len(h.c)
	fmt.Printf("当前连接数: %d\n", count)
}

// 处理连接
func HandleConnection(conn *websocket.Conn) {
	// 设置读写超时
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 设置 Keep-Alive
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	// 日志记录连接关闭
	defer func() {
		fmt.Printf("连接已关闭：%s\n", conn.RemoteAddr().String())
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("读取消息时发生错误:", err)
			return
		}

		// 在这里添加用户名和消息内容到 Message 结构体
		msg.Username = "你的用户名"   // 替换成从前端获取的用户名
		msg.Type = "userMessage" // 设置消息类型为用户消息

		// 将消息内容发送到 hub 的 b 通道
		h.b <- []byte(msg.Content)

		// 打印用户名和消息内容到控制台
		fmt.Printf("收到消息：%s: %s\n", msg.Username, msg.Content)
	}
}

// 用于路由器调用
func RunHub() {
	go h.run()
}

// ToJSON 方法，将 Message 结构体转换为 JSON 字符串
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)

}

// 广播连接数信息给所有连接
func (h *hub) broadcastConnectionCount() {
	count := len(h.c)
	data := Message{Type: "connectionCount", Content: fmt.Sprintf("%d", count)}

	for c := range h.c {
		jsonMsg, err := data.ToJSON()
		if err != nil {
			fmt.Println("Error converting message to JSON:", err)
			continue
		}

		select {
		case c.sc <- jsonMsg: // 使用 c.sc 发送连接数消息
		default:
			// 处理发送失败的情况
			fmt.Println("Failed to send message to client")
		}
	}
}

// 广播消息
// func BroadcastMessage(sender *websocket.Conn, messageType int, message []byte) {
// 	for client := range Clients { //遍历所有客户端
// 		if client != sender {
// 			err := client.WriteMessage(messageType, message) //向客户端发送消息
// 			if err != nil {
// 				fmt.Println("向客户端广播消息时出错:", err)
// 			}
// 		}
// 	}
// }
