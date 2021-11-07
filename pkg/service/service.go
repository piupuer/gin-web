package service

import (
	"context"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/query"
)

type MysqlService struct {
	Q query.MySql
}

func New(ctx context.Context) MysqlService {
	ops := []func(*query.MysqlOptions){
		query.WithMysqlLogger(global.Log),
		query.WithMysqlCtx(ctx),
		query.WithMysqlDb(global.Mysql),
		query.WithMysqlCasbinEnforcer(global.CasbinEnforcer),
	}
	if global.Conf.Redis.Enable {
		ops = append(ops, query.WithMysqlRedis(global.Redis))
	}
	my := MysqlService{
		Q: query.NewMySql(ops...),
	}
	return my
}
