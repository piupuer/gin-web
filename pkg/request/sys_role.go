package request

import "go-shipment-api/pkg/response"

// 获取角色列表结构体
type RoleListRequestStruct struct {
	Name              string `json:"name" form:"name"` // 角色名称
	response.PageInfo        // 分页参数
}

// 创建角色结构体
type CreateRoleRequestStruct struct {
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	Desc    string `json:"desc"`
	Status  *bool  `json:"status"`
	Creator string `json:"creator"`
}
