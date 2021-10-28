package api

import (
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 检查服务器是否通畅
func Ping(c *gin.Context) {
	resp.SuccessWithData("pong")
}
