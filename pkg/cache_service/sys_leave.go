package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

// 获取所有请假(当前用户)
func (s *RedisService) GetLeaves(req *request.LeaveListRequestStruct) ([]models.SysLeave, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetLeaves(req)
	}
	var err error
	list := make([]models.SysLeave, 0)
	query := s.redis.
		Table(new(models.SysLeave).TableName()).
		Where("user_id", "=", req.UserId).
		Order("created_at DESC")
	statusVal, statusFlag := req.Status.Uint()
	if statusFlag {
		query = query.Where("status", "=", statusVal)
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		query = query.Where("desc", "contains", desc)
	}
	// TODO 按id逆序
	// query = query.SortBy("id", "desc")
	// 查询条数
	err = query.Count(&req.PageInfo.Total).Error
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
func (s *RedisService) GetLeaveApprovalLogs(leaveId uint) ([]models.SysWorkflowLog, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetLeaveApprovalLogs(leaveId)
	}
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
