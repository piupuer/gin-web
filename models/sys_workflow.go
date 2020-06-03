package models

// 流程相关的常量
const (
	// 流程类别
	SysWorkflowCategoryOnlyOneApproval    uint   = 1 // 每个节点有一个人通过
	SysWorkflowCategoryAllApproval        uint   = 2 // 每个节点必须所有人审批通过
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
	Model
	Uuid              string `gorm:"unique;comment:'唯一标识'" json:"uuid"`
	Category          uint   `gorm:"default:1;comment:'类别(1:每个节点有一个人通过 2:每个节点必须所有人审批通过(指定了Users) 其他自行扩展)'" json:"category"`
	SubmitUserConfirm *bool  `gorm:"type:tinyint(1);default:0;comment:'是否需要提交人确认'" json:"submitUserConfirm"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	TargetCategory    uint   `gorm:"default:1;comment:'目标类别(1:请假(需要关联SysUser表) 其他自行扩展)'" json:"targetCategory"`
	Self              *bool  `gorm:"type:tinyint(1);default:0;comment:'是否可以自我审批(当前节点角色与可能提交人角色一致)'" json:"self"`
	Name              string `gorm:"comment:'名称'" json:"name"`
	Desc              string `gorm:"comment:'说明'" json:"desc"`
	Creator           string `gorm:"comment:'创建人'" json:"creator"`
}

func (m SysWorkflow) TableName() string {
	return m.Model.TableName("sys_workflow")
}

// 流程节点
type SysWorkflowNode struct {
	Model
	FlowId  uint        `gorm:"comment:'流程编号'" json:"flowId"`
	Flow    SysWorkflow `gorm:"foreignkey:FlowId" json:"flow"`
	RoleId  uint        `gorm:"comment:'审批人角色编号(拥有该角色才能审批)'" json:"roleId"`
	Role    SysRole     `gorm:"foreignkey:RoleId" json:"role"`
	Users   []SysUser   `gorm:"many2many:relation_user_workflow_node;comment:'审批人列表(指定了具体审批人, 则不再使用角色判断)'" json:"users"`
	Edit    *bool       `gorm:"type:tinyint(1);default:1;comment:'是否有编辑权限'" json:"edit"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Name    string      `gorm:"comment:'名称'" json:"name"`
	Creator string      `gorm:"comment:'创建人'" json:"creator"`
}

func (m SysWorkflowNode) TableName() string {
	return m.Model.TableName("sys_workflow_node")
}

// 用户与工作流节点关联关系
type RelationUserWorkflowNode struct {
	SysUserId         uint `json:"sysUserId"`
	SysWorkflowNodeId uint `json:"sysWorkflowNodeId"`
}

func (m RelationUserWorkflowNode) TableName() string {
	// 多对多关系表在tag中写死, 不能加自定义表前缀
	return "relation_user_workflow_node"
}

// 流程流水线
type SysWorkflowLine struct {
	Model
	FlowId uint              `gorm:"comment:'流程编号'" json:"flowId"`
	Flow   SysWorkflow       `gorm:"foreignkey:FlowId" json:"flow"`
	Sort   uint              `gorm:"comment:'排序'" json:"sort"`
	End    *bool             `gorm:"default:0;comment:'是否到达末尾'" json:"end"`
	Nodes  []SysWorkflowNode `gorm:"many2many:relation_workflow_line_node;comment:'审批节点列表(可能同一节点需多人无序审批)'" json:"nodes"`
}

func (m SysWorkflowLine) TableName() string {
	return m.Model.TableName("sys_workflow_line")
}

// 流水线与节点多对多关系
type RelationWorkflowLineNode struct {
	SysWorkflowLineId uint `json:"sysWorkflowLineId"`
	SysWorkflowNodeId uint `json:"sysWorkflowNodeId"`
}

func (m RelationWorkflowLineNode) TableName() string {
	// 多对多关系表在tag中写死, 不能加自定义表前缀
	return "relation_workflow_line_node"
}

// 流程日志: 任何一种工作流程都会关联到某一张表, 需要targetId
type SysWorkflowLog struct {
	Model
	FlowId          uint            `gorm:"comment:'流程编号'" json:"flowId"`
	Flow            SysWorkflow     `gorm:"foreignkey:FlowId" json:"flow"`
	TargetId        uint            `gorm:"comment:'目标表编号'" json:"targetId"`
	CurrentLineId   uint            `gorm:"comment:'当前审批线编号'" json:"currentLineId"`
	CurrentLine     SysWorkflowLine `gorm:"foreignkey:CurrentLineId" json:"currentLine"`
	Status          *uint           `gorm:"default:0;comment:'状态(0:提交 1:批准 2:拒绝 3:取消 4:重启 5:结束)'" json:"status"`
	SubmitUserId    uint            `gorm:"comment:'提交人编号'" json:"submitUserId"`
	SubmitUser      SysUser         `gorm:"foreignkey:SubmitUserId" json:"submitUser"`
	ApprovalId      uint            `gorm:"comment:'审批人编号'" json:"approvalId"`
	ApprovalUser    SysUser         `gorm:"foreignkey:ApprovalId" json:"approvalId"`
	ApprovalOpinion string          `gorm:"comment:'审批意见'" json:"approvalOpinion"`
}

func (m SysWorkflowLog) TableName() string {
	return m.Model.TableName("sys_workflow_log")
}
