package cache_service

import (
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/query"
)

type RedisService struct {
	Q      query.Redis
	mysql  service.MysqlService
	binlog bool
}

func New(c *gin.Context) RedisService {
	rd := RedisService{
		mysql:  service.New(c),
		binlog: global.Conf.Redis.EnableBinlog,
	}
	if global.Conf.Redis.EnableBinlog {
		ops := []func(*query.RedisOptions){
			query.WithRedisLogger(global.Log),
			query.WithRedisCtx(c),
			query.WithRedisClient(global.Redis),
			query.WithRedisDatabase(global.Conf.Mysql.DSN.DBName),
			query.WithRedisNamingStrategy(global.Mysql.NamingStrategy),
		}
		rd.Q = query.NewRedis(ops...)
	}
	return rd
}
