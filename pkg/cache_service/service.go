package cache_service

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/resp"
	"reflect"
)

// 所有的查询可以走redis, 但数据的更新还是走mysql
type RedisService struct {
	ctx   *gin.Context         // 上下文
	mysql service.MysqlService // 保留mysql, 如果没开启redis可以走mysql
	redis *redis.QueryRedis    // redis对象实例
}

// 初始化服务
func New(c *gin.Context) RedisService {
	nc := gin.Context{}
	if c != nil {
		nc = *c
	}
	s := RedisService{
		ctx: &nc,
	}
	rc := s.RequestIdContext("")
	s.ctx.Set(global.RequestIdContextKey, rc.Value(global.RequestIdContextKey))
	s.mysql = service.New(s.ctx)
	s.redis = redis.New(s.ctx)
	return s
}

// 获取携带request id的上下文
func (s RedisService) RequestIdContext(requestId string) context.Context {
	if s.ctx != nil && requestId == "" {
		requestId = s.ctx.GetString(global.RequestIdContextKey)
	}
	return global.RequestIdContext(requestId)
}

// 查询, model需使用指针, 否则可能无法绑定数据
func (s RedisService) Find(query *redis.QueryRedis, page *resp.Page, model interface{}) (err error) {
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Slice) {
		return fmt.Errorf("model必须是非空指针数组类型")
	}

	if !page.NoPagination {
		// 查询条数
		err = query.Count(&page.Total).Error
		if err == nil && page.Total > 0 {
			// 获取分页参数
			limit, offset := page.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(model).Error
		}
	} else {
		// 不使用分页
		err = query.Find(model).Error
		if err == nil {
			page.Total = int64(rv.Elem().Len())
			// 获取分页参数
			page.GetLimit()
		}
	}
	return
}
