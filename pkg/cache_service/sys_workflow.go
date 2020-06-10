package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"strings"
)

// 获取所有工作流
func (s *RedisService) GetWorkflows(req *request.WorkflowListRequestStruct) ([]models.SysWorkflow, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetWorkflows(req)
	}
	var err error
	list := make([]models.SysWorkflow, 0)
	// 查询接口表所有缓存
	jsonWorkflows := s.GetListFromCache(nil, new(models.SysWorkflow).TableName())
	query := s.JsonQuery().FromString(jsonWorkflows)
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name", "contains", name)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Category > 0 {
		// redis存的json转换为int, 因此这里转一下类型
		query = query.Where("category", "=", int(req.Category))
	}
	if req.TargetCategory > 0 {
		// redis存的json转换为int, 因此这里转一下类型
		query = query.Where("targetCategory", "=", int(req.TargetCategory))
	}
	if req.Self != nil {
		query = query.Where("self", "=", *req.Self)
	}
	if req.SubmitUserConfirm != nil {
		query = query.Where("submitUserConfirm", "=", *req.SubmitUserConfirm)
	}

	// 查询条数
	req.PageInfo.Total = uint(query.Count())
	var res interface{}
	if req.PageInfo.NoPagination {
		// 不使用分页
		res = query.Get()
	} else {
		// 获取分页参数
		limit, offset := req.GetLimit()
		res = query.Limit(int(limit)).Offset(int(offset)).Get()
	}
	// 转换为结构体
	utils.Struct2StructByJson(res, &list)
	return list, err
}

// 获取所有流水线
func (s *RedisService) GetWorkflowLines(req *request.WorkflowLineListRequestStruct) ([]models.SysWorkflowLine, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetWorkflowLines(req)
	}
	var err error
	list := make([]models.SysWorkflowLine, 0)
	// 查询接口表所有缓存
	jsonWorkflowLines := s.GetListFromCache(nil, new(models.SysWorkflowLine).TableName())
	query := s.JsonQuery().FromString(jsonWorkflowLines)
	if req.FlowId > 0 {
		// redis存的json转换为int, 因此这里转一下类型
		query = query.Where("flowId", "=", int(req.FlowId))
	}

	// 查询条数
	req.PageInfo.Total = uint(query.Count())
	var res interface{}
	if req.PageInfo.NoPagination {
		// 不使用分页
		res = query.Get()
	} else {
		// 获取分页参数
		limit, offset := req.GetLimit()
		res = query.Limit(int(limit)).Offset(int(offset)).Get()
	}
	// 转换为结构体
	utils.Struct2StructByJson(res, &list)
	// 查询所有节点
	nodeList := make([]models.SysWorkflowNode, 0)
	jsonWorkflowNodes := s.GetListFromCache(nil, new(models.SysWorkflowNode).TableName())
	nodeRes := s.JsonQuery().FromString(jsonWorkflowNodes).Get()
	utils.Struct2StructByJson(nodeRes, &nodeList)
	for i, line := range list {
		var node models.SysWorkflowNode
		// 查找节点
		for _, item := range nodeList {
			if item.Id == line.NodeId {
				node = item
				break
			}
		}
		// 查找节点用户
		node.Users = s.getWorkflowUsersByNodeId(line.NodeId)
		list[i].Node = node
	}
	return list, err
}

// 获取工作流(指定审批单类型)
func (s *RedisService) GetWorkflowByTargetCategory(targetCategory uint) (models.SysWorkflow, error) {
	var flow models.SysWorkflow
	// 查询第一条记录即可
	// id在JSONQ中以int存在
	jsonWorkflows := s.GetListFromCache(nil, new(models.SysWorkflow).TableName())
	res, err := JsonQueryFindOne(s.JsonQuery().FromString(jsonWorkflows).Where("targetCategory", "=", int(targetCategory)))
	if err != nil {
		return flow, err
	}
	utils.Struct2StructByJson(res, &flow)
	return flow, err
}

