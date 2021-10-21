package v1

import (
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindFsm(c *gin.Context) {
	var r req.FsmMachine
	req.ShouldBind(c, &r)
	s := service.New(c)
	list, err := s.FindFsm(r)
	resp.CheckErr(err)
	resp.SuccessWithData(list)
}

func CreateFsm(c *gin.Context) {
	var r req.FsmCreateMachine
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.CreateFsm(r)
	resp.CheckErr(err)
	resp.Success()
}

func FindFsmApprovingLog(c *gin.Context) {
	var r req.FsmPendingLog
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.ApprovalRoleId = user.RoleId
	r.ApprovalUserId = user.Id
	s := service.New(c)
	list, err := s.FindFsmApprovingLog(r)
	resp.CheckErr(err)
	resp.SuccessWithData(list)
}
