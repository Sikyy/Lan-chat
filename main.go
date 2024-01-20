package main

import (
	"Lan-chat/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// 启动 hub 的运行
	service.RunHub()

	router := gin.Default()

	// 指定 HTML 文件的路径
	router.LoadHTMLGlob("templates/*")

	// 静态资源处理
	router.Static("static", "./static")

	// 首页
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	router.GET("/connectionCount", func(c *gin.Context) {
		count := service.H.GetConnectionCount()
		c.JSON(http.StatusOK, gin.H{"connectionCount": count})
	})

	// 升级到 WebSocket
	router.GET("/ws", func(c *gin.Context) {
		conn, err := service.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("升级到 WebSocket 时出错:", err)
			return
		}

		// 用协程处理连接，传递连接信息
		go service.HandleConnection(conn)
	})

	router.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// 启动服务
	router.Run("0.0.0.0:8880")

}
