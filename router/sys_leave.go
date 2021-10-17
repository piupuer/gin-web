package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitLeaveRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/leave")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/leave")
	{
		router1.GET("/list", v1.GetLeaves)
		router1.GET("/approving/track/:leaveId", v1.FindLeaveFsmTrack)
		router2.POST("/create", v1.CreateLeave)
		router1.PATCH("/update/:leaveId", v1.UpdateLeaveById)
		router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
	}
	return r
}
