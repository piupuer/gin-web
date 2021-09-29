package request

import (
	"gin-web/models"
	"gin-web/pkg/response"
)

type LeaveReq struct {
	UserId            uint   `json:"-"`
	Status            *uint  `json:"status" form:"status"`
	ApprovalOpinion   string `json:"approvalOpinion" form:"approvalOpinion"`
	Desc              string `json:"desc" form:"desc"`
	response.PageInfo        // 分页参数
}

type CreateLeaveReq struct {
	User models.SysUser `json:"user"`
	Desc string         `json:"desc" validate:"required"`
}

func (s CreateLeaveReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "说明"
	return m
}

type UpdateLeaveReq struct {
	Desc *string `json:"desc"`
}
