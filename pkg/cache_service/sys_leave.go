package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"strings"
)

// 获取所有请假
func (s *RedisService) GetLeaves(req *request.LeaveListRequestStruct) ([]models.SysLeave, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetLeaves(req)
	}
	var err error
	list := make([]models.SysLeave, 0)
	// 查询请假表所有缓存
	jsonLeaves := s.GetListFromCache(nil, new(models.SysLeave).TableName())
	query := s.JsonQuery().FromString(jsonLeaves)
	if req.Status != nil {
		// redis存的json转换为int, 因此这里转一下类型
		query = query.Where("status", "=", int(*req.Status))
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		query = query.Where("desc", "contains", desc)
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
