package redis

import (
	"gin-web/pkg/utils"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"sync"
)

// Statement statement
type Statement struct {
	DB           *QueryRedis
	Schema       *schema.Schema
	ReflectValue reflect.Value
	Model        interface{}
	// 表名称
	Table string
	// 输出值
	Dest interface{}
	// 记录需要preload的所有字段信息
	preloads []searchPreload
	// where条件
	whereConditions []whereCondition
	// 排序条件
	orderConditions []orderCondition
	// 分页
	limit int
	// 偏移量
	offset int
	// 是否只取一条数据
	first bool
	// 是否只查询数据条数
	count bool
	// 不查询表, 直接给定了json字符串
	json    bool
	jsonStr string
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

// 预加载
func (stmt *Statement) Preload(schema string) *Statement {
	var preloads []searchPreload
	for _, preload := range stmt.preloads {
		if preload.schema != schema {
			preloads = append(preloads, preload)
		}
	}
	preloads = append(preloads, searchPreload{schema})
	stmt.preloads = preloads
	return stmt
}

// 指定json字符串
func (stmt *Statement) FromString(str string) *Statement {
	stmt.jsonStr = str
	stmt.json = true
	return stmt
}

// Where
func (stmt *Statement) Where(key, cond string, val interface{}) *Statement {
	// 可能通过小数点分隔子表多条件
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
	for _, condition := range stmt.whereConditions {
		if condition.key != key {
			whereConditions = append(whereConditions, condition)
		}
	}
	// 添加新数据
	whereConditions = append(whereConditions, whereCondition{key, cond, val})
	stmt.whereConditions = whereConditions
	return stmt
}

// order条件
func (stmt *Statement) Order(key string) *Statement {
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
	for _, condition := range stmt.orderConditions {
		if condition.property != key {
			orderConditions = append(orderConditions, condition)
		}
	}
	// 添加新数据
	orderConditions = append(orderConditions, orderCondition{property, asc})
	stmt.orderConditions = orderConditions
	return stmt
}

func (stmt *Statement) Parse(value interface{}) (err error) {
	if stmt.DB.cacheStore == nil {
		stmt.DB.cacheStore = &sync.Map{}
	}
	if stmt.DB.NamingStrategy == nil {
		stmt.DB.NamingStrategy = schema.NamingStrategy{}
	}
	if stmt.Schema, err = schema.Parse(value, stmt.DB.cacheStore, stmt.DB.NamingStrategy); err == nil && stmt.Table == "" {
		if tables := strings.Split(stmt.Schema.Table, "."); len(tables) == 2 {
			stmt.Table = tables[1]
			return
		}

		stmt.Table = stmt.Schema.Table
	}
	return err
}
