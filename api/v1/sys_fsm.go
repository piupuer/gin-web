package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
)

// find finite state machine
func FindFsm(c *gin.Context) {
	var r req.FsmMachine
	request.ShouldBind(c, &r)
	s := service.New(c)
	list, err := s.FindFsm(r)
	response.CheckErr(err)
	response.SuccessWithData(list)
}

// create finite state machine
func CreateFsm(c *gin.Context) {
	var r req.FsmCreateMachine
	request.ShouldBind(c, &r)
	s := service.New(c)
	err := s.CreateFsm(r)
	response.CheckErr(err)
	response.Success()
}

// find waiting approve log
func FindFsmApprovingLog(c *gin.Context) {
	var r req.FsmPendingLog
	request.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.ApprovalRoleId = user.RoleId
	r.ApprovalUserId = user.Id
	s := service.New(c)
	list, err := s.FindFsmApprovingLog(r)
	response.CheckErr(err)
	response.SuccessWithData(list)
}
