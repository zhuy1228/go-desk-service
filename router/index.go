package router

import "github.com/gin-gonic/gin"

var routerTest Test
var websocks Websocks
var user User

func Init(app *gin.Engine) {
	routerTest.Init(app)
	websocks.Init(app)
	user.Init(app)
}
