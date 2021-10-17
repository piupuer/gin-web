package response

import (
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/resp"
)

type LeaveResp struct {
	resp.Base
	Status *uint  `json:"status"`
	Desc   string `json:"desc"`
}

type LeaveLogResp struct {
	LeaveId uint    `json:"leaveId"`
	Log     fsm.Log `json:"log"`
}
