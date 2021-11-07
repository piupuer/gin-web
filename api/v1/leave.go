package v1

import (
	"context"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindLeave(c *gin.Context) {
	var r request.Leave
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.UserId = user.Id
	s := service.New(c)
	list := s.FindLeave(&r)
	resp.SuccessWithPageData(list, &[]response.Leave{}, r.Page)
}

func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateLeave
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	r.User = user
	s := service.New(c)
	err := s.CreateLeave(&r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateLeaveById(c *gin.Context) {
	var r request.UpdateLeave
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateLeaveById(id, r, user)
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteLeaveByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.DeleteLeaveByIds(r.Uints(), user)
	resp.CheckErr(err)
	resp.Success()
}

func LeaveTransition(ctx context.Context, logs ...resp.FsmApprovalLog) error {
	s := service.New(ctx)
	return s.LeaveTransition(logs...)
}

func GetLeaveFsmDetail(c *gin.Context, detail req.FsmSubmitterDetail) []resp.FsmSubmitterDetail {
	s := service.New(c)
	return s.GetLeaveFsmDetail(detail)
}

func UpdateLeaveFsmDetail(c *gin.Context, detail req.UpdateFsmSubmitterDetail) error {
	s := service.New(c)
	return s.UpdateLeaveFsmDetail(detail)
}
