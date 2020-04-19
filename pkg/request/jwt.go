package request

import "github.com/dgrijalva/jwt-go"

// 保留登录用户信息
type JwtClaimsStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}
