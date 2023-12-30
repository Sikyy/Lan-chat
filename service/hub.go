package service

import (
	"Lan-chat/data"
	"encoding/json"
	"fmt"
	"time"
)

// 初始化 Hub
var H = Hub{
	c: make(map[*connection]bool), // 连接映射表，记录连接状态
	b: make(chan []byte),          // 消息广播通道
	r: make(chan *connection),     // 新连接通道
	u: make(chan *connection),     // 连接断开通道
}

var User_list = []string{} // 声明并初始化全局的 User_list

type Hub struct {
	c map[*connection]bool // 连接映射表，记录连接状态
	b chan []byte          // 消息广播通道
	r chan *connection     // 新连接通道
	u chan *connection     // 连接断开通道
}

// run 方法，用于启动 hub 的运行，监听各个通道的数据
func (H *Hub) RunHub() {

	// 定时器，每五秒触发一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	printConnectionCount := func() {
		count := len(H.c)
		fmt.Printf("当前连接数：%d\n", count)
	}

	for {
		select {
		case c := <-H.r: // 有新连接加入
			H.c[c] = true                          // 将连接状态设置为 true，表示连接已建立
			c.data.Ip = c.ws.RemoteAddr().String() // 获取连接的远程地址，并设置到连接数据中的 Ip 字段
			c.data.Type = "handshake"              // 设置连接数据的类型为 "handshake"
			c.data.UserList = User_list            // 将全局的 User_list 设置到连接数据中的 UserList 字段
			data_b, _ := json.Marshal(c.data)      // 将连接数据转换为 JSON 格式
			c.sc <- data_b                         // 将 JSON 数据发送到连接的消息通道
			printConnectionCount()                 // 打印当前连接数
		case c := <-H.u: // 有连接断开
			if _, ok := H.c[c]; ok {
				delete(H.c, c)
				close(c.sc)
				printConnectionCount()
			}
		case data := <-H.b: // 有消息广播
			for c := range H.c {
				select {
				case c.sc <- data:
				default:
					delete(H.c, c)
					close(c.sc)
					printConnectionCount()
				}
			}
		case <-ticker.C:
			// 每五秒触发一次的统计和打印操作
			printConnectionCount()
		}
	}
}

// 广播连接数信息给所有连接
func (h *Hub) BroadcastConnectionCount() {
	count := len(h.c)
	data := data.Data{Type: "connectionCount", Content: fmt.Sprintf("%d", count)}

	for c := range h.c {
		jsonMsg, err := data.ToJSON()
		if err != nil {
			fmt.Println("Error converting message to JSON:", err)
			continue
		}

		select {
		case c.sc <- jsonMsg:
		default:
			fmt.Println("Failed to send message to client")
		}
	}
}

// 启动 Hub 运行的函数
func RunHub() {
	go H.RunHub()
}
