package response

import "github.com/piupuer/go-helper/pkg/resp"

type LoginResp struct {
	Username  string `json:"username"`  // 登录用户名
	Token     string `json:"token"`     // jwt令牌
	ExpiresAt int64  `json:"expiresAt"` // 过期时间, 秒
}

type UserInfoResp struct {
	Id           uint     `json:"id"`
	Username     string   `json:"username"`
	Mobile       string   `json:"mobile"`
	Avatar       string   `json:"avatar"`
	Nickname     string   `json:"nickname"`
	Introduction string   `json:"introduction"`
	Roles        []string `json:"roles"`
	RoleSort     uint     `json:"roleSort"`
	Keyword      string   `json:"keyword"`
}

type UserResp struct {
	resp.Base
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       *uint  `json:"status"`
	RoleId       uint   `json:"roleId"`
}
