package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// 获取所有工作流
func (s *MysqlService) GetWorkflows(req *request.WorkflowListRequestStruct) ([]models.SysWorkflow, error) {
	var err error
	list := make([]models.SysWorkflow, 0)
	query := s.tx
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

// 获取所有流水线
func (s *MysqlService) GetWorkflowLines(req *request.WorkflowLineListRequestStruct) ([]models.SysWorkflowLine, error) {
	var err error
	list := make([]models.SysWorkflowLine, 0)
	query := s.tx
	if req.FlowId > 0 {
		query = query.Where("flow_id = ?", req.FlowId)
	}

	// 查询条数
	err = query.Preload("Node").Find(&list).Count(&req.PageInfo.Total).Error
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

// 查询待审批目标列表(指定用户)
func (s *MysqlService) GetWorkflowApprovingList(flowId uint, targetId uint, approvalId uint) ([]models.SysWorkflowLog, error) {
	// 查询需要审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	list := make([]models.SysWorkflowLog, 0)
	if approvalId == 0 {
		return list, fmt.Errorf("用户不存在, approvalId=%d", approvalId)
	}
	// 查询审批人
	approval, err := s.GetUserById(approvalId)
	if err != nil {
		return list, err
	}
	// 查询日志
	err = s.tx.
		Preload("CurrentLine").
		Preload("CurrentLine.Node").
		Preload("CurrentLine.Node.Users").
		Preload("CurrentLine.Node.Role").
		Preload("CurrentLine.Node.Role.Users").
		Where(&models.SysWorkflowLog{
			FlowId:   flowId,   // 流程号一致
			TargetId: targetId, // 目标一致
		}).Where(
		"status = ?", models.SysWorkflowLogStateSubmit, // 状态已提交
	).Find(&logs).Error
	if err != nil {
		return list, err
	}

	for _, log := range logs {
		// 获取当前待审批人
		userIds := s.getApprovingUsers(log.FlowId, log.TargetId, log.CurrentLineId, log.CurrentLine.Node)
		// 包含当前审批人
		if utils.ContainsUint(userIds, approval.Id) {
			list = append(list, log)
		}
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
		Preload("CurrentLine.Node").
		Preload("CurrentLine.Node.Users").
		Preload("CurrentLine.Node.Role").
		Preload("CurrentLine.Node.Role.Users").
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
	userIds := s.getApprovingUsers(log.FlowId, log.TargetId, log.CurrentLineId, log.CurrentLine.Node)
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
	}).Where(
		"status > ?", models.SysWorkflowLogStateSubmit, // 状态非提交
	).Find(&logs).Error
	return logs, err
}

// 查询下一流水线
func (s *MysqlService) GetNextWorkflowLine(flowId uint, currentSort uint) (models.SysWorkflowLine, error) {
	return s.GetPrevWorkflowLineBySort(flowId, currentSort+1)
}

// 查询上一流水线
func (s *MysqlService) GetPrevWorkflowLine(flowId uint, currentSort uint) (models.SysWorkflowLine, error) {
	return s.GetPrevWorkflowLineBySort(flowId, currentSort-1)
}

// 查询指定sort流水线
func (s *MysqlService) GetPrevWorkflowLineBySort(flowId uint, sort uint) (models.SysWorkflowLine, error) {
	// 查询流程线
	var line models.SysWorkflowLine
	if sort <= 0 {
		return line, gorm.ErrRecordNotFound
	}
	err := s.tx.Preload("Node").Preload("Node.Users").Where(&models.SysWorkflowLine{
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

// 获取多个节点
func (s *MysqlService) GetWorkflowNodesByIds(ids []uint) ([]models.SysWorkflowNode, error) {
	var nodes []models.SysWorkflowNode
	var err error
	err = s.tx.Preload("Users").Where("id IN (?)", ids).Find(&nodes).Error
	return nodes, err
}

// 创建工作流节点
func (s *MysqlService) CreateWorkflowNode(req *request.UpdateWorkflowNodeRequestStruct) (err error) {
	var node models.SysWorkflowNode
	if len(req.UserIds) > 0 {
		// 查询所有用户
		users := make([]models.SysUser, 0)
		err = s.tx.Where("id IN (?)", req.UserIds).Find(&users).Error
		if err != nil {
			return
		}
		node.Users = users
	}
	utils.Struct2StructByJson(req, &node)
	// 创建数据
	err = s.tx.Create(&node).Error
	// 获取id
	req.Id = node.Id
	return
}

// 更新流程流水线
func (s *MysqlService) UpdateWorkflowLineByNodes(req *request.UpdateWorkflowLineRequestStruct) (err error) {
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
	// 查询增改的所有用户/节点
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
	users, err := s.GetUsersByIds(userIds)
	if err != nil {
		return
	}

	// 获取历史节点
	nodeIds := make([]uint, 0)
	for _, line := range oldLines {
		if line.NodeId > 0 && !utils.ContainsUint(nodeIds, line.NodeId) {
			nodeIds = append(nodeIds, line.NodeId)
		}
	}
	// 删除/更新/新增需要有序, 否则可能会打破流水线顺序
	// 1. 删除节点
	for _, item := range req.Delete {
		// 删除节点
		err = s.tx.Where("id = ?", item.Id).Delete(models.SysWorkflowNode{}).Error
		if err != nil {
			return
		}
		// 删除流水线
		err = s.tx.Where("node_id = ?", item.Id).Delete(models.SysWorkflowLine{}).Error
		if err != nil {
			return
		}
		// 删除节点
		deleteIndex := utils.ContainsUintIndex(nodeIds, item.Id)
		if deleteIndex > -1 {
			if deleteIndex < len(nodeIds)-1 {
				nodeIds = append(nodeIds[:deleteIndex], nodeIds[deleteIndex+1:]...)
			} else {
				nodeIds = append(nodeIds[:deleteIndex])
			}
		}
	}
	// 2. 更新节点
	for _, item := range req.Update {
		var node models.SysWorkflowNode
		utils.Struct2StructByJson(item, &node)
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
		query := s.tx.Model(&node)
		if item.RoleId != nil {
			// 需要强制更新roleId
			query = query.Update("role_id", item.RoleId)
		}
		// 更新数据, 替换users
		err = query.Update(&node).Association("Users").Replace(&us).Error
		if err != nil {
			return
		}
	}
	// 2. 创建节点
	for _, item := range req.Create {
		var node models.SysWorkflowNode
		utils.Struct2StructByJson(item, &node)
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
		node.Users = us
		// 设置flowId
		node.FlowId = oldFlow.Id
		node.Creator = req.Creator
		// 创建数据
		err = s.tx.Create(&node).Error
		if err != nil {
			return
		}
		// 保留id
		nodeIds = append(nodeIds, node.Id)
	}
	nodes, err := s.GetWorkflowNodesByIds(nodeIds)
	if err != nil {
		return
	}

	lines := make([]models.SysWorkflowLine, 0)
	count := len(nodeIds)
	endPtr := true
	for i, id := range nodeIds {
		// 寻找当前level的节点
		var currentNode models.SysWorkflowNode
		// 所有节点必须属于同一流程号
		for _, node := range nodes {
			if node.FlowId != oldFlow.Id {
				return fmt.Errorf("第%d级节点所属流程[%s]与其他节点不一致", i+1, oldFlow.Name)
			}
			if node.Id == id {
				currentNode = node
				break
			}
		}
		// 构造流程线
		var line models.SysWorkflowLine
		needCreate := false
		if i < len(oldLines) {
			line = oldLines[i]
		} else {
			line.FlowId = oldFlow.Id
			needCreate = true
		}
		// 替换节点
		line.Node = currentNode
		// 排序, 从1开始
		line.Sort = uint(i + 1)
		// 结束标识
		if i == count-1 {
			line.End = &endPtr
		}
		if needCreate {
			// gorm 1.x不支持批量插入, 这里单条插入
			err = s.tx.Create(&line).Error
		} else {
			err = s.tx.Model(&line).Update(&line).Error
		}
		if err != nil {
			return
		}
		lines = append(lines, line)
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
	notFound := s.tx.
		Preload("CurrentLine").
		Preload("CurrentLine.Node").
		Preload("CurrentLine.Node.Users").
		Preload("CurrentLine.Node.Role").
		Preload("CurrentLine.Node.Role.Users").
		Preload("Flow").
		Where(&models.SysWorkflowLog{
			TargetId: req.TargetId,
			FlowId:   req.FlowId,
		}).Last(&lastLog).RecordNotFound()
	if notFound {
		// 走提交逻辑
		return s.first(req)
	} else {
		// 走审批逻辑
		return s.next(req, lastLog)
	}
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
	approvalStatus := models.SysWorkflowLogStateApproval
	// 状态为自己批准
	firstLog.Status = &approvalStatus
	// 当前节点为开始节点
	firstLog.SubmitUserId = submitUser.Id
	firstLog.ApprovalId = submitUser.Id
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
	// 状态为提交, 当前节点指向下一节点, 创建新日志
	err = s.newLog(models.SysWorkflowLogStateSubmit, nextLine.Id, firstLog)
	return err
}

// 第二次提交流程工单
func (s *MysqlService) next(req *request.WorkflowTransitionRequestStruct, lastLog models.SysWorkflowLog) error {
	if *lastLog.Status == models.SysWorkflowLogStateEnd {
		return fmt.Errorf("流程已结束")
	}
	if req.ApprovalId == 0 {
		return fmt.Errorf("审批人不存在, approvalId=%d", req.ApprovalId)
	}
	// 查询审批人是否存在
	approval, err := s.GetUserById(req.ApprovalId)
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
			req.SubmitUserId = req.ApprovalId
			if strings.TrimSpace(req.ApprovalOpinion) == "" {
				req.ApprovalOpinion = "提交人再次重启"
			}
			return s.first(req)
		}
	}
	var nextLine models.SysWorkflowLine
	// 开启提交人确认是没有下一节点的
	if !*lastLog.Flow.SubmitUserConfirm && !*lastLog.CurrentLine.End {
		// 判断是否末尾节点
		nextLine, err = s.GetNextWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
		if err != nil {
			return err
		}
	}
	if nextLine.NodeId == 0 {
		// 1. 下一节点为空
		// 上一条没有被拒绝, 否则不允许自己通过
		if *lastLog.Status != models.SysWorkflowLogStateDeny {
			if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
				// 结束
				return s.end(req.ApprovalOpinion, approval, lastLog)
			} else if !*lastLog.Flow.SubmitUserConfirm {
				// 回退到上一节点(提交人确认是不允许被拒绝的)
				return s.deny(req, approval, lastLog)
			}
		}
	} else if *lastLog.Flow.Self && s.checkPermission(approval.Id, lastLog) {
		// 2. 开启自我审批 且 有权限审批
		if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
			// 通过
			return s.approval(req, approval, lastLog)
		} else {
			// 回退到上一节点
			return s.deny(req, approval, lastLog)
		}
	}
	return fmt.Errorf("无权限审批或审批流程未创建")
}

// 开始正常审批
func (s *MysqlService) start(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 当前状态已提交 且 必须有权限审批
	if *lastLog.Status == models.SysWorkflowLogStateSubmit && s.checkPermission(approval.Id, lastLog) {
		if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
			// 通过
			return s.approval(req, approval, lastLog)
		} else {
			// 回退到上一节点
			return s.deny(req, approval, lastLog)
		}
	}
	return fmt.Errorf("无权限审批或审批流程未创建")
}

