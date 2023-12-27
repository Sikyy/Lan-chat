package router

import (
	"Lan-chat/service"

	"github.com/gin-gonic/gin"
)

func Rrouter() {
	//启动hub，监听各个通道的数据
	service.RunHub()

	//创建路由
	r := gin.Default()

	// 指定 HTML 文件的路径
	r.LoadHTMLGlob("templates/*.html")
	// 静态资源处理
	r.Static("static", "./static")

	//websocket连接相关
	// r.GET("/ws", service.HandleConnection())
}
