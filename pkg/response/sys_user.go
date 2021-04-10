package response

// User login response structure
type LoginResponseStruct struct {
	Username  string `json:"username"`  // 登录用户名
	Token     string `json:"token"`     // jwt令牌
	ExpiresAt int64  `json:"expiresAt"` // 过期时间, 秒
}

// 用户信息响应
type UserInfoResponseStruct struct {
	Id           uint     `json:"id"`
	Username     string   `json:"username"`
	Mobile       string   `json:"mobile"`
	Avatar       string   `json:"avatar"`
	Nickname     string   `json:"nickname"`
	Introduction string   `json:"introduction"`
	Roles        []string `json:"roles"`
	RoleSort     uint     `json:"roleSort"`
}

// 用户信息响应, 字段含义见models.SysUser
type UserListResponseStruct struct {
	BaseData
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       *uint  `json:"status"`
	RoleId       uint   `json:"roleId"`
	Creator      string `json:"creator"`
}
