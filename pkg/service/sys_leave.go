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

// 获取所有请假(当前用户)
func (s MysqlService) GetLeaves(r *request.LeaveReq) ([]models.SysLeave, error) {
	var err error
	list := make([]models.SysLeave, 0)
	query := s.tx.
		Model(&models.SysLeave{}).
		Order("created_at DESC").
		Where("user_id = ?", r.UserId)
	if r.Status != nil {
		query = query.Where("status = ?", *r.Status)
	}
	desc := strings.TrimSpace(r.Desc)
	if desc != "" {
		query = query.Where("desc LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	// 查询列表
	err = s.Find(query, &r.Page, &list)
	return list, err
}

// query leave fsm track
func (s MysqlService) FindLeaveApprovalLog(leaveId uint) ([]fsm.Log, error) {
	fsmUuid := s.GetLeaveFsmUuid(leaveId)
	f := fsm.New(s.tx)
	return f.FindLog(req.FsmLog{
		Category: req.NullUint(global.FsmCategoryLeave),
		Uuid:     fsmUuid,
	})
}

// query leave fsm track
func (s MysqlService) FindLeaveFsmTrack(leaveId uint) ([]resp.FsmLogTrack, error) {
	fsmUuid := s.GetLeaveFsmUuid(leaveId)
	f := fsm.New(s.tx)
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
func (s MysqlService) CreateLeave(r *request.CreateLeaveReq) error {
	f := fsm.New(s.tx)
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
	var leave models.SysLeave
	// save fsm uuid
	leave.FsmUuid = fsmUuid
	leave.Desc = r.Desc
	err = s.tx.Create(&leave).Error
	return err
}

// query leave fsm uuid by id
func (s MysqlService) GetLeaveFsmUuid(leaveId uint) string {
	// create leave to db
	var leave models.SysLeave
	s.tx.
		Model(&models.SysLeave{}).
		Where("id = ?", leaveId).
		First(&leave)
	return leave.FsmUuid
}
