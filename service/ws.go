package service

import (
	"encoding/json"
	"fmt"
	"net/http"
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
)

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

func HandleConnection(conn *websocket.Conn) {

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
		H.b <- []byte(msg.Content)

		// 打印用户名和消息内容到控制台
		fmt.Printf("收到消息：%s: %s\n", msg.Username, msg.Content)
	}
}

// ToJSON 方法，将 Message 结构体转换为 JSON 字符串
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)

}

// 广播连接数信息给所有连接
func (h *Hub) broadcastConnectionCount() {
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
