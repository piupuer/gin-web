package request

import "go-shipment-api/pkg/response"

// User login structure
type RegisterAndLoginRequestStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 修改密码结构体
type ChangePwdRequestStruct struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// 获取用户列表结构体
type UserListRequestStruct struct {
	Id                uint   `json:"id" form:"id"`
	Username          string `json:"username" form:"username"`
	Mobile            string `json:"mobile" form:"mobile"`
	Avatar            string `json:"avatar" form:"avatar"`
	Nickname          string `json:"nickname" form:"nickname"`
	Introduction      string `json:"introduction" form:"introduction"`
	Status            *bool  `json:"status" form:"status"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 创建用户结构体
type CreateUserRequestStruct struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password"`
	Mobile       string `json:"mobile" validate:"required"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       *bool  `json:"status"`
	Creator      string `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateUserRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["username"] = "用户名"
	m["mobile"] = "手机号"
	return m
}