// 查询审批日志(指定目标)
func (s *RedisService) GetWorkflowLogs(flowId uint, targetId uint) ([]models.SysWorkflowLog, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetWorkflowLogs(flowId, targetId)
	}
	var err error
	// 查询所有用户
	userList := make([]models.SysUser, 0)
	jsonWorkflowUsers := s.GetListFromCache(nil, new(models.SysUser).TableName())
	userRes := s.JsonQuery().FromString(jsonWorkflowUsers).Get()
	utils.Struct2StructByJson(userRes, &userList)

	// 查询所有工作流
	workflowList := make([]models.SysWorkflow, 0)
	jsonWorkflows := s.GetListFromCache(nil, new(models.SysWorkflow).TableName())
	workflowRes := s.JsonQuery().FromString(jsonWorkflows).Get()
	utils.Struct2StructByJson(workflowRes, &workflowList)

	// 查询所有工作流日志
	workflowLogList := make([]models.SysWorkflowLog, 0)
	jsonWorkflowWorkflowLogs := s.GetListFromCache(nil, new(models.SysWorkflowLog).TableName())
	workflowLogRes := s.JsonQuery().
		FromString(jsonWorkflowWorkflowLogs).
		// 流程号一致
		Where("flowId", "=", int(flowId)).
		// 目标一致
		Where("targetId", "=", int(targetId)).
		Get()
	utils.Struct2StructByJson(workflowLogRes, &workflowLogList)

	newLogs := make([]models.SysWorkflowLog, 0)
	for _, workflowLog := range workflowLogList {
		newLog := workflowLog
		// 查找审批人/提交人
		for _, user := range userList {
			if workflowLog.ApprovalUserId == user.Id {
				newLog.ApprovalUser = user
			}
			if workflowLog.SubmitUserId == user.Id {
				newLog.SubmitUser = user
			}
		}
		// 查找工作流
		for _, flow := range workflowList {
			if workflowLog.FlowId == flow.Id {
				newLog.Flow = flow
			}
		}
		newLogs = append(newLogs, newLog)
	}
	return newLogs, err
}

// 查询待审批目标列表(指定用户)
func (s *RedisService) GetWorkflowApprovings(req *request.WorkflowApprovingListRequestStruct) ([]models.SysWorkflowLog, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetWorkflowApprovings(req)
	}
	list := make([]models.SysWorkflowLog, 0)
	if req.ApprovalUserId == 0 {
		return list, fmt.Errorf("用户不存在, approvalUserId=%d", req.ApprovalUserId)
	}
	// 查询审批人
	approval, err := s.GetUserById(req.ApprovalUserId)
	if err != nil {
		return list, err
	}
	// 查询所有日志
	workflowLogList := make([]models.SysWorkflowLog, 0)
	jsonWorkflowLogs := s.GetListFromCache(nil, new(models.SysWorkflowLog).TableName())
	// 由于还需判断是否包含当前审批人, 因此无法直接分页
	workflowLogRes := s.JsonQuery().FromString(jsonWorkflowLogs).
		Where("status", "=", int(models.SysWorkflowLogStateSubmit)). // 状态已提交
		Get()
	utils.Struct2StructByJson(workflowLogRes, &workflowLogList)

	// 查询所有工作流
	workflowList := make([]models.SysWorkflow, 0)
	jsonWorkflows := s.GetListFromCache(nil, new(models.SysWorkflow).TableName())
	workflowRes := s.JsonQuery().FromString(jsonWorkflows).Get()
	utils.Struct2StructByJson(workflowRes, &workflowList)

	// 查询所有流水线
	workflowLineList := make([]models.SysWorkflowLine, 0)
	jsonWorkflowLines := s.GetListFromCache(nil, new(models.SysWorkflowLine).TableName())
	workflowLineRes := s.JsonQuery().FromString(jsonWorkflowLines).Get()
	utils.Struct2StructByJson(workflowLineRes, &workflowLineList)

	// 查询所有节点
	workflowNodeList := make([]models.SysWorkflowNode, 0)
	jsonWorkflowNodes := s.GetListFromCache(nil, new(models.SysWorkflowNode).TableName())
	workflowNodeRes := s.JsonQuery().FromString(jsonWorkflowNodes).Get()
	utils.Struct2StructByJson(workflowNodeRes, &workflowNodeList)

	// 查询所有角色
	roleList := make([]models.SysRole, 0)
	jsonRoles := s.GetListFromCache(nil, new(models.SysRole).TableName())
	roleRes := s.JsonQuery().FromString(jsonRoles).Get()
	utils.Struct2StructByJson(roleRes, &roleList)

	// 查询所有用户
	userList := make([]models.SysUser, 0)
	jsonUsers := s.GetListFromCache(nil, new(models.SysUser).TableName())
	userRes := s.JsonQuery().FromString(jsonUsers).Get()
	utils.Struct2StructByJson(userRes, &userList)

	// 加载日志的关联对象
	for i, log := range workflowLogList {
		if log.CurrentLineId > 0 {
			for _, line := range workflowLineList {
				if line.Id == log.CurrentLineId {
					// 查找节点
					var currentNode models.SysWorkflowNode
					for _, node := range workflowNodeList {
						if node.Id == line.NodeId {
							// 查找节点角色
							for _, role := range roleList {
								if node.RoleId == role.Id {
									// 加载角色关联的用户
									users := make([]models.SysUser, 0)
									for _, user := range userList {
										if role.Id == user.RoleId {
											users = append(users, user)
											break
										}
									}
									role.Users = users
									node.Role = role
									break
								}
							}
							currentNode = node
							break
						}
					}
					// 查找节点用户
					currentNode.Users = s.getWorkflowUsersByNodeId(line.NodeId)
					line.Node = currentNode
					log.CurrentLine = line
					break
				}
			}
		}
		// 查找提交人
		if log.SubmitUserId > 0 {
			for _, user := range userList {
				if log.SubmitUserId == user.Id {
					log.SubmitUser = user
					break
				}
			}
		}
		// 查找审批人
		if log.ApprovalUserId > 0 {
			for _, user := range userList {
				if log.ApprovalUserId == user.Id {
					log.ApprovalUser = user
					break
				}
			}
		}
		// 查找工作流
		if log.FlowId > 0 {
			for _, flow := range workflowList {
				if log.FlowId == flow.Id {
					log.Flow = flow
					break
				}
			}
		}
		workflowLogList[i] = log
	}

	for _, log := range workflowLogList {
		// 获取当前待审批人
		userIds := s.getApprovingUsers(log.FlowId, log.TargetId, log.CurrentLineId, log.CurrentLine.Node)
		log.ApprovingUserIds = userIds
		// 包含当前审批人
		if utils.ContainsUint(userIds, approval.Id) {
			list = append(list, log)
		}
	}
	// 处理分页(转为json)
	query := s.JsonQuery().FromString(utils.Struct2Json(list))
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

