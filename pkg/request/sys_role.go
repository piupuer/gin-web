package request

import (
	"gin-web/pkg/response"
)

// 获取角色列表结构体
type RoleListRequestStruct struct {
	Name              string `json:"name" form:"name"`
	Keyword           string `json:"keyword" form:"keyword"`
	CurrentRoleSort   uint   `json:"currentRoleSort"`
	Status            *uint  `json:"status" form:"status"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 创建角色结构体
type CreateRoleRequestStruct struct {
	CurrentRoleSort uint     `json:"currentRoleSort"`
	Name            string   `json:"name" validate:"required"`
	Keyword         string   `json:"keyword" validate:"required"`
	Sort            *ReqUint `json:"sort" validate:"required"`
	Desc            string   `json:"desc"`
	Status          *ReqUint `json:"status"`
	Creator         string   `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateRoleRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "角色名称"
	m["Keyword"] = "角色关键字"
	m["Sort"] = "角色排序"
	return m
}
