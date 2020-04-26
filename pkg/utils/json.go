package utils

import (
	"encoding/json"
)

// 结构体转为json
func Struct2Json(obj interface{}) string {
	str, _ := json.Marshal(obj)
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	_ = json.Unmarshal([]byte(str), obj)
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
