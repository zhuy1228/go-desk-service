package router

import "github.com/gin-gonic/gin"

var routerTest Test

func Init(app *gin.Engine) {
	routerTest.Init(app)
}