// 通过审批, 流转到下一节点
func (s *MysqlService) approval(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 默认节点不变
	lineId := lastLog.CurrentLineId
	if s.checkNextLineSort(approval.Id, lastLog) {
		// 流转到下一节点
		var err error
		var nextLine models.SysWorkflowLine
		if !*lastLog.CurrentLine.End {
			// 获取下一节点
			nextLine, err = s.GetNextWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
			if err != nil {
				return err
			}
		}
		// 下一节点为空, 直接结束
		if nextLine.NodeId == 0 {
			return s.end(req.ApprovalOpinion, approval, lastLog)
		}
		// 更新日志
		err = s.updateLog(models.SysWorkflowLogStateApproval, req.ApprovalOpinion, approval, lastLog)
		if err != nil {
			return err
		}
		// 当前节点指向下一节点
		lineId = nextLine.Id
	} else {
		// 保留当前节点, 还需其他人继续审批
		err := s.updateLog(models.SysWorkflowLogStateApproval, req.ApprovalOpinion, approval, lastLog)
		if err != nil {
			return err
		}
	}
	// 状态为提交, 创建新日志
	err := s.newLog(models.SysWorkflowLogStateSubmit, lineId, lastLog)
	return err
}

// 拒绝审批, 回退到上一节点
func (s *MysqlService) deny(req *request.WorkflowTransitionRequestStruct, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 流转到上一节点
	// 获取上一节点
	if lastLog.CurrentLine.Sort <= 1 {
		// 上一节点不存在, 说明拒绝到最初提交状态
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
	// 状态为提交,当前节点指向上一节点, 创建新日志
	err = s.newLog(models.SysWorkflowLogStateSubmit, prevLine.Id, lastLog)
	return err
}

// 结束审批, 末尾节点
func (s *MysqlService) end(approvalOpinion string, approval models.SysUser, lastLog models.SysWorkflowLog) error {
	// 结束
	status := models.SysWorkflowLogStateEnd
	if *lastLog.Flow.SubmitUserConfirm && approval.Id != lastLog.SubmitUserId {
		// 开启了提交人确认则不能直接结束
		status = models.SysWorkflowLogStateApproval
	}
	// 提交人确认节点
	if approval.Id == lastLog.SubmitUserId && strings.TrimSpace(approvalOpinion) == "" {
		approvalOpinion = "提交人已确认"
	}
	err := s.updateLog(status, approvalOpinion, approval, lastLog)
	if err != nil {
		return err
	}
	if status == models.SysWorkflowLogStateEnd {
		return nil
	}
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
	// 提交人
	updateLog.SubmitUserId = lastLog.SubmitUserId
	// 审批人以及意见
	updateLog.ApprovalId = approval.Id
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
	// 提交人
	newLog.SubmitUserId = lastLog.SubmitUserId
	// 创建数据
	err := s.tx.Create(&newLog).Error
	return err
}

// 检查当前审批人是否有权限
func (s *MysqlService) checkPermission(approvalId uint, lastLog models.SysWorkflowLog) bool {
	// 获取当前待审批人
	userIds := s.getApprovingUsers(lastLog.FlowId, lastLog.TargetId, lastLog.CurrentLineId, lastLog.CurrentLine.Node)
	return utils.ContainsUint(userIds, approvalId)
}

// 检查是否可以切换流水线到下一个(通过审批会使用)
func (s *MysqlService) checkNextLineSort(approvalId uint, lastLog models.SysWorkflowLog) bool {
	// 获取当前待审批人
	userIds := s.getApprovingUsers(lastLog.FlowId, lastLog.TargetId, lastLog.CurrentLineId, lastLog.CurrentLine.Node)
	// 判断流程类别
	switch lastLog.Flow.Category {
	case models.SysWorkflowCategoryOnlyOneApproval:
		// 只需要1人通过: 当前审批人在待审批列表中
		return utils.ContainsUint(userIds, approvalId)
	case models.SysWorkflowCategoryAllApproval:
		// 查询全部审批人数
		allUserIds := s.getAllApprovalUsers(lastLog.CurrentLine.Node)
		// 查询历史审批人数
		historyUserIds := s.getHistoryApprovalUsers(lastLog.FlowId, lastLog.TargetId, lastLog.CurrentLineId)
		// 需要全部人通过: 当前审批人在待审批列表中 且 历史审批人+当前审批人刚好等于全部审批人
		return utils.ContainsUint(userIds, approvalId) && len(historyUserIds) >= len(allUserIds)-1
	}
	return false
}

// 获取待审批人(当前节点)
func (s *MysqlService) getApprovingUsers(flowId uint, targetId uint, currentLineId uint, currentNode models.SysWorkflowNode) []uint {
	userIds := make([]uint, 0)
	allUserIds := s.getAllApprovalUsers(currentNode)
	historyUserIds := s.getHistoryApprovalUsers(flowId, targetId, currentLineId)
	for _, allUserId := range allUserIds {
		// 不在历史列表中
		if !utils.ContainsUint(historyUserIds, allUserId) {
			userIds = append(userIds, allUserId)
		}
	}
	return userIds
}

// 获取历史审批人(最后一个节点, 主要用于判断是否审批完成)
func (s *MysqlService) getHistoryApprovalUsers(flowId uint, targetId uint, currentLineId uint) []uint {
	historyUserIds := make([]uint, 0)
	// 查询已审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	err := s.tx.Preload("ApprovalUser").Where(&models.SysWorkflowLog{
		FlowId:   flowId,   // 流程号一致
		TargetId: targetId, // 目标一致
	}).Where(
		"status > ?", models.SysWorkflowLogStateSubmit, // 状态非提交
	).Order("id DESC").Find(&logs).Error
	if err != nil {
		global.Log.Warn("[getHistoryApprovalUsers]", err)
	}
	// 保留连续审核通过记录
	l := len(logs)
	for i := 0; i < l; i++ {
		log := logs[i]
		// 如果不是通过立即结束, 必须保证连续的通过 或 当前节点不一致
		if *log.Status != models.SysWorkflowLogStateApproval || log.CurrentLineId != currentLineId {
			break
		}
		// 审批人为配置中的一人
		if !utils.ContainsUint(historyUserIds, log.ApprovalId) {
			historyUserIds = append(historyUserIds, log.ApprovalId)
		}
	}
	return historyUserIds
}

// 获取全部审批人(当前节点)
func (s *MysqlService) getAllApprovalUsers(currentNode models.SysWorkflowNode) []uint {
	userIds := make([]uint, 0)
	for _, user := range currentNode.Users {
		if user.Id > 0 && !utils.ContainsUint(userIds, user.Id) {
			userIds = append(userIds, user.Id)
		}
	}
	if currentNode.RoleId > 0 {
		for _, user := range currentNode.Role.Users {
			if user.Id > 0 && !utils.ContainsUint(userIds, user.Id) {
				userIds = append(userIds, user.Id)
			}
		}
	}
	return userIds
}
