package request

import (
	"gin-web/pkg/response"
)

// 获取流程列表结构体
type WorkflowListRequestStruct struct {
	Uuid              string `json:"uuid" form:"uuid"`
	Category          uint   `json:"category" form:"category"`
	SubmitUserConfirm *bool  `json:"submitUserConfirm" form:"submitUserConfirm"`
	TargetCategory    uint   `json:"targetCategory" form:"targetCategory"`
	Self              *bool  `json:"self" form:"self"`
	Name              string `json:"name" form:"name"`
	Desc              string `json:"desc" form:"desc"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 创建流程结构体
type CreateWorkflowRequestStruct struct {
	Category          uint   `json:"category"`
	SubmitUserConfirm *bool  `json:"submitUserConfirm"`
	TargetCategory    uint   `json:"targetCategory"`
	Self              *bool  `json:"self"`
	Name              string `json:"name" validate:"required"`
	Desc              string `json:"desc"`
	Creator           string `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateWorkflowRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "流程名称"
	return m
}

// 创建流程节点结构体
type CreateWorkflowNodeRequestStruct struct {
	FlowId  uint   `json:"flowId" validate:"required"`
	RoleId  uint   `json:"roleId" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Desc    string `json:"desc"`
	Creator string `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateWorkflowNodeRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["FlowId"] = "流程"
	m["RoleId"] = "审批人角色"
	m["Name"] = "节点名称"
	return m
}

// 工作流转移结构体
type WorkflowTransitionRequestStruct struct {
	FlowId          uint   `json:"targetId"`
	TargetCategory  uint   `json:"targetCategory" validate:"required"`
	TargetId        uint   `json:"targetId" validate:"required"`
	SubmitUserId    uint   `json:"submitUserId"`
	ApprovalId      uint   `json:"approvalId"`
	ApprovalOpinion string `json:"approvalOpinion" validate:"required"`
	ApprovalStatus  *uint  `json:"status" validate:"required||in=0,1,2,3,4"`
}

// 翻译需要校验的字段名称
func (s WorkflowTransitionRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["TargetCategory"] = "目标类别"
	m["TargetId"] = "目标编号"
	m["ApprovalOpinion"] = "审批意见"
	m["ApprovalStatus"] = "审批状态"
	return m
}
