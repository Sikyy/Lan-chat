package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	// hub实例
	H = Hub{
		c: make(map[*connection]bool), //表示连接是否存在
		b: make(chan []byte),          //用于接收要广播的消息数据
		r: make(chan *connection),     //用于接收要移除的连接
		u: make(chan *connection),     //用于接收新连接
	}
)

type Hub struct {
	c  map[*connection]bool //连接池--当前连接的集合（数）
	b  chan []byte          //用于接收要广播的消息数据
	r  chan *connection     //用于接收要移除的连接
	u  chan *connection     //用于接收新连接
	mu sync.Mutex           // 互斥锁
}

// 运行hub，监听各个通道的数据
func (h *Hub) Run() {
	// 启动定时器，每隔5秒执行一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case c := <-h.r: // 有新连接加入
			fmt.Println("新连接加入:", c.ws.RemoteAddr().String())
			h.c[c] = true
			c.data.Type = "handshake"
			c.data.Content = fmt.Sprintf("新用户已连接：%s", c.ws.RemoteAddr().String())
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
			h.broadcastConnectionCount()
		case c := <-h.u: // 有连接断开
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
				fmt.Println("连接已关闭:", c.ws.RemoteAddr().String())
			}
			h.broadcastConnectionCount()
		case data := <-h.b: // 有消息广播
			fmt.Println("广播消息:", string(data))
			for c := range h.c {
				select {
				case c.sc <- data:
				default:
					delete(h.c, c)
					close(c.sc)
					fmt.Println("无法向客户端发送消息。连接已关闭:", c.ws.RemoteAddr().String())
					h.broadcastConnectionCount()
				}
			}
		case <-ticker.C:
			// 定时器触发时打印当前连接数
			h.PrintConnectionCount()
		}
	}
}

// 打印当前连接数
func (h *Hub) PrintConnectionCount() {
	h.mu.Lock()
	defer h.mu.Unlock()

	count := len(h.c)
	fmt.Printf("当前连接数: %d\n", count)
}

// 用于路由器调用
func RunHub() {
	go H.Run()
}
