package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/response"
	"runtime/debug"
)

// 处理异常
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("未知异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			// 响应服务器异常
			response.FailWithCode(c, response.EXCEPTION)
			c.Abort()
			return
		}
	}()
	c.Next()
}
