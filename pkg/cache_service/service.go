package cache_service

import (
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/thedevsaddam/gojsonq/v2"
)

// 所有的查询可以走redis, 但数据的更新还是走mysql
type RedisService struct {
	mysql service.MysqlService // 保留mysql, 如果没开启redis可以走mysql
	json  *gojsonq.JSONQ       // json查询实例, 由于读取redis以json数据格式保存, 该实例可实现类似sql语句的查询
	redis *redis.Client        // redis对象实例
}

// 初始化服务
func New(c *gin.Context) RedisService {
	return RedisService{
		mysql: service.New(c),
		json:  gojsonq.New(),
		redis: global.Redis,
	}
}
