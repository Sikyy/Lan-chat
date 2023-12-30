package service

//ws.go
import (
	"Lan-chat/data"
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

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	for {
		var msg data.Data
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("读取消息时发生错误:", err)
			return
		}

		msg.Username = "你的用户名"
		msg.Type = "userMessage"
		H.b <- []byte(msg.Content)
		fmt.Printf("收到消息：%s: %s\n", msg.Username, msg.Content)
	}
}
