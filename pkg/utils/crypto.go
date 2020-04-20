package utils

import "golang.org/x/crypto/bcrypt"

// 生成密码, 由于使用自适应hash算法, 不可逆
func GenPwd(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hash)
}

// 通过比较两个字符串hash判断是否出自同一个明文
// str 明文
// pwd 需要对比的密文
func ComparePwd(str string, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(str)); err != nil {
		return false
	}
	return true
}
