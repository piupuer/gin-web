package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 字符串转uint数组, 默认逗号分割
func Str2UintArr(str string) (ids []uint) {
	idArr := strings.Split(str, ",")
	for _, v := range idArr {
		ids = append(ids, Str2Uint(v))
	}
	return
}

// 字符串转uint
func Str2Uint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(num)
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

// 使用gizp压缩字符串
func Str2BytesByGzip(str string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(str))
	gz.Close()
	return b.Bytes()
}

// 使用gizp压缩字符串
func Bytes2StrByGzip(b []byte) string {
	data := bytes.NewReader(b)
	r, _ := gzip.NewReader(data)
	s, _ := ioutil.ReadAll(r)
	return string(s)
}