// 获取所有用户(根据节点编号)
func (s *RedisService) getWorkflowUsersByNodeId(nodeId uint) []models.SysUser {
	// 查询所有用户
	userList := make([]models.SysUser, 0)
	jsonWorkflowUsers := s.GetListFromCache(nil, new(models.SysUser).TableName())
	userRes := s.JsonQuery().FromString(jsonWorkflowUsers).Get()
	utils.Struct2StructByJson(userRes, &userList)
	// 查询所有关联关系
	relationList := make([]models.RelationUserWorkflowNode, 0)
	// 查询所有节点
	jsonWorkflowRelations := s.GetListFromCache(nil, new(models.RelationUserWorkflowNode).TableName())
	relationRes := s.JsonQuery().FromString(jsonWorkflowRelations).Get()
	utils.Struct2StructByJson(relationRes, &relationList)
	users := make([]models.SysUser, 0)
	for _, relation := range relationList {
		for _, user := range userList {
			if nodeId == relation.SysWorkflowNodeId && user.Id == relation.SysUserId {
				users = append(users, user)
				break
			}
		}
	}
	return users
}

// 获取待审批人(当前节点)
func (s *RedisService) getApprovingUsers(flowId uint, targetId uint, currentLineId uint, currentNode models.SysWorkflowNode) []uint {
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
func (s *RedisService) getHistoryApprovalUsers(flowId uint, targetId uint, currentLineId uint) []uint {
	historyUserIds := make([]uint, 0)
	// 查询已审核的日志
	logs := make([]models.SysWorkflowLog, 0)
	jsonWorkflowLogs := s.GetListFromCache(nil, new(models.SysWorkflowLog).TableName())
	workflowLogRes := s.JsonQuery().FromString(jsonWorkflowLogs).
		Where("flowId", "=", int(flowId)).
		Where("targetId", "=", int(targetId)).
		Where("status", ">", models.SysWorkflowLogStateSubmit). // 状态非提交
		Get()
	utils.Struct2StructByJson(workflowLogRes, &logs)

	// 保留连续审核通过记录
	l := len(logs)
	for i := 0; i < l; i++ {
		log := logs[i]
		// 如果不是通过立即结束, 必须保证连续的通过 或 当前节点不一致
		if *log.Status != models.SysWorkflowLogStateApproval || log.CurrentLineId != currentLineId {
			break
		}
		// 审批人为配置中的一人
		if !utils.ContainsUint(historyUserIds, log.ApprovalUserId) {
			historyUserIds = append(historyUserIds, log.ApprovalUserId)
		}
	}
	return historyUserIds
}

// 获取全部审批人(当前节点)
func (s *RedisService) getAllApprovalUsers(currentNode models.SysWorkflowNode) []uint {
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
