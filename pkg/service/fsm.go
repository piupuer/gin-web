package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
)

func (my MysqlService) FsmTransition(logs ...resp.FsmApprovalLog) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FsmTransition"))
	defer span.End()
	m := make(map[uint][]string)
	for _, item := range logs {
		switch item.Category {
		case global.FsmCategoryLeave:
			if item.Resubmit == constant.One {
				arr := make([]string, 0)
				if v, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = v
				}
				m[models.LevelStatusRefused] = append(arr, item.Uuid)
			} else if item.Cancel == constant.One {
				arr := make([]string, 0)
				if v, ok := m[models.LevelStatusCancelled]; ok {
					arr = v
				}
				m[models.LevelStatusCancelled] = append(arr, item.Uuid)
			} else if item.Confirm == constant.One {
				arr := make([]string, 0)
				if v, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = v
				}
				m[models.LevelStatusWaitingConfirm] = append(arr, item.Uuid)
			} else if item.End == constant.One {
				arr := make([]string, 0)
				if v, ok := m[models.LevelStatusApproved]; ok {
					arr = v
				}
				m[models.LevelStatusApproved] = append(arr, item.Uuid)
			} else {
				arr := make([]string, 0)
				if v, ok := m[models.LevelStatusApproving]; ok {
					arr = v
				}
				m[models.LevelStatusApproving] = append(arr, item.Uuid)
			}
		}
	}
	for status, uuids := range m {
		my.Q.Tx.
			Model(&models.Leave{}).
			Where("fsm_uuid IN (?)", uuids).
			Update("status", status)
	}
	return
}

func (my MysqlService) GetFsmLogDetail(detail req.FsmLogSubmitterDetail) (rp []resp.FsmLogSubmitterDetail) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "GetFsmLogDetail"))
	defer span.End()
	rp = make([]resp.FsmLogSubmitterDetail, 0)
	switch uint(detail.Category) {
	case global.FsmCategoryLeave:
		var leave models.Leave
		my.Q.Tx.
			Model(&models.Leave{}).
			Where("fsm_uuid = ?", detail.Uuid).
			First(&leave)
		if leave.Id > constant.Zero {
			rp = append(rp, resp.FsmLogSubmitterDetail{
				Name: "leave desc",
				Key:  "desc",
				Val:  leave.Desc,
			})
			if !leave.StartTime.IsZero() {
				rp = append(rp, resp.FsmLogSubmitterDetail{
					Name: "leave start time",
					Key:  "startTime",
					Val:  leave.StartTime.String(),
				})
			}
			if !leave.EndTime.IsZero() {
				rp = append(rp, resp.FsmLogSubmitterDetail{
					Name: "leave end time",
					Key:  "endTime",
					Val:  leave.EndTime.String(),
				})
			}
		}
	}
	return
}

func (my MysqlService) UpdateFsmLogDetail(detail req.UpdateFsmLogSubmitterDetail) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UpdateFsmLogDetail"))
	defer span.End()
	switch uint(detail.Category) {
	case global.FsmCategoryLeave:
		detail.Parse()
		m := make(map[string]interface{})
		for i, key := range detail.Keys {
			m[utils.SnakeCase(key)] = detail.Vals[i]
		}
		var leave models.Leave
		q := my.Q.Tx.
			Model(&models.Leave{}).
			Where("fsm_uuid = ?", detail.Uuid)
		q.First(&leave)
		if leave.Id > constant.Zero {
			q.Updates(&m)
			return
		}
	}
	return
}
