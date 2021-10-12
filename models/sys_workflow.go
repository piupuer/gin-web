package models

import (
	"github.com/piupuer/go-helper/models"
)

// 流程相关的常量
const (
	// 流程类别
	SysWorkflowCategoryOnlyOneApproval    uint   = 1 // 每个流水线有一个人通过
	SysWorkflowCategoryAllApproval        uint   = 2 // 每个流水线必须所有人审批通过
	SysWorkflowCategoryOnlyOneApprovalStr string = "只需一人通过"
	SysWorkflowCategoryAllApprovalStr     string = "必须全部通过"

	// 流程目标类别
	SysWorkflowTargetCategoryLeave    uint   = 1 // 请假
	SysWorkflowTargetCategoryLeaveStr string = "请假流程"

	// 流程日志状态
	SysWorkflowLogStateSubmit      uint   = 0 // 已提交
	SysWorkflowLogStateApproval    uint   = 1 // 通过
	SysWorkflowLogStateDeny        uint   = 2 // 拒绝
	SysWorkflowLogStateCancel      uint   = 3 // 取消
	SysWorkflowLogStateRestart     uint   = 4 // 重启
	SysWorkflowLogStateEnd         uint   = 5 // 结束
	SysWorkflowLogStateSubmitStr   string = "已提交"
	SysWorkflowLogStateApprovalStr string = "通过"
	SysWorkflowLogStateDenyStr     string = "拒绝"
	SysWorkflowLogStateCancelStr   string = "取消"
	SysWorkflowLogStateRestartStr  string = "重启"
	SysWorkflowLogStateEndStr      string = "结束"
)

// 定义map方便取值
var SysWorkflowCategoryConst = map[uint]string{
	SysWorkflowCategoryOnlyOneApproval: SysWorkflowCategoryOnlyOneApprovalStr,
	SysWorkflowCategoryAllApproval:     SysWorkflowCategoryAllApprovalStr,
}

var SysWorkflowTargetCategoryConst = map[uint]string{
	SysWorkflowTargetCategoryLeave: SysWorkflowTargetCategoryLeaveStr,
}

var SysWorkflowLogStateConst = map[uint]string{
	SysWorkflowLogStateSubmit:   SysWorkflowLogStateSubmitStr,
	SysWorkflowLogStateApproval: SysWorkflowLogStateApprovalStr,
	SysWorkflowLogStateDeny:     SysWorkflowLogStateDenyStr,
	SysWorkflowLogStateCancel:   SysWorkflowLogStateCancelStr,
	SysWorkflowLogStateRestart:  SysWorkflowLogStateRestartStr,
	SysWorkflowLogStateEnd:      SysWorkflowLogStateEndStr,
}

// 流程
type SysWorkflow struct {
	models.Model
	Uuid              string `gorm:"index:idx_uuid,unique;comment:'唯一标识'" json:"uuid"`
	Category          uint   `gorm:"default:1;comment:'类别(1:每个流水线有一个人通过 2:每个流水线必须所有人审批通过(指定了Users) 其他自行扩展)'" json:"category"`
	SubmitUserConfirm *uint  `gorm:"type:tinyint(1);default:0;comment:'是否需要提交人确认'" json:"submitUserConfirm"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	TargetCategory    uint   `gorm:"default:1;comment:'目标类别(1:请假(需要关联SysUser表) 其他自行扩展)'" json:"targetCategory"`
	Self              *uint  `gorm:"type:tinyint(1);default:0;comment:'是否可以自我审批(当前流水线角色与可能提交人角色一致)'" json:"self"`
	Name              string `gorm:"comment:'名称'" json:"name"`
	Desc              string `gorm:"comment:'说明'" json:"desc"`
	Creator           string `gorm:"comment:'创建人'" json:"creator"`
}

// 流程流水线
type SysWorkflowLine struct {
	models.Model
	FlowId uint        `gorm:"comment:'流程编号'" json:"flowId"`
	Flow   SysWorkflow `gorm:"foreignKey:FlowId" json:"flow"`
	Sort   uint        `gorm:"comment:'排序'" json:"sort"`
	End    *uint       `gorm:"default:0;comment:'是否到达末尾'" json:"end"`
	RoleId uint        `gorm:"comment:'审批人角色编号(拥有该角色才能审批)'" json:"roleId"`
	Role   SysRole     `gorm:"foreignKey:RoleId" json:"role"`
	Users  []SysUser   `gorm:"many2many:sys_workflow_line_user_relation;comment:'审批人列表(指定了具体审批人, 则不再使用角色判断)'" json:"users"`
	Edit   *uint       `gorm:"type:tinyint(1);default:1;comment:'是否有编辑权限'" json:"edit"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Name   string      `gorm:"comment:'名称'" json:"name"`
}

// 用户与流水线关联关系
type SysWorkflowLineUserRelation struct {
	SysUserId         uint `json:"sysUserId"`
	SysWorkflowLineId uint `json:"sysWorkflowLineId"`
}

// 流程日志: 任何一种工作流程都会关联到某一张表, 需要targetId
type SysWorkflowLog struct {
	models.Model
	FlowId           uint            `gorm:"comment:'流程编号'" json:"flowId"`
	Flow             SysWorkflow     `gorm:"foreignKey:FlowId" json:"flow"`
	TargetId         uint            `gorm:"comment:'目标表编号'" json:"targetId"`
	CurrentLineId    uint            `gorm:"comment:'当前审批线编号'" json:"currentLineId"`
	CurrentLine      SysWorkflowLine `gorm:"foreignKey:CurrentLineId" json:"currentLine"`
	Status           *uint           `gorm:"default:0;comment:'状态(0:提交 1:批准 2:拒绝 3:取消 4:重启 5:结束)'" json:"status"`
	End              *uint           `gorm:"default:0;comment:'是否到达末尾'" json:"end"`
	SubmitUserId     uint            `gorm:"comment:'提交人编号'" json:"submitUserId"`
	SubmitUser       SysUser         `gorm:"foreignKey:SubmitUserId" json:"submitUser"`
	SubmitDetail     string          `gorm:"comment:'提交明细(待审批可避免二次查询)'" json:"submitDetail"`
	ApprovalUserId   uint            `gorm:"comment:'审批人编号'" json:"approvalUserId"`
	ApprovalUser     SysUser         `gorm:"foreignKey:ApprovalUserId" json:"approvalUser"`
	ApprovalOpinion  string          `gorm:"comment:'审批意见'" json:"approvalOpinion"`
	ApprovingUserIds []uint          `gorm:"-" json:"approvingUserIds"` // status为0提交时有效, 表示审批人列表, 无需保存到数据库
}
