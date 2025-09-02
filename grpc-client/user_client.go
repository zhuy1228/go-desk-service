package grpcClient

import (
	"fmt"
	"go-desk-service/config"
	userpb "go-desk-service/proto/gen"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserClient userpb.UserServiceClient

func UserClientInit() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Println("读取配置失败")
	}

	conn, err := grpc.NewClient(
		appConfig.GrpcUrl, // gRPC服务器地址
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 使用非安全连接（测试用）
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)
	UserClient = client
}

func GetUserClient() userpb.UserServiceClient {
	return UserClient
}
