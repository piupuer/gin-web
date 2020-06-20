package redis

import (
	"database/sql/driver"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/thedevsaddam/gojsonq/v2"
	"reflect"
	"strings"
)

// 以redis client为基础, 加入一些类似gorm的查询方法

// 自定义redis查询结构
type QueryRedis struct {
	// 错误信息
	Error error
	// redis对象
	redis *redis.Client
	// mysql对象, 主要作用是映射相关字段
	mysql *gorm.DB
	// 表名称
	tableName string
	// 记录需要preload的所有字段信息
	preloads []preload
	// where条件
	whereConditions []whereCondition
	// 分页
	limit int
	// 偏移量
	offset int
	// 输出值
	out interface{}
}

// 单次查询域
type QueryRedisScope struct {
	scope *gorm.Scope
}

// 预加载
type preload struct {
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

// 初始化服务
func New() *QueryRedis {
	return &QueryRedis{
		redis: global.Redis,
		mysql: global.Mysql,
	}
}

// 指定表名称
func (s *QueryRedis) Table(name string) *QueryRedis {
	s.tableName = name
	return s
}

// 预先加载某一列
func (s *QueryRedis) Preload(column string) *QueryRedis {
	s.preload(column)
	return s
}

// 查询条件
func (s *QueryRedis) Where(key, cond string, val interface{}) *QueryRedis {
	s.where(key, cond, val)
	return s
}

// 查询列表
func (s *QueryRedis) Find(out interface{}) *QueryRedis {
	if !s.check() {
		return s
	}
	// 记录输出值
	s.out = out
	// 获取数据
	s.get(s.tableName)
	return s
}

// 查询一条
func (s *QueryRedis) First(out interface{}) *QueryRedis {
	s.limit = 1
	s.Find(out)
	return s
}

// 获取总数
func (s *QueryRedis) Count(out *int) *QueryRedis {
	s.out = out
	s.count(out)
	return s
}

// 分页
func (s *QueryRedis) Limit(limit int) *QueryRedis {
	s.limit = limit
	return s
}

func (s *QueryRedis) Offset(offset int) *QueryRedis {
	s.offset = offset
	return s
}

// //////////////////////////////////////////////////////////////////////////////
// Private Methods For QueryRedis
// //////////////////////////////////////////////////////////////////////////////

// 预加载
func (s *QueryRedis) preload(schema string) *QueryRedis {
	var preloads []preload
	// 保留旧数据
	for _, preload := range s.preloads {
		if preload.schema != schema {
			preloads = append(preloads, preload)
		}
	}
	// 添加新数据
	preloads = append(preloads, preload{schema: schema})
	s.preloads = preloads
	return s
}

// where条件
func (s *QueryRedis) where(key, cond string, val interface{}) *QueryRedis {
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

// 获取总数
func (s *QueryRedis) count(value *int) *QueryRedis {
	if !s.check() {
		*value = 0
	}
	// 读取某个表的数据总数
	*value = s.get(s.tableName).Count()
	return s
}

// 从缓存中获取model全部数据, 返回json字符串
func (s QueryRedis) get(tableName string) *gojsonq.JSONQ {
	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, tableName)
	res, err := s.redis.Get(cacheKey).Result()
	global.Log.Debug(fmt.Sprintf("[QueryRedis.get]读取redis缓存: %s", tableName))
	if err != nil {
		global.Log.Debug(fmt.Sprintf("[QueryRedis.get]读取redis缓存异常: %v", err))
	}
	if res == "" {
		// 如果是空字符串, 将其设置为空数组, 否则list会被转为nil
		res = "[]"
	}
	query := s.jsonQuery(res)
	fmt.Println(query.Get())
	// 获取json结果并转换为结构体
	utils.Struct2StructByJson(query.Get(), s.out)

	// 处理预加载
	s.processPreload()
	return query
}

// 新建域
func (s *QueryRedis) NewScope(value interface{}) *QueryRedisScope {
	scope := s.mysql.NewScope(value)
	return &QueryRedisScope{scope: scope}
}

// New create a new Scope without search information
func (s *QueryRedisScope) New(value interface{}) *QueryRedisScope {
	scope := global.Mysql.NewScope(value)
	return &QueryRedisScope{scope: scope}
}

// jsonq需要每一次新建不同实例, 否则可能查询条件重叠, 找不到记录
func (s QueryRedis) jsonQuery(str string) *gojsonq.JSONQ {
	// 使用jsonq
	query := gojsonq.New().FromString(str)
	// 添加where条件
	for _, condition := range s.whereConditions {
		query = query.Where(condition.key, condition.cond, condition.val)
	}
	// 添加limit/offset
	query.Limit(s.limit)
	query.Offset(s.offset)
	return query
}

// 校验表名是否正常
func (s QueryRedis) check() bool {
	if strings.TrimSpace(s.tableName) == "" {
		s.Error = fmt.Errorf("invalid table name: '%s'", s.tableName)
		return false
	}
	return true
}

// //////////////////////////////////////////////////////////////////////////////
// Private Methods For QueryRedisScope. From gorm.scope
// //////////////////////////////////////////////////////////////////////////////

// 处理预加载
func (s *QueryRedis) processPreload() {
	// 获取当前域
	scope := s.NewScope(s.out)
	var (
		preloadedMap = map[string]bool{}
		fields       = scope.scope.Fields()
	)
	// preload其他关联表
	for _, preload := range s.preloads {
		var (
			// 由.分隔子项目
			preloadFields = strings.Split(preload.schema, ".")
			currentScope  = scope
			currentFields = fields
		)
		for idx, preloadField := range preloadFields {
			if currentScope == nil {
				continue
			}

			// if not preloaded
			if preloadKey := strings.Join(preloadFields[:idx+1], "."); !preloadedMap[preloadKey] {
				for _, field := range currentFields {
					if field.Name != preloadField || field.Relationship == nil {
						continue
					}
					fmt.Println(field.Relationship.Kind)
					// 根据关联字段不同类型处理
					switch field.Relationship.Kind {
					// TODO 项目中暂时没有用到has_one/many_to_many
					// case "has_one":
					// 	currentScope.handleHasOnePreload(field)
					case "has_many":
						currentScope.handleHasManyPreload(field)
					case "belongs_to":
						currentScope.handleBelongsToPreload(field)
					// case "many_to_many":
					// 	currentScope.handleManyToManyPreload(field)
					default:
						s.Error = fmt.Errorf("unsupported relation: %s", field.Relationship.Kind)
					}

					// 记录已经加载过的数据
					preloadedMap[preloadKey] = true
					break
				}

				if !preloadedMap[preloadKey] {
					currentScope.scope.Err(fmt.Errorf("can't preload field %s for %s", preloadField, currentScope.scope.GetModelStruct().ModelType))
					return
				}
			}
			// preload next level
			if idx < len(preloadFields)-1 {
				currentScope = currentScope.getColumnAsScope(preloadField)
				if currentScope != nil {
					currentFields = currentScope.scope.Fields()
				}
			}
		}
	}
}

// handleHasManyPreload used to preload has many associations
func (s *QueryRedisScope) handleHasManyPreload(field *gorm.Field) {
	relation := field.Relationship

	// get relations's primary keys
	primaryKeys := s.getColumnAsArray(relation.AssociationForeignFieldNames, s.scope.Value)
	if len(primaryKeys) == 0 {
		return
	}

	// 新建数据库查询对象
	preloadDB := New()

	results := makeSlice(field.Struct.Type)
	// 获取目标表名
	tableName := s.scope.New(results).GetModelStruct().TableName(s.scope.DB())

	preloadDB.Table(tableName).Where(
		"id",
		"in",
		toQueryValues(primaryKeys),
	).Find(results)

	// assign find results
	var (
		resultsValue       = indirect(reflect.ValueOf(results))
		indirectScopeValue = s.scope.IndirectValue()
	)

	if indirectScopeValue.Kind() == reflect.Slice {
		preloadMap := make(map[string][]reflect.Value)
		for i := 0; i < resultsValue.Len(); i++ {
			result := resultsValue.Index(i)
			foreignValues := getValueFromFields(result, relation.ForeignFieldNames)
			preloadMap[toString(foreignValues)] = append(preloadMap[toString(foreignValues)], result)
		}

		for j := 0; j < indirectScopeValue.Len(); j++ {
			object := indirect(indirectScopeValue.Index(j))
			objectRealValue := getValueFromFields(object, relation.AssociationForeignFieldNames)
			f := object.FieldByName(field.Name)
			if results, ok := preloadMap[toString(objectRealValue)]; ok {
				f.Set(reflect.Append(f, results...))
			} else {
				f.Set(reflect.MakeSlice(f.Type(), 0, 0))
			}
		}
	} else {
		s.scope.Err(field.Set(resultsValue))
	}
}

// handleBelongsToPreload used to preload belongs to associations
func (s *QueryRedisScope) handleBelongsToPreload(field *gorm.Field) {
	relation := field.Relationship

	// 新建数据库查询对象
	preloadDB := New()

	// get relations's primary keys
	primaryKeys := s.getColumnAsArray(relation.ForeignFieldNames, s.scope.Value)
	if len(primaryKeys) == 0 {
		return
	}

	// 查询关联关系
	results := makeSlice(field.Struct.Type)
	// 获取目标表名
	tableName := s.scope.New(results).GetModelStruct().TableName(s.scope.DB())

	preloadDB.Table(tableName).Where(
		"id",
		"in",
		toQueryValues(primaryKeys),
	).Find(results)

	// assign find results
	var (
		resultsValue       = indirect(reflect.ValueOf(results))
		indirectScopeValue = s.scope.IndirectValue()
	)

	foreignFieldToObjects := make(map[string][]*reflect.Value)
	if indirectScopeValue.Kind() == reflect.Slice {
		for j := 0; j < indirectScopeValue.Len(); j++ {
			object := indirect(indirectScopeValue.Index(j))
			valueString := toString(getValueFromFields(object, relation.ForeignFieldNames))
			foreignFieldToObjects[valueString] = append(foreignFieldToObjects[valueString], &object)
		}
	}

	for i := 0; i < resultsValue.Len(); i++ {
		result := resultsValue.Index(i)
		if indirectScopeValue.Kind() == reflect.Slice {
			valueString := toString(getValueFromFields(result, relation.AssociationForeignFieldNames))
			if objects, found := foreignFieldToObjects[valueString]; found {
				for _, object := range objects {
					object.FieldByName(field.Name).Set(result)
				}
			}
		} else {
			s.scope.Err(field.Set(result))
		}
	}

	fmt.Println(results, preloadDB)

}

func (s *QueryRedisScope) getColumnAsArray(columns []string, values ...interface{}) (results [][]interface{}) {
	resultMap := make(map[string][]interface{})
	for _, value := range values {
		indirectValue := indirect(reflect.ValueOf(value))

		switch indirectValue.Kind() {
		case reflect.Slice:
			for i := 0; i < indirectValue.Len(); i++ {
				var result []interface{}
				var object = indirect(indirectValue.Index(i))
				var hasValue = false
				for _, column := range columns {
					field := object.FieldByName(column)
					if hasValue || !isBlank(field) {
						hasValue = true
					}
					result = append(result, field.Interface())
				}

				if hasValue {
					h := fmt.Sprint(result...)
					if _, exist := resultMap[h]; !exist {
						resultMap[h] = result
					}
				}
			}
		case reflect.Struct:
			var result []interface{}
			var hasValue = false
			for _, column := range columns {
				field := indirectValue.FieldByName(column)
				if hasValue || !isBlank(field) {
					hasValue = true
				}
				result = append(result, field.Interface())
			}

			if hasValue {
				h := fmt.Sprint(result...)
				if _, exist := resultMap[h]; !exist {
					resultMap[h] = result
				}
			}
		}
	}
	for _, v := range resultMap {
		results = append(results, v)
	}
	return
}

func (s *QueryRedisScope) getColumnAsScope(column string) *QueryRedisScope {
	indirectScopeValue := s.scope.IndirectValue()

	switch indirectScopeValue.Kind() {
	case reflect.Slice:
		if fieldStruct, ok := s.scope.GetModelStruct().ModelType.FieldByName(column); ok {
			fieldType := fieldStruct.Type
			if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			resultsMap := map[interface{}]bool{}
			results := reflect.New(reflect.SliceOf(reflect.PtrTo(fieldType))).Elem()

			for i := 0; i < indirectScopeValue.Len(); i++ {
				result := indirect(indirect(indirectScopeValue.Index(i)).FieldByName(column))

				if result.Kind() == reflect.Slice {
					for j := 0; j < result.Len(); j++ {
						if elem := result.Index(j); elem.CanAddr() && resultsMap[elem.Addr()] != true {
							resultsMap[elem.Addr()] = true
							results = reflect.Append(results, elem.Addr())
						}
					}
				} else if result.CanAddr() && resultsMap[result.Addr()] != true {
					resultsMap[result.Addr()] = true
					results = reflect.Append(results, result.Addr())
				}
			}
			return s.New(results.Interface())
		}
	case reflect.Struct:
		if field := indirectScopeValue.FieldByName(column); field.CanAddr() {
			return s.New(field.Addr().Interface())
		}
	}
	return nil
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}

	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func makeSlice(elemType reflect.Type) interface{} {
	if elemType.Kind() == reflect.Slice {
		elemType = elemType.Elem()
	}
	sliceType := reflect.SliceOf(elemType)
	slice := reflect.New(sliceType)
	slice.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))
	return slice.Interface()
}

