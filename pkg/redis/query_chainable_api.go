package redis

import (
	"gin-web/pkg/global"
	"regexp"
	"strings"
)

var tableRegexp = regexp.MustCompile(`(?i).+? AS (\w+)\s*(?:$|,)`)

// 指定json字符串
func (s QueryRedis) FromString(str string) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.FromString(str)
	ins.Statement.Dest = nil
	return ins
}

// 指定表名称
func (s QueryRedis) Table(name string, args ...interface{}) (ins *QueryRedis) {
	ins = s.getInstance()
	// 降低复杂度, 只支持单表
	if strings.Contains(name, " ") || strings.Contains(name, "`") || len(args) > 0 {
		if results := tableRegexp.FindStringSubmatch(name); len(results) == 2 {
			ins.Statement.Table = global.Mysql.NamingStrategy.TableName(results[1])
			return
		}
	}

	ins.Statement.Table = global.Mysql.NamingStrategy.TableName(name)
	return
}

// 预先加载某一列
func (s QueryRedis) Preload(column string) *QueryRedis {
	return s.getInstance().Statement.Preload(column).DB
}

// Where条件
func (s QueryRedis) Where(key, cond string, val interface{}) *QueryRedis {
	return s.getInstance().Statement.Where(key, cond, val).DB
}

// 排序
func (s QueryRedis) Order(key string) *QueryRedis {
	return s.getInstance().Statement.Order(key).DB
}

// 分页
func (s QueryRedis) Limit(limit int) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.limit = limit
	return ins
}

func (s QueryRedis) Offset(offset int) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.offset = offset
	return ins
}
