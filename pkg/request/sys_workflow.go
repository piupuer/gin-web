package request

import (
	"gin-web/pkg/response"
)

type WorkflowReq struct {
	Uuid              string `json:"uuid" form:"uuid"`
	Category          *uint  `json:"category" form:"category"`
	SubmitUserConfirm *uint  `json:"submitUserConfirm" form:"submitUserConfirm"`
	TargetCategory    *uint  `json:"targetCategory" form:"targetCategory"`
	Self              *uint  `json:"self" form:"self"`
	Name              string `json:"name" form:"name"`
	Desc              string `json:"desc" form:"desc"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

type WorkflowLineReq struct {
	FlowId            uint `json:"flowId" form:"flowId"`
	response.PageInfo      // 分页参数
}

type WorkflowApprovingReq struct {
	ApprovalUserId    uint `json:"approvalUserId"`
	response.PageInfo      // 分页参数
}

type CreateWorkflowReq struct {
	Uuid              string  `json:"uuid"`
	Category          ReqUint `json:"category"`
	SubmitUserConfirm ReqUint `json:"submitUserConfirm"`
	TargetCategory    ReqUint `json:"targetCategory"`
	Self              ReqUint `json:"self"`
	Name              string  `json:"name" validate:"required"`
	Desc              string  `json:"desc"`
	Creator           string  `json:"creator"`
}

func (s CreateWorkflowReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "流程名称"
	return m
}

type UpdateWorkflowReq struct {
	Category          *ReqUint `json:"category"`
	SubmitUserConfirm *ReqUint `json:"submitUserConfirm"`
	TargetCategory    *ReqUint `json:"targetCategory"`
	Self              *ReqUint `json:"self"`
	Name              *string  `json:"name"`
	Desc              *string  `json:"desc"`
}

type UpdateWorkflowLineReq struct {
	Id      uint    `json:"id"`
	FlowId  ReqUint `json:"flowId" validate:"required"`
	RoleId  ReqUint `json:"roleId"`
	UserIds []uint  `json:"userIds"`
	Name    string  `json:"name" validate:"required"`
	Edit    ReqUint `json:"edit"`
}

func (s UpdateWorkflowLineReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Id"] = "流水线"
	m["FlowId"] = "流程"
	m["RoleId"] = "审批人所属角色"
	m["UserIds"] = "审批人"
	m["Name"] = "流水线名称"
	m["Edit"] = "编辑权限"
	return m
}

type UpdateWorkflowLineIncrementalReq struct {
	FlowId uint                    `json:"flowId" validate:"required"`
	Create []UpdateWorkflowLineReq `json:"create"` // 需要新增的流水线编号集合
	Update []UpdateWorkflowLineReq `json:"update"` // 需要新增的流水线编号集合
	Delete []UpdateWorkflowLineReq `json:"delete"` // 需要删除的流水线编号集合
}

func (s UpdateWorkflowLineIncrementalReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["FlowId"] = "流程号"
	return m
}

type WorkflowTransitionReq struct {
	FlowId          uint   `json:"flowId" validate:"required"`
	TargetCategory  uint   `json:"targetCategory" validate:"required"`
	TargetId        uint   `json:"targetId" validate:"required"`
	SubmitUserId    uint   `json:"submitUserId"`
	SubmitDetail    string `json:"submitDetail"`
	ApprovalUserId  uint   `json:"approvalUserId"`
	ApprovalOpinion string `json:"approvalOpinion"`
	ApprovalStatus  *uint  `json:"approvalStatus" validate:"required,min=1,max=4"`
}

func (s WorkflowTransitionReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["FlowId"] = "流程"
	m["TargetCategory"] = "目标类别"
	m["TargetId"] = "目标编号"
	m["ApprovalStatus"] = "审批状态"
	return m
}
