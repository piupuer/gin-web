package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 消息机器路由
func InitMachineRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/machine")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/machine")
	{
		router1.GET("/shell/ws", v1.MachineShellWs)
		router1.GET("/list", v1.GetMachines)
		router2.POST("/create", v1.CreateMachine)
		router1.PATCH("/update/:machineId", v1.UpdateMachineById)
		router1.PATCH("/connect/:machineId", v1.ConnectMachineById)
		router1.DELETE("/delete/batch", v1.BatchDeleteMachineByIds)
	}
	return r
}
