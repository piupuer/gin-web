package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

func (rd RedisService) FindRoleIdBySort(currentRoleSort uint) ([]uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindRoleIdBySort(currentRoleSort)
	}
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	err := rd.Q.Table("sys_role").Where("sort", ">=", currentRoleSort).Find(&roles).Error
	if err != nil {
		return roleIds, err
	}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds, nil
}

func (rd RedisService) FindRole(req *request.RoleReq) ([]models.SysRole, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindRole(req)
	}
	var err error
	list := make([]models.SysRole, 0)
	query := rd.Q.
		Table("sys_role").
		Order("created_at DESC").
		Where("sort", ">=", req.CurrentRoleSort)
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name", "contains", name)
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		query = query.Where("keyword", "contains", keyword)
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status", "=", 1)
		} else {
			query = query.Where("status", "=", 0)
		}
	}
	err = rd.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}
