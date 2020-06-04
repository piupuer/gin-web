package response

import (
	"gin-web/models"
)

// 请假信息响应, 字段含义见models.SysRole
type LeaveListResponseStruct struct {
	Id        uint             `json:"id"`
	Status    *uint            `json:"status"`
	Desc      string           `json:"desc"`
	CreatedAt models.LocalTime `json:"createdAt"`
}