// getValueFromFields return given fields's value
func getValueFromFields(value reflect.Value, fieldNames []string) (results []interface{}) {
	// If value is a nil pointer, Indirect returns a zero Value!
	// Therefor we need to check for a zero value,
	// as FieldByName could panic
	if indirectValue := reflect.Indirect(value); indirectValue.IsValid() {
		for _, fieldName := range fieldNames {
			if fieldValue := reflect.Indirect(indirectValue.FieldByName(fieldName)); fieldValue.IsValid() {
				result := fieldValue.Interface()
				if r, ok := result.(driver.Valuer); ok {
					result, _ = r.Value()
				}
				results = append(results, result)
			}
		}
	}
	return
}

func toQueryValues(values [][]interface{}) (results []int) {
	for _, value := range values {
		for _, v := range value {
			// redis存的json转换为int, 因此这里转一下类型
			if item, ok := v.(uint); ok {
				results = append(results, int(item))
			} else if item, ok := v.(int); ok {
				results = append(results, item)
			} else {
				fmt.Println("当前关联表主键类型不兼容")
			}
		}
	}
	return
}

func toString(str interface{}) string {
	if values, ok := str.([]interface{}); ok {
		var results []string
		for _, value := range values {
			results = append(results, toString(value))
		}
		return strings.Join(results, "_")
	} else if bytes, ok := str.([]byte); ok {
		return string(bytes)
	} else if reflectValue := reflect.Indirect(reflect.ValueOf(str)); reflectValue.IsValid() {
		return fmt.Sprintf("%v", reflectValue.Interface())
	}
	return ""
}
