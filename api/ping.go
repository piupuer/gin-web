package api

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/response"
)

// 检查服务器是否通畅
func Ping(c *gin.Context) {
	response.Success(c)
}
