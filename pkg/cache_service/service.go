package cache_service

import (
	"context"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"github.com/piupuer/go-helper/pkg/query"
)

type RedisService struct {
	Q      query.Redis
	mysql  service.MysqlService
	binlog bool
}

func New(ctx context.Context) RedisService {
	rd := RedisService{
		mysql:  service.New(ctx),
		binlog: global.Conf.Redis.EnableBinlog,
	}
	if global.Conf.Redis.EnableBinlog {
		ops := []func(*query.RedisOptions){
			query.WithRedisCtx(ctx),
			query.WithRedisClient(global.Redis),
			query.WithRedisDatabase(global.Conf.Mysql.DSN.DBName),
			query.WithRedisNamingStrategy(global.Mysql.NamingStrategy),
			query.WithRedisCasbinEnforcer(global.CasbinEnforcer),
		}
		rd.Q = query.NewRedis(ops...)
	}
	return rd
}
