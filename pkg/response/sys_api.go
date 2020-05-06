package response

import "time"

// 接口信息响应, 字段含义见models.SysRole
type ApiListResponseStruct struct {
	Id        uint      `json:"id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Category  string    `json:"category"`
	Creator   string    `json:"creator"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
}

// 权限接口信息响应, 字段含义见models.SysRole
type RoleApiListResponseStruct struct {
	Id     uint   `json:"id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Desc   string `json:"desc"`
	Access bool   `json:"access"` // 是否有权限访问
}
