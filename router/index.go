package router

import "github.com/gin-gonic/gin"

var routerTest Test
var websocks Websocks
var user User
var device Device

func Init(app *gin.Engine) {
	routerTest.Init(app)
	websocks.Init(app)
	user.Init(app)
	device.Init(app)
}
