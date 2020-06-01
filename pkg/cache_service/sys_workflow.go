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
