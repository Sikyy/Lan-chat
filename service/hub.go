package service

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// 初始化 Hub
var H = Hub{
	connections:    make(map[*connection]bool),   // 连接映射表，记录连接状态
	idToConnection: make(map[string]*connection), // id映射表，记录对应状态
	broadcast:      make(chan []byte),            // 消息广播通道
	register:       make(chan *connection),       // 新连接通道
	unregister:     make(chan *connection),       // 连接断开通道
	message:        make(chan string),            // 测试消息通道
}

var User_list = []string{} // 声明并初始化全局的 User_list

type Hub struct {
	connections    map[*connection]bool   // 连接映射表，记录连接状态
	idToConnection map[string]*connection // id映射表，记录对应连接
	broadcast      chan []byte            // 消息广播通道
	register       chan *connection       // 新连接通道
	unregister     chan *connection       // 连接断开通道
	message        chan string            // 测试消息通道
}

var message struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
	SenderIP string `json:"senderIP"`
}

// run 方法，用于启动 hub 的运行，监听各个通道的数据，进行处理
func (H *Hub) RunHub() {

	for {
		select {
		case c := <-H.register: // 有新连接加入
			H.connections[c] = true // 将连接状态设置为 true，表示连接已建立
			c.data.Ip = c.ws.RemoteAddr().String()
		case c := <-H.unregister: // 有连接断开
			if _, ok := H.connections[c]; ok {
				delete(H.connections, c) // 从连接映射表中删除连接
				close(c.sc)              // 关闭连接的消息通道
			}
		case data := <-H.broadcast: // 有消息广播
			fmt.Println("Broadcasting message:", string(data))
			//打印发送者的IP
			// 解析 JSON 数据
			err := json.Unmarshal([]byte(data), &message)
			if err != nil {
				fmt.Println("解析JSON时发生错误:", err)
				return
			}
			fmt.Println("senderIP:", message.SenderIP)
			for c := range H.connections { // 遍历连接映射表
				//打印遍历的连接IP
				fmt.Println(c.ws.RemoteAddr().String())
				// 给除发送方以外的所有连接发送消息
				// 获取发送方的地址
				sendip := message.SenderIP
				if sendip != c.ws.RemoteAddr().String() {
					select {
					case c.sc <- data: // 将消息发送到连接的消息通道
					default:
						delete(H.connections, c)
						close(c.sc)

					}
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

// 获取当前本机连接
func GetWsLocal(conn *websocket.Conn) string {
	return conn.RemoteAddr().String()
}
