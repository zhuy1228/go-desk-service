package main

import (
	"fmt"
	"go-desk-service/config"
	"go-desk-service/router"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Println("读取配置失败")
	}

	app := gin.Default()
	router.Init(app)
	app.Run(":" + strconv.Itoa(appConfig.Port))
}
