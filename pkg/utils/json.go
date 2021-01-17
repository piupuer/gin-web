package utils

import (
	"encoding/json"
	"fmt"
	"gin-web/models"
	"github.com/shopspring/decimal"
	"reflect"
)

// 结构体转为json
func Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) {
	// 将json interface转为string
	jsonStr, _ := str.(string)
	Json2Struct(jsonStr, obj)
}

// 结构体转结构体, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
func Struct2StructByJson(struct1 interface{}, struct2 interface{}) {
	// 转换为响应结构体, 隐藏部分字段
	jsonStr := Struct2Json(struct1)
	Json2Struct(jsonStr, struct2)
}

// 两结构体比对不同的字段, 不同时将取newStruct中的字段返回, json为中间桥梁
func CompareDifferenceStructByJson(oldStruct interface{}, newStruct interface{}, update *map[string]interface{}) {
	// 通过json先将其转为map集合
	m1 := make(map[string]interface{}, 0)
	m2 := make(map[string]interface{}, 0)
	m3 := make(map[string]interface{}, 0)
	Struct2StructByJson(newStruct, &m1)
	Struct2StructByJson(oldStruct, &m2)
	for k1, v1 := range m1 {
		for k2, v2 := range m2 {
			switch v1.(type) {
			// 复杂结构不做对比
			case map[string]interface{}:
				continue
			}
			rv := reflect.ValueOf(v1)
			// 值类型必须有效
			if rv.Kind() != reflect.Invalid {
				// key相同, 值不同
				if k1 == k2 && v1 != v2 {
					t := reflect.TypeOf(oldStruct)
					key := CamelCase(k1)
					var fieldType reflect.Type
					oldStructV := reflect.ValueOf(oldStruct)
					// map与结构体取值方式不同
					if oldStructV.Kind() == reflect.Map {
						mapV := oldStructV.MapIndex(reflect.ValueOf(k1))
						if !mapV.IsValid() {
							break
						}
						fieldType = mapV.Type()
					} else if oldStructV.Kind() == reflect.Struct {
						structField, ok := t.FieldByName(key)
						if !ok {
							break
						}
						fieldType = structField.Type
					} else {
						// oldStruct类型不对, 直接跳过不处理
						break
					}
					// 取到结构体对应字段
					realT := fieldType
					// 指针类型需要剥掉一层获取真实类型
					if fieldType.Kind() == reflect.Ptr {
						realT = fieldType.Elem()
					}
					// 获得元素
					e := reflect.New(realT).Elem()
					// 不同类型不一定可以强制转换
					switch e.Interface().(type) {
					case decimal.Decimal:
						d, _ := decimal.NewFromString(rv.String())
						m3[k1] = d
					case models.LocalTime:
						t := new(models.LocalTime).SetString(rv.String())
						// 时间过滤空值
						if !t.IsZero() {
							m3[k1] = *t
						}
					default:
						// 强制转换rv赋值给e
						e.Set(rv.Convert(realT))
						m3[k1] = e.Interface()
					}
					break
				}
			}
		}
	}
	*update = m3
}

// 两结构体比对不同的字段, 将key转为蛇形
func CompareDifferenceStruct2SnakeKeyByJson(oldStruct interface{}, newStruct interface{}, update *map[string]interface{}) {
	m1 := make(map[string]interface{}, 0)
	m2 := make(map[string]interface{}, 0)
	CompareDifferenceStructByJson(oldStruct, newStruct, &m1)
	for key, item := range m1 {
		m2[SnakeCase(key)] = item
	}
	*update = m2
}
