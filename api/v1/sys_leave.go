package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取请假列表
func GetLeaves(c *gin.Context) {
	// 绑定参数
	var req request.LeaveListRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取当前登录用户
	user := GetCurrentUser(c)
	req.UserId = user.Id
	// 创建服务
	s := cache_service.New(c)
	leaves, err := s.GetLeaves(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.LeaveListResponseStruct
	utils.Struct2StructByJson(leaves, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 获取请假列表
func GetLeaveApprovalLogs(c *gin.Context) {
	// 绑定参数
	var req request.LeaveListRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := cache_service.New(c)
	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	leaves, err := s.GetLeaveApprovalLogs(leaveId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 将日志包装下
	respStruct := make([]response.LeaveLogListResponseStruct, 0)
	for _, log := range leaves {
		respStruct = append(respStruct, response.LeaveLogListResponseStruct{
			LeaveId: leaveId,
			Log: response.WorkflowLogsListResponseStruct{
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
				CreatedAt:             log.Model.CreatedAt,
				UpdatedAt:             log.Model.UpdatedAt,
			},
		})
	}
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建请假
func CreateLeave(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateLeaveRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 参数校验
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前用户
	req.User = user
	// 创建服务
	s := service.New(c)
	err = s.CreateLeave(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新请假
func UpdateLeaveById(c *gin.Context) {
	// 绑定参数
	var req models.SysLeave
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的leaveId
	leaveId := utils.Str2Uint(c.Param("leaveId"))
	if leaveId == 0 {
		response.FailWithMsg("请假编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateLeaveById(leaveId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除请假
func BatchDeleteLeaveByIds(c *gin.Context) {
	var req request.Req
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteLeaveByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
