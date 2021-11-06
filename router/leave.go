package router

import (
	v1 "gin-web/api/v1"
	"github.com/piupuer/go-helper/router"
)

func InitLeaveRouter(r *router.Router) {
	router1 := r.Casbin("/leave")
	router2 := r.CasbinAndIdempotence("/leave")
	router1.GET("/list", v1.FindLeave)
	router1.GET("/approving/track/:id", v1.FindLeaveFsmTrack)
	router2.POST("/create", v1.CreateLeave)
	router1.PATCH("/update/:id", v1.UpdateLeaveById)
	router1.PATCH("/resubmit/:id", v1.ResubmitLeaveById)
	router1.PATCH("/confirm/:id", v1.ConfirmLeaveById)
	router1.PATCH("/cancel/:id", v1.CancelLeaveById)
	router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
}
