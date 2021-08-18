package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(c, "[Exception]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack()))
			// 服务器异常
			resp := response.Resp{
				Code:      response.InternalServerError,
				Data:      map[string]interface{}{},
				Msg:       response.CustomError[response.InternalServerError],
				RequestId: c.GetString(global.RequestIdContextKey),
			}
			// 以json方式写入响应
			response.JSON(c, http.StatusOK, resp)
			c.Abort()
			return
		}
	}()
	c.Next()
}
