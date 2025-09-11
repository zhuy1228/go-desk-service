package middleware

import (
	"context"
	grpcClient "go-desk-service/grpc-client"
	"go-desk-service/libs"
	userpb "go-desk-service/proto/gen"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用当前grpc的token校验接口
		client := grpcClient.GetUserClient()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tokenExpired := libs.ErrorCode["TokenExpired"]
			c.JSON(http.StatusUnauthorized, gin.H{"code": tokenExpired.Code, "data": tokenExpired.Data, "msg": tokenExpired.Msg})
			c.Abort()
			return
		}
		TokenStatus, err := client.ValidateToken(ctx, &userpb.ValidateTokenRequest{
			AccessToken: authHeader,
		})
		if err != nil {
			tokenExpired := libs.ErrorCode["TokenExpired"]
			c.JSON(http.StatusUnauthorized, gin.H{"code": tokenExpired.Code, "data": tokenExpired.Data, "msg": tokenExpired.Msg})
			c.Abort()
			return
		}
		if !TokenStatus.IsValid {
			tokenExpired := libs.ErrorCode["TokenExpired"]
			c.JSON(http.StatusUnauthorized, gin.H{"code": tokenExpired.Code, "data": tokenExpired.Data, "msg": tokenExpired.Msg})
			c.Abort()
			return
		}
		c.Set("userID", TokenStatus.UserId)
		c.Next()
	}
}
