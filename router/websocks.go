package router

import (
	"go-desk-service/api"

	"github.com/gin-gonic/gin"
)

type Websocks struct {
}

var apiWebsocks api.Websocks

func (*Websocks) Init(app *gin.Engine) {
	app.GET("/ws", apiWebsocks.Init)
}
