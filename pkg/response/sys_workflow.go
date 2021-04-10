package response

// 工作流日志信息响应, 字段含义见models.WorkflowLog
type WorkflowLogsListResponseStruct struct {
	BaseData
	FlowName              string `json:"name"`
	FlowId                uint   `json:"flowId"`
	FlowUuid              string `json:"flowUuid"`
	FlowCategory          uint   `json:"flowCategory"`
	FlowCategoryStr       string `json:"flowCategoryStr"`
	FlowTargetCategory    uint   `json:"flowTargetCategory"`
	FlowTargetCategoryStr string `json:"flowTargetCategoryStr"`
	TargetId              uint   `json:"targetId"`
	Status                *uint  `json:"status"`
	StatusStr             string `json:"statusStr"`
	SubmitUsername        string `json:"submitUsername"`
	SubmitUserNickname    string `json:"submitUserNickname"`
	SubmitDetail          string `json:"submitDetail"`
	ApprovalUsername      string `json:"approvalUsername"`
	ApprovalUserNickname  string `json:"approvalUserNickname"`
	ApprovalOpinion       string `json:"approvalOpinion"`
	ApprovingUserIds      []uint `json:"approvingUserIds"`
}

// 工作流信息响应, 字段含义见models.Workflow
type WorkflowListResponseStruct struct {
	BaseData
	Uuid              string `json:"uuid"`
	Category          uint   `json:"category"`
	SubmitUserConfirm *uint  `json:"submitUserConfirm"`
	TargetCategory    uint   `json:"targetCategory"`
	Self              *uint  `json:"self"`
	Name              string `json:"name"`
	Desc              string `json:"desc"`
	Creator           string `json:"creator"`
}

// 流水线信息响应, 字段含义见models.WorkflowLine
type WorkflowLineListResponseStruct struct {
	BaseData
	FlowId  uint   `json:"flowId"`
	Sort    uint   `json:"sort"`
	End     *uint  `json:"end"`
	RoleId  uint   `json:"roleId"`
	UserIds []uint `json:"userIds"`
	Edit    *uint  `json:"edit"`
	Name    string `json:"name"`
}
