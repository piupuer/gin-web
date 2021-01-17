package request

import "gin-web/pkg/response"

// 获取菜单列表结构体
type MenuListRequestStruct struct {
	Name              string `json:"name" form:"name"`
	Title             string `json:"title" form:"title"`
	Path              string `json:"path" form:"path"`
	Component         string `json:"component" form:"component"`
	Redirect          string `json:"redirect"`
	Status            *uint  `json:"status" form:"status"`
	Visible           *uint  `json:"visible" form:"visible"`
	Breadcrumb        *uint  `json:"breadcrumb" form:"breadcrumb"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 创建菜单结构体
type CreateMenuRequestStruct struct {
	Name       string  `json:"name" validate:"required"`
	Title      string  `json:"title"`
	Icon       string  `json:"icon"`
	Path       string  `json:"path"`
	Redirect   string  `json:"redirect"`
	Component  string  `json:"component"`
	Permission string  `json:"permission"`
	Sort       ReqUint `json:"sort"`
	Status     ReqUint `json:"status"`
	Visible    ReqUint `json:"visible"`
	Breadcrumb ReqUint `json:"breadcrumb"`
	ParentId   ReqUint `json:"parentId"`
	Creator    string  `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateMenuRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "菜单名称"
	return m
}

// 更新菜单结构体
type UpdateMenuRequestStruct struct {
	Name       *string  `json:"name"`
	Title      *string  `json:"title"`
	Icon       *string  `json:"icon"`
	Path       *string  `json:"path"`
	Redirect   *string  `json:"redirect"`
	Component  *string  `json:"component"`
	Permission *string  `json:"permission"`
	Sort       *ReqUint `json:"sort"`
	Status     *ReqUint `json:"status"`
	Visible    *ReqUint `json:"visible"`
	Breadcrumb *ReqUint `json:"breadcrumb"`
	ParentId   *ReqUint `json:"parentId"`
}
