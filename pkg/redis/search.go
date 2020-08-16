package redis

import (
	"gin-web/pkg/utils"
	"strings"
)

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
	// 排序条件
	orderConditions []orderCondition
	// 是否只取一条数据
	first bool
	// 不查询表, 直接给定了json字符串
	json    bool
	jsonStr string
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

// order条件(与jsonq SortBy参数一致)
type orderCondition struct {
	// 键名称
	property string
	// 是否正向排序
	asc bool
}

// 指定json字符串
func (s *search) FromString(str string) *search {
	s.jsonStr = str
	s.json = true
	return s
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
	// 可能通过小数点分隔多条件
	keys := strings.Split(key, ".")
	newKeys := make([]string, 0)
	for _, item := range keys {
		// redis存的key均为驼峰
		newKeys = append(newKeys, utils.CamelCaseLowerFirst(item))
	}
	key = strings.Join(newKeys, ".")
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

// order条件
func (s *search) Order(key string) *search {
	key = strings.ToLower(key)
	// 通过空格拆分
	fields := strings.Split(key, " ")
	property := key
	asc := true
	// 刚好2个数据说明指定顺序
	if len(fields) == 2 && strings.TrimSpace(fields[1]) == "desc" {
		property = fields[0]
		asc = false
	}
	// redis存的key均为驼峰
	property = utils.CamelCaseLowerFirst(property)

	var orderConditions []orderCondition
	// 保留旧数据
	for _, condition := range s.orderConditions {
		if condition.property != key {
			orderConditions = append(orderConditions, condition)
		}
	}
	// 添加新数据
	orderConditions = append(orderConditions, orderCondition{
		property: property,
		asc:      asc,
	})
	s.orderConditions = orderConditions
	return s
}
