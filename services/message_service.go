package services

import (
	"log"

	"github.com/gorilla/websocket"
)

type MessageService struct{}

// websockt消息处理
func (*MessageService) WebSocketMessage(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取错误: %v (客户端: %s)", err, conn.RemoteAddr())
			break
		}
		log.Printf("收到来自 %s 的消息: %s", conn.RemoteAddr(), string(message))
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
