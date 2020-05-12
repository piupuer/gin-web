package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/response"
	"net/http"
	"runtime/debug"
)

// 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		// 获取事务对象
		tx := global.GetTx(c)
		if err := recover(); err != nil {
			var resp response.Resp
			// 判断是否自定义响应结果
			resp, ok := err.(response.Resp)
			if ok {
				if global.Conf.System.Transaction {
					if resp.Code == response.Ok {
						// 有效的请求, 提交事务
						tx.Commit()
					} else {
						// 回滚事务
						tx.Rollback()
					}
				}
			} else {
				// 将其他异常写入日志
				global.Log.Error(fmt.Sprintf("未知异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
				// 服务器异常
				resp.Code = response.InternalServerError
				resp.Msg = response.CustomError[response.InternalServerError]
				resp.Data = map[string]interface{}{}
				if global.Conf.System.Transaction {
					// 回滚事务
					tx.Rollback()
				}
			}
			// 以json方式写入响应
			c.JSON(http.StatusOK, resp)
		} else {
			if global.Conf.System.Transaction {
				// 没有异常, 提交事务
				tx.Commit()
			}
		}
		// 结束请求, 避免二次调用
		c.Abort()
	}()
	if global.Conf.System.Transaction {
		fmt.Println("准备开启事务")
		// 开启事务, 写入当前请求
		tx := global.Mysql.Begin()
		fmt.Println("开启事务成功", tx.Error)
		c.Set("tx", tx)
	}
	// 处理请求
	c.Next()
}
