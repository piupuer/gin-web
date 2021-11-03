package response

import (
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/resp"
)

type Leave struct {
	resp.Base
	Status *uint  `json:"status"`
	Desc   string `json:"desc"`
}

type LeaveLog struct {
	LeaveId uint    `json:"leaveId"`
	Log     fsm.Log `json:"log"`
}
