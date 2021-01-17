package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/thedevsaddam/gojsonq/v2"
	"strings"
)

// 获取工作流(指定审批单类型)
func (s *RedisService) GetWorkflowByTargetCategory(targetCategory uint) (models.SysWorkflow, error) {
	var flow models.SysWorkflow
	err := s.redis.Table(flow.TableName()).Where("target_category", "=", targetCategory).First(&flow).Error
	return flow, err
}

// 获取所有工作流
func (s *RedisService) GetWorkflows(req *request.WorkflowListRequestStruct) ([]models.SysWorkflow, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetWorkflows(req)
	}
	var err error
	list := make([]models.SysWorkflow, 0)
	query := s.redis.Table(new(models.SysWorkflow).TableName())
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name", "contains", name)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Category != nil {
		query = query.Where("category", "=", *req.Category)
	}
	if req.TargetCategory != nil {
		query = query.Where("target_category", "=", *req.TargetCategory)
	}
	if req.Self != nil {
		query = query.Where("self", "=", *req.Self)
	}
	if req.SubmitUserConfirm != nil {
		query = query.Where("submit_user_confirm", "=", *req.SubmitUserConfirm)
	}

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

// 获取所有流水线
func (s *RedisService) GetWorkflowLines(req *request.WorkflowLineListRequestStruct) ([]models.SysWorkflowLine, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetWorkflowLines(req)
	}
	var err error
	list := make([]models.SysWorkflowLine, 0)
	query := s.redis.Table(new(models.SysWorkflowLine).TableName()).Preload("Users")
	if req.FlowId > 0 {
		query = query.Where("flow_id", "=", req.FlowId)
	}

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

// 查询审批日志(指定目标)
func (s *RedisService) GetWorkflowLogs(flowId uint, targetId uint) ([]models.SysWorkflowLog, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetWorkflowLogs(flowId, targetId)
	}
	// 查询已审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	err := s.redis.Table(new(models.SysWorkflowLog).TableName()).
		Preload("ApprovalUser").
		Preload("SubmitUser").
		Preload("Flow").
		Where("flow_id", "=", flowId). // 流程号一致
		Where("target_id", "=", targetId). // 目标一致
		Find(&logs).Error
	return logs, err
}

// 查询待审批目标列表(指定用户)
func (s *RedisService) GetWorkflowApprovings(req *request.WorkflowApprovingListRequestStruct) ([]models.SysWorkflowLog, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetWorkflowApprovings(req)
	}
	// 查询需要审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	list := make([]models.SysWorkflowLog, 0)
	if req.ApprovalUserId == 0 {
		return list, fmt.Errorf("用户不存在, approvalUserId=%d", req.ApprovalUserId)
	}
	// 查询审批人
	approval, err := s.GetUserById(req.ApprovalUserId)
	if err != nil {
		return list, err
	}
	// 由于还需判断是否包含当前审批人, 因此无法直接分页
	err = s.redis.Table(new(models.SysWorkflowLog).TableName()).
		Preload("Flow").
		Preload("CurrentLine").
		Preload("CurrentLine.Users").
		Preload("CurrentLine.Role").
		Preload("CurrentLine.Role.Users").
		Where("status", "=", models.SysWorkflowLogStateSubmit). // 状态已提交
		Find(&logs).Error
	if err != nil {
		return list, err
	}

	for _, log := range logs {
		// 获取当前待审批人
		userIds := s.mysql.GetApprovingUsers(log)
		log.ApprovingUserIds = userIds
		// 包含当前审批人
		if utils.ContainsUint(userIds, approval.Id) {
			list = append(list, log)
		}
	}
	// 处理分页(转为json)
	query := gojsonq.New().FromString(utils.Struct2Json(list))
	req.PageInfo.Total = int64(query.Count())
	if !req.PageInfo.NoPagination {
		// 获取分页参数
		limit, offset := req.GetLimit()
		res := query.Limit(int(limit)).Offset(int(offset)).Get()
		// 转换为结构体
		utils.Struct2StructByJson(res, &list)
	}
	return list, err
}
