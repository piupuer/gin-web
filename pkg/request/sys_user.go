package request

import (
	"gin-web/models"
	"gin-web/pkg/response"
)

// User login structure
type RegisterAndLoginRequestStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 修改密码结构体
type ChangePwdRequestStruct struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

// 获取用户列表结构体
type UserListRequestStruct struct {
	CurrentRole       models.SysRole `json:"currentRole"`
	Username          string         `json:"username" form:"username"`
	Mobile            string         `json:"mobile" form:"mobile"`
	Avatar            string         `json:"avatar" form:"avatar"`
	Nickname          string         `json:"nickname" form:"nickname"`
	Introduction      string         `json:"introduction" form:"introduction"`
	Status            *uint          `json:"status" form:"status"`
	RoleId            uint           `json:"roleId" form:"roleId"`
	Creator           string         `json:"creator" form:"creator"`
	response.PageInfo                // 分页参数
}

// 创建用户结构体
type CreateUserRequestStruct struct {
	Username     string   `json:"username" validate:"required"`
	InitPassword string   `json:"initPassword" validate:"required"` // 不使用SysUser的Password字段, 避免请求劫持绕过系统校验
	NewPassword  string   `json:"newPassword"`
	Mobile       string   `json:"mobile" validate:"required"`
	Avatar       string   `json:"avatar"`
	Nickname     string   `json:"nickname"`
	Introduction string   `json:"introduction"`
	Status       *ReqUint `json:"status"`
	RoleId       uint     `json:"roleId" validate:"required"`
	Creator      string   `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateUserRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Username"] = "用户名"
	m["InitPassword"] = "初始密码"
	m["Mobile"] = "手机号"
	m["RoleId"] = "角色"
	return m
}
