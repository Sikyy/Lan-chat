package service

//ws.go
import (
	"Lan-chat/data"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection结构体 用于保存每个连接的信息
type connection struct {
	ws   *websocket.Conn //表示 WebSocket 连接的实例
	sc   chan []byte     //表示用于向客户端发送消息的通道
	data *data.Data      //表示用于存储与连接相关的消息数据的结构体
}

func HandleConnection(conn *websocket.Conn) {

	// 设置 Keep-Alive
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	// 日志记录连接建立
	fmt.Printf("WebSocket 连接已建立：%s\n", conn.RemoteAddr().String())

	// defer 用于连接关闭时的日志记录
	defer func() {
		fmt.Printf("WebSocket 连接已关闭：%s\n", conn.RemoteAddr().String())
	}()

	// 创建 connection 实例，用于管理该 WebSocket 连接的信息
	c := &connection{sc: make(chan []byte, 256), ws: conn, data: &data.Data{}}
	H.register <- c

	// 启动 connection 的写入和读取协程
	go c.writer()
	c.reader()

	// 在连接关闭时执行清理操作
	defer func() {
		// 向服务器发送注销消息
		c.data.Type = "logout"
		User_list = del(User_list, c.data.User)
		c.data.UserList = User_list
		c.data.Content = c.data.User
		data_b, _ := json.Marshal(c.data)
		H.broadcast <- data_b
		H.register <- c
	}()
}

// writer 方法用于向 WebSocket 连接写入消息
func (c *connection) writer() {
	for message := range c.sc {
		// 将消息通过 WebSocket 发送给客户端
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	// 关闭 WebSocket 连接
	c.ws.Close()
}

// reader 方法用于从 WebSocket 连接读取消息reader
func (c *connection) reader() {
	defer func() {
		H.unregister <- c
		c.ws.Close()
	}()

	for {
		_, rawMessage, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			// 处理解析错误
			continue
		}

		messageType, ok := msg["type"].(string)
		if !ok {
			// 处理缺少类型字段的错误
			continue
		}

		content, ok := msg["content"].(string)
		if !ok {
			// 处理缺少内容字段的错误
			continue
		}

		switch messageType {
		case "message":
			// 在这里处理收到的消息，例如广播给其他连接
			H.broadcast <- []byte(fmt.Sprintf("siky: %s", content))
			H.message <- fmt.Sprintf("siky: %s", content)
			// 添加其他消息类型的处理
		}
	}
}

// del 方法用于从切片中删除指定的元素
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}

// 广播接收到的消息给其他客户端
func BroadcastMessage(message []byte) {
	for c := range H.connections {
		select {
		case c.sc <- message:
		default:
			delete(H.connections, c)
			close(c.sc)
		}
	}
}

func (h *Hub) HandleMessage(message string) {
	// 在这里处理收到的消息，这里简单地打印消息内容
	fmt.Println("Received message:", message)
}
