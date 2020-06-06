package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 请假路由
func InitLeaveRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("leave").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetLeaves)
		router.GET("/approval/list/:leaveId", v1.GetLeaveApprovalLogs)
		router.POST("/create", v1.CreateLeave)
		router.PATCH("/update/:leaveId", v1.UpdateLeaveById)
		router.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
	}
	return router
}
