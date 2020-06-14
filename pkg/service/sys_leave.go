package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// 获取所有请假(当前用户)
func (s *MysqlService) GetLeaves(req *request.LeaveListRequestStruct) ([]models.SysLeave, error) {
	var err error
	list := make([]models.SysLeave, 0)
	query := s.tx.Where("user_id = ?", req.UserId)
	desc := strings.TrimSpace(req.Desc)
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if desc != "" {
		query = query.Where("desc LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	// 按id逆序
	query = query.Order("id DESC")
	// 查询条数
	err = query.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
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
	var leave models.SysLeave
	utils.Struct2StructByJson(req, &leave)
	// 创建数据
	err = s.tx.Create(&leave).Error
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
			"请假单[申请人: %s(%s), 申请说明: %s]",
			req.User.Nickname,
			req.User.Username,
			leave.Desc,
		), // 提交明细
	})
	return
}

// 更新请假
func (s *MysqlService) UpdateLeaveById(id uint, req gin.H) (err error) {
	var leave models.SysLeave
	query := s.tx.Table(leave.TableName()).Where("id = ?", id).First(&leave)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(leave, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除请假
func (s *MysqlService) DeleteLeaveByIds(ids []uint) (err error) {
	var list []models.SysLeave
	query := s.tx.Where("id IN (?)", ids).Find(&list)
	if query.Error != nil {
		return
	}
	return query.Delete(models.SysLeave{}).Error
}
