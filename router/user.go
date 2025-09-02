package router

import (
	"go-desk-service/api"

	"github.com/gin-gonic/gin"
)

type User struct{}

var ApiUser api.User

func (*User) Init(app *gin.Engine) {
	group := app.Group("/user")
	group.POST("/login", ApiUser.Login)
}
