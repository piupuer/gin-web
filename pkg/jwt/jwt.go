package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"go-shipment-api/pkg/request"
)

var jwtSecret = []byte("jwt-secret")

// 根据用户名和密码生成token
func GenerateToken(claims *request.JwtClaimsStruct) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

// 解析token
func ParseToken(token string) (*request.JwtClaimsStruct, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &request.JwtClaimsStruct{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*request.JwtClaimsStruct); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
