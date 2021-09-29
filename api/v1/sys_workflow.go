package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// 获取全部审批日志(目前有1: 请假)
func GetWorkflowApprovings(c *gin.Context) {
	var req request.WorkflowApprovingRequestStruct
	request.ShouldBind(c, &req)
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	// 绑定当前用户
	req.ApprovalUserId = user.Id
	approvings, err := s.GetWorkflowApprovings(&req)
	response.CheckErr(err)
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
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 获取工作流列表
func GetWorkflows(c *gin.Context) {
	var req request.WorkflowRequestStruct
	request.ShouldBind(c, &req)
	s := cache_service.New(c)
	workflows, err := s.GetWorkflows(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.WorkflowListResponseStruct
	utils.Struct2StructByJson(workflows, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 获取工作流列表
func GetWorkflowLines(c *gin.Context) {
	var req request.WorkflowLineRequestStruct
	request.ShouldBind(c, &req)
	s := cache_service.New(c)
	workflowLines, err := s.GetWorkflowLines(&req)
	response.CheckErr(err)
	// 隐藏部分字段
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
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建工作流
func CreateWorkflow(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateWorkflowRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	req.Uuid = uuid.NewV4().String()
	err := s.Create(req, new(models.SysWorkflow))
	response.CheckErr(err)
	response.Success()
}

// 更新工作流流水线
func UpdateWorkflowLineIncremental(c *gin.Context) {
	var req request.UpdateWorkflowLineIncrementalRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	s := service.New(c)
	// 更新流程线
	err := s.UpdateWorkflowLineByIncremental(&req)
	response.CheckErr(err)
	response.Success()
}

// 更新工作流
func UpdateWorkflowById(c *gin.Context) {
	var req request.UpdateWorkflowRequestStruct
	request.ShouldBind(c, &req)
	workflowId := utils.Str2Uint(c.Param("workflowId"))
	if workflowId == 0 {
		response.CheckErr("工作流编号不正确")
	}
	s := service.New(c)
	err := s.UpdateById(workflowId, req, new(models.SysWorkflow))
	response.CheckErr(err)
	response.Success()
}

// 更新工作流日志: 审批通过
func UpdateWorkflowLogApproval(c *gin.Context) {
	var req request.WorkflowTransitionRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	user := GetCurrentUser(c)
	// 记录审批人
	req.ApprovalUserId = user.Id
	s := service.New(c)
	// 工作流状态转移
	err := s.WorkflowTransition(&req)
	response.CheckErr(err)
	response.Success()
}

// 批量删除工作流
func BatchDeleteWorkflowByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysWorkflow))
	response.CheckErr(err)
	response.Success()
}
