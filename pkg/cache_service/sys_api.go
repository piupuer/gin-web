package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"strings"
)

// 获取所有接口
func (s *RedisService) GetApis(req *request.ApiListRequestStruct) ([]models.SysApi, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetApis(req)
	}
	var err error
	list := make([]models.SysApi, 0)
	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, new(models.SysApi).TableName())
	// 查询接口表所有缓存
	jsonApis, err := s.redis.Get(cacheKey).Result()
	if err != nil {
		return list, err
	}
	query := s.json.FromString(jsonApis)
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method", "contains", method)
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path", "contains", path)
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category", "contains", category)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
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
