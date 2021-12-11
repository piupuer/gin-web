package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
	"strings"
)

// user login check
func (rd RedisService) LoginCheck(user *models.SysUser) (u models.SysUser, err error) {
	if !rd.binlog {
		return rd.mysql.LoginCheck(user)
	}
	// Does the user exist
	err = rd.Q.Table("sys_user").Preload("Role").Where("username", "=", user.Username).First(&u).Error
	if err != nil {
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	// Verify password
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		err = rd.mysql.UserWrongPwd(u)
		if err != nil {
			return
		}
		err = errors.Errorf(resp.LoginCheckErrorMsg)
		return
	}
	err = rd.mysql.UserLastLogin(user.Id)
	return
}

func (rd RedisService) FindUser(r *request.User) []models.SysUser {
	if !rd.binlog {
		return rd.mysql.FindUser(r)
	}
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

func (rd RedisService) GetUserById(id uint) (models.SysUser, error) {
	if !rd.binlog {
		return rd.mysql.GetUserById(id)
	}
	var user models.SysUser
	var err error
	err = rd.Q.
		Table("sys_user").
		Preload("Role").
		Where("id", "=", id).
		Where("status", "=", models.SysUserStatusEnable).
		First(&user).Error
	return user, err
}

func (rd RedisService) FindUserByIds(ids []uint) []models.SysUser {
	if !rd.binlog {
		return rd.mysql.FindUserByIds(ids)
	}
	list := make([]models.SysUser, 0)
	rd.Q.
		Table("sys_user").
		Where("id", "in", ids).
		Where("status", "=", models.SysUserStatusEnable).
		Find(&list)
	return list
}
