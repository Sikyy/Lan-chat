package service

import "fmt"

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
			H.connections[c] = true // 将连接状态设置为 true，表示连接已建立

		case c := <-H.unregister: // 有连接断开
			if _, ok := H.connections[c]; ok {
				delete(H.connections, c) // 从连接映射表中删除连接
				close(c.sc)              // 关闭连接的消息通道
			}
		case data := <-H.broadcast: // 有消息广播
			fmt.Println("Broadcasting message:", string(data))
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

// GetConnectionCount 返回当前连接数量
func (h *Hub) GetConnectionCount() int {
	return len(h.connections)
}
