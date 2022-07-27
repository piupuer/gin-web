package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func (my MysqlService) LoginCheck(r req.LoginCheck) (u models.SysUser, err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "LoginCheck"))
	defer span.End()
	err = my.Q.Tx.Preload("Role").Where("username = ?", r.Username).First(&u).Error
	if err != nil {
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	flag := my.Q.UserNeedCaptcha(req.UserNeedCaptcha{
		Wrong: u.Wrong,
	})
	if flag {
		if !my.Q.VerifyCaptcha(r) {
			err = errors.Errorf(resp.InvalidCaptchaMsg)
			return
		}
	}
	timestamp := time.Now().Unix()
	if u.Locked == constant.One && (u.LockExpire == 0 || timestamp < u.LockExpire) {
		err = errors.Errorf(resp.UserLockedMsg)
		return
	}
	if ok := utils.ComparePwd(r.Password, u.Password); !ok {
		err = my.UserWrongPwd(u)
		if err != nil {
			return
		}
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	err = my.UserLastLogin(u.Id)
	return
}

func (my MysqlService) UserWrongPwd(user models.SysUser) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UserWrongPwd"))
	defer span.End()
	// do not use transaction
	q := my.Q.Db.
		Model(&models.SysUser{}).
		Where("id = ?", user.Id)
	m := make(map[string]interface{})
	newWrong := user.Wrong + 1
	if newWrong >= 10 {
		m["locked"] = constant.One
		if newWrong == 10 {
			m["lock_expire"] = carbon.Now().AddDuration("10m").Carbon2Time().Unix()
		} else if newWrong == 20 {
			m["lock_expire"] = carbon.Now().AddDuration("60m").Carbon2Time().Unix()
		} else if newWrong >= 30 {
			m["lock_expire"] = 0
		}
	}
	m["wrong"] = newWrong
	err = q.Updates(&m).Error
	return
}

func (my MysqlService) UserLastLogin(id uint) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UserLastLogin"))
	defer span.End()
	m := make(map[string]interface{})
	m["wrong"] = constant.Zero
	m["last_login"] = carbon.Now()
	m["locked"] = constant.Zero
	m["lock_expire"] = constant.Zero
	err = my.Q.Tx.
		Model(&models.SysUser{}).
		Where("id = ?", id).
		Updates(&m).Error
	return
}

func (my MysqlService) FindUser(r *request.User) []models.SysUser {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindUser"))
	defer span.End()
	list := make([]models.SysUser, 0)
	q := my.Q.Tx.
		Model(&models.SysUser{}).
		Order("created_at DESC")
	if *r.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds := my.FindRoleIdBySort(*r.CurrentRole.Sort)
		q.Where("role_id IN (?)", roleIds)
	}
	username := strings.TrimSpace(r.Username)
	if username != "" {
		q.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(r.Mobile)
	if mobile != "" {
		q.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	nickname := strings.TrimSpace(r.Nickname)
	if nickname != "" {
		q.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if r.Status != nil {
		if *r.Status > 0 {
			q.Where("status = ?", 1)
		} else {
			q.Where("status = ?", 0)
		}
	}
	my.Q.FindWithPage(q, &r.Page, &list)
	return list
}

func (my MysqlService) GetUserById(id uint) (rp models.SysUser) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "GetUserById"))
	defer span.End()
	my.Q.Tx.Preload("Role").
		Where("id = ?", id).
		Where("status = ?", models.SysUserStatusEnable).
		First(&rp)
	return
}

func (my MysqlService) GetUserByUsername(username string) (models.SysUser, error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "GetUserByUsername"))
	defer span.End()
	var user models.SysUser
	var err error
	err = my.Q.Tx.Preload("Role").
		Where("username = ?", username).
		Where("status = ?", models.SysUserStatusEnable).
		First(&user).Error
	return user, err
}

func (my MysqlService) FindUserByIds(ids []uint) []models.SysUser {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindUserByIds"))
	defer span.End()
	list := make([]models.SysUser, 0)
	my.Q.Tx.
		Model(&models.SysUser{}).
		Where("id IN (?)", ids).
		Where("status = ?", models.SysUserStatusEnable).
		Find(&list)
	return list
}
