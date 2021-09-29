package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取请假列表
func GetLeaves(c *gin.Context) {
	var req request.LeaveRequestStruct
	request.ShouldBind(c, &req)
	// 获取当前登录用户
	user := GetCurrentUser(c)
	req.UserId = user.Id
	s := cache_service.New(c)
	leaves, err := s.GetLeaves(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.LeaveListResponseStruct
	utils.Struct2StructByJson(leaves, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 获取请假列表
func GetLeaveApprovalLogs(c *gin.Context) {
	var req request.LeaveRequestStruct
	request.ShouldBind(c, &req)
	s := cache_service.New(c)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	leaves, err := s.GetLeaveApprovalLogs(leaveId)
	response.CheckErr(err)

	// 将日志包装下
	respStruct := make([]response.LeaveLogListResponseStruct, 0)
	for _, log := range leaves {
		respStruct = append(respStruct, response.LeaveLogListResponseStruct{
			LeaveId: leaveId,
			Log: response.WorkflowLogsListResponseStruct{
				BaseData: response.BaseData{
					CreatedAt: log.CreatedAt,
					UpdatedAt: log.UpdatedAt,
				},
				FlowName:              log.Flow.Name,
				FlowUuid:              log.Flow.Uuid,
				FlowCategoryStr:       models.SysWorkflowCategoryConst[log.Flow.Category],
				FlowTargetCategoryStr: models.SysWorkflowTargetCategoryConst[log.Flow.TargetCategory],
				Status:                log.Status,
				StatusStr:             models.SysWorkflowLogStateConst[*log.Status],
				SubmitUsername:        log.SubmitUser.Username,
				SubmitUserNickname:    log.SubmitUser.Nickname,
				ApprovalUsername:      log.ApprovalUser.Username,
				ApprovalUserNickname:  log.ApprovalUser.Nickname,
				ApprovalOpinion:       log.ApprovalOpinion,
			},
		})
	}
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建请假
func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateLeaveRequestStruct
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
	var req request.UpdateLeaveRequestStruct
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
