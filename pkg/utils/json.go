package utils

import (
	"encoding/json"
	"fmt"
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

// 两结构体比对不同的字段, 不同时将取struct1中的字段返回, json为中间桥梁, update必须以指针方式传递, 否则可能获取到空数据
func CompareDifferenceStructByJson(oldStruct interface{}, newStruct interface{}, update interface{}) {
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
			// key相同, 值不同
			if k1 == k2 && v1 != v2 {
				m3[k1] = v1
				break
			}
		}
	}
	Struct2StructByJson(m3, &update)
}
