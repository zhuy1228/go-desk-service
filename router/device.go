package router

import (
	"go-desk-service/api"
	"go-desk-service/middleware"

	"github.com/gin-gonic/gin"
)

type Device struct{}

var DeviceApi api.DeviceApi

func (*Device) Init(app *gin.Engine) {
	group := app.Group("/device", middleware.TokenAuth())
	group.POST("/seva", DeviceApi.SevaDeviceInfo)
	group.POST("/all", DeviceApi.GetAll)
}
