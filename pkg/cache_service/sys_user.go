package cache_service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/pkg/resp"
	"strings"
)

// user login check
func (rd RedisService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
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

func (rd RedisService) FindUser(req *request.UserReq) []models.SysUser {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindUser(req)
	}
	list := make([]models.SysUser, 0)
	query := rd.Q.
		Table("sys_user").
		Order("created_at DESC")
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds := rd.FindRoleIdBySort(*req.CurrentRole.Sort)
		query = query.Where("role_id", "in", roleIds)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		query = query.Where("username", "contains", username)
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		query = query.Where("mobile", "contains", mobile)
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		query = query.Where("nickname", "contains", nickname)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	rd.Q.FindWithPage(query, &req.Page, &list)
	return list
}

func (rd RedisService) GetUserById(id uint) (models.SysUser, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
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
