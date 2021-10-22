package request

import (
	"gin-web/models"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type LeaveReq struct {
	UserId          uint          `json:"-"`
	Status          *req.NullUint `json:"status" form:"status"`
	ApprovalOpinion string        `json:"approvalOpinion" form:"approvalOpinion"`
	Desc            string        `json:"desc" form:"desc"`
	resp.Page
}

type CreateLeaveReq struct {
	User models.SysUser `json:"user"`
	Desc string         `json:"desc" validate:"required"`
}

func (s CreateLeaveReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "description"
	return m
}

type UpdateLeaveReq struct {
	Desc *string `json:"desc"`
}
