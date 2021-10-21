package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitMessageRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	ops := &middleware.IdempotenceOptions{}
	for _, f := range GetIdempotenceMiddlewareOps() {
		f(ops)
	}
	v1.StartMessageHub(ops)
	router1 := GetCasbinRouter(r, jwtOptions, "/message")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/message")
	{
		router1.GET("/ws", v1.MessageWs)
		router1.GET("/all", v1.FindMessage)
		router1.GET("/unRead/count", v1.GetUnReadMessageCount)
		router2.POST("/push", v1.PushMessage)
		router1.PATCH("/read/batch", v1.BatchUpdateMessageRead)
		router1.PATCH("/deleted/batch", v1.BatchUpdateMessageDeleted)
		router1.PATCH("/read/all", v1.UpdateAllMessageRead)
		router1.PATCH("/deleted/all", v1.UpdateAllMessageDeleted)
	}
	return r
}
