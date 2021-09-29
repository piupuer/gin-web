package response

type LeaveResp struct {
	BaseData
	Status *uint  `json:"status"`
	Desc   string `json:"desc"`
}

type LeaveLogResp struct {
	LeaveId uint            `json:"leaveId"`
	Log     WorkflowLogResp `json:"log"`
}
