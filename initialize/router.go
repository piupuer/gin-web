package initialize

import (
	"fmt"
	"gin-web/api"
	v1 "gin-web/api/v1"
	"gin-web/docs/swagger"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	hv1 "github.com/piupuer/go-helper/api/v1"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/query"
	"github.com/piupuer/go-helper/pkg/req"
	hr "github.com/piupuer/go-helper/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			middleware.WithRateRedis(global.Redis),
			middleware.WithRateMaxLimit(global.Conf.System.RateLimitMax),
		),
		middleware.Cors,
		middleware.SecurityHeader,
		middleware.RequestId(),
		middleware.Sign(
			middleware.WithSignLogger(global.Log),
			middleware.WithSignCheckScope(false),
			middleware.WithSignGetSignUser(func(c *gin.Context, appId string) ms.SignUser {
				return ms.SignUser{
					AppSecret: "gin-web",
					Status:    constant.One,
				}
			}),
			middleware.WithSignFindSkipPath(func(c *gin.Context) []string {
				return []string{
					fmt.Sprintf("%s/ping", global.Conf.System.UrlPrefix),
					fmt.Sprintf("%s/message/ws", global.Conf.System.Base),
					fmt.Sprintf("%s/upload/file", global.Conf.System.Base),
				}
			}),
		),
		middleware.AccessLog(
			middleware.WithAccessLogLogger(global.Log),
			middleware.WithAccessLogUrlPrefix(global.Conf.System.UrlPrefix),
		),
		middleware.Exception(),
		middleware.Transaction(
			middleware.WithTransactionDbNoTx(global.Mysql),
			middleware.WithTransactionTxCtxKey(constant.MiddlewareTransactionTxCtxKey),
		),
		middleware.OperationLog(
			middleware.WithOperationLogLogger(global.Log),
			middleware.WithOperationLogRedis(global.Redis),
			middleware.WithOperationLogCachePrefix(fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, constant.MiddlewareOperationLogApiCacheKey)),
			middleware.WithOperationLogUrlPrefix(global.Conf.System.UrlPrefix),
			middleware.WithOperationLogRealIpKey(global.Conf.System.AmapKey),
			middleware.WithOperationLogFindSkipPath(v1.OperationLogFindSkipPath),
			middleware.WithOperationLogSaveMaxCount(50),
			middleware.WithOperationLogGetCurrentUser(v1.GetCurrentUserAndRole),
			middleware.WithOperationLogSave(v1.OperationLogSave),
			middleware.WithOperationLogFindApi(v1.OperationLogFindApi),
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
		middleware.WithJwtLoginPwdCheck(func(c *gin.Context, r req.LoginCheck) (userId int64, err error) {
			s := cache_service.New(c)
			user, err := s.LoginCheck(r)
			return int64(user.Id), err
		}),
	}

	// set path prefix
	v1Group := apiGroup.Group(global.Conf.System.ApiVersion)
	// set swagger
	swagger.SwaggerInfo.Version = global.Conf.System.ApiVersion
	swagger.SwaggerInfo.BasePath = v1Group.BasePath()
	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DocExpansion("none"),
		),
	)

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
		hr.WithIdempotenceOps(
			middleware.WithIdempotenceCachePrefix(fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, constant.MiddlewareIdempotencePrefix)),
		),
		hr.WithV1Ops(
			hv1.WithCachePrefix(fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, "v1")),
			hv1.WithDbOps(
				query.WithMysqlDb(global.Mysql),
				query.WithMysqlFsmTransition(v1.LeaveTransition),
				query.WithMysqlCachePrefix(fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, constant.QueryCachePrefix)),
			),
			hv1.WithBinlogOps(
				query.WithRedisDatabase(global.Conf.Mysql.DSN.DBName),
				query.WithRedisNamingStrategy(global.Mysql.NamingStrategy),
			),
			hv1.WithGetCurrentUser(v1.GetCurrentUserAndRole),
			hv1.WithFindRoleKeywordByRoleIds(v1.RouterFindRoleKeywordByRoleIds),
			hv1.WithFindRoleByIds(v1.RouterFindRoleByIds),
			hv1.WithFindUserByIds(v1.RouterFindUserByIds),
			hv1.WithGetUserLoginStatus(v1.GetUserLoginStatus),
			hv1.WithFsmGetFsmSubmitterDetail(v1.GetLeaveFsmDetail),
			hv1.WithFsmUpdateFsmSubmitterDetail(v1.UpdateLeaveFsmDetail),
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
