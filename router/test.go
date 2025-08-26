package router

import (
	"go-desk-service/api"

	"github.com/gin-gonic/gin"
)

type Test struct{}

var apiTest api.Test

func (*Test) Init(app *gin.Engine) {
	test := app.Group("/test")
	test.GET("/status", apiTest.Status)
}
