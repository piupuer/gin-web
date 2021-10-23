package v1

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindLeave(c *gin.Context) {
	var r request.LeaveReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.UserId = user.Id
	s := service.New(c)
	list := s.FindLeave(&r)
	resp.SuccessWithPageData(list, []response.LeaveResp{}, r.Page)
}

func FindLeaveFsmTrack(c *gin.Context) {
	var r request.LeaveReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	list, err := s.FindLeaveFsmTrack(id)
	resp.CheckErr(err)
	resp.SuccessWithData(list)
}

func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateLeaveReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	r.User = user
	s := service.New(c)
	err := s.CreateLeave(&r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateLeaveById(c *gin.Context) {
	var r request.UpdateLeaveReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(models.Leave))
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteLeaveByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.Uints(), new(models.Leave))
	resp.CheckErr(err)
	resp.Success()
}
