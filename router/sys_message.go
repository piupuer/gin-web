package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 消息中心路由
func InitMessageRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/message")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/message")
	{
		router1.GET("/all", v1.GetAllMessages)
		router1.GET("/unRead/count", v1.GetUnReadMessageCount)
		router2.POST("/push", v1.PushMessage)
		router1.PATCH("/read/batch", v1.BatchUpdateMessageRead)
		router1.PATCH("/deleted/batch", v1.BatchUpdateMessageDeleted)
		router1.PATCH("/read/all", v1.UpdateAllMessageRead)
		router1.PATCH("/deleted/all", v1.UpdateAllMessageDeleted)
	}
	return r
}
