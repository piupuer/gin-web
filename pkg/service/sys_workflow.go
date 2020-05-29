package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"strings"
)

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
	err := s.tx.Preload("Nodes").Preload("Nodes.Users").Where(&models.SysWorkflowLine{
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

// 创建工作流节点
func (s *MysqlService) CreateWorkflowNode(req *request.CreateWorkflowNodeRequestStruct) (err error) {
	var node models.SysWorkflowNode
	utils.Struct2StructByJson(req, &node)
	// 创建数据
	err = s.tx.Create(&node).Error
	return
}

// 创建流程流水线, 传入多个节点, 二维数组, 外层表示审批层级(顺序排列), 内层表示同级存在多个子节点
func (s *MysqlService) CreateWorkflowLine(nodeIds [][]uint) (err error) {
	// 拆分所有id
	ids := make([]uint, 0)
	for i, item := range nodeIds {
		if len(item) == 0 {
			return fmt.Errorf("第%d个节点为空", i)
		}
		for _, id := range item {
			if !utils.ContainsUint(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	if len(ids) == 0 {
		return fmt.Errorf("节点至少有一个")
	}
	// 查询所有节点明细
	nodes := make([]models.SysWorkflowNode, 0)
	err = s.tx.Preload("Flow").Where("id IN (?)", ids).Find(&nodes).Error
	if err != nil {
		return
	}
	if len(nodes) == 0 {
		return gorm.ErrRecordNotFound
	}
	flowId := nodes[0].FlowId
	flow := nodes[0].Flow
	// 一个流程只能有一条线
	var oldLine models.SysWorkflowLine
	oldLine.FlowId = flowId
	notFound := s.tx.Where(&oldLine).First(&oldLine).RecordNotFound()
	if !notFound {
		return fmt.Errorf("流程%s已存在一条流水线, 无法再创建新的, flowId=%d, line.id=%d", flow.Name, flowId, oldLine.Id)
	}
	lines := make([]models.SysWorkflowLine, 0)
	count := len(nodeIds)
	endPtr := true
	for i, item := range nodeIds {
		// 寻找当前level的所有节点
		currentNodes := make([]models.SysWorkflowNode, 0)
		for j, id := range item {
			for _, node := range nodes {
				// 所有节点必须属于同一流程号
				if node.FlowId != flowId {
					return fmt.Errorf("第%d级第%d个节点所属流程[%s]与其他节点不一致", i+1, j+1, flow.Name)
				}
				if node.Id == id {
					currentNodes = append(currentNodes, node)
				}
			}
		}
		// 构造流程线
		var line models.SysWorkflowLine
		// 流程号
		line.FlowId = flowId
		// 节点
		line.Nodes = currentNodes
		// 排序, 从1开始
		line.Sort = uint(i + 1)
		// 结束标识
		if i == count-1 {
			line.End = &endPtr
		}
		// gorm 1.x不支持批量插入, 这里单条插入
		s.tx.Create(&line)
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
	notFound := s.tx.Preload("CurrentLine").Preload("CurrentLine.Nodes").Preload("CurrentLine.Nodes.Users").Preload("Flow").Where(&models.SysWorkflowLog{
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
	if req.SubmitterId == 0 {
		return fmt.Errorf("提交人不存在, submitterId=%d", req.SubmitterId)
	}
	// 查询提交人是否存在
	submitUser, err := s.GetUserById(req.SubmitterId)
	if err != nil {
		return err
	}

	// 初次创建
	var firstLog models.SysWorkflowLog
	firstLog.FlowId = req.FlowId
	firstLog.TargetId = req.TargetId
	// 状态为自己批准
	firstLog.Status = &models.SysWorkflowLogStateApproval
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
			req.SubmitterId = req.ApprovalId
			if strings.TrimSpace(req.ApprovalOpinion) == "" {
				req.ApprovalOpinion = "提交人再次重启"
			}
			return s.first(req)
		}
	}
	var nextLine models.SysWorkflowLine
	// 开启提交人确认是没有下一节点的
	if !*lastLog.Flow.SubmitterConfirm && !*lastLog.CurrentLine.End {
		// 判断是否末尾节点
		nextLine, err = s.GetNextWorkflowLine(req.FlowId, lastLog.CurrentLine.Sort)
		if err != nil {
			return err
		}
	}
	if nextLine.Nodes == nil {
		// 1. 下一节点为空
		// 上一条没有被拒绝, 否则不允许自己通过
		if *lastLog.Status != models.SysWorkflowLogStateDeny {
			if *req.ApprovalStatus == models.SysWorkflowLogStateApproval {
				// 结束
				return s.end(req.ApprovalOpinion, approval, lastLog)
			} else if !*lastLog.Flow.SubmitterConfirm {
				// 回退到上一节点(提交人确认是不允许被拒绝的)
				return s.deny(req, approval, lastLog)
			}
		}
	} else if *lastLog.Flow.Self && s.checkPermission(approval, lastLog.CurrentLine.Nodes) {
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
	if *lastLog.Status == models.SysWorkflowLogStateSubmit && s.checkPermission(approval, lastLog.CurrentLine.Nodes) {
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
	if s.checkNextLineSort(lastLog) {
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
		if nextLine.Nodes == nil {
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
	if *lastLog.Flow.SubmitterConfirm && approval.Id != lastLog.SubmitUserId {
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
func (s *MysqlService) checkPermission(approval models.SysUser, nodes []models.SysWorkflowNode) bool {
	// 获取当前节点审批人
	userIds, roleIds := s.getApprovalUsers(nodes)
	// 配置了用户, 优先以用户为准
	if len(userIds) > 0 {
		return utils.ContainsUint(userIds, approval.Id)
	}
	return utils.ContainsUint(roleIds, approval.RoleId)
}

// 检查是否可以切换流水线到下一个(通过审批会使用)
func (s *MysqlService) checkNextLineSort(lastLog models.SysWorkflowLog) bool {
	// 判断流程类别
	switch lastLog.Flow.Category {
	case models.SysWorkflowCategoryOnlyOneApproval:
		// 只需要1人通过
		return true
	case models.SysWorkflowCategoryAllApproval:
		// 需要全部人通过
		// 查询当前流水线当前节点总共配置了多少人审核
		userIds, roleIds := s.getApprovalUsers(lastLog.CurrentLine.Nodes)
		// 查询已审核的日志
		logs := make([]models.SysWorkflowLog, 0)
		err := s.tx.Preload("ApprovalUser").Where(&models.SysWorkflowLog{
			FlowId:   lastLog.FlowId,   // 流程号一致
			TargetId: lastLog.TargetId, // 目标一致
		}).Where(
			"status > ?", models.SysWorkflowLogStateSubmit, // 状态非提交
		).Order("id DESC").Find(&logs).Error
		if err != nil {
			return false
		}
		// 保留连续审核通过记录
		l := len(logs)
		historyUserIds := make([]uint, 0)
		for i := 0; i < l; i++ {
			log := logs[i]
			// 如果不是通过立即结束, 必须保证连续的通过
			if *log.Status != models.SysWorkflowLogStateApproval {
				break
			}
			// 审批人为配置中的一人
			if utils.ContainsUint(userIds, log.ApprovalId) && !utils.ContainsUint(historyUserIds, log.ApprovalId) {
				historyUserIds = append(historyUserIds, log.ApprovalId)
			}
		}
		// 配置了用户, 优先以用户为准(checkNextLineSort方法在更新审批前, 因此-1)
		if len(userIds) > 0 {
			return len(historyUserIds) >= len(userIds)-1
		}
		return len(historyUserIds) >= len(roleIds)-1
	}
	return false
}

// 从节点中获取审批人
func (s *MysqlService) getApprovalUsers(nodes []models.SysWorkflowNode) ([]uint, []uint) {
	userIds := make([]uint, 0)
	roleIds := make([]uint, 0)
	for _, node := range nodes {
		for _, user := range node.Users {
			if user.Id > 0 && !utils.ContainsUint(userIds, user.Id) {
				userIds = append(userIds, user.Id)
			}
		}
		if node.RoleId > 0 && !utils.ContainsUint(roleIds, node.RoleId) {
			roleIds = append(roleIds, node.RoleId)
		}
	}
	return userIds, roleIds
}

// // 工作流转移(从一个状态转移到另一个状态)
// func (s *MysqlService) WorkflowTransition(req *request.WorkflowTransitionRequestStruct) error {
// 	if req.FlowId == 0 {
// 		return fmt.Errorf("流程号不存在, flowId=%d, 请检查参数", req.FlowId)
// 	}
// 	// 查询最后一条审批日志
// 	var lastLog models.SysWorkflowLog
// 	var newLog models.SysWorkflowLog
// 	notFound := s.tx.Preload("CurrentNode").Preload("CurrentNode.Users").Preload("Flow").Where(&models.SysWorkflowLog{TargetId: req.TargetId, FlowId: req.FlowId}).Last(&lastLog).RecordNotFound()
// 	if notFound {
// 		if req.SubmitterId == 0 {
// 			return fmt.Errorf("提交人不存在, submitterId=%d, 请检查参数", req.SubmitterId)
// 		}
// 		// 查询提交人是否存在
// 		submitter, err := s.GetUserById(req.SubmitterId)
// 		if err != nil {
// 			return err
// 		}
//
// 		// 获取下一节点, 当前节点为开始节点=0
// 		nextNode, err := s.GetNextWorkflowLine(req.FlowId, 0)
// 		if err != nil {
// 			return err
// 		}
//
// 		// 初次创建
// 		var firstLog models.SysWorkflowLog
// 		firstLog.FlowId = req.FlowId
// 		firstLog.TargetId = req.TargetId
// 		// 状态为自己批准
// 		firstLog.Status = &models.SysWorkflowLogStateApproval
// 		// 当前节点为开始节点
// 		// lastLog.CurrentNodeId = 0
// 		firstLog.SubmitterId = submitter.Id
// 		firstLog.SubmitterName = submitter.Nickname
// 		firstLog.ApprovalId = submitter.Id
// 		firstLog.ApprovalName = submitter.Nickname
// 		firstLog.ApprovalOpinion = "提交"
// 		// 创建首条日志
// 		s.tx.Create(&firstLog)
// 		// 创建下一日志
// 		newLog.FlowId = req.FlowId
// 		newLog.TargetId = req.TargetId
// 		// 状态为提交
// 		newLog.Status = &models.SysWorkflowLogStateSubmit
// 		newLog.SubmitterId = submitter.Id
// 		newLog.SubmitterName = submitter.Nickname
// 		// 当前节点指向下一节点
// 		newLog.CurrentNodeId = &nextNode.Id
// 	} else {
// 		if req.ApprovalId == 0 {
// 			return fmt.Errorf("审批人不存在, approvalId=%d, 请检查参数", req.ApprovalId)
// 		}
// 		if *lastLog.Status == models.SysWorkflowLogStateEnd {
// 			return fmt.Errorf("流程已结束")
// 		}
// 		// 查询审批人是否存在
// 		start, err := s.GetUserById(req.ApprovalId)
// 		if err != nil {
// 			return err
// 		}
// 		if lastLog.SubmitterId == start.Id {
// 			var updateLog models.SysWorkflowLog
// 			// 自我取消
// 			if *req.ApprovalStatus == models.SysWorkflowLogStateCancel {
// 				// 记录审批人以及审批意见
// 				updateLog.FlowId = req.FlowId
// 				updateLog.TargetId = req.TargetId
// 				// 状态: 取消
// 				updateLog.Status = req.ApprovalStatus
// 				updateLog.SubmitterId = lastLog.SubmitterId
// 				updateLog.SubmitterName = lastLog.SubmitterName
// 				// 取消: 审批人是自己
// 				updateLog.ApprovalId = lastLog.SubmitterId
// 				updateLog.ApprovalName = lastLog.SubmitterName
// 				approvalOpinion := req.ApprovalOpinion
// 				if strings.TrimSpace(approvalOpinion) == "" {
// 					approvalOpinion = "提交人主动取消"
// 				}
// 				updateLog.ApprovalOpinion = approvalOpinion
// 				err = s.tx.Table(updateLog.TableName()).Where("id = ?", lastLog.Id).Update(&updateLog).Error
// 				return err
// 			}
// 			// 自我重启
// 			if *req.ApprovalStatus == models.SysWorkflowLogStateRestart {
// 				if *lastLog.Status != models.SysWorkflowLogStateCancel {
// 					return fmt.Errorf("流程正常, 无需重启")
// 				}
// 				// 获取下一节点, 当前节点为开始节点=0
// 				nextNode, err := s.GetNextWorkflowLine(req.FlowId, 0)
// 				if err != nil {
// 					return err
// 				}
// 				// 流程从头开始
// 				var firstLog models.SysWorkflowLog
// 				firstLog.FlowId = req.FlowId
// 				firstLog.TargetId = req.TargetId
// 				// 状态为自己批准
// 				firstLog.Status = &models.SysWorkflowLogStateApproval
// 				firstLog.SubmitterId = lastLog.SubmitterId
// 				firstLog.SubmitterName = lastLog.SubmitterName
// 				firstLog.ApprovalId = lastLog.SubmitterId
// 				firstLog.ApprovalName = lastLog.SubmitterName
// 				firstLog.ApprovalOpinion = "重新提交"
// 				// 创建首条日志
// 				s.tx.Create(&firstLog)
// 				// 创建下一日志
// 				newLog.FlowId = req.FlowId
// 				newLog.TargetId = req.TargetId
// 				// 状态为提交
// 				newLog.Status = &models.SysWorkflowLogStateSubmit
// 				newLog.SubmitterId = lastLog.SubmitterId
// 				newLog.SubmitterName = lastLog.SubmitterName
// 				// 当前节点指向下一节点
// 				newLog.CurrentNodeId = &nextNode.Id
// 				// 创建新记录
// 				err = s.tx.Create(&newLog).Error
// 				return err
// 			}
// 		} else {
// 			if *lastLog.Status == models.SysWorkflowLogStateCancel {
// 				return fmt.Errorf("流程已被提交人主动取消")
// 			}
// 		}
// 		// 是否有权限审批
// 		if !checkNodePermission(start, lastLog) {
// 			return fmt.Errorf("无权限审批或审批流程未创建")
// 		}
//
// 		// 查询第一条审批日志
// 		var firstLog models.SysWorkflowLog
// 		notFound = s.tx.Where(&models.SysWorkflowLog{TargetId: req.TargetId}).First(&firstLog).RecordNotFound()
// 		if notFound {
// 			return fmt.Errorf("第一条流程日志不存在, targetId=%d, 请检查参数", req.TargetId)
// 		}
//
// 		end := *lastLog.CurrentNodeIsEnd && *req.ApprovalStatus == models.SysWorkflowLogStateApproval
// 		// 更新日志
// 		var status *uint
// 		if req.ApprovalStatus == nil {
// 			// 默认批准
// 			status = &models.SysWorkflowLogStateApproval
// 		} else {
// 			status = req.ApprovalStatus
// 		}
// 		if end {
// 			// 结束节点, 批准
// 			*status = models.SysWorkflowLogStateEnd
// 		}
// 		var updateLog models.SysWorkflowLog
// 		// 记录审批人以及审批意见
// 		updateLog.FlowId = req.FlowId
// 		updateLog.TargetId = req.TargetId
// 		// 状态为参数给定的状态
// 		updateLog.Status = status
// 		updateLog.SubmitterId = firstLog.SubmitterId
// 		updateLog.SubmitterName = firstLog.SubmitterName
// 		updateLog.ApprovalId = start.Id
// 		updateLog.ApprovalName = start.Nickname
// 		approvalOpinion := req.ApprovalOpinion
// 		if end && strings.TrimSpace(approvalOpinion) == "" {
// 			approvalOpinion = "提交人已确认"
// 		}
// 		updateLog.ApprovalOpinion = approvalOpinion
// 		err = s.tx.Table(updateLog.TableName()).Where("id = ?", lastLog.Id).Update(&updateLog).Error
// 		if err != nil {
// 			return err
// 		}
// 		if end {
// 			return nil
// 		}
// 		if *lastLog.CurrentNodeIsEnd {
//
// 		}
// 		if *status == models.SysWorkflowLogStateDeny {
// 			// 拒绝
// 			newLog.Status = &models.SysWorkflowLogStateSubmit
// 			// 当前节点指向上一节点
// 			// 获取上一节点
// 			prevNode, err := s.GetPrevWorkflowLine(req.FlowId, *lastLog.CurrentNodeId)
// 			if err != nil {
// 				return err
// 			}
// 			newLog.CurrentNodeId = &prevNode.Id
// 		} else {
// 			// 下一节点
// 			var nextNode models.SysWorkflowNode
// 			if !*lastLog.CurrentNodeIsEnd {
// 				// 非末尾节点则获取下一节点
// 				nextNode, err = s.GetNextWorkflowLine(req.FlowId, *lastLog.CurrentNodeId)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			// 同意
// 			if nextNode.Id == 0 {
// 				// 下一节点为空
// 				if *lastLog.Flow.SubmitterConfirm {
// 					// 需要提交人确认
// 					newLog.Status = &models.SysWorkflowLogStateSubmit
// 					newLog.CurrentNodeIsEnd = lastLog.Flow.SubmitterConfirm
// 				} else {
// 					// 不需要提交人确认, 结束
// 					newLog.Status = &models.SysWorkflowLogStateEnd
// 				}
// 			} else {
// 				// 当前节点指向下一节点
// 				newLog.CurrentNodeId = &nextNode.Id
// 			}
// 		}
// 		// 创建新日志
// 		newLog.FlowId = req.FlowId
// 		newLog.TargetId = req.TargetId
// 		// 状态为提交
// 		newLog.SubmitterId = firstLog.SubmitterId
// 		newLog.SubmitterName = firstLog.SubmitterName
// 	}
// 	// 创建新记录
// 	err := s.tx.Create(&newLog).Error
// 	return err
// }
//
// // 校验节点是否有权限审批
// func checkNodePermission(start models.SysUser, lastLog models.SysWorkflowLog) bool {
// 	if !*lastLog.CurrentNodeIsEnd {
// 		// 不是末尾节点
// 		// 以用户优先
// 		if len(lastLog.CurrentNode.Users) > 0 {
// 			userIds := make([]uint, 0)
// 			for _, user := range lastLog.CurrentNode.Users {
// 				userIds = append(userIds, user.Id)
// 			}
// 			if utils.Contains(userIds, start.Id) {
// 				// 审批人在列表中
// 				if lastLog.SubmitterId == start.Id {
// 					// 提交人也是审批人, 需要判断是否开启自我审批
// 					return *lastLog.Flow.Self
// 				}
// 				return true
// 			}
// 		} else if lastLog.CurrentNode.RoleId == start.RoleId {
// 			if lastLog.SubmitterId == start.Id {
// 				// 提交人也是审批人, 需要判断是否开启自我审批
// 				return *lastLog.Flow.Self
// 			}
// 			return true
// 		}
// 	} else {
// 		// 末尾节点需要自己审批
// 		return lastLog.SubmitterId == start.Id
// 	}
// 	return false
// }
