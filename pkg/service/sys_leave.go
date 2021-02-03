package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"strings"
)

// 获取所有请假(当前用户)
func (s *MysqlService) GetLeaves(req *request.LeaveListRequestStruct) ([]models.SysLeave, error) {
	var err error
	list := make([]models.SysLeave, 0)
	query := s.tx.
		Table(new(models.SysLeave).TableName()).
		Order("created_at DESC").
		Where("user_id = ?", req.UserId)
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		query = query.Where("desc LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	// 查询条数
	err = query.Count(&req.PageInfo.Total).Error
	if err == nil && req.PageInfo.Total > 0 {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 获取请假审批日志(指定请假编号)
func (s *MysqlService) GetLeaveApprovalLogs(leaveId uint) ([]models.SysWorkflowLog, error) {
	list := make([]models.SysWorkflowLog, 0)

	// 获取请假对应的工作流
	flow, err := s.GetWorkflowByTargetCategory(models.SysWorkflowTargetCategoryLeave)
	if err != nil {
		return list, err
	}
	// 获取工作流日志
	list, err = s.GetWorkflowLogs(flow.Id, leaveId)
	return list, err
}

// 创建请假
func (s *MysqlService) CreateLeave(req *request.CreateLeaveRequestStruct) (err error) {
	leave := new(models.SysLeave)
	err = s.Create(req, &leave)
	if err != nil {
		return
	}
	// 获取请假对应的工作流
	flow, err := s.GetWorkflowByTargetCategory(models.SysWorkflowTargetCategoryLeave)
	if err != nil {
		return
	}
	// 创建请假工作流日志
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       leave.Id,                              // 请假编号
		SubmitUserId:   req.User.Id,                           // 提交人编号
		SubmitDetail: fmt.Sprintf(
			"请假条[申请人: %s(%s), 申请说明: %s]",
			req.User.Nickname,
			req.User.Username,
			leave.Desc,
		), // 提交明细
	})
	return
}
