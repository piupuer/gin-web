package utils

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 是否空字符串
func StrIsEmpty(str string) bool {
	return str == "null" || strings.TrimSpace(str) == ""
}

// 字符串转uint数组, 默认逗号分割
func Str2UintArr(str string) (ids []uint) {
	idArr := strings.Split(str, ",")
	for _, v := range idArr {
		ids = append(ids, Str2Uint(v))
	}
	return
}

// 字符串转int
func Str2Int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

// 字符串转uint
func Str2Uint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(num)
}

// 字符串转uint
func Str2Uint32(str string) uint32 {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(num)
}

// 字符串转uint
func Str2Bool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return b
}

// 字符串转float64
func Str2Float64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return num
}

// 字符串转uint数组, 默认逗号分割
func UintArr2IntArr(arr []uint) (newArr []int) {
	for _, v := range arr {
		newArr = append(newArr, int(v))
	}
	return
}

var (
	camelRe = regexp.MustCompile("(_)([a-zA-Z]+)")
	snakeRe = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// 字符串转为驼峰
func CamelCase(str string) string {
	camel := camelRe.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

// 字符串转为驼峰(首字母小写)
func CamelCaseLowerFirst(str string) string {
	camel := CamelCase(str)
	for i, v := range camel {
		return string(unicode.ToLower(v)) + camel[i+1:]
	}
	return camel
}

// 驼峰式写法转为下划线蛇形写法
func SnakeCase(str string) string {
	snake := snakeRe.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
