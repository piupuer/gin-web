package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// LoginCheck user login check
func (rd RedisService) LoginCheck(r req.LoginCheck) (u models.SysUser, err error) {
	if !rd.binlog {
		return rd.mysql.LoginCheck(r)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "LoginCheck"))
	defer span.End()
	// Does the user exist
	err = rd.Q.Table("sys_user").Preload("Role").Where("username", "=", r.Username).First(&u).Error
	if err != nil {
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	// Need captcha
	flag := rd.mysql.Q.UserNeedCaptcha(req.UserNeedCaptcha{
		Wrong: u.Wrong,
	})
	if flag {
		if !rd.mysql.Q.VerifyCaptcha(r) {
			err = errors.Errorf(resp.InvalidCaptchaMsg)
			return
		}
	}
	// Is locked
	timestamp := time.Now().Unix()
	if u.Locked == constant.One && (u.LockExpire == 0 || timestamp < u.LockExpire) {
		err = errors.Errorf(resp.UserLockedMsg)
		return
	}
	// Verify password
	if ok := utils.ComparePwd(r.Password, u.Password); !ok {
		err = rd.mysql.UserWrongPwd(u)
		if err != nil {
			return
		}
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	err = rd.mysql.UserLastLogin(u.Id)
	return
}

func (rd RedisService) FindUser(r *request.User) []models.SysUser {
	if !rd.binlog {
		return rd.mysql.FindUser(r)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "FindUser"))
	defer span.End()
	list := make([]models.SysUser, 0)
	q := rd.Q.
		Table("sys_user").
		Order("created_at DESC")
	if *r.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds := rd.FindRoleIdBySort(*r.CurrentRole.Sort)
		q.Where("role_id", "in", roleIds)
	}
	username := strings.TrimSpace(r.Username)
	if username != "" {
		q.Where("username", "contains", username)
	}
	mobile := strings.TrimSpace(r.Mobile)
	if mobile != "" {
		q.Where("mobile", "contains", mobile)
	}
	nickname := strings.TrimSpace(r.Nickname)
	if nickname != "" {
		q.Where("nickname", "contains", nickname)
	}
	if r.Status != nil {
		q.Where("status", "=", *r.Status)
	}
	rd.Q.FindWithPage(q, &r.Page, &list)
	return list
}

func (rd RedisService) GetUserById(id uint) (rp models.SysUser) {
	if !rd.binlog {
		return rd.mysql.GetUserById(id)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "GetUserById"))
	defer span.End()
	rd.Q.
		Table("sys_user").
		Preload("Role").
		Where("id", "=", id).
		Where("status", "=", models.SysUserStatusEnable).
		First(&rp)
	return
}

func (rd RedisService) FindUserByIds(ids []uint) []models.SysUser {
	if !rd.binlog {
		return rd.mysql.FindUserByIds(ids)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "FindUserByIds"))
	defer span.End()
	list := make([]models.SysUser, 0)
	rd.Q.
		Table("sys_user").
		Where("id", "in", ids).
		Where("status", "=", models.SysUserStatusEnable).
		Find(&list)
	return list
}
