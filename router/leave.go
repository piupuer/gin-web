package router

import (
	v1 "gin-web/api/v1"
	"github.com/piupuer/go-helper/router"
)

func InitLeaveRouter(r *router.Router) {
	router1 := r.Casbin("/leave")
	router2 := r.CasbinAndIdempotence("/leave")
	router1.GET("/list", v1.FindLeave)
	router2.POST("/create", v1.CreateLeave)
	router1.PATCH("/update/:id", v1.UpdateLeaveById)
	router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
}
