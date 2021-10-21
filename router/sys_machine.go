package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitMachineRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/machine")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/machine")
	{
		router1.GET("/shell/ws", v1.MachineShellWs)
		router1.GET("/list", v1.FindMachine)
		router2.POST("/create", v1.CreateMachine)
		router1.PATCH("/update/:id", v1.UpdateMachineById)
		router1.PATCH("/connect/:id", v1.ConnectMachineById)
		router1.DELETE("/delete/batch", v1.BatchDeleteMachineByIds)
	}
	return r
}
