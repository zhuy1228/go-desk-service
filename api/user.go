package api

import (
	"context"
	"fmt"
	"go-desk-service/config"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	"go-desk-service/models"
	userpb "go-desk-service/proto/gen"
	"go-desk-service/services"
	"go-desk-service/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
}

type LoginRequestParams struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	DeviceId string `form:"device_id" json:"device_id" uri:"device_id" xml:"device_id" binding:"required"`
}

type RegisterRequestParams struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	Email    string `form:"email" json:"email" uri:"email" xml:"email" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" uri:"nickname" xml:"nickname" binding:"required"`
	Phone    string `form:"phone" json:"phone" uri:"phone" xml:"phone" binding:"required"`
}

func (*User) Login(ctx *gin.Context) {
	var res *libs.ErrorInfo
	var paramsJson LoginRequestParams
	if err := ctx.ShouldBindJSON(&paramsJson); err != nil {
		res = libs.ErrorCode["ParamsError"]
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		return
	}
	// 获取当前的grpc连接
	client := grpcClient.GetUserClient()
	ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	loginResp, err := client.Login(ctx1, &userpb.LoginRequest{
		Username: paramsJson.Username,
		Password: paramsJson.Password,
	})
	if err != nil {
		res = libs.ErrorCode["LoginError"]
		// log.Fatalf("登录失败: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		return
	}
	ProfileService := services.InitProfileService()
	log.Println(loginResp.UserId)
	ProfileService.UpdateAccount(loginResp.UserId, &models.Profile{
		LoginStatus: 1,
	})
	// 登录成功之后查看设备是否是当前设备
	DeviceService := services.InitDeviceService()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": loginResp,
	})
}

func (*User) Register(ctx *gin.Context) {
	var res *libs.ErrorInfo
	var paramsJson RegisterRequestParams
	if err := ctx.ShouldBindJSON(&paramsJson); err != nil {
		res = libs.ErrorCode["ParamsError"]
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		return
	}
	client := grpcClient.GetUserClient()
	ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	registerResp, err := client.Register(ctx1, &userpb.RegisterRequest{
		Username: paramsJson.Username,
		Password: paramsJson.Password,
		Email:    paramsJson.Email,
		Nickname: paramsJson.Nickname,
		Phone:    paramsJson.Phone,
	})
	if err != nil {
		res = libs.ErrorCode["RegistrationFailed"]

		ctx.JSON(http.StatusOK, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		// log.Fatalf("注册失败: %v", err)
		return
	}
	// 信息映射到当前表
	appConfig, _ := config.LoadConfig()
	Snowflake, _ := utils.NewSnowflake(appConfig.WorkerID, appConfig.DatacenterID)
	ProfileId := Snowflake.NextID()
	ProfileService := services.InitProfileService()
	ProfileService.Create(&models.Profile{
		ID:          ProfileId,
		UserId:      registerResp.UserId,
		Nickname:    paramsJson.Nickname,
		Email:       paramsJson.Email,
		Phone:       paramsJson.Phone,
		Status:      1,
		LoginStatus: 0,
	})
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": registerResp,
	})
}
