package response

import "time"

// 角色信息响应, 字段含义见models.SysRole
type RoleListResponseStruct struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Keyword   string    `json:"keyword"`
	Desc      string    `json:"desc"`
	Status    *bool     `json:"status"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
}
