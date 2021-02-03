package cache_service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

// 登录校验
func (s *RedisService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.LoginCheck(user)
	}
	var u models.SysUser
	// 查询用户及其角色
	err := s.redis.Table(new(models.SysUser).TableName()).Preload("Role").Where("username", "=", user.Username).First(&u).Error
	if err != nil {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	// 校验密码
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	return &u, err
}

func (s *RedisService) GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetUsers(req)
	}
	var err error
	list := make([]models.SysUser, 0)
	query := s.redis.
		Table(new(models.SysUser).TableName()).
		Order("created_at DESC")
	// 非超级管理员
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds, err := s.GetRoleIdsBySort(*req.CurrentRole.Sort)
		if err != nil {
			return list, err
		}
		query = query.Where("roleId", "in", roleIds)
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
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	err = query.Count(&req.PageInfo.Total).Error
	if err == nil && req.PageInfo.Total > 0 {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 获取单个用户
func (s *RedisService) GetUserById(id uint) (models.SysUser, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetUserById(id)
	}
	var user models.SysUser
	var err error
	err = s.redis.
		Table(new(models.SysUser).TableName()).
		Preload("Role").
		Where("id", "=", id).
		// 状态为正常
		Where("status", "=", models.SysUserStatusNormal).
		First(&user).Error
	return user, err
}
