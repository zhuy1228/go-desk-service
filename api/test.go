package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Test struct {
}

func (*Test) Status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": "",
	})
}
