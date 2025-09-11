package api

import (
	"context"
	"fmt"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	userpb "go-desk-service/proto/gen"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		log.Fatalf("登录失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		log.Fatalf("注册失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": res.Code,
			"msg":  res.Msg,
			"data": res.Data,
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": registerResp,
	})
}
