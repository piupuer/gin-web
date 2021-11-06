package initialize

import (
	"gin-web/api"
	v1 "gin-web/api/v1"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	hv1 "github.com/piupuer/go-helper/api/v1"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/query"
	hr "github.com/piupuer/go-helper/router"
)

func Routers() *gin.Engine {
	// use custom router not default
	// r := gin.Default()
	r := gin.New()

	// replace default router
	gin.DebugPrintRouteFunc = middleware.PrintRouter(
		middleware.WithPrintRouterLogger(global.Log),
		middleware.WithPrintRouterCtx(ctx),
	)
	
	// set middleware
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
			middleware.WithOperationLogSaveMaxCount(50),
			middleware.WithOperationLogSave(v1.OperationLogSave),
			middleware.WithOperationLogFindApi(v1.OperationLogFindApi),
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

	// init default routers
	nr := hr.NewRouter(
		hr.WithLogger(global.Log),
		hr.WithRedis(global.Redis),
		hr.WithRedisBinlog(global.Conf.Redis.EnableBinlog),
		hr.WithGroup(v1Group),
		hr.WithJwt(true),
		hr.WithJwtOps(jwtOps...),
		hr.WithCasbin(true),
		hr.WithCasbinOps(
			middleware.WithCasbinEnforcer(global.CasbinEnforcer),
			middleware.WithCasbinGetCurrentUser(v1.GetCurrentUserAndRole),
		),
		hr.WithIdempotence(true),
		hr.WithV1Ops(
			hv1.WithDbOps(
				query.WithMysqlDb(global.Mysql),
			),
			hv1.WithBinlogOps(
				query.WithRedisCasbinEnforcer(global.CasbinEnforcer),
				query.WithRedisDatabase(global.Conf.Mysql.DSN.DBName),
				query.WithRedisNamingStrategy(global.Mysql.NamingStrategy),
			),
			hv1.WithGetCurrentUser(v1.GetCurrentUserAndRole),
			hv1.WithFindRoleKeywordByRoleIds(v1.RouterFindRoleKeywordByRoleIds),
			hv1.WithFindRoleByIds(v1.RouterFindRoleByIds),
			hv1.WithFindUserByIds(v1.RouterFindUserByIds),
			hv1.WithFsmTransition(v1.LeaveTransition),
			hv1.WithFsmGetFsmSubmitterDetail(v1.GetLeaveFsmDetail),
			hv1.WithUploadSaveDir(global.Conf.Upload.SaveDir),
			hv1.WithUploadSingleMaxSize(global.Conf.Upload.SingleMaxSize),
			hv1.WithUploadMergeConcurrentCount(global.Conf.Upload.MergeConcurrentCount),
			hv1.WithUploadMinio(global.Minio),
			hv1.WithUploadMinioBucket(global.Conf.Upload.Minio.Bucket),
		),
	)
	nr.Api()
	nr.Base()
	nr.Dict()
	nr.Fsm()
	nr.Machine()
	nr.Menu()
	nr.Message()
	nr.OperationLog()
	nr.Upload()

	// init custom routers
	router.InitLeaveRouter(nr)
	router.InitRoleRouter(nr)
	router.InitUserRouter(nr)

	global.Log.Info(ctx, "initialize router success")
	return r
}
