package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"github.com/jinzhu/copier"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// find leave by current user id
func (my MysqlService) FindLeave(r *request.Leave) []models.Leave {
	list := make([]models.Leave, 0)
	q := my.Q.Tx.
		Model(&models.Leave{}).
		Order("created_at DESC").
		Where("user_id = ?", r.UserId)
	if r.Status != nil {
		q.Where("status = ?", *r.Status)
	}
	desc := strings.TrimSpace(r.Desc)
	if desc != "" {
		q.Where("desc LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	my.Q.FindWithPage(q, &r.Page, &list)
	return list
}

// query leave fsm track
func (my MysqlService) FindLeaveApprovalLog(leaveId uint) ([]fsm.Log, error) {
	fsmUuid := my.GetLeaveFsmUuid(leaveId)
	f := fsm.New(my.Q.Tx)
	return f.FindLog(req.FsmLog{
		Category: req.NullUint(global.FsmCategoryLeave),
		Uuid:     fsmUuid,
	})
}

// query leave fsm track
func (my MysqlService) FindLeaveFsmTrack(leaveId uint) ([]resp.FsmLogTrack, error) {
	fsmUuid := my.GetLeaveFsmUuid(leaveId)
	f := fsm.New(my.Q.Tx)
	logs, err := f.FindLog(req.FsmLog{
		Category: req.NullUint(global.FsmCategoryLeave),
		Uuid:     fsmUuid,
	})
	if err != nil {
		return nil, err
	}
	return f.FindLogTrack(logs)
}

// create leave
func (my MysqlService) CreateLeave(r *request.CreateLeave) error {
	f := fsm.New(my.Q.Tx)
	fsmUuid := uuid.NewV4().String()
	// submit fsm log
	_, err := f.SubmitLog(req.FsmCreateLog{
		Category:        req.NullUint(global.FsmCategoryLeave),
		Uuid:            fsmUuid,
		SubmitterUserId: r.User.Id,
		SubmitterRoleId: r.User.RoleId,
	})
	if err != nil {
		return err
	}

	// create leave to db
	var leave models.Leave
	copier.Copy(&leave, r)
	// save fsm uuid
	leave.FsmUuid = fsmUuid
	leave.Status = models.LevelStatusWaiting
	leave.UserId = r.User.Id
	err = my.Q.Tx.Create(&leave).Error
	return err
}

// query leave fsm uuid by id
func (my MysqlService) GetLeaveFsmUuid(leaveId uint) string {
	// create leave to db
	var leave models.Leave
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("id = ?", leaveId).
		First(&leave)
	return leave.FsmUuid
}

// query leave by fsm uuids
func (my MysqlService) FindLevelByFsmUuids(uuids []string) []models.Leave {
	// create leave to db
	leaves := make([]models.Leave, 0)
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("uuid IN (?)", uuids).
		Find(&leaves)
	return leaves
}

func (my MysqlService) ApprovedLeaveById(r request.ApproveLeave) (err error) {
	var leave models.Leave
	q := my.Q.Tx.
		Model(&models.Leave{}).
		Where("id = ?", r.Id)
	err = q.First(&leave).Error
	if err != nil {
		return err
	}
	f := fsm.New(my.Q.Tx)
	var log *resp.FsmApprovalLog
	log, err = f.ApproveLog(req.FsmApproveLog{
		Category:       req.NullUint(global.FsmCategoryLeave),
		Uuid:           leave.FsmUuid,
		ApprovalRoleId: r.User.RoleId,
		ApprovalUserId: r.User.Id,
		Approved:       req.NullUint(r.Approved),
	})
	if err != nil {
		return
	}
	// log status transition
	return my.LeaveTransition(*log)
}

func (my MysqlService) DeleteLeaveByIds(ids []uint) (err error) {
	list := make([]string, 0)
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("id IN (?)", ids).
		Pluck("fsm_uuid", &list)
	if len(list) > 0 {
		f := fsm.New(my.Q.Tx)
		_, err = f.CancelLogByUuids(list)
		if err != nil {
			return
		}
	}
	return my.Q.DeleteByIds(ids, new(models.Leave))
}

func (my MysqlService) LeaveTransition(logs ...resp.FsmApprovalLog) (err error) {
	m := make(map[uint][]string)
	for _, log := range logs {
		if log.Category == global.FsmCategoryLeave {
			if log.Resubmit {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = item
				}
				m[models.LevelStatusRefused] = append(arr, log.Uuid)
			} else if log.Cancel {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusCancelled]; ok {
					arr = item
				}
				m[models.LevelStatusCancelled] = append(arr, log.Uuid)
			} else if log.Confirm {
				arr := make([]string, 0)
				if item, ok := m[models.LevelStatusWaitingConfirm]; ok {
					arr = item
				}
				m[models.LevelStatusWaitingConfirm] = append(arr, log.Uuid)
			} else if log.End {
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
			return
		}
	}
	return nil
}

func (my MysqlService) GetLeaveFsmDetail(detail req.FsmSubmitterDetail) []resp.FsmSubmitterDetail {
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

func (my MysqlService) UpdateLeaveFsmDetail(detail req.UpdateFsmSubmitterDetail) (err error) {
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
			return
		}
	}
	return nil
}
