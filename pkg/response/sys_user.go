package response

import "github.com/piupuer/go-helper/pkg/resp"

type Login struct {
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type UserInfo struct {
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

type User struct {
	resp.Base
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       *uint  `json:"status"`
	Locked       uint   `json:"locked"`
	RoleId       uint   `json:"roleId"`
}
