package response

// 菜单树信息响应基础字段, 字段含义见models.SysMenu
type MenuTreeBaseResponseStruct struct {
	Id         uint   `json:"id"`
	ParentId   uint   `json:"parentId"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Redirect   string `json:"redirect"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Creator    string `json:"creator"`
	Sort       int    `json:"sort"`
	Status     bool   `json:"status"`
	Visible    bool   `json:"visible"`
	Breadcrumb bool   `json:"breadcrumb"`
}

// 菜单树信息响应
type MenuTreeResponseStruct struct {
	MenuTreeBaseResponseStruct
	Children []MenuTreeResponseStruct `json:"children"`
}

// 菜单树信息响应, 包含access字段
type MenuTreeWithAccessResponseStruct struct {
	MenuTreeBaseResponseStruct
	Children []MenuTreeWithAccessResponseStruct `json:"children"`
	Access   bool                               `json:"access"` // 是否有权限访问
}
