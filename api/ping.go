package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/response"
)

// 检查服务器是否通畅
func Ping(c *gin.Context) {
	a := 1000
	i := 1000 - a
	fmt.Println(1 / i)
	response.Success(c)
}
