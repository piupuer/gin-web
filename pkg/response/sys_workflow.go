package response

import "gin-web/models"

// 工作流日志信息响应, 字段含义见models.WorkflowLog
type WorkflowLogsListResponseStruct struct {
	FlowName              string           `json:"name"`
	FlowUuid              string           `json:"flowUuid"`
	FlowCategoryStr       string           `json:"flowCategoryStr"`
	FlowTargetCategoryStr string           `json:"flowTargetCategoryStr"`
	Status                *uint            `json:"status"`
	StatusStr             string           `json:"statusStr"`
	SubmitUsername        string           `json:"submitUsername"`
	SubmitUserNickname    string           `json:"submitUserNickname"`
	ApprovalUsername      string           `json:"approvalUsername"`
	ApprovalUserNickname  string           `json:"approvalUserNickname"`
	ApprovalOpinion       string           `json:"approvalOpinion"`
	CreatedAt             models.LocalTime `json:"createdAt"`
	UpdatedAt             models.LocalTime `json:"updatedAt"`
}

// 工作流信息响应, 字段含义见models.Workflow
type WorkflowListResponseStruct struct {
	Id                uint             `json:"id"`
	Uuid              string           `json:"uuid"`
	Category          uint             `json:"category"`
	SubmitUserConfirm *bool            `json:"submitUserConfirm"`
	TargetCategory    uint             `json:"targetCategory"`
	Self              *bool            `json:"self"`
	Name              string           `json:"name"`
	Desc              string           `json:"desc"`
	Creator           string           `json:"creator"`
	CreatedAt         models.LocalTime `json:"createdAt"`
}
