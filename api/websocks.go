package api

import (
	"context"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	userpb "go-desk-service/proto/gen"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Websocks struct {
}

type Message struct {
}

var Clients = make(map[int64]map[string]*websocket.Conn) // 已连接的客户端

var ClientsMu sync.Mutex // 保护 clients 映射的互斥锁

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境中应严格限制来源
	},
}

func (*Websocks) Init(ctx *gin.Context) {
	WebsocktFailed := libs.ErrorCode["WebsocktFailed"]
	// 判断当前连接是否合法
	tokenStr := ctx.Request.Header.Get("Sec-WebSocket-Protocol")
	if tokenStr == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": WebsocktFailed.Code, "data": WebsocktFailed.Data, "msg": WebsocktFailed.Msg})
		return
	}
	client := grpcClient.GetUserClient()
	ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	TokenStatus, err1 := client.ValidateToken(ctx1, &userpb.ValidateTokenRequest{
		AccessToken: tokenStr,
	})
	// 验证当前token
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": WebsocktFailed.Code, "data": WebsocktFailed.Data, "msg": WebsocktFailed.Msg})
		return
	}
	if !TokenStatus.IsValid {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": WebsocktFailed.Code, "data": WebsocktFailed.Data, "msg": WebsocktFailed.Msg})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"code": WebsocktFailed.Code, "data": WebsocktFailed.Data, "msg": WebsocktFailed.Msg})
		return
	}
	defer conn.Close()

	// 将保存连接状态
	ClientsMu.Lock()
	Clients[TokenStatus.UserId] = conn
	ClientsMu.Unlock()
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
