package utils

import (
	"fmt"
	"testing"
)

func TestGenPwd(t *testing.T) {
	var s1 = "123456"
	var s2 = "654321"
	var s3 = "123456"
	pw1 := GenPwd(s1)
	pw2 := GenPwd(s2)
	pw3 := GenPwd(s3)
	fmt.Println(fmt.Sprintf("明文1: %s, 加密为密文1: %s", s1, pw1))
	fmt.Println(fmt.Sprintf("明文2: %s, 加密为密文2: %s", s2, pw2))
	fmt.Println(fmt.Sprintf("比较明文1: %s, 密文2: %s, 是否出自同一明文: %v", s1, pw2, ComparePwd(s1, pw2)))

	fmt.Println(fmt.Sprintf("明文3: %s, 加密结果: %s", s3, pw3))
	fmt.Println(fmt.Sprintf("比较明文1: %s, 密文3: %s, 是否出自同一明文: %v", s1, pw3, ComparePwd(s1, pw3)))
}
