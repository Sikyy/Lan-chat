package service

import (
	"encoding/json"
)

// 初始化 Hub
var H = Hub{
	connections: make(map[*connection]bool), // 连接映射表，记录连接状态
	broadcast:   make(chan []byte),          // 消息广播通道
	register:    make(chan *connection),     // 新连接通道
	unregister:  make(chan *connection),     // 连接断开通道
	message:     make(chan string),          // 测试消息通道
}

var User_list = []string{} // 声明并初始化全局的 User_list

type Hub struct {
	connections map[*connection]bool // 连接映射表，记录连接状态
	broadcast   chan []byte          // 消息广播通道
	register    chan *connection     // 新连接通道
	unregister  chan *connection     // 连接断开通道
	message     chan string          // 测试消息通道
}

// run 方法，用于启动 hub 的运行，监听各个通道的数据，进行处理
func (H *Hub) RunHub() {

	for {
		select {
		case c := <-H.register: // 有新连接加入
			H.connections[c] = true                // 将连接状态设置为 true，表示连接已建立
			c.data.Ip = c.ws.RemoteAddr().String() // 获取连接的远程地址，并设置到连接数据中的 Ip 字段
			c.data.Type = "handshake"              // 设置连接数据的类型为 "handshake"
			c.data.UserList = User_list            // 获取用户列表，并设置到连接数据中的 UserList 字段
			data_b, _ := json.Marshal(c.data)      // 将连接数据转换为 JSON 格式
			c.sc <- data_b                         // 将 JSON 数据发送到连接的消息通道
		case c := <-H.unregister: // 有连接断开
			if _, ok := H.connections[c]; ok {
				delete(H.connections, c) // 从连接映射表中删除连接
				close(c.sc)              // 关闭连接的消息通道
			}
		case data := <-H.broadcast: // 有消息广播
			for c := range H.connections { // 遍历连接映射表
				select {
				case c.sc <- data: // 将消息发送到连接的消息通道
				default:
					delete(H.connections, c)
					close(c.sc)
				}
			}
		case message := <-H.message:
			// 处理收到的消息
			H.HandleMessage(message)

		}
	}
}

// 启动 Hub 运行的函数
func RunHub() {
	go H.RunHub()
}
