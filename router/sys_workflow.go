package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 工作流路由
func InitWorkflowRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/workflow")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/workflow")
	{
		router1.GET("/list", v1.GetWorkflows)
		router1.GET("/line/list", v1.GetWorkflowLines)
		router1.GET("/approving/list", v1.GetWorkflowApprovings)
		router2.POST("/create", v1.CreateWorkflow)
		router1.PATCH("/update/:workflowId", v1.UpdateWorkflowById)
		router1.PATCH("/log/approval", v1.UpdateWorkflowLogApproval)
		router1.DELETE("/delete/batch", v1.BatchDeleteWorkflowByIds)
		router1.PATCH("/line/update", v1.UpdateWorkflowLineIncremental)
	}
	return r
}
