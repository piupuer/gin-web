package response

// User login response structure
type LoginResponseStruct struct {
	Username  string `json:"username"`  // 登录用户名
	Token     string `json:"token"`     // jwt令牌
	ExpiresAt int64  `json:"expiresAt"` // 过期时间, 秒
}
