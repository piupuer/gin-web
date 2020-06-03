package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 工作流路由
func InitWorkflowRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("workflow").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetWorkflows)
		router.GET("/line/list", v1.GetWorkflowLines)
		router.POST("/create", v1.CreateWorkflow)
		router.PATCH("/update/:workflowId", v1.UpdateWorkflowById)
		router.DELETE("/delete/batch", v1.BatchDeleteWorkflowByIds)
		router.PATCH("/line/update", v1.UpdateWorkflowLineByNodes)
	}
	return router
}
