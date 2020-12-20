package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 消息中心路由
func InitMessageRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	// 初始化消息中心仓库
	v1.StartMessageHub(middleware.CheckIdempotenceToken)
	router := r.Group("/message").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/ws", v1.MessageWs)
		router.GET("/all", v1.GetAllMessages)
		router.GET("/unRead/count", v1.GetUnReadMessageCount)
		router.
			Use(middleware.Idempotence).
			POST("/push", v1.PushMessage)
		router.PATCH("/read/batch", v1.BatchUpdateMessageRead)
		router.PATCH("/deleted/batch", v1.BatchUpdateMessageDeleted)
		router.PATCH("/read/all", v1.UpdateAllMessageRead)
		router.PATCH("/deleted/all", v1.UpdateAllMessageDeleted)
	}
	return router
}
