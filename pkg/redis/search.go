package redis

import "gin-web/pkg/utils"

type search struct {
	query *QueryRedis
	// 表名称
	tableName string
	// 记录需要preload的所有字段信息
	preload []searchPreload
	// where条件
	whereConditions []whereCondition
	// 分页
	limit int
	// 偏移量
	offset int
	// 输出值
	out interface{}
	// 是否只取一条数据
	first bool
}

func (s *search) clone() *search {
	clone := *s
	return &clone
}

// 预加载
type searchPreload struct {
	schema string
}

// where条件(与jsonq Where参数一致)
type whereCondition struct {
	// 键名称
	key string
	// 条件
	cond string
	// 值
	val interface{}
}

// 指定表名称
func (s *search) Table(name string) *search {
	s.tableName = name
	return s
}

// 预加载
func (s *search) Preload(schema string) *search {
	var preloads []searchPreload
	for _, preload := range s.preload {
		if preload.schema != schema {
			preloads = append(preloads, preload)
		}
	}
	preloads = append(preloads, searchPreload{
		schema: schema,
	})
	s.preload = preloads
	return s
}

// where条件
func (s *search) Where(key, cond string, val interface{}) *search {
	// redis存的key均为驼峰
	key = utils.CamelCaseLowerFirst(key)
	// 如果参数为uint, redis存的json转换为int, 因此这里转一下类型
	if item, ok := val.(uint); ok {
		val = int(item)
	}
	var whereConditions []whereCondition
	// 保留旧数据
	for _, condition := range s.whereConditions {
		if condition.key != key {
			whereConditions = append(whereConditions, condition)
		}
	}
	// 添加新数据
	whereConditions = append(whereConditions, whereCondition{
		key:  key,
		cond: cond,
		val:  val,
	})
	s.whereConditions = whereConditions
	return s
}
