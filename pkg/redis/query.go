package redis

import (
	"database/sql/driver"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/thedevsaddam/gojsonq/v2"
	"reflect"
	"sort"
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
	// 查询条件
	search *search
}

// 单次查询域
type QueryRedisScope struct {
	scope *gorm.Scope
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
	clone := s.clone()
	clone.search.Table(name)
	clone.search.out = nil
	return clone
}

// 预先加载某一列
func (s *QueryRedis) Preload(column string) *QueryRedis {
	return s.clone().search.Preload(column).query
}

// 查询条件
func (s *QueryRedis) Where(key, cond string, val interface{}) *QueryRedis {
	return s.clone().search.Where(key, cond, val).query
}

// 查询列表
func (s *QueryRedis) Find(out interface{}) *QueryRedis {
	clone := s.clone()
	if !clone.check() {
		return clone
	}
	// 记录输出值
	clone.search.out = out
	// 获取数据
	clone.get(clone.search.tableName)
	return s
}

// 查询一条
func (s *QueryRedis) First(out interface{}) *QueryRedis {
	clone := s.clone()
	clone.search.limit = 1
	clone.search.first = true
	clone.Find(out)
	return clone
}

// 获取总数
func (s *QueryRedis) Count(out *uint) *QueryRedis {
	clone := s.clone()
	clone.search.out = out
	clone.count(out)
	return clone
}

// 分页
func (s *QueryRedis) Limit(limit uint) *QueryRedis {
	clone := s.clone()
	clone.search.limit = int(limit)
	return clone
}

func (s *QueryRedis) Offset(offset uint) *QueryRedis {
	clone := s.clone()
	clone.search.offset = int(offset)
	return clone
}

// //////////////////////////////////////////////////////////////////////////////
// Private Methods For QueryRedis
// //////////////////////////////////////////////////////////////////////////////

func (s *QueryRedis) clone() *QueryRedis {
	query := &QueryRedis{
		redis: s.redis,
		mysql: s.mysql,
	}

	if s.search == nil {
		query.search = &search{limit: -1, offset: -1}
	} else {
		query.search = s.search.clone()
	}

	query.search.query = query
	return query
}

// 获取总数
func (s *QueryRedis) count(value *uint) *QueryRedis {
	if !s.check() {
		*value = 0
	}
	// 读取某个表的数据总数
	*value = uint(s.get(s.search.tableName).Count())
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
	var nullList interface{}
	list := query.Get()
	if s.search.first {
		// 取第一条数据
		switch list.(type) {
		case []interface{}:
			v, _ := list.([]interface{})
			if len(v) > 0 {
				list = v[0]
			} else {
				// 设置为空元素
				list = nullList
			}
		}
	}
	// 获取json结果并转换为结构体
	utils.Struct2StructByJson(list, s.search.out)

	if list != nil {
		// 处理预加载
		s.processPreload()
	}
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
	for _, condition := range s.search.whereConditions {
		query = query.Where(condition.key, condition.cond, condition.val)
	}
	// 添加limit/offset
	query.Limit(s.search.limit)
	query.Offset(s.search.offset)
	return query
}

// 校验表名是否正常
func (s QueryRedis) check() bool {
	if strings.TrimSpace(s.clone().search.tableName) == "" {
		s.Error = fmt.Errorf("invalid table name: '%s'", s.clone().search.tableName)
		return false
	}
	return true
}

// 通过反射获取元素真实种类, 支持指针元素
func getRealKind(v interface{}) (reflect.Value, reflect.Kind) {
	rv := reflect.ValueOf(v)
	// 指针类型/接口类型需要继续获取元素值再判断
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv, rv.Kind()
}

// 类型转换(redis查询如果指定interface会将数字类型默认转为float64)
func convertTypeFromKind(source reflect.Value, target reflect.Value) interface{} {
	switch source.Kind() {
	case reflect.Float64:
		var v float64
		v = source.Float()
		return convertFloat64TypeFormKind(v, target)
	case reflect.Int:
		var v int64
		v = source.Int()
		return convertIntTypeFormKind(int(v), target)
	case reflect.Uint:
		var v uint64
		v = source.Uint()
		return convertUintTypeFormKind(uint(v), target)
	case reflect.String:
		var v string
		v = source.String()
		return convertStringTypeFormKind(v, target)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertTypeFromKind]类型转换失败, %v", target.Kind()))
	return nil
}

