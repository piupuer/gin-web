package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"gorm.io/gorm"
	"strings"
)

// FindLeave find leave by current user id
func (my MysqlService) FindLeave(r *request.Leave) []models.Leave {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindLeave"))
	defer span.End()
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
		q.Where("`desc` LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	my.Q.FindWithPage(q, &r.Page, &list)
	return list
}

// FindLeaveFsmTrack query leave fsm track
func (my MysqlService) FindLeaveFsmTrack(leaveId uint) []resp.FsmLogTrack {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindLeaveFsmTrack"))
	defer span.End()
	fsmUuid := my.GetLeaveFsmUuid(leaveId)
	return my.Q.FindFsmLogTrack(req.FsmLog{
		Category: req.NullUint(global.FsmCategoryLeave),
		Uuid:     fsmUuid,
	})
}

// CreateLeave create leave
func (my MysqlService) CreateLeave(r *request.CreateLeave) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "CreateLeave"))
	defer span.End()
	fsmUuid := uuid.NewString()
	// submit fsm log
	err = my.Q.FsmSubmitLog(req.FsmCreateLog{
		Category:        req.NullUint(global.FsmCategoryLeave),
		Uuid:            fsmUuid,
		SubmitterUserId: r.User.Id,
		SubmitterRoleId: r.User.RoleId,
	})
	if err != nil {
		return
	}

	// create leave to db
	var leave models.Leave
	copier.Copy(&leave, r)
	// save fsm uuid
	leave.FsmUuid = fsmUuid
	leave.Status = models.LevelStatusWaiting
	leave.UserId = r.User.Id
	my.Q.Tx.Create(&leave)
	return
}

func (my MysqlService) UpdateLeaveById(id uint, r request.UpdateLeave, u models.SysUser) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UpdateLeaveById"))
	defer span.End()
	var leave models.Leave
	my.Q.Tx.
		Where("id = ?", id).
		First(&leave)
	if leave.Id == constant.Zero {
		err = gorm.ErrRecordNotFound
		return
	}
	// check edit permission
	err = my.Q.FsmCheckEditLogDetailPermission(req.FsmCheckEditLogDetailPermission{
		Category:       req.NullUint(global.FsmCategoryLeave),
		Uuid:           leave.FsmUuid,
		ApprovalRoleId: u.RoleId,
		ApprovalUserId: u.Id,
		Fields:         []string{"desc", "start_time", "end_time"},
	})
	if err != nil {
		return
	}
	// update
	my.Q.UpdateById(id, r, new(models.Leave))
	return
}

// GetLeaveFsmUuid query leave fsm uuid by id
func (my MysqlService) GetLeaveFsmUuid(leaveId uint) string {
	// create leave to db
	var leave models.Leave
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("id = ?", leaveId).
		First(&leave)
	return leave.FsmUuid
}

// FindLevelByFsmUuids query leave by fsm uuids
func (my MysqlService) FindLevelByFsmUuids(uuids []string) (rp []models.Leave) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindLevelByFsmUuids"))
	defer span.End()
	// create leave to db
	rp = make([]models.Leave, 0)
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("uuid IN (?)", uuids).
		Find(&rp)
	return
}

func (my MysqlService) ApprovedLeaveById(r request.ApproveLeave) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "ApprovedLeaveById"))
	defer span.End()
	var leave models.Leave
	q := my.Q.Tx.
		Model(&models.Leave{}).
		Where("id = ?", r.Id)
	q.First(&leave)
	err = my.Q.FsmApproveLog(req.FsmApproveLog{
		Category:       req.NullUint(global.FsmCategoryLeave),
		Uuid:           leave.FsmUuid,
		ApprovalRoleId: r.User.RoleId,
		ApprovalUserId: r.User.Id,
		Approved:       req.NullUint(r.Approved),
	})
	return
}

func (my MysqlService) DeleteLeaveByIds(ids []uint, u models.SysUser) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "DeleteLeaveByIds"))
	defer span.End()
	list := make([]string, 0)
	my.Q.Tx.
		Model(&models.Leave{}).
		Where("id IN (?)", ids).
		Pluck("fsm_uuid", &list)
	if len(list) > 0 {
		my.Q.FsmCancelLogByUuids(req.FsmCancelLog{
			ApprovalRoleId: u.RoleId,
			ApprovalUserId: u.Id,
			Uuids:          list,
		})
	}
	err = my.Q.DeleteByIds(ids, new(models.Leave))
	return
}
