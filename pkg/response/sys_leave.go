package response

import (
	"gin-web/models"
)

// 请假信息响应, 字段含义见models.SysLeave
type LeaveListResponseStruct struct {
	Id        uint             `json:"id"`
	Status    *uint            `json:"status"`
	Desc      string           `json:"desc"`
	CreatedAt models.LocalTime `json:"createdAt"`
}

// 请假日志信息响应
type LeaveLogListResponseStruct struct {
	LeaveId uint                           `json:"leaveId"`
	Log     WorkflowLogsListResponseStruct `json:"log"`
}
