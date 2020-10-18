package redis

// 查询列表
func (s *QueryRedis) Find(dest interface{}) *QueryRedis {
	ins := s.getInstance()
	if !ins.check() {
		return ins
	}
	// 记录输出值
	ins.Statement.Dest = dest
	// 重新指定model
	ins.Statement.Model = dest
	// 获取数据
	ins.beforeQuery(ins).get(ins.Statement.Table)
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
func (s *QueryRedis) Count(count *int64) *QueryRedis {
	ins := s.getInstance()
	ins.Statement.Dest = count
	if !ins.check() {
		*count = 0
	}
	// 记录输出值
	ins.Statement.Dest = count
	// 读取某个表的数据总数
	*count = int64(ins.beforeQuery(ins).get(ins.Statement.Table).Count())
	return ins
}
