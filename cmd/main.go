package main

import (
	"fmt"
	"go-desk-service/config"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	"go-desk-service/router"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Println("读取配置失败")
	}

	go handleShutdownSignals()
	libs.Connect()
	grpcClient.UserClientInit()
	// 创建 STUN 服务器实例
	stunServer := libs.NewSTUNServer(appConfig.StunPort)
	// 启动 STUN 服务器
	if err1 := stunServer.Start(); err1 != nil {
		log.Fatalf("Failed to start STUN server: %v", err1)
	}
	defer stunServer.Close()
	app := gin.Default()
	router.Init(app)
	app.Run(":" + strconv.Itoa(appConfig.Port))
}

// 处理关闭信号的函数
func handleShutdownSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 阻塞等待信号
	sig := <-sigChan
	log.Printf("接收到信号: %v, 系统关闭", sig)

	// 强制退出程序
	os.Exit(0)
}
