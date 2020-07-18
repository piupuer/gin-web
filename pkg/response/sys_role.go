package response

import (
	"gin-web/models"
)

// 角色信息响应, 字段含义见models
type RoleListResponseStruct struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Keyword   string           `json:"keyword"`
	Desc      string           `json:"desc"`
	Status    *bool            `json:"status"`
	Creator   string           `json:"creator"`
	CreatedAt models.LocalTime `json:"createdAt"`
}
