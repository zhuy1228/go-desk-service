package api

import (
	"fmt"
	"go-desk-service/libs"
	"go-desk-service/models"
	"go-desk-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceApi struct{}

type DeviceInfo struct {
	Hostname        string `form:"hostname" json:"hostname" uri:"hostname" xml:"hostname" binding:"required"`
	Platform        string `form:"platform" json:"platform" uri:"platform" xml:"platform" binding:"required"`
	PlatformVersion string `form:"platform_version" json:"platform_version" uri:"platform_version" xml:"platform_version" binding:"required"`
	Mac             string `form:"mac" json:"mac" uri:"mac" xml:"mac" binding:"required"`
	CPU             string `form:"cpu" json:"cpu" uri:"cpu" xml:"cpu" binding:"required"`
	Mem             string `form:"mem" json:"mem" uri:"mem" xml:"mem" binding:"required"`
	Disk            string `form:"disk" json:"disk" uri:"disk" xml:"disk" binding:"required"`
}

func (*DeviceApi) SevaDeviceInfo(ctx *gin.Context) {
	var res *libs.ErrorInfo
	var paramsJson DeviceInfo
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
	// 获取当前的用户ID跟token信息
	token, _ := ctx.Get("token")
	IP := ctx.ClientIP()
	// 保存当前的设备
	DeviceService := services.InitDeviceService()
	DeviceService.UpdateByToken(token.(string), &models.Device{
		Hostname:        paramsJson.Hostname,
		Platform:        paramsJson.Platform,
		PlatformVersion: paramsJson.PlatformVersion,
		Mac:             paramsJson.Mac,
		Cpu:             paramsJson.CPU,
		Mem:             paramsJson.Mem,
		Disk:            paramsJson.Disk,
		Ip:              IP,
	})
	res = libs.ErrorCode["SevaSuccessful"]
	ctx.JSON(http.StatusOK, gin.H{
		"code": res.Code,
		"msg":  res.Msg,
		"data": res.Data,
	})
}

func (*DeviceApi) GetAll(ctx *gin.Context) {
	userId, _ := ctx.Get("userID")
	DeviceService := services.InitDeviceService()
	devices := DeviceService.GetAll(userId.(int64))
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": devices,
	})
}
