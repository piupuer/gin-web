package response

import (
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/resp"
)

type Leave struct {
	resp.Base
	Status    *uint           `json:"status"`
	FsmUuid   string          `json:"fsmUuid"`
	Desc      string          `json:"desc"`
	StartTime carbon.DateTime `json:"startTime"`
	EndTime   carbon.DateTime `json:"endTime"`
}

type LeaveLog struct {
	LeaveId uint    `json:"leaveId"`
	Log     fsm.Log `json:"log"`
}
