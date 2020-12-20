package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 消息机器路由
func InitMachineRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/machine").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/shell/ws", v1.MachineShellWs)
		router.GET("/list", v1.GetMachines)
		router.
			Use(middleware.Idempotence).
			POST("/create", v1.CreateMachine)
		router.PATCH("/update/:machineId", v1.UpdateMachineById)
		router.PATCH("/connect/:machineId", v1.ConnectMachineById)
		router.DELETE("/delete/batch", v1.BatchDeleteMachineByIds)
	}
	return router
}
