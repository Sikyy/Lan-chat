package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader 用于升级 HTTP 连接到 WebSocket 连接
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// 保存所有客户端
	clients = make(map[*websocket.Conn]struct{})
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("index.html")

	// 静态资源处理
	router.Static("/css", "./css")
	router.Static("/images", "./images")
	router.Static("/js", "./js")

	// 首页
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 升级到 WebSocket
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("升级到 WebSocket 时出错:", err)
			return
		}

		// 将客户端保存到 clients 中
		clients[conn] = struct{}{}

		// 用协程处理连接
		go handleConnection(conn)
	})

	// 启动服务
	router.Run("0.0.0.0:8880")
}

// 处理连接
func handleConnection(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage() //返回三个值，分别是消息类型，消息内容，错误信息
		if err != nil {
			fmt.Println("读取消息时出错:", err)
			delete(clients, conn) //删除客户端
			return
		}

		message := string(p)

		fmt.Printf("收到来自 %s 的消息: %s\n", conn.RemoteAddr(), message)

		broadcastMessage(conn, messageType, p)
	}
}

// 广播消息
func broadcastMessage(sender *websocket.Conn, messageType int, message []byte) {
	for client := range clients { //遍历所有客户端
		if client != sender {
			err := client.WriteMessage(messageType, message) //向客户端发送消息
			if err != nil {
				fmt.Println("向客户端广播消息时出错:", err)
			}
		}
	}
}