// float64转为target类型
func convertFloat64TypeFormKind(source float64, target reflect.Value) interface{} {
	switch target.Kind() {
	case reflect.Float64:
		return source
	case reflect.Uint:
		var v uint
		return convertFloat64Type(source, v)
	case reflect.Int:
		var v int
		return convertFloat64Type(source, v)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertFloat64TypeFormKind]类型转换失败, %v", target.Kind()))
	return nil
}

// int转为target类型
func convertIntTypeFormKind(source int, target reflect.Value) interface{} {
	switch target.Kind() {
	case reflect.Float64:
		return source
	case reflect.Uint:
		var v uint
		return convertIntType(source, v)
	case reflect.Int:
		var v int
		return convertIntType(source, v)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertIntTypeFormKind]类型转换失败, %v", target.Kind()))
	return nil
}

// int转为target类型
func convertUintTypeFormKind(source uint, target reflect.Value) interface{} {
	switch target.Kind() {
	case reflect.Float64:
		return source
	case reflect.Uint:
		var v uint
		return convertUintType(source, v)
	case reflect.Int:
		var v int
		return convertUintType(source, v)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertUintTypeFormKind]类型转换失败, %v", target.Kind()))
	return nil
}

// string转为target类型
func convertStringTypeFormKind(source string, target reflect.Value) interface{} {
	switch target.Kind() {
	case reflect.Struct:
		switch target.Interface().(type) {
		// 自定义日期转换
		case models.LocalTime:
			var time models.LocalTime
			_ = time.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", source)))
			return time
		}
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertStringTypeFormKind]类型转换失败, %v", target.Kind()))
	return nil
}

// 类型转换(redis查询如果指定interface会将数字类型默认转为float64)
func convertType(source interface{}, target interface{}) interface{} {
	switch source.(type) {
	case float64:
		return convertFloat64Type(source.(float64), target)
	case int:
		return convertIntType(source.(int), target)
	case uint:
		return convertUintType(source.(uint), target)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertType]类型转换失败, %v", target))
	return nil
}

// float64转为target类型
func convertFloat64Type(source float64, target interface{}) interface{} {
	switch target.(type) {
	case float64:
		return source
	case uint:
		return uint(source)
	case int:
		return int(source)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertFloat64Type]类型转换失败, %v", target))
	return nil
}

// int转为target类型
func convertIntType(source int, target interface{}) interface{} {
	switch target.(type) {
	case int:
		return source
	case uint:
		return uint(source)
	case float64:
		return float64(source)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertIntType]类型转换失败, %v", target))
	return nil
}

// uint转为target类型
func convertUintType(source uint, target interface{}) interface{} {
	switch target.(type) {
	case uint:
		return source
	case int:
		return int(source)
	case float64:
		return float64(source)
	}
	global.Log.Warn(fmt.Sprintf("[QueryRedis.convertUintType]类型转换失败, %v", target))
	return nil
}

// //////////////////////////////////////////////////////////////////////////////
// Private Methods For QueryRedisScope. From gorm.scope
// //////////////////////////////////////////////////////////////////////////////

// 处理预加载
func (s *QueryRedis) processPreload() {
	// 获取当前域
	scope := s.NewScope(s.search.out)
	var (
		preloadedMap = map[string]bool{}
		fields       = scope.scope.Fields()
	)
	// preload其他关联表
	for _, preload := range s.search.preload {
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
					// 根据关联字段不同类型处理
					switch field.Relationship.Kind {
					// TODO 项目中暂时没有用到has_one
					// case "has_one":
					// 	currentScope.handleHasOnePreload(field)
					case "has_many":
						currentScope.handleHasManyPreload(field)
					case "belongs_to":
						currentScope.handleBelongsToPreload(field)
					case "many_to_many":
						currentScope.handleManyToManyPreload(field)
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
			// searchPreload next level
			if idx < len(preloadFields)-1 {
				currentScope = currentScope.getColumnAsScope(preloadField)
				if currentScope != nil {
					currentFields = currentScope.scope.Fields()
				}
			}
		}
	}
}

// handleHasManyPreload used to searchPreload has many associations
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

// handleBelongsToPreload used to searchPreload belongs to associations
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
}

