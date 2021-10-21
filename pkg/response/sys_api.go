package response

import "github.com/piupuer/go-helper/pkg/resp"

type ApiResp struct {
	resp.Base
	Method   string `json:"method"`
	Path     string `json:"path"`
	Category string `json:"category"`
	Desc     string `json:"desc"`
	Title    string `json:"title"`
}

type ApiGroupByCategoryResp struct {
	Title    string    `json:"title"`    // 标题
	Category string    `json:"category"` // 分组名称
	Children []ApiResp `json:"children"` // 前端以树形图结构展示, 这里用children表示
}

type ApiTreeWithAccessResp struct {
	List      []ApiGroupByCategoryResp `json:"list"`
	AccessIds []uint                   `json:"accessIds"`
}
