package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/tracing"
	"strings"
)

func (rd RedisService) FindRoleIdBySort(currentRoleSort uint) []uint {
	if !rd.binlog {
		return rd.mysql.FindRoleIdBySort(currentRoleSort)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "FindRoleIdBySort"))
	defer span.End()
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
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "FindRole"))
	defer span.End()
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

func (rd RedisService) GetRoleById(id uint) (rp models.SysRole) {
	if !rd.binlog {
		return rd.mysql.GetRoleById(id)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "GetRoleById"))
	defer span.End()
	rd.Q.
		Table("sys_role").
		Where("id", "=", id).
		Where("status", "=", models.SysRoleStatusNormal).
		First(&rp)
	return
}

func (rd RedisService) FindRoleByIds(ids []uint) []models.SysRole {
	if !rd.binlog {
		return rd.mysql.FindRoleByIds(ids)
	}
	_, span := tracer.Start(rd.Q.Ctx, tracing.Name(tracing.Cache, "FindRoleByIds"))
	defer span.End()
	roles := make([]models.SysRole, 0)
	rd.Q.
		Table("sys_role").
		Where("id", "in", ids).
		Where("status", "=", models.SysRoleStatusNormal).
		Find(&roles)
	return roles
}
