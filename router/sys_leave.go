package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 请假路由
func InitLeaveRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/leave")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/leave")
	{
		router1.GET("/list", v1.GetLeaves)
		router1.GET("/approval/list/:leaveId", v1.GetLeaveApprovalLogs)
		router2.POST("/create", v1.CreateLeave)
		router1.PATCH("/update/:leaveId", v1.UpdateLeaveById)
		router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
	}
	return r
}
