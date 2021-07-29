package initialize

import (
	"fmt"
	"gin-web/api"
	"gin-web/middleware"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/gin-gonic/gin"
)

// 初始化总路由
func Routers() *gin.Engine {
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()

	// 添加速率访问中间件
	r.Use(middleware.RateLimiter())
	// 添加访问记录
	r.Use(middleware.AccessLog)
	// 添加操作日志
	r.Use(middleware.OperationLog)
	// 添加全局异常处理中间件
	r.Use(middleware.Exception)
	// 添加全局事务处理中间件
	r.Use(middleware.Transaction)
	// 添加跨域中间件, 让请求支持跨域
	r.Use(middleware.Cors())
	global.Log.Info("请求已支持跨域")

	// 初始化jwt auth中间件
	authMiddleware, err := middleware.InitAuth()

	if err != nil {
		panic(fmt.Sprintf("初始化jwt auth中间件失败: %v", err))
	}
	global.Log.Info("初始化jwt auth中间件完成")

	apiGroup := r.Group(global.Conf.System.UrlPathPrefix)
	// ping
	apiGroup.GET("/ping", api.Ping)

	// 方便统一添加路由前缀
	v1Group := apiGroup.Group(global.Conf.System.ApiVersion)
	router.InitPublicRouter(v1Group)                       // 注册公共路由
	router.InitBaseRouter(v1Group, authMiddleware)         // 注册基础路由
	router.InitUserRouter(v1Group, authMiddleware)         // 注册用户路由
	router.InitMenuRouter(v1Group, authMiddleware)         // 注册菜单路由
	router.InitRoleRouter(v1Group, authMiddleware)         // 注册角色路由
	router.InitApiRouter(v1Group, authMiddleware)          // 注册接口路由
	router.InitUploadRouter(v1Group, authMiddleware)       // 注册文件上传路由
	router.InitOperationLogRouter(v1Group, authMiddleware) // 注册操作日志路由
	router.InitMessageRouter(v1Group, authMiddleware)      // 注册消息中心路由
	router.InitMachineRouter(v1Group, authMiddleware)      // 注册机器路由
	router.InitDictRouter(v1Group, authMiddleware)         // 注册数据字典路由

	global.Log.Info("初始化路由完成")
	return r
}
