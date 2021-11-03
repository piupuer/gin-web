package request

import (
	"gin-web/models"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type Leave struct {
	UserId          uint          `json:"-"`
	Status          *req.NullUint `json:"status" form:"status"`
	ApprovalOpinion string        `json:"approvalOpinion" form:"approvalOpinion"`
	Desc            string        `json:"desc" form:"desc"`
	resp.Page
}

type CreateLeave struct {
	User models.SysUser `json:"user"`
	Desc string         `json:"desc" validate:"required"`
}

func (s CreateLeave) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "description"
	return m
}

type UpdateLeave struct {
	Desc *string `json:"desc"`
}
