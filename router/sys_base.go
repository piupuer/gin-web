package router

import (
	v1 "gin-web/api/v1"
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitBaseRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router := r.Group("/base")
	{
		router.POST("/login", middleware.JwtLogin(jwtOptions...))
		router.POST("/logout", middleware.JwtLogout(jwtOptions...))
		router.POST("/refreshToken", middleware.JwtRefresh(jwtOptions...))
		// need login
		router.
			Use(middleware.Jwt(jwtOptions...)).
			Use(GetCasbinMiddleware()).
			GET("/idempotenceToken", middleware.GetIdempotenceToken(
				middleware.WithIdempotenceRedis(global.Redis),
			))
	}
	return r
}

// get casbin middleware router
func GetCasbinRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions), path string) gin.IRoutes {
	return r.Group(path).
		Use(middleware.Jwt(jwtOptions...)).
		Use(GetCasbinMiddleware())
}

// get casbin and idempotence middleware router
func GetCasbinAndIdempotenceRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions), path string) gin.IRoutes {
	return GetCasbinRouter(r, jwtOptions, path).
		Use(
			middleware.Idempotence(GetIdempotenceMiddlewareOps()...),
		)
}

func GetCasbinMiddleware() gin.HandlerFunc {
	return middleware.Casbin(
		middleware.WithCasbinEnforcer(global.CasbinEnforcer),
		middleware.WithCasbinRoleKey(func(c *gin.Context) string {
			user := v1.GetCurrentUser(c)
			return user.Role.Keyword
		}),
	)
}

func GetIdempotenceMiddlewareOps() []func(*middleware.IdempotenceOptions) {
	return []func(*middleware.IdempotenceOptions){
		middleware.WithIdempotenceRedis(global.Redis),
		middleware.WithIdempotencePrefix(constant.MiddlewareIdempotencePrefix),
	}
}
