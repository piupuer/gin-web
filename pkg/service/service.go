package service

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/query"
)

type MysqlService struct {
	Q query.MySql
}

func New(c *gin.Context) MysqlService {
	ops := []func(*query.MysqlOptions){
		query.WithMysqlLogger(global.Log),
		query.WithMysqlCtx(c),
	}
	if global.Conf.Redis.Enable {
		ops = append(ops, query.WithMysqlRedis(global.Redis))
	}
	my := MysqlService{
		Q: query.NewMySql(global.Mysql, ops...),
	}
	return my
}