// handleManyToManyPreload used to searchPreload many to many associations
func (s *QueryRedisScope) handleManyToManyPreload(field *gorm.Field) {
	var (
		relation         = field.Relationship
		joinTableHandler = relation.JoinTableHandler
		fieldType        = field.Struct.Type.Elem()
		linkHash         = map[string][]reflect.Value{}
		isPtr            bool
	)

	if fieldType.Kind() == reflect.Ptr {
		isPtr = true
		fieldType = fieldType.Elem()
	}

	var sourceKeys []string
	for _, key := range joinTableHandler.SourceForeignKeys() {
		sourceKeys = append(sourceKeys, key.DBName)
	}
	var destinationKeys []string
	var destinationCamelKeys []string
	for _, key := range joinTableHandler.DestinationForeignKeys() {
		destinationKeys = append(destinationKeys, key.DBName)
		destinationCamelKeys = append(destinationCamelKeys, utils.CamelCaseLowerFirst(key.DBName))
	}

	// 查询关系表
	relationRows := make([]map[string]interface{}, 0)
	if many2many, _ := field.TagSettingsGet("MANY2MANY"); many2many != "" {
		many2manyDB := New().Table(many2many)
		// 找到对应字段名, 组成查询条件
		for _, key := range sourceKeys {
			// get relations's primary keys
			primaryKeys := s.getColumnAsArray(relation.ForeignFieldNames, s.scope.Value)
			if len(primaryKeys) == 0 {
				return
			}

			many2manyDB = many2manyDB.Where(
				// 转为驼峰命名
				utils.CamelCaseLowerFirst(key),
				"in",
				toQueryValues(primaryKeys),
			)
		}
		err := many2manyDB.Find(&relationRows).Error
		if s.scope.Err(err) != nil {
			return
		}
	}

	sourceIds := make([][]interface{}, 0)
	// 取出符合关系的关联表数据id
	for _, row := range relationRows {
		for k, v := range row {
			for _, key := range destinationKeys {
				if k == utils.CamelCaseLowerFirst(key) {
					// redis查询如果指定interface会将数字类型默认转为float64
					var target int
					sourceIds = append(sourceIds, []interface{}{convertType(v, target)})
				}
			}
		}
	}
	if len(sourceIds) == 0 {
		return
	}

	// 查询需要关联的表
	foreignScope := s.scope.New(reflect.New(fieldType).Interface())
	foreignRows := make([]map[string]interface{}, 0)
	foreignDB := New().Table(foreignScope.TableName())
	// 找到对应字段名, 组成查询条件
	for _, key := range relation.AssociationForeignFieldNames {
		foreignDB = foreignDB.Where(
			// 转为驼峰命名
			utils.CamelCaseLowerFirst(key),
			"in",
			toQueryValues(sourceIds),
		)
	}
	err := foreignDB.Find(&foreignRows).Error

	if s.scope.Err(err) != nil {
		return
	}

	// 获取外键对应的全部值
	foreignKeys := s.getColumnAsArray(relation.ForeignFieldNames, s.scope.Value)
	hashedSourceKeys := toString(foreignKeys[0])

	// 将数据暂存至linkHash
	for _, row := range foreignRows {
		var (
			elem   = reflect.New(fieldType).Elem()
			fields = s.scope.New(elem.Addr().Interface()).Fields()
		)

		// 将每一行的数据写入结果集
		s.scan(row, append(fields))

		// 暂存hash
		if isPtr {
			linkHash[hashedSourceKeys] = append(linkHash[hashedSourceKeys], elem.Addr())
		} else {
			linkHash[hashedSourceKeys] = append(linkHash[hashedSourceKeys], elem)
		}
	}

	// assign find results
	var (
		indirectScopeValue = s.scope.IndirectValue()
		fieldsSourceMap    = map[string][]reflect.Value{}
		foreignFieldNames  []string
	)

	for _, dbName := range relation.ForeignFieldNames {
		if field, ok := s.scope.FieldByName(dbName); ok {
			foreignFieldNames = append(foreignFieldNames, field.Name)
		}
	}

	if indirectScopeValue.Kind() == reflect.Slice {
		for j := 0; j < indirectScopeValue.Len(); j++ {
			object := indirect(indirectScopeValue.Index(j))
			key := toString(getValueFromFields(object, foreignFieldNames))
			fieldsSourceMap[key] = append(fieldsSourceMap[key], object.FieldByName(field.Name))
		}
	} else if indirectScopeValue.IsValid() {
		key := toString(getValueFromFields(indirectScopeValue, foreignFieldNames))
		fieldsSourceMap[key] = append(fieldsSourceMap[key], indirectScopeValue.FieldByName(field.Name))
	}

	for source, fields := range fieldsSourceMap {
		for _, f := range fields {
			// If not 0 this means Value is a pointer and we already added preloaded models to it
			if f.Len() != 0 {
				continue
			}

			v := reflect.MakeSlice(f.Type(), 0, 0)
			if len(linkHash[source]) > 0 {
				v = reflect.Append(f, linkHash[source]...)
			}

			f.Set(v)
		}
	}
}

