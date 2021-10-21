package response

import "github.com/piupuer/go-helper/pkg/resp"

type MenuTreeResp struct {
	resp.Base
	ParentId   uint           `json:"parentId"`
	Name       string         `json:"name"`
	Title      string         `json:"title"`
	Icon       string         `json:"icon"`
	Path       string         `json:"path"`
	Redirect   string         `json:"redirect"`
	Component  string         `json:"component"`
	Permission string         `json:"permission"`
	Sort       int            `json:"sort"`
	Status     uint           `json:"status"`
	Visible    uint           `json:"visible"`
	Breadcrumb uint           `json:"breadcrumb"`
	Children   []MenuTreeResp `json:"children"`
}

type MenuTreeWithAccessResp struct {
	List      []MenuTreeResp `json:"list"`
	AccessIds []uint         `json:"accessIds"`
}
