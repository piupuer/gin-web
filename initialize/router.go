package initialize

import (
	"gin-web/api"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func Routers() *gin.Engine {
	// use custom router not default
	// r := gin.Default()
	r := gin.New()

	// replace default router
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		global.Log.Debug(ctx, "[gin-route] %-6s %-40s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	r.Use(
		middleware.Rate(
			middleware.WithRateMaxLimit(global.Conf.System.RateLimitMax),
		),
		middleware.Cors,
		middleware.RequestId(),
		middleware.Exception(),
		middleware.AccessLog(
			middleware.WithAccessLogLogger(global.Log),
			middleware.WithAccessLogUrlPrefix(global.Conf.System.UrlPrefix),
		),
		middleware.OperationLog(
			middleware.WithOperationLogLogger(global.Log),
			middleware.WithOperationLogRedis(global.Redis),
			middleware.WithOperationLogUrlPrefix(global.Conf.System.UrlPrefix),
			middleware.WithOperationLogRealIpKey(global.Conf.System.AmapKey),
			middleware.WithOperationLogSkipPaths(global.Conf.Logs.OperationDisabledPathArr...),
			middleware.WithOperationLogSave(func(c *gin.Context, list []middleware.OperationRecord) {
			}),
			middleware.WithOperationLogFindApi(func(c *gin.Context) []middleware.OperationApi {
				return []middleware.OperationApi{}
			}),
		),
		middleware.Transaction(
			middleware.WithTransactionDbNoTx(global.Mysql),
			middleware.WithTransactionTxCtxKey(constant.MiddlewareTransactionTxCtxKey),
		),
	)

	apiGroup := r.Group(global.Conf.System.UrlPrefix)
	// ping
	apiGroup.GET("/ping", api.Ping)

	jwtOps := []func(*middleware.JwtOptions){
		middleware.WithJwtLogger(global.Log),
		middleware.WithJwtRealm(global.Conf.Jwt.Realm),
		middleware.WithJwtKey(global.Conf.Jwt.Key),
		middleware.WithJwtTimeout(global.Conf.Jwt.Timeout),
		middleware.WithJwtMaxRefresh(global.Conf.Jwt.MaxRefresh),
		middleware.WithJwtPrivateBytes(global.Conf.Jwt.RSAPrivateBytes),
		middleware.WithJwtLoginPwdCheck(func(c *gin.Context, username, password string) (userId int64, pass bool) {
			s := cache_service.New(c)
			user, err := s.LoginCheck(&models.SysUser{
				Username: username,
				Password: password,
			})
			if err != nil {
				return 0, false
			}
			return int64(user.Id), true
		}),
	}

	// set path prefix
	v1Group := apiGroup.Group(global.Conf.System.ApiVersion)
	// init routers
	router.InitPublicRouter(v1Group)
	router.InitBaseRouter(v1Group, jwtOps)
	router.InitUserRouter(v1Group, jwtOps)
	router.InitMenuRouter(v1Group, jwtOps)
	router.InitRoleRouter(v1Group, jwtOps)
	router.InitApiRouter(v1Group, jwtOps)
	router.InitFsmRouter(v1Group, jwtOps)
	router.InitLeaveRouter(v1Group, jwtOps)
	router.InitUploadRouter(v1Group, jwtOps)
	router.InitOperationLogRouter(v1Group, jwtOps)
	router.InitMessageRouter(v1Group, jwtOps)
	router.InitMachineRouter(v1Group, jwtOps)
	router.InitDictRouter(v1Group, jwtOps)

	global.Log.Info(ctx, "initialize router success")
	return r
}
