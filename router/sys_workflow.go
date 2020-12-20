package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 工作流路由
func InitWorkflowRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/workflow").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetWorkflows)
		router.GET("/line/list", v1.GetWorkflowLines)
		router.GET("/approving/list", v1.GetWorkflowApprovings)
		router.
			Use(middleware.Idempotence).
			POST("/create", v1.CreateWorkflow)
		router.PATCH("/update/:workflowId", v1.UpdateWorkflowById)
		router.PATCH("/log/approval", v1.UpdateWorkflowLogApproval)
		router.DELETE("/delete/batch", v1.BatchDeleteWorkflowByIds)
		router.PATCH("/line/update", v1.UpdateWorkflowLineIncremental)
	}
	return router
}
