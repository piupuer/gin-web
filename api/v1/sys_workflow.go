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

// 获取全部审批日志(目前有1: 请假)
func GetWorkflowApprovings(c *gin.Context) {
	// 绑定参数
	var req request.WorkflowApprovingListRequestStruct
	_ = c.Bind(&req)
	user := GetCurrentUser(c)
	// 创建服务
	s := cache_service.New(c)
	// 绑定当前用户
	req.ApprovalUserId = user.Id
	approvings, err := s.GetWorkflowApprovings(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 将日志包装下
	respStruct := make([]response.WorkflowLogsListResponseStruct, 0)
	for _, log := range approvings {
		respStruct = append(respStruct, response.WorkflowLogsListResponseStruct{
			Id:                    log.Id,
			FlowName:              log.Flow.Name,
			FlowId:                log.Flow.Id,
			FlowUuid:              log.Flow.Uuid,
			FlowCategory:          log.Flow.Category,
			FlowCategoryStr:       models.SysWorkflowCategoryConst[log.Flow.Category],
			FlowTargetCategory:    log.Flow.TargetCategory,
			FlowTargetCategoryStr: models.SysWorkflowTargetCategoryConst[log.Flow.TargetCategory],
			TargetId:              log.TargetId,
			Status:                log.Status,
			StatusStr:             models.SysWorkflowLogStateConst[*log.Status],
			SubmitUsername:        log.SubmitUser.Username,
			SubmitUserNickname:    log.SubmitUser.Nickname,
			SubmitDetail:          log.SubmitDetail,
			ApprovalUsername:      log.ApprovalUser.Username,
			ApprovalUserNickname:  log.ApprovalUser.Nickname,
			ApprovalOpinion:       log.ApprovalOpinion,
			ApprovingUserIds:      log.ApprovingUserIds,
			CreatedAt:             log.Model.CreatedAt,
			UpdatedAt:             log.Model.UpdatedAt,
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

// 获取工作流列表
func GetWorkflows(c *gin.Context) {
	// 绑定参数
	var req request.WorkflowListRequestStruct
	_ = c.Bind(&req)
	// 创建服务
	s := cache_service.New(c)
	workflows, err := s.GetWorkflows(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.WorkflowListResponseStruct
	utils.Struct2StructByJson(workflows, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 获取工作流列表
func GetWorkflowLines(c *gin.Context) {
	// 绑定参数
	var req request.WorkflowLineListRequestStruct
	_ = c.Bind(&req)
	// 创建服务
	s := cache_service.New(c)
	workflowLines, err := s.GetWorkflowLines(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.WorkflowLineListResponseStruct
	utils.Struct2StructByJson(workflowLines, &respStruct)
	// 绑定流水线userIds
	for i, line := range workflowLines {
		userIds := make([]uint, 0)
		if len(line.Node.Users) > 0 {
			for _, user := range line.Node.Users {
				userIds = append(userIds, user.Id)
			}
		}
		respStruct[i].Node.UserIds = userIds
	}
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建工作流
func CreateWorkflow(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateWorkflowRequestStruct
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	err = s.CreateWorkflow(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新工作流流水线
func UpdateWorkflowLineByNodes(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.UpdateWorkflowLineRequestStruct
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	// 更新流程线
	err = s.UpdateWorkflowLineByNodes(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新工作流
func UpdateWorkflowById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	_ = c.Bind(&req)
	// 获取path中的workflowId
	workflowId := utils.Str2Uint(c.Param("workflowId"))
	if workflowId == 0 {
		response.FailWithMsg("工作流编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err := s.UpdateWorkflowById(workflowId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新工作流日志: 审批通过
func UpdateWorkflowLogApproval(c *gin.Context) {
	// 绑定参数
	var req request.WorkflowTransitionRequestStruct
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	user := GetCurrentUser(c)
	// 记录审批人
	req.ApprovalUserId = user.Id
	// 创建服务
	s := service.New(c)
	// 工作流状态转移
	err = s.WorkflowTransition(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除工作流
func BatchDeleteWorkflowByIds(c *gin.Context) {
	var req request.Req
	_ = c.Bind(&req)
	// 创建服务
	s := service.New(c)
	// 删除数据
	err := s.DeleteWorkflowByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
