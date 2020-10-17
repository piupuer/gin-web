package redis

import (
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
)

// 查询列表
func (s *QueryRedis) Find(dest interface{}) *QueryRedis {
	ins := s.getInstance()
	if !ins.check() {
		return ins
	}
	// 记录输出值
	ins.Statement.Dest = dest

	stmt := ins.Statement

	if stmt.Model == nil {
		stmt.Model = stmt.Dest
	} else if stmt.Dest == nil {
		stmt.Dest = stmt.Model
	}

	if stmt.Model != nil {
		if err := stmt.Parse(stmt.Model); err != nil && (!errors.Is(err, schema.ErrUnsupportedDataType) || (stmt.Table == "")) {
			ins.AddError(err)
		}
	}

	if stmt.Dest != nil {
		stmt.ReflectValue = reflect.ValueOf(stmt.Dest)
		for stmt.ReflectValue.Kind() == reflect.Ptr {
			stmt.ReflectValue = stmt.ReflectValue.Elem()
		}
		if !stmt.ReflectValue.IsValid() {
			ins.AddError(fmt.Errorf("invalid value"))
		}
	}

	// 获取数据
	ins.get(ins.Statement.Table)
	return ins
}

// 查询一条
func (s *QueryRedis) First(dest interface{}) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.limit = 1
	ins.Statement.first = true
	ins.Find(dest)
	return ins
}

// 获取总数
func (s *QueryRedis) Count(dest *int64) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.Dest = dest
	if !ins.check() {
		*dest = 0
	}
	// 读取某个表的数据总数
	*dest = int64(ins.get(s.Statement.Table).Count())
	return ins
}
