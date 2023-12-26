package main

import (
	"Lan-chat/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("index.html")

	// 静态资源处理
	router.Static("static", "./static")

	// 首页
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 升级到 WebSocket
	router.GET("/ws", func(c *gin.Context) {
		conn, err := service.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("升级到 WebSocket 时出错:", err)
			return
		}

		// 将客户端保存到 clients 中
		service.Clients[conn] = struct{}{}

		// 用协程处理连接
		go service.HandleConnection(conn)
	})

	// 启动服务
	router.Run("0.0.0.0:8880")
}
