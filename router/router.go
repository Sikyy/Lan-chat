package router

import "github.com/gin-gonic/gin"

func Rrouter() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	// 设置静态文件服务
	r.Static("/static", "./static")

	//websocket连接相关
	// r.GET("/ws", service.Ws)
}
