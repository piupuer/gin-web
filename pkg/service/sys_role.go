package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/ms"
	"github.com/pkg/errors"
	"strings"
)

func (my MysqlService) FindRoleIdBySort(currentRoleSort uint) []uint {
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	my.Q.Tx.
		Model(&models.SysRole{}).
		Where("sort >= ?", currentRoleSort).
		Find(&roles)
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds
}

func (my MysqlService) FindRole(r *request.Role) []models.SysRole {
	list := make([]models.SysRole, 0)
	q := my.Q.Tx.
		Model(&models.SysRole{}).
		Order("created_at DESC").
		Where("sort >= ?", r.CurrentRoleSort)
	name := strings.TrimSpace(r.Name)
	if name != "" {
		q.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(r.Keyword)
	if keyword != "" {
		q.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
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

func (my MysqlService) DeleteRoleByIds(ids []uint) (err error) {
	var roles []models.SysRole
	my.Q.Tx.
		Preload("Users").
		Where("id IN (?)", ids).
		Find(&roles)
	newIds := make([]uint, 0)
	oldCasbins := make([]ms.SysRoleCasbin, 0)
	for _, v := range roles {
		if len(v.Users) > 0 {
			return errors.Errorf("role %s has %d associated users, please delete the user before deleting the role", v.Name, len(v.Users))
		}
		oldCasbins = append(oldCasbins, my.Q.FindRoleCasbin(ms.SysRoleCasbin{
			Keyword: v.Keyword,
		})...)
		newIds = append(newIds, v.Id)
	}
	if len(oldCasbins) > 0 {
		my.Q.BatchDeleteRoleCasbin(oldCasbins)
	}
	if len(newIds) > 0 {
		err = my.Q.DeleteByIds(newIds, new(models.SysRole))
		err = errors.WithStack(err)
	}
	return
}

func (my MysqlService) GetRoleById(id uint) (models.SysRole, error) {
	var role models.SysRole
	var err error
	err = my.Q.Tx.
		Where("id = ?", id).
		Where("status = ?", models.SysRoleStatusNormal).
		First(&role).Error
	return role, err
}

func (my MysqlService) FindRoleByIds(ids []uint) []models.SysRole {
	roles := make([]models.SysRole, 0)
	my.Q.Tx.
		Model(&models.SysRole{}).
		Where("id IN (?)", ids).
		Where("status = ?", models.SysRoleStatusNormal).
		Find(&roles)
	return roles
}
