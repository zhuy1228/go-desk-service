package router

import (
	"go-desk-service/api"
	"go-desk-service/middleware"

	"github.com/gin-gonic/gin"
)

type Test struct{}

var apiTest api.Test

func (*Test) Init(app *gin.Engine) {
	test := app.Group("/test", middleware.TokenAuth())
	test.GET("/status", apiTest.Status)
}
