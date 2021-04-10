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
	uuid "github.com/satori/go.uuid"
)

// 获取全部审批日志(目前有1: 请假)
func GetWorkflowApprovings(c *gin.Context) {
	// 绑定参数
	var req request.WorkflowApprovingRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}
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
			BaseData: response.BaseData{
				Id:        log.Id,
				CreatedAt: log.CreatedAt,
				UpdatedAt: log.UpdatedAt,
			},
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
	var req request.WorkflowRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

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
	var req request.WorkflowLineRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

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
		if len(line.Users) > 0 {
			for _, user := range line.Users {
				userIds = append(userIds, user.Id)
			}
		}
		respStruct[i].UserIds = userIds
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
	err := c.ShouldBind(&req)
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
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	req.Uuid = uuid.NewV4().String()
	err = s.Create(req, new(models.SysWorkflow))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新工作流流水线
func UpdateWorkflowLineIncremental(c *gin.Context) {
	// 绑定参数
	var req request.UpdateWorkflowLineIncrementalRequestStruct
	err := c.ShouldBind(&req)
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
	// 创建服务
	s := service.New(c)
	// 更新流程线
	err = s.UpdateWorkflowLineByIncremental(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新工作流
func UpdateWorkflowById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateWorkflowRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的workflowId
	workflowId := utils.Str2Uint(c.Param("workflowId"))
	if workflowId == 0 {
		response.FailWithMsg("工作流编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateById(workflowId, req, new(models.SysWorkflow))
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
	err := c.ShouldBind(&req)
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
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysWorkflow))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
