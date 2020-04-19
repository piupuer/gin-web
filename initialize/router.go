package initialize

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/api"
	"go-shipment-api/pkg/global"
	"go-shipment-api/router"
)

// 初始化总路由
func Routers() *gin.Engine {
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()

	// ping
	r.GET("/ping", api.Ping)

	// 方便统一添加路由前缀
	group := r.Group("")
	router.InitPublicRouter(group) // 注册公共路由

	global.Log.Debug("初始化路由完成")
	return r
}
