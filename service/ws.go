package service

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// upgrader 用于升级 HTTP 连接到 WebSocket 连接
var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// 保存所有客户端
	Clients = make(map[*websocket.Conn]struct{})
)

// 处理连接
func HandleConnection(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage() //返回三个值，分别是消息类型，消息内容，错误信息
		if err != nil {
			fmt.Println("读取消息时出错:", err)
			delete(Clients, conn) //删除客户端
			return
		}

		message := string(p)
		// 打印消息
		fmt.Printf("收到来自 %s 的消息: %s\n", conn.RemoteAddr(), message)

		BroadcastMessage(conn, messageType, p)
	}
}

// 广播消息
func BroadcastMessage(sender *websocket.Conn, messageType int, message []byte) {
	for client := range Clients { //遍历所有客户端
		if client != sender {
			err := client.WriteMessage(messageType, message) //向客户端发送消息
			if err != nil {
				fmt.Println("向客户端广播消息时出错:", err)
			}
		}
	}
}
