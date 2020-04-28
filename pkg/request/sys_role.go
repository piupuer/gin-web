package request

import "go-shipment-api/pkg/response"

// 获取角色列表结构体
type RoleListRequestStruct struct {
	Name              string `json:"name" form:"name"` // 角色名称
	response.PageInfo        // 分页参数
}
