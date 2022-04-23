package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/ms"
)

const (
	LevelStatusWaiting        uint = iota // waiting approval
	LevelStatusApproved                   // approved
	LevelStatusRefused                    // approval rejection
	LevelStatusCancelled                  // approval cancelled
	LevelStatusApproving                  // pending approval but not end
	LevelStatusWaitingConfirm             // approval waiting confirm
)

type Leave struct {
	ms.M
	UserId          uint            `gorm:"comment:user id(SysUser.Id)" json:"userId"`
	User            SysUser         `gorm:"foreignKey:UserId" json:"user"`
	FsmUuid         string          `gorm:"size:100;comment:finite state machine uuid" json:"fsmUuid"`
	Status          uint            `gorm:"default:0;comment:status(0:submitted 1:approved 2:refused 3:cancel 4:approving 5:waiting confirm)" json:"status"`
	ApprovalOpinion string          `gorm:"comment:approval opinion or remark" json:"approvalOpinion"`
	Desc            string          `gorm:"comment:submitter description" json:"desc"`
	StartTime       carbon.DateTime `gorm:"comment:start time" json:"startTime"`
	EndTime         carbon.DateTime `gorm:"comment:end time" json:"endTime"`
}
