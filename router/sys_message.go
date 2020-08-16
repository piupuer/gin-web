package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 消息中心路由
func InitMessageRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("message").Use(authMiddleware.MiddlewareFunc())
	// .Use(middleware.CasbinMiddleware)
	{
		router.GET("/all", v1.GetAllMessages)
		router.GET("/unRead/count", v1.GetUnReadMessageCount)
		router.PATCH("/read/batch", v1.BatchUpdateMessageRead)
		router.PATCH("/deleted/batch", v1.BatchUpdateMessageDeleted)
		router.PATCH("/read/all", v1.UpdateAllMessageRead)
		router.PATCH("/deleted/all", v1.UpdateAllMessageDeleted)
	}
	return router
}
