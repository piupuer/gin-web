package models

import "github.com/piupuer/go-helper/models"

// 请假记录表
type SysLeave struct {
	models.M
	UserId          uint    `gorm:"comment:'用户编号'" json:"userId"`
	User            SysUser `gorm:"foreignKey:UserId" json:"user"`
	FsmUuid         string  `gorm:"size:100;comment:'状态机uuid'" json:"fsmUuid"`
	Status          uint    `gorm:"default:0;comment:'状态(0:提交 1:批准 2:拒绝 3:取消 4:结束)'" json:"status"`
	ApprovalOpinion string  `gorm:"comment:'审批意见'" json:"approvalOpinion"`
	Desc            string  `gorm:"comment:'说明'" json:"desc"`
}
