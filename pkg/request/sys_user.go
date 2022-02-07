package request

import (
	"gin-web/models"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type RegisterAndLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s RegisterAndLogin) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Username"] = "username"
	m["Password"] = "password"
	return m
}

type ChangePwd struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

type User struct {
	CurrentRole  models.SysRole `json:"currentRole" swaggerignore:"true"`
	Username     string         `json:"username" form:"username"`
	Mobile       string         `json:"mobile" form:"mobile"`
	Avatar       string         `json:"avatar" form:"avatar"`
	Nickname     string         `json:"nickname" form:"nickname"`
	UsernameOr   string         `json:"usernameOr" form:"usernameOr"`
	MobileOr     string         `json:"mobileOr" form:"mobileOr"`
	NicknameOr   string         `json:"nicknameOr" form:"nicknameOr"`
	Introduction string         `json:"introduction" form:"introduction"`
	Status       *uint          `json:"status" form:"status"`
	RoleId       uint           `json:"roleId" form:"roleId"`
	resp.Page
}

type CreateUser struct {
	Username     string        `json:"username" validate:"required"`
	Password     string        `json:"password"`
	InitPassword string        `json:"initPassword" validate:"required"`
	NewPassword  string        `json:"newPassword"`
	Mobile       string        `json:"mobile" validate:"required"`
	Avatar       string        `json:"avatar"`
	Nickname     string        `json:"nickname"`
	Introduction string        `json:"introduction"`
	Status       *req.NullUint `json:"status"`
	RoleId       uint          `json:"roleId" validate:"required"`
}

func (s CreateUser) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Username"] = "username"
	m["InitPassword"] = "initialize password"
	m["Mobile"] = "mobile number"
	m["RoleId"] = "role id"
	return m
}

type UpdateUser struct {
	Username     *string       `json:"username"`
	Password     *string       `json:"password"`
	InitPassword *string       `json:"initPassword"`
	NewPassword  *string       `json:"newPassword"`
	Mobile       *string       `json:"mobile"`
	Avatar       *string       `json:"avatar"`
	Nickname     *string       `json:"nickname"`
	Introduction *string       `json:"introduction"`
	Status       *req.NullUint `json:"status"`
	Locked       *req.NullUint `json:"locked"`
	LockExpire   *int64        `json:"lockExpire"`
	Wrong        *int          `json:"wrong"`
	RoleId       *uint         `json:"roleId"`
}
