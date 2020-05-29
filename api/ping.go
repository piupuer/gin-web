package api

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

// 检查服务器是否通畅
func Ping(c *gin.Context) {
	response.SuccessWithData("pong")
}
