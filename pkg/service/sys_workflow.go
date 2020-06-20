package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/service/strategy"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/thedevsaddam/gojsonq/v2"
	"strings"
)

// 获取工作流(指定审批单类型)
func (s *MysqlService) GetWorkflowByTargetCategory(targetCategory uint) (models.SysWorkflow, error) {
	var flow models.SysWorkflow
	err := s.tx.Where("target_category = ?", targetCategory).First(&flow).Error
	return flow, err
}

// 获取所有工作流
func (s *MysqlService) GetWorkflows(req *request.WorkflowListRequestStruct) ([]models.SysWorkflow, error) {
	var err error
	list := make([]models.SysWorkflow, 0)
	query := s.tx.Table(new(models.SysWorkflow).TableName())
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Category > 0 {
		query = query.Where("category = ?", req.Category)
	}
	if req.TargetCategory > 0 {
		query = query.Where("targetCategory = ?", req.TargetCategory)
	}
	if req.Self != nil {
		query = query.Where("self = ?", *req.Self)
	}
	if req.SubmitUserConfirm != nil {
		query = query.Where("submitUserConfirm = ?", *req.SubmitUserConfirm)
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

// 获取所有流水线
func (s *MysqlService) GetWorkflowLines(req *request.WorkflowLineListRequestStruct) ([]models.SysWorkflowLine, error) {
	var err error
	list := make([]models.SysWorkflowLine, 0)
	query := s.tx.Preload("Users").Model(&models.SysWorkflowLine{})
	if req.FlowId > 0 {
		query = query.Where("flow_id = ?", req.FlowId)
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

// 查询待审批目标列表(指定用户)
func (s *MysqlService) GetWorkflowApprovings(req *request.WorkflowApprovingListRequestStruct) ([]models.SysWorkflowLog, error) {
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
	err = s.tx.
		Preload("Flow").
		Preload("CurrentLine").
		Preload("CurrentLine.Users").
		Preload("CurrentLine.Role").
		Preload("CurrentLine.Role.Users").
		Where("status = ?", models.SysWorkflowLogStateSubmit). // 状态已提交
		Find(&logs).Error
	if err != nil {
		return list, err
	}

	for _, log := range logs {
		// 获取当前待审批人
		userIds := s.GetApprovingUsers(log)
		log.ApprovingUserIds = userIds
		// 包含当前审批人
		if utils.ContainsUint(userIds, approval.Id) {
			list = append(list, log)
		}
	}
	// 处理分页(转为json)
	query := gojsonq.New().FromString(utils.Struct2Json(list))
	req.PageInfo.Total = uint(query.Count())
	if !req.PageInfo.NoPagination {
		// 获取分页参数
		limit, offset := req.GetLimit()
		res := query.Limit(int(limit)).Offset(int(offset)).Get()
		// 转换为结构体
		utils.Struct2StructByJson(res, &list)
	}
	return list, err
}

// 查询下一审批人(指定目标)
func (s *MysqlService) GetWorkflowNextApprovingUsers(flowId uint, targetId uint) ([]models.SysUser, error) {
	// 查询需要审核的日志
	var log models.SysWorkflowLog
	users := make([]models.SysUser, 0)
	err := s.tx.
		Preload("CurrentLine").
		Preload("CurrentLine.Users").
		Preload("CurrentLine.Role").
		Preload("CurrentLine.Role.Users").
		Where(&models.SysWorkflowLog{
			FlowId:   flowId,   // 流程号一致
			TargetId: targetId, // 目标一致
		}).Where(
		"status = ?", models.SysWorkflowLogStateSubmit, // 状态已提交
	).First(&log).Error
	if err != nil {
		return users, err
	}

	// 获取当前待审批人
	userIds := s.GetApprovingUsers(log)
	err = s.tx.Where("id IN (?)", userIds).Find(&users).Error
	return users, err
}

// 查询审批日志(指定目标)
func (s *MysqlService) GetWorkflowLogs(flowId uint, targetId uint) ([]models.SysWorkflowLog, error) {
	// 查询已审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	err := s.tx.Preload("ApprovalUser").Preload("SubmitUser").Preload("Flow").Where(&models.SysWorkflowLog{
		FlowId:   flowId,   // 流程号一致
		TargetId: targetId, // 目标一致
	}).Find(&logs).Error
	return logs, err
}

// 查询下一流水线
func (s *MysqlService) GetNextWorkflowLine(flowId uint, currentSort uint) (models.SysWorkflowLine, error) {
	return s.GetWorkflowLineBySort(flowId, currentSort+1)
}

// 查询上一流水线
func (s *MysqlService) GetPrevWorkflowLine(flowId uint, currentSort uint) (models.SysWorkflowLine, error) {
	return s.GetWorkflowLineBySort(flowId, currentSort-1)
}

// 查询指定sort流水线
func (s *MysqlService) GetWorkflowLineBySort(flowId uint, sort uint) (models.SysWorkflowLine, error) {
	// 查询流程线
	var line models.SysWorkflowLine
	if sort <= 0 {
		return line, gorm.ErrRecordNotFound
	}
	err := s.tx.Preload("Users").Where(&models.SysWorkflowLine{
		FlowId: flowId,
		Sort:   sort,
	}).First(&line).Error
	return line, err
}

// 创建工作流
func (s *MysqlService) CreateWorkflow(req *request.CreateWorkflowRequestStruct) (err error) {
	var flow models.SysWorkflow
	utils.Struct2StructByJson(req, &flow)
	// 生成uuid
	flow.Uuid = uuid.NewV4().String()
	// 创建数据
	err = s.tx.Create(&flow).Error
	return
}

// 更新工作流
func (s *MysqlService) UpdateWorkflowById(id uint, req gin.H) (err error) {
	var oldWorkflow models.SysWorkflow
	query := s.tx.Table(oldWorkflow.TableName()).Where("id = ?", id).First(&oldWorkflow)
	if query.RecordNotFound() {
		return fmt.Errorf("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldWorkflow, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除工作流
func (s *MysqlService) DeleteWorkflowByIds(ids []uint) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(models.SysWorkflow{}).Error
}

// 更新流程流水线
func (s *MysqlService) UpdateWorkflowLineByIncremental(req *request.UpdateWorkflowLineIncrementalRequestStruct) (err error) {
	// 查询流程以及流水线
	var oldFlow models.SysWorkflow
	oldLines := make([]models.SysWorkflowLine, 0)
	noFlow := s.tx.Where("id = ?", req.FlowId).First(&oldFlow).RecordNotFound()
	if noFlow {
		return fmt.Errorf("流程不存在")
	}
	err = s.tx.Where(&models.SysWorkflowLine{FlowId: req.FlowId}).Find(&oldLines).Error
	if err != nil {
		return
	}
	// 查询增改的所有用户/流水线
	userIds := make([]uint, 0)
	for _, item := range req.Create {
		for _, userId := range item.UserIds {
			if userId > 0 && !utils.ContainsUint(userIds, userId) {
				userIds = append(userIds, userId)
			}
		}
	}
	for _, item := range req.Update {
		for _, userId := range item.UserIds {
			if userId > 0 && !utils.ContainsUint(userIds, userId) {
				userIds = append(userIds, userId)
			}
		}
	}
	// 批量查找所有用户
	users, err := s.GetUsersByIds(userIds)
	if err != nil {
		return
	}

	// 删除/更新/新增需要有序, 否则可能会打破流水线顺序
	// 1. 删除流水线
	for _, item := range req.Delete {
		// 删除流水线
		err = s.tx.Where("id = ?", item.Id).Delete(models.SysWorkflowLine{}).Error
		if err != nil {
			return
		}
	}
	// 2. 更新流水线
	for _, item := range req.Update {
		var line models.SysWorkflowLine
		utils.Struct2StructByJson(item, &line)
		// 获取用户
		us := make([]models.SysUser, 0)
		if len(item.UserIds) > 0 {
			for _, userId := range item.UserIds {
				for _, user := range users {
					if userId == user.Id {
						us = append(us, user)
						break
					}
				}
			}
		}
		query := s.tx.Model(&line)
		if item.RoleId != nil {
			// 需要强制更新roleId
			query = query.Update("role_id", item.RoleId)
		}
		// 更新数据, 替换users
		err = query.Update(&line).Association("Users").Replace(&us).Error
		if err != nil {
			return
		}
	}
	// 2. 创建流水线
	for _, item := range req.Create {
		var line models.SysWorkflowLine
		utils.Struct2StructByJson(item, &line)
		// 获取用户
		us := make([]models.SysUser, 0)
		if len(item.UserIds) > 0 {
			for _, userId := range item.UserIds {
				for _, user := range users {
					if userId == user.Id {
						us = append(us, user)
						break
					}
				}
			}
		}
		line.Users = us
		// 设置flowId
		line.FlowId = oldFlow.Id
		// 创建数据
		err = s.tx.Create(&line).Error
		if err != nil {
			return
		}
	}
	newLines := make([]models.SysWorkflowLine, 0)
	err = s.tx.Where(&models.SysWorkflowLine{FlowId: req.FlowId}).Find(&newLines).Error
	if err != nil {
		return
	}
	// 序号重排
	count := len(newLines)
	endPtr := true
	sort := uint(1)
	for i, line := range newLines {
		var attr models.SysWorkflowLine
		// 结束标识
		if i == count-1 {
			attr.End = &endPtr
		}
		// 序号
		attr.Sort = sort
		err = s.tx.Model(&line).Update(&attr).Error
		if err != nil {
			return
		}
		sort += 1
	}
	return
}

// 工作流流转(从一个状态转移到另一个状态)
func (s *MysqlService) WorkflowTransition(req *request.WorkflowTransitionRequestStruct) error {
	if req.FlowId == 0 {
		return fmt.Errorf("流程号不存在, flowId=%d", req.FlowId)
	}
	if req.TargetId == 0 {
		return fmt.Errorf("目标表编号不存在, flowId=%d", req.TargetId)
	}
	// 查询最后一条审批日志, 判断是否存在
	var lastLog models.SysWorkflowLog
	var err error
	notFound := s.tx.
		Preload("CurrentLine").
		Preload("CurrentLine.Users").
		Preload("CurrentLine.Role").
		Preload("CurrentLine.Role.Users").
		Preload("Flow").
		Where(&models.SysWorkflowLog{
			TargetId: req.TargetId,
			FlowId:   req.FlowId,
		}).Last(&lastLog).RecordNotFound()
	if notFound {
		// 走提交逻辑
		err = s.first(req)
		if err != nil {
			return err
		}
	} else {
		// 走审批逻辑
		err = s.next(req, lastLog)
		if err != nil {
			return err
		}
	}

	// 查询最后一条审批日志
	var newLastLog models.SysWorkflowLog
	err = s.tx.
		Where(&models.SysWorkflowLog{
			TargetId: req.TargetId,
			FlowId:   req.FlowId,
		}).Last(&newLastLog).Error
	if err != nil {
		return err
	}

	// 查询最后一条日志状态, 回写到对应的目标表中
	ctx, err := strategy.NewAfterTransitionContext(s.tx, req.TargetCategory, req.TargetId, newLastLog)
	if err != nil {
		return err
	}
	// 执行更新策略
	return ctx.Strategy.UpdateTarget()
}

// 初次提交流程工单
func (s *MysqlService) first(req *request.WorkflowTransitionRequestStruct) error {
	if req.SubmitUserId == 0 {
		return fmt.Errorf("提交人不存在, submitUserId=%d", req.SubmitUserId)
	}
	// 查询提交人是否存在
	submitUser, err := s.GetUserById(req.SubmitUserId)
	if err != nil {
		return err
	}

	// 初次创建
	var firstLog models.SysWorkflowLog
	firstLog.FlowId = req.FlowId
	firstLog.TargetId = req.TargetId
	firstLog.SubmitDetail = req.SubmitDetail
	approvalStatus := models.SysWorkflowLogStateApproval
	// 状态为自己批准
	firstLog.Status = &approvalStatus
	// 当前流水线为开始流水线
	firstLog.SubmitUserId = submitUser.Id
	firstLog.ApprovalUserId = submitUser.Id
	approvalOpinion := req.ApprovalOpinion
	if strings.TrimSpace(approvalOpinion) == "" {
		approvalOpinion = "初次提交"
	}
	firstLog.ApprovalOpinion = approvalOpinion
	// 创建首条日志
	s.tx.Create(&firstLog)
	// 获取下一流水线
	nextLine, err := s.GetNextWorkflowLine(req.FlowId, 0)
	if err != nil {
		return err
	}
	// 状态为提交, 当前流水线指向下一流水线, 创建新日志
	err = s.newLog(models.SysWorkflowLogStateSubmit, nextLine.Id, firstLog)
	return err
}

// 第二次提交流程工单
func (s *MysqlService) next(req *request.WorkflowTransitionRequestStruct, lastLog models.SysWorkflowLog) error {
	if *lastLog.Status == models.SysWorkflowLogStateEnd {
		return fmt.Errorf("流程已结束")
	}
	if req.ApprovalUserId == 0 {
		return fmt.Errorf("审批人不存在, approvalUserId=%d", req.ApprovalUserId)
	}
	// 查询审批人是否存在
	approval, err := s.GetUserById(req.ApprovalUserId)
	if err != nil {
		return err
	}
	if lastLog.SubmitUserId == approval.Id {
		// 自我审批
		return s.selfStart(req, approval, lastLog)
	} else {
		// 正常审批
		return s.start(req, approval, lastLog)
	}
}

// 开始自我审批
func (s *MysqlService) selfStart(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	var err error
	if *req.ApprovalStatus == models.SysWorkflowLogStateCancel {
		if *lastLog.Status == models.SysWorkflowLogStateCancel {
			return fmt.Errorf("流程已被取消, 无需重复操作")
		}
		// 提交人主动取消
		approvalOpinion := req.ApprovalOpinion
		if strings.TrimSpace(approvalOpinion) == "" {
			approvalOpinion = "提交人主动取消"
		}
		err := s.updateLog(models.SysWorkflowLogStateCancel, approvalOpinion, approval, lastLog)
		return err
	} else if *lastLog.Status == models.SysWorkflowLogStateCancel {
		// 提交人再次重启
		if *req.ApprovalStatus == models.SysWorkflowLogStateRestart {
			// 从头开始创建新的
			req.SubmitUserId = req.ApprovalUserId
			if strings.TrimSpace(req.ApprovalOpinion) == "" {
				req.ApprovalOpinion = "提交人再次重启"
			}
			return s.first(req)
		}
	}
	var nextLine models.SysWorkflowLine
	// 获取下一流水线
	if lastLog.CurrentLineId > 0 && !*lastLog.CurrentLine.End {
		nextLine, err = s.GetNextWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
		if err != nil {
			return err
		}
	}
	// 1. 未结束 且 开启自我审批 且 有权限审批
	if !*lastLog.End && *lastLog.Flow.Self && s.checkPermission(approval.Id, lastLog) {
		if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
			if nextLine.Id == 0 {
				return s.end(false, req.ApprovalOpinion, approval, lastLog)
			}
			// 通过
			return s.approval(req, approval, lastLog)
		} else {
			// 回退到上一流水线
			return s.deny(req, approval, lastLog)
		}
	}
	// 2.开启提交人确认
	if *lastLog.Flow.SubmitUserConfirm {
		// 下一流水线为空
		if nextLine.Id == 0 {
			var log models.SysWorkflowLog
			err = s.tx.
				Where(&models.SysWorkflowLog{
					FlowId:   req.FlowId,   // 流程号一致
					TargetId: req.TargetId, // 目标一致
				}).
				Where("status > ?", models.SysWorkflowLogStateSubmit). // 状态非提交
				Last(&log).Error

			if err != nil {
				return err
			}
			// 最后审批为通过 且 当前也为通过, 直接结束
			if *log.Status == models.SysWorkflowLogStateApproval && *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
				return s.end(true, req.ApprovalOpinion, approval, lastLog)
			}
		}
	}
	return fmt.Errorf("无权限自我审批, 请查看是否开启自我审批或提交人确认")
}

// 开始正常审批
func (s *MysqlService) start(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 当前状态已提交 且 必须有权限审批
	if *lastLog.Status == models.SysWorkflowLogStateSubmit && s.checkPermission(approval.Id, lastLog) {
		if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
			// 通过
			return s.approval(req, approval, lastLog)
		} else {
			// 回退到上一流水线
			return s.deny(req, approval, lastLog)
		}
	}
	return fmt.Errorf("无权限审批或审批流程未创建")
}

// 通过审批, 流转到下一流水线
func (s *MysqlService) approval(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 默认流水线不变
	lineId := lastLog.CurrentLineId
	if s.checkNextLineSort(approval.Id, lastLog) {
		// 流转到下一流水线
		var err error
		var nextLine models.SysWorkflowLine
		if !*lastLog.CurrentLine.End {
			// 获取下一流水线
			nextLine, err = s.GetNextWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
			if err != nil {
				return err
			}
		}
		// 下一流水线为空, 直接结束
		if nextLine.Id == 0 {
			return s.end(false, req.ApprovalOpinion, approval, lastLog)
		}
		// 更新日志
		err = s.updateLog(models.SysWorkflowLogStateApproval, req.ApprovalOpinion, approval, lastLog)
		if err != nil {
			return err
		}
		// 当前流水线指向下一流水线
		lineId = nextLine.Id
	} else {
		// 保留当前流水线, 还需其他人继续审批
		err := s.updateLog(models.SysWorkflowLogStateApproval, req.ApprovalOpinion, approval, lastLog)
		if err != nil {
			return err
		}
	}
	// 状态为提交, 创建新日志
	err := s.newLog(models.SysWorkflowLogStateSubmit, lineId, lastLog)
	return err
}

// 拒绝审批, 回退到上一流水线
func (s *MysqlService) deny(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 流转到上一流水线
	// 获取上一流水线
	if lastLog.CurrentLine.Sort <= 1 {
		// 上一流水线不存在, 说明拒绝到最初提交状态
		err := s.updateLog(models.SysWorkflowLogStateDeny, req.ApprovalOpinion, approval, lastLog)
		return err
	}
	prevLine, err := s.GetPrevWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
	if err != nil {
		return err
	}
	// 更新日志
	err = s.updateLog(models.SysWorkflowLogStateDeny, req.ApprovalOpinion, approval, lastLog)
	if err != nil {
		return err
	}
	// 状态为提交,当前流水线指向上一流水线, 创建新日志
	err = s.newLog(models.SysWorkflowLogStateSubmit, prevLine.Id, lastLog)
	return err
}

// 结束审批, 末尾流水线
func (s *MysqlService) end(submitConfirm bool, approvalOpinion string, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 结束
	status := models.SysWorkflowLogStateEnd
	// 状态为结束
	end := true
	if *lastLog.Flow.SubmitUserConfirm && !submitConfirm {
		// 开启了提交人确认则不能直接结束
		status = models.SysWorkflowLogStateApproval
		end = false
	}
	lastLog.End = &end
	// 提交人确认流水线
	if end && lastLog.SubmitUserId == approval.Id && strings.TrimSpace(approvalOpinion) == "" {
		approvalOpinion = "提交人已确认"
	}
	err := s.updateLog(status, approvalOpinion, approval, lastLog)
	if err != nil {
		return err
	}
	if end {
		return nil
	}
	end = true
	lastLog.End = &end
	// 创建新记录
	return s.newLog(models.SysWorkflowLogStateSubmit, 0, lastLog)
}

// 更新日志
func (s *MysqlService) updateLog(status uint, approvalOpinion string, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	var updateLog models.SysWorkflowLog
	// 流程信息
	updateLog.FlowId = lastLog.FlowId
	updateLog.TargetId = lastLog.TargetId
	// 状态
	updateLog.Status = &status
	// 是否结束
	updateLog.End = lastLog.End
	// 提交人
	updateLog.SubmitUserId = lastLog.SubmitUserId
	// 提交明细
	updateLog.SubmitDetail = lastLog.SubmitDetail
	// 审批人以及意见
	updateLog.ApprovalUserId = approval.Id
	updateLog.ApprovalOpinion = approvalOpinion
	err := s.tx.Table(updateLog.TableName()).Where("id = ?", lastLog.Id).Update(&updateLog).Error
	return err
}

// 创建新的日志
func (s *MysqlService) newLog(status uint, lineId uint, lastLog models.SysWorkflowLog) error {
	var newLog models.SysWorkflowLog
	// 流程信息
	newLog.FlowId = lastLog.FlowId
	newLog.TargetId = lastLog.TargetId
	// 当前流水线
	newLog.CurrentLineId = lineId
	// 状态
	newLog.Status = &status
	// 是否结束
	newLog.End = lastLog.End
	// 提交人
	newLog.SubmitUserId = lastLog.SubmitUserId
	// 提交人
	newLog.SubmitDetail = lastLog.SubmitDetail
	// 创建数据
	err := s.tx.Create(&newLog).Error
	return err
}

// 检查当前审批人是否有权限
func (s *MysqlService) checkPermission(approvalUserId uint, lastLog models.SysWorkflowLog) bool {
	// 获取当前待审批人
	userIds := s.GetApprovingUsers(lastLog)
	return utils.ContainsUint(userIds, approvalUserId)
}

// 检查是否可以切换流水线到下一个(通过审批会使用)
func (s *MysqlService) checkNextLineSort(approvalUserId uint, lastLog models.SysWorkflowLog) bool {
	// 获取当前待审批人
	userIds := s.GetApprovingUsers(lastLog)
	// 判断流程类别
	switch lastLog.Flow.Category {
	case models.SysWorkflowCategoryOnlyOneApproval:
		// 只需要1人通过: 当前审批人在待审批列表中
		return utils.ContainsUint(userIds, approvalUserId)
	case models.SysWorkflowCategoryAllApproval:
		// 查询全部审批人数
		allUserIds := s.getAllApprovalUsers(lastLog)
		// 查询历史审批人数
		historyUserIds := s.getHistoryApprovalUsers(lastLog)
		// 需要全部人通过: 当前审批人在待审批列表中 且 历史审批人+当前审批人刚好等于全部审批人
		return utils.ContainsUint(userIds, approvalUserId) && len(historyUserIds) >= len(allUserIds)-1
	}
	return false
}

// 获取待审批人(当前流水线)
func (s *MysqlService) GetApprovingUsers(log models.SysWorkflowLog) []uint {
	userIds := make([]uint, 0)
	allUserIds := s.getAllApprovalUsers(log)
	historyUserIds := s.getHistoryApprovalUsers(log)
	for _, allUserId := range allUserIds {
		// 不在历史列表中
		if !utils.ContainsUint(historyUserIds, allUserId) {
			userIds = append(userIds, allUserId)
		}
	}
	return userIds
}

// 获取历史审批人(最后一个流水线, 主要用于判断是否审批完成)
func (s *MysqlService) getHistoryApprovalUsers(log models.SysWorkflowLog) []uint {
	historyUserIds := make([]uint, 0)
	// 查询已审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	err := s.tx.Where(&models.SysWorkflowLog{
		FlowId:   log.FlowId,   // 流程号一致
		TargetId: log.TargetId, // 目标一致
	}).Where(
		"status > ?", models.SysWorkflowLogStateSubmit, // 状态非提交
	).Order("id DESC").Find(&logs).Error
	if err != nil {
		global.Log.Warn("[getHistoryApprovalUsers]", err)
	}
	// 保留连续审核通过记录
	l := len(logs)
	for i := 0; i < l; i++ {
		item := logs[i]
		// 如果不是通过立即结束, 必须保证连续的通过 或 当前流水线不一致
		if *item.Status != models.SysWorkflowLogStateApproval || item.CurrentLineId != log.CurrentLineId {
			break
		}
		// 审批人为配置中的一人
		if !utils.ContainsUint(historyUserIds, item.ApprovalUserId) {
			historyUserIds = append(historyUserIds, item.ApprovalUserId)
		}
	}
	return historyUserIds
}

// 获取全部审批人(当前流水线)
func (s *MysqlService) getAllApprovalUsers(log models.SysWorkflowLog) []uint {
	userIds := make([]uint, 0)
	if log.CurrentLineId == 0 && *log.Flow.SubmitUserConfirm {
		// 末尾节点 且 开启提交人确认
		userIds = append(userIds, log.SubmitUserId)
	} else {
		for _, user := range log.CurrentLine.Users {
			if user.Id > 0 && !utils.ContainsUint(userIds, user.Id) {
				userIds = append(userIds, user.Id)
			}
		}
		if log.CurrentLine.RoleId > 0 {
			for _, user := range log.CurrentLine.Role.Users {
				if user.Id > 0 && !utils.ContainsUint(userIds, user.Id) {
					userIds = append(userIds, user.Id)
				}
			}
		}
	}
	return userIds
}
