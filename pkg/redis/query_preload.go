package redis

import (
	"fmt"
	localUtils "gin-web/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
	"reflect"
	"sort"
	"strings"
)

// 处理预加载
func (s *QueryRedis) processPreload() {
	if s.Error == nil && len(s.Statement.preloads) > 0 {
		preloadMap := map[string][]string{}
		// 组装需要预加载的字段列表
		for _, preload := range s.Statement.preloads {
			preloadFields := strings.Split(preload.schema, ".")
			for idx := range preloadFields {
				preloadMap[strings.Join(preloadFields[:idx+1], ".")] = preloadFields[:idx+1]
			}
		}

		preloadNames := make([]string, len(preloadMap))
		idx := 0
		for key := range preloadMap {
			preloadNames[idx] = key
			idx++
		}
		sort.Strings(preloadNames)

		for _, name := range preloadNames {
			var (
				curSchema     = s.Statement.Schema
				preloadFields = preloadMap[name]
				rels          = make([]*schema.Relationship, len(preloadFields))
			)

			for idx, preloadField := range preloadFields {
				if rel := curSchema.Relationships.Relations[preloadField]; rel != nil {
					rels[idx] = rel
					curSchema = rel.FieldSchema
				} else {
					s.AddError(fmt.Errorf("%v: %w", name, gorm.ErrUnsupportedRelation))
				}
			}

			preload(s, rels)
		}
	}
}

