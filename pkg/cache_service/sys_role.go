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

func (rd RedisService) FindRole(r *request.Role) []models.SysRole {
	if !rd.binlog {
		return rd.mysql.FindRole(r)
	}
	list := make([]models.SysRole, 0)
	q := rd.Q.
		Table("sys_role").
		Order("created_at DESC").
		Where("sort", ">=", r.CurrentRoleSort)
	name := strings.TrimSpace(r.Name)
	if name != "" {
		q.Where("name", "contains", name)
	}
	keyword := strings.TrimSpace(r.Keyword)
	if keyword != "" {
		q.Where("keyword", "contains", keyword)
	}
	if r.Status != nil {
		if *r.Status > 0 {
			q.Where("status", "=", 1)
		} else {
			q.Where("status", "=", 0)
		}
	}
	rd.Q.FindWithPage(q, &r.Page, &list)
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
