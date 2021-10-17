package v1

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取请假列表
func GetLeaves(c *gin.Context) {
	var r request.LeaveReq
	req.ShouldBind(c, &r)
	// 获取当前登录用户
	user := GetCurrentUser(c)
	r.UserId = user.Id
	s := service.New(c)
	list, err := s.GetLeaves(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.LeaveResp{}, r.Page)
}

// 获取请假列表
func FindLeaveFsmTrack(c *gin.Context) {
	var r request.LeaveReq
	req.ShouldBind(c, &r)
	s := service.New(c)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	logs, err := s.FindLeaveFsmTrack(leaveId)
	resp.CheckErr(err)
	resp.SuccessWithData(logs)
}

// 创建请假
func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateLeaveReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	// 记录当前用户
	r.User = user
	s := service.New(c)
	err := s.CreateLeave(&r)
	resp.CheckErr(err)
	resp.Success()
}

// 更新请假
func UpdateLeaveById(c *gin.Context) {
	var r request.UpdateLeaveReq
	req.ShouldBind(c, &r)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	if leaveId == 0 {
		resp.CheckErr("请假编号不正确")
	}
	s := service.New(c)
	err := s.UpdateById(leaveId, r, new(models.SysLeave))
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除请假
func BatchDeleteLeaveByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.DeleteByIds(r.GetUintIds(), new(models.SysLeave))
	resp.CheckErr(err)
	resp.Success()
}