func preload(db *QueryRedis, rels []*schema.Relationship) {
	var (
		reflectValue     = db.Statement.ReflectValue
		rel              = rels[len(rels)-1]
		relForeignKeys   []string
		relForeignFields []*schema.Field
		foreignFields    []*schema.Field
		foreignValues    [][]interface{}
		identityMap      = map[string][]reflect.Value{}
	)

	if len(rels) > 1 {
		reflectValue = schema.GetRelationsValues(reflectValue, rels[:len(rels)-1])
	}

	// 全新的实例
	ins := New()

	if rel.JoinTable != nil {
		// 多对多关系
		var joinForeignFields, joinRelForeignFields []*schema.Field
		var joinForeignKeys []string
		for _, ref := range rel.References {
			if ref.OwnPrimaryKey {
				joinForeignKeys = append(joinForeignKeys, ref.ForeignKey.DBName)
				joinForeignFields = append(joinForeignFields, ref.ForeignKey)
				foreignFields = append(foreignFields, ref.PrimaryKey)
			} else if ref.PrimaryValue != "" {
				// key需要转为驼峰
				ins = ins.Where(localUtils.CamelCaseLowerFirst(ref.ForeignKey.DBName), "=", ref.PrimaryValue)
			} else {
				joinRelForeignFields = append(joinRelForeignFields, ref.ForeignKey)
				relForeignKeys = append(relForeignKeys, ref.PrimaryKey.DBName)
				relForeignFields = append(relForeignFields, ref.PrimaryKey)
			}
		}

		joinIdentityMap, joinForeignValues := schema.GetIdentityFieldValuesMap(reflectValue, foreignFields)
		if len(joinForeignValues) == 0 {
			return
		}

		// 如果用rel.JoinTable.MakeSlice().Elem()键不是驼峰, 导致无法赋值, 改为map
		// joinResults是一个带指针的slice形如[]*users
		joinResults := rel.JoinTable.MakeSlice().Elem()
		joinResultsType := joinResults.Type()
		// slice是指针, 内层还有指针
		joinResultType := joinResultsType.Elem().Elem()
		relatedRows := make([]map[string]interface{}, 0)
		column, values := schema.ToQueryValues(rel.JoinTable.Table, joinForeignKeys, joinForeignValues)
		// column可能的值有2种: 1.clause.Column 2.[]clause.Column
		cols := make([]clause.Column, 0)
		if col1, ok := column.(clause.Column); ok {
			cols = append(cols, col1)
		} else if col2, ok := column.([]clause.Column); ok {
			cols = append(cols, col2...)
		}
		for _, col := range cols {
			// key需要转为驼峰
			ins.AddError(ins.getInstance().Table(rel.JoinTable.Table).Where(localUtils.CamelCaseLowerFirst(col.Name), "in", toQueryValues(values)).Find(&relatedRows).Error)
		}

		// convert join identity map to relation identity map
		fieldValues := make([]interface{}, len(joinForeignFields))
		joinFieldValues := make([]interface{}, len(joinRelForeignFields))

		for _, result := range relatedRows {
			for idx, field := range joinForeignFields {
				fieldValues[idx] = result[localUtils.CamelCaseLowerFirst(field.Name)]
			}

			for idx, field := range joinRelForeignFields {
				joinFieldValues[idx] = result[localUtils.CamelCaseLowerFirst(field.Name)]
			}

			if results, ok := joinIdentityMap[utils.ToStringKey(fieldValues...)]; ok {
				joinKey := utils.ToStringKey(joinFieldValues...)
				identityMap[joinKey] = append(identityMap[joinKey], results...)
			}
			// new创建的元素为指针类型
			itemPtr := reflect.New(joinResultType)
			item := itemPtr.Elem()
			for _, field := range rel.JoinTable.PrimaryFields {
				// 由于redis存的json整数类型一律为float64, 这里需要做Convert强制转换
				// 创建数组项
				item.FieldByName(field.Name).Set(reflect.ValueOf(result[localUtils.CamelCaseLowerFirst(field.DBName)]).Convert(item.FieldByName(field.Name).Type()))
				joinResults = reflect.Append(joinResults, itemPtr)
			}
		}

		_, foreignValues = schema.GetIdentityFieldValuesMap(joinResults, joinRelForeignFields)
	} else {
		// 一对多或多对一关系
		for _, ref := range rel.References {
			if ref.OwnPrimaryKey {
				relForeignKeys = append(relForeignKeys, ref.ForeignKey.DBName)
				relForeignFields = append(relForeignFields, ref.ForeignKey)
				foreignFields = append(foreignFields, ref.PrimaryKey)
			} else if ref.PrimaryValue != "" {
				// key需要转为驼峰
				ins = ins.Where(localUtils.CamelCaseLowerFirst(ref.ForeignKey.DBName), "=", ref.PrimaryValue)
			} else {
				relForeignKeys = append(relForeignKeys, ref.PrimaryKey.DBName)
				relForeignFields = append(relForeignFields, ref.PrimaryKey)
				foreignFields = append(foreignFields, ref.ForeignKey)
			}
		}

		identityMap, foreignValues = schema.GetIdentityFieldValuesMap(reflectValue, foreignFields)
		if len(foreignValues) == 0 {
			return
		}
	}

	reflectResults := rel.FieldSchema.MakeSlice().Elem()
	column, values := schema.ToQueryValues(clause.CurrentTable, relForeignKeys, foreignValues)

	// column可能的值有2种: 1.clause.Column 2.[]clause.Column
	cols := make([]clause.Column, 0)
	if col1, ok := column.(clause.Column); ok {
		cols = append(cols, col1)
	} else if col2, ok := column.([]clause.Column); ok {
		cols = append(cols, col2...)
	}
	for _, col := range cols {
		ins.AddError(ins.getInstance().Table(rel.FieldSchema.Table).Where(localUtils.CamelCaseLowerFirst(col.Name), "in", toQueryValues(values)).Find(reflectResults.Addr().Interface()).Error)
	}
	fieldValues := make([]interface{}, len(relForeignFields))

	// clean up old values before preloading
	switch reflectValue.Kind() {
	case reflect.Struct:
		switch rel.Type {
		case schema.HasMany, schema.Many2Many:
			rel.Field.Set(reflectValue, reflect.MakeSlice(rel.Field.IndirectFieldType, 0, 0).Interface())
		default:
			rel.Field.Set(reflectValue, reflect.New(rel.Field.FieldType).Interface())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < reflectValue.Len(); i++ {
			switch rel.Type {
			case schema.HasMany, schema.Many2Many:
				rel.Field.Set(reflectValue.Index(i), reflect.MakeSlice(rel.Field.IndirectFieldType, 0, 0).Interface())
			default:
				rel.Field.Set(reflectValue.Index(i), reflect.New(rel.Field.FieldType).Interface())
			}
		}
	}

	for i := 0; i < reflectResults.Len(); i++ {
		elem := reflectResults.Index(i)
		for idx, field := range relForeignFields {
			fieldValues[idx], _ = field.ValueOf(elem)
		}

		for _, data := range identityMap[utils.ToStringKey(fieldValues...)] {
			reflectFieldValue := rel.Field.ReflectValueOf(data)
			if reflectFieldValue.Kind() == reflect.Ptr && reflectFieldValue.IsNil() {
				reflectFieldValue.Set(reflect.New(rel.Field.FieldType.Elem()))
			}

			reflectFieldValue = reflect.Indirect(reflectFieldValue)
			switch reflectFieldValue.Kind() {
			case reflect.Struct:
				rel.Field.Set(data, reflectResults.Index(i).Interface())
			case reflect.Slice, reflect.Array:
				if reflectFieldValue.Type().Elem().Kind() == reflect.Ptr {
					rel.Field.Set(data, reflect.Append(reflectFieldValue, elem).Interface())
				} else {
					rel.Field.Set(data, reflect.Append(reflectFieldValue, elem.Elem()).Interface())
				}
			}
		}
	}
}

// 值转换
func toQueryValues(values []interface{}) (results []int) {
	for _, v := range values {
		// redis存的json转换为int, 因此这里转一下类型
		if item, ok := v.(uint); ok {
			results = append(results, int(item))
		} else if item, ok := v.(int); ok {
			results = append(results, item)
		} else {
			fmt.Println("当前关联表主键类型不兼容")
		}
	}
	return
}
