package cache_service

import (
	"gin-web/pkg/redis"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
)

// 所有的查询可以走redis, 但数据的更新还是走mysql
type RedisService struct {
	mysql service.MysqlService // 保留mysql, 如果没开启redis可以走mysql
	redis *redis.QueryRedis    // redis对象实例
}

// 初始化服务
func New(c *gin.Context) RedisService {
	return RedisService{
		mysql: service.New(c),
		redis: redis.New(),
	}
}
