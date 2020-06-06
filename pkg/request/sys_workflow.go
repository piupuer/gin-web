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

// 获取流水线结构体
type WorkflowLineListRequestStruct struct {
	FlowId            uint `json:"flowId" form:"flowId"`
	response.PageInfo      // 分页参数
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

// 更新流程节点结构体
type UpdateWorkflowNodeRequestStruct struct {
	Id      uint   `json:"id"`
	FlowId  uint   `json:"flowId" validate:"required"`
	RoleId  *uint  `json:"roleId"`
	UserIds []uint `json:"userIds"`
	Name    string `json:"name" validate:"required"`
	Edit    *bool  `json:"edit"`
	Creator string `json:"creator"`
}

// 翻译需要校验的字段名称
func (s UpdateWorkflowNodeRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Id"] = "节点"
	m["FlowId"] = "流程"
	m["RoleId"] = "审批人所属角色"
	m["UserIds"] = "审批人"
	m["Name"] = "节点名称"
	m["Edit"] = "编辑权限"
	return m
}

// 更新流水线结构体
type UpdateWorkflowLineRequestStruct struct {
	FlowId  uint                              `json:"flowId" validate:"required"`
	Create  []UpdateWorkflowNodeRequestStruct `json:"create"` // 需要新增的节点编号集合
	Update  []UpdateWorkflowNodeRequestStruct `json:"update"` // 需要新增的节点编号集合
	Delete  []UpdateWorkflowNodeRequestStruct `json:"delete"` // 需要删除的节点编号集合
	Creator string                            `json:"creator"`
}

// 翻译需要校验的字段名称
func (s UpdateWorkflowLineRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["FlowId"] = "流程号"
	return m
}

// 工作流转移结构体
type WorkflowTransitionRequestStruct struct {
	FlowId          uint   `json:"targetId"`
	TargetCategory  uint   `json:"targetCategory" validate:"required"`
	TargetId        uint   `json:"targetId" validate:"required"`
	SubmitUserId    uint   `json:"submitUserId"`
	ApprovalUserId  uint   `json:"approvalUserId"`
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
