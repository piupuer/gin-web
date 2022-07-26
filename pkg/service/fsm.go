package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
)

func (my MysqlService) FsmTransition(logs ...resp.FsmApprovalLog) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FsmTransition"))
	defer span.End()
	m := make(map[uint][]string)
	for _, log := range logs {
		if log.Category == global.FsmCategoryLeave {
			if log.Resubmit == constant.One {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = item
				}
				m[models.LevelStatusRefused] = append(arr, log.Uuid)
			} else if log.Cancel == constant.One {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusCancelled]; ok {
					arr = item
				}
				m[models.LevelStatusCancelled] = append(arr, log.Uuid)
			} else if log.Confirm == constant.One {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = item
				}
				m[models.LevelStatusWaitingConfirm] = append(arr, log.Uuid)
			} else if log.End == constant.One {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusApproved]; ok {
					arr = item
				}
				m[models.LevelStatusApproved] = append(arr, log.Uuid)
			} else {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusApproving]; ok {
					arr = item
				}
				m[models.LevelStatusApproving] = append(arr, log.Uuid)
			}
		}
	}
	for status, uuids := range m {
		err = my.Q.Tx.
			Model(&models.Leave{}).
			Where("fsm_uuid IN (?)", uuids).
			Update("status", status).Error
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (my MysqlService) GetFsmDetail(detail req.FsmSubmitterDetail) []resp.FsmSubmitterDetail {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "GetFsmDetail"))
	defer span.End()
	arr := make([]resp.FsmSubmitterDetail, 0)
	switch uint(detail.Category) {
	case global.FsmCategoryLeave:
		var leave models.Leave
		my.Q.Tx.
			Model(&models.Leave{}).
			Where("fsm_uuid = ?", detail.Uuid).
			First(&leave)
		if leave.Id > 0 {
			arr = append(arr, resp.FsmSubmitterDetail{
				Name: "leave desc",
				Key:  "desc",
				Val:  leave.Desc,
			})
			if !leave.StartTime.IsZero() {
				arr = append(arr, resp.FsmSubmitterDetail{
					Name: "leave start time",
					Key:  "startTime",
					Val:  leave.StartTime.String(),
				})
			}
			if !leave.EndTime.IsZero() {
				arr = append(arr, resp.FsmSubmitterDetail{
					Name: "leave end time",
					Key:  "endTime",
					Val:  leave.EndTime.String(),
				})
			}
		}
	}
	return arr
}

func (my MysqlService) UpdateFsmDetail(detail req.UpdateFsmSubmitterDetail) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UpdateFsmDetail"))
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
		if leave.Id > 0 {
			err = q.Updates(&m).Error
			return errors.WithStack(err)
		}
	}
	return nil
}
