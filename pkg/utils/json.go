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
