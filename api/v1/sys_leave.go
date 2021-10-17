package v1

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取请假列表
func GetLeaves(c *gin.Context) {
	var req request.LeaveReq
	request.ShouldBind(c, &req)
	// 获取当前登录用户
	user := GetCurrentUser(c)
	req.UserId = user.Id
	s := service.New(c)
	leaves, err := s.GetLeaves(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.LeaveResp
	utils.Struct2StructByJson(leaves, &respStruct)
	// 返回分页数据
	var rp resp.PageData
	rp.Page = req.Page
	rp.List = respStruct
	response.SuccessWithData(rp)
}

// 获取请假列表
func FindLeaveFsmTrack(c *gin.Context) {
	var req request.LeaveReq
	request.ShouldBind(c, &req)
	s := service.New(c)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	logs, err := s.FindLeaveFsmTrack(leaveId)
	response.CheckErr(err)
	response.SuccessWithData(logs)
}

// 创建请假
func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateLeaveReq
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	// 记录当前用户
	req.User = user
	s := service.New(c)
	err := s.CreateLeave(&req)
	response.CheckErr(err)
	response.Success()
}

// 更新请假
func UpdateLeaveById(c *gin.Context) {
	var req request.UpdateLeaveReq
	request.ShouldBind(c, &req)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	if leaveId == 0 {
		response.CheckErr("请假编号不正确")
	}
	s := service.New(c)
	err := s.UpdateById(leaveId, req, new(models.SysLeave))
	response.CheckErr(err)
	response.Success()
}

// 批量删除请假
func BatchDeleteLeaveByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysLeave))
	response.CheckErr(err)
	response.Success()
}
