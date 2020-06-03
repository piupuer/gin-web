package cache_service

import (
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
