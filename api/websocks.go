package api

import (
	"context"
	"encoding/json"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	"go-desk-service/models"
	userpb "go-desk-service/proto/gen"
	"go-desk-service/services"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Websocks struct {
}

var DeviceList = make(map[int64]models.Device)

type Message struct {
	Type      string `form:"type" json:"type" uri:"type" xml:"type"`
	Data      any    `form:"data" json:"data" uri:"data" xml:"data"`
	Recipient int64  `form:"recipient" json:"recipient" uri:"recipient" xml:"recipient"`
	Sender    int64  `form:"sender" json:"sender" uri:"sender" xml:"sender"`
}

var Clients = make(map[int64]map[int64]*websocket.Conn) // 已连接的客户端

var ClientsMu sync.Mutex // 保护 clients 映射的互斥锁

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境中应严格限制来源
	},
}

func (w *Websocks) Init(ctx *gin.Context) {
	WebsocktFailed := libs.ErrorCode["WebsocktFailed"]
	// 判断当前连接是否合法
	key := http.CanonicalHeaderKey("sec-websocket-protocol")
	protos := ctx.Request.Header[key]
	if len(protos) <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": WebsocktFailed.Code, "data": WebsocktFailed.Data, "msg": WebsocktFailed.Msg})
		return
	}
	client := grpcClient.GetUserClient()
	ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tokenStr := protos[0]
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
	deviceService := services.InitDeviceService()
	deviceLoginInfo := deviceService.Login(tokenStr)
	DeviceList[deviceLoginInfo.ID] = deviceLoginInfo

	// 将保存连接状态
	ClientsMu.Lock()
	if _, ok := Clients[TokenStatus.UserId]; !ok {
		Clients[TokenStatus.UserId] = make(map[int64]*websocket.Conn)
	}
	Clients[TokenStatus.UserId][deviceLoginInfo.ID] = conn
	ClientsMu.Unlock()
	log.Printf("客户端已连接: %s", conn.RemoteAddr())
	w.MessageHandle(conn, deviceLoginInfo)
}

func (*Websocks) MessageHandle(conn *websocket.Conn, deviceLoginInfo models.Device) {
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("读取错误: %v (客户端: %s)", err, conn.RemoteAddr())
			log.Println("客户端断开", deviceLoginInfo.Token)
			deviceService := services.InitDeviceService()
			deviceService.Logout(deviceLoginInfo.Token)
			break
		}

		log.Printf("收到来自 %s 的消息: %s", conn.RemoteAddr(), string(message))
		var data Message
		// 结构体解析
		if err := json.Unmarshal(message, &data); err == nil && data.Type != "" {
			// 说明是结构体(JSON)
			log.Printf("收到结构体: %+v\n", data)
			var RecipientDevice models.Device
			if data.Recipient != 0 {
				RecipientDevice = DeviceList[data.Recipient]
			}
			log.Println("推送用户", RecipientDevice.UserId)
			log.Println("推送设备", RecipientDevice.ID)

			switch data.Type {
			case "message":
				if data.Recipient != 0 {
					// 将数据推送给对应设备
					userConn := Clients[RecipientDevice.UserId][RecipientDevice.ID]
					if userConn != nil {
						userConn.WriteJSON(&Message{
							Type:      "message",
							Data:      data.Data,
							Sender:    deviceLoginInfo.ID,
							Recipient: data.Recipient,
						})
					} else {
						log.Println("未获取到设备ID：", RecipientDevice.ID)
					}

				}
			}
		} else {
			log.Printf("收到普通消息: %s\n", string(message))
		}
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}
