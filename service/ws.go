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
	H.r <- c

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
		H.b <- data_b
		H.r <- c
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

// reader 方法用于从 WebSocket 连接读取消息
func (c *connection) reader() {
	for {
		// 从 WebSocket 连接中读取消息
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			// 发生错误时，将 connection 加入到 hub 的注册通道中，以便清理资源
			H.r <- c
			break
		}

		// 解析消息并根据消息类型执行相应的操作
		json.Unmarshal(message, &c.data)
		switch c.data.Type {
		case "login":
			// 处理用户登录消息
			c.data.User = c.data.Content
			c.data.From = c.data.User
			User_list = append(User_list, c.data.User)
			c.data.UserList = User_list
			data_b, _ := json.Marshal(c.data)
			H.b <- data_b
		case "user":
			// 处理用户发送的普通消息
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			H.b <- data_b
		case "logout":
			// 处理用户注销消息
			c.data.Type = "logout"
			User_list = del(User_list, c.data.User)
			data_b, _ := json.Marshal(c.data)
			H.b <- data_b
			H.r <- c
		default:
			// 处理其他未知类型的消息
			fmt.Print("========default================")
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
