package response

// 请假信息响应, 字段含义见models.SysLeave
type LeaveListResponseStruct struct {
	BaseData
	Status *uint  `json:"status"`
	Desc   string `json:"desc"`
}

// 请假日志信息响应
type LeaveLogListResponseStruct struct {
	LeaveId uint                           `json:"leaveId"`
	Log     WorkflowLogsListResponseStruct `json:"log"`
}
