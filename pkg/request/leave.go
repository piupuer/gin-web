package request

import (
	"gin-web/models"
	"github.com/golang-module/carbon/v2"
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
	User      models.SysUser  `json:"user" swaggerignore:"true"`
	Desc      string          `json:"desc" validate:"required"`
	StartTime carbon.DateTime `json:"startTime" swaggertype:"string"`
	EndTime   carbon.DateTime `json:"endTime" swaggertype:"string"`
}

func (s CreateLeave) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "description"
	return m
}

type UpdateLeave struct {
	Desc      *string          `json:"desc"`
	StartTime *carbon.DateTime `json:"startTime" swaggertype:"string"`
	EndTime   *carbon.DateTime `json:"endTime" swaggertype:"string"`
}

type ApproveLeave struct {
	Id          uint           `json:"id"`
	AfterStatus uint           `json:"afterStatus"`
	Approved    uint           `json:"approved"`
	User        models.SysUser `json:"user" swaggerignore:"true"`
}
