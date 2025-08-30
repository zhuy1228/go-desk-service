package router

import "github.com/gin-gonic/gin"

var routerTest Test
var websocks Websocks

func Init(app *gin.Engine) {
	routerTest.Init(app)
	websocks.Init(app)
}
