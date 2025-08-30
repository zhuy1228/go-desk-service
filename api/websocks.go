package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Websocks struct {
}

var clients = make(map[*websocket.Conn]bool) // 已连接的客户端
var clientsMu sync.Mutex                     // 保护 clients 映射的互斥锁

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境中应严格限制来源
	},
}

func (*Websocks) Init(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法建立 WebSocket 连接"})
		return
	}
	defer conn.Close()
	// 将保存连接状态
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Printf("客户端已连接: %s", conn.RemoteAddr())
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