func (s *QueryRedisScope) scan(row map[string]interface{}, fields []*gorm.Field) {
	// 从row中读取列
	columns := make([]string, 0)
	for key := range row {
		columns = append(columns, key)
	}
	// 统一排序
	sort.Strings(columns)
	var (
		ignored            interface{}
		values             = make([]interface{}, len(columns))
		selectFields       []*gorm.Field
		selectedColumnsMap = map[string]int{}
		resetFields        = map[int]*gorm.Field{}
	)

	for index, column := range columns {
		values[index] = &ignored

		selectFields = fields
		offset := 0
		if idx, ok := selectedColumnsMap[column]; ok {
			offset = idx + 1
			selectFields = selectFields[offset:]
		}

		for fieldIndex, field := range selectFields {
			// 转为驼峰命名
			if utils.CamelCaseLowerFirst(field.DBName) == column {
				if field.Field.Kind() == reflect.Ptr {
					values[index] = field.Field.Addr().Interface()
				} else {
					reflectValue := reflect.New(reflect.PtrTo(field.Struct.Type))
					reflectValue.Elem().Set(field.Field.Addr())
					values[index] = reflectValue.Interface()
					resetFields[index] = field
				}

				selectedColumnsMap[column] = offset + fieldIndex

				if field.IsNormal {
					break
				}
			}
		}
	}

	for i, v := range values {
		// 通过反射赋值
		// 获取元素真实类型
		reflectV, reflectVKind := getRealKind(v)
		item := row[columns[i]]
		reflectRow, reflectRowKind := getRealKind(item)

		// row数据为无效直接跳过
		if reflectRowKind == reflect.Invalid {
			continue
		}

		// 类型相同 或 类型不同但v值有效(说明不是空指针)
		if reflectVKind == reflectRowKind || reflectV.IsValid() {
			if reflectVKind != reflectRowKind {
				// redis查询如果指定interface会将数字类型默认转为float64
				newVal := convertTypeFromKind(reflectRow, reflectV)
				if newVal != nil {
					reflectV.Set(reflect.ValueOf(newVal))
				}
			} else {
				// 类型一致
				reflectV.Set(reflectRow)
			}
		} else {
			// v的类型为无效, 可能当前为空指针
			if reflectVKind == reflect.Invalid {
				// 寻址只需1层Elem(v默认2层Elem才能找到最后的值类型)
				elem := reflect.ValueOf(v).Elem()
				kind := elem.Kind()
				if kind == reflect.Ptr || kind == reflect.Interface {
					// 以当前row的值类型为基础创建新对象, 指针类型, 类似于new(xxx)
					rowPtr := reflect.New(reflectRow.Type()).Elem()
					// 将当前row值写入rowPtr
					rowPtr.Set(reflectRow)
					// 将rowPtr写入v
					elem.Set(rowPtr.Addr())
				}
			} else {
				global.Log.Warn(fmt.Sprintf("[QueryRedisScope.scan]类型不匹配, row type: %v, value type: %v", reflectRowKind, reflectVKind))
			}
		}
	}

	for index, field := range resetFields {
		if v := reflect.ValueOf(values[index]).Elem().Elem(); v.IsValid() {
			field.Field.Set(v)
		}
	}
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
					// 忽略字段大小写
					field := object.FieldByNameFunc(func(n string) bool {
						return strings.ToLower(n) == strings.ToLower(column)
					})
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
				// 忽略字段大小写
				field := indirectValue.FieldByNameFunc(func(n string) bool {
					return strings.ToLower(n) == strings.ToLower(column)
				})
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
