package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
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
		MachineId:       1,
		SubmitterUserId: r.User.Id,
		SubmitterRoleId: r.User.RoleId,
	})
	if err != nil {
		return err
	}

	// create leave to db
	var leave models.Leave
	// save fsm uuid
	leave.FsmUuid = fsmUuid
	leave.Desc = r.Desc
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
