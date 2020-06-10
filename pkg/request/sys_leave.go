package request

import (
	"gin-web/models"
	"gin-web/pkg/response"
)

// 获取请假列表结构体
type LeaveListRequestStruct struct {
	UserId            uint   `json:"userId"`
	Status            *uint  `json:"status" form:"status"`
	ApprovalOpinion   string `json:"approvalOpinion" form:"approvalOpinion"`
	Desc              string `json:"desc" form:"desc"`
	response.PageInfo        // 分页参数
}

// 创建请假结构体
type CreateLeaveRequestStruct struct {
	User models.SysUser `json:"user"`
	Desc string         `json:"desc" validate:"required"`
}

// 翻译需要校验的字段名称
func (s CreateLeaveRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "说明"
	return m
}
