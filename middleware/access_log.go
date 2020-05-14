package middleware

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"time"
)

// 访问日志
func AccessLog(c *gin.Context) {
	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	execTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := c.Request.Method

	// 请求路由
	reqUri := c.Request.RequestURI

	// 状态码
	statusCode := c.Writer.Status()

	// 请求IP
	clientIP := c.ClientIP()

	global.Log.Info(
		fmt.Sprintf(
			"%s %s %d %s %s",
			reqMethod,
			reqUri,
			statusCode,
			execTime.String(),
			clientIP,
		),
	)
}
