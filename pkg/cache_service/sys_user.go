package cache_service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
	"strings"
)

// user login check
func (rd RedisService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	if !rd.binlog {
		return rd.mysql.LoginCheck(user)
	}
	var u models.SysUser
	// Does the user exist
	err := rd.Q.Table("sys_user").Preload("Role").Where("username", "=", user.Username).First(&u).Error
	if err != nil {
		return nil, errors.New(resp.LoginCheckErrorMsg)
	}
	// Verify password
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New(resp.LoginCheckErrorMsg)
	}
	return &u, err
}

func (rd RedisService) FindUser(req *request.User) []models.SysUser {
	if !rd.binlog {
		return rd.mysql.FindUser(req)
	}
	list := make([]models.SysUser, 0)
	q := rd.Q.
		Table("sys_user").
		Order("created_at DESC")
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds := rd.FindRoleIdBySort(*req.CurrentRole.Sort)
		q.Where("role_id", "in", roleIds)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		q.Where("username", "contains", username)
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		q.Where("mobile", "contains", mobile)
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		q.Where("nickname", "contains", nickname)
	}
	if req.Status != nil {
		q.Where("status", "=", *req.Status)
	}
	rd.Q.FindWithPage(q, &req.Page, &list)
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
