package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"strings"
)

func (rd RedisService) FindRoleIdBySort(currentRoleSort uint) []uint {
	if !rd.binlog {
		return rd.mysql.FindRoleIdBySort(currentRoleSort)
	}
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	rd.Q.
		Table("sys_role").
		Where("sort", ">=", currentRoleSort).
		Find(&roles)
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds
}

func (rd RedisService) FindRole(req *request.RoleReq) []models.SysRole {
	if !rd.binlog {
		return rd.mysql.FindRole(req)
	}
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
	rd.Q.FindWithPage(query, &req.Page, &list)
	return list
}

func (rd RedisService) GetRoleById(id uint) (models.SysRole, error) {
	if !rd.binlog {
		return rd.mysql.GetRoleById(id)
	}
	var role models.SysRole
	var err error
	err = rd.Q.
		Table("sys_role").
		Where("id", "=", id).
		Where("status", "=", models.SysRoleStatusNormal).
		First(&role).Error
	return role, err
}

func (rd RedisService) FindRoleByIds(ids []uint) []models.SysRole {
	if !rd.binlog {
		return rd.mysql.FindRoleByIds(ids)
	}
	roles := make([]models.SysRole, 0)
	rd.Q.
		Table("sys_role").
		Where("id", "in", ids).
		Where("status", "=", models.SysRoleStatusNormal).
		Find(&roles)
	return roles
}
