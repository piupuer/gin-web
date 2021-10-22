package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gorm.io/gorm"
	"strings"
)

func (my MysqlService) FindRoleIdBySort(currentRoleSort uint) ([]uint, error) {
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	err := my.Q.Tx.Model(&models.SysRole{}).Where("sort >= ?", currentRoleSort).Find(&roles).Error
	if err != nil {
		return roleIds, err
	}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds, nil
}

func (my MysqlService) FindRole(req *request.RoleReq) ([]models.SysRole, error) {
	var err error
	list := make([]models.SysRole, 0)
	query := my.Q.Tx.
		Model(&models.SysRole{}).
		Order("created_at DESC").
		Where("sort >= ?", req.CurrentRoleSort)
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		query = query.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status = ?", 1)
		} else {
			query = query.Where("status = ?", 0)
		}
	}
	err = my.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}

func (my MysqlService) UpdateRoleMenuById(currentRole models.SysRole, id uint, req request.UpdateMenuIncrementalIdsReq) (err error) {
	allMenu := my.FindMenu(currentRole)
	roleMenus := my.findMenuByRoleId(id)
	menuIds := make([]uint, 0)
	for _, menu := range roleMenus {
		menuIds = append(menuIds, menu.Id)
	}
	incremental := FindIncremental(req, menuIds, allMenu)
	incrementalMenus := make([]models.SysMenu, 0)
	err = my.Q.Tx.
		Model(&models.SysMenu{}).
		Where("id in (?)", incremental).
		Find(&incrementalMenus).Error
	if err != nil {
		return
	}
	newIncrementalMenus := make([]models.SysMenu, 0)
	for _, menu := range incrementalMenus {
		if !hasChildrenMenu(menu.Id, allMenu) {
			newIncrementalMenus = append(newIncrementalMenus, menu)
		}
	}
	var role models.SysRole
	err = my.Q.Tx.Where("id = ?", id).First(&role).Error
	if err != nil {
		return
	}
	err = my.Q.Tx.Model(&role).Association("Menus").Replace(&incrementalMenus)
	return
}

func (my MysqlService) UpdateRoleApiById(id uint, req request.UpdateMenuIncrementalIdsReq) (err error) {
	var oldRole models.SysRole
	query := my.Q.Tx.Model(&oldRole).Where("id = ?", id).First(&oldRole)
	if query.Error == gorm.ErrRecordNotFound {
		return gorm.ErrRecordNotFound
	}
	if len(req.Delete) > 0 {
		deleteApis := make([]models.SysApi, 0)
		err = my.Q.Tx.Where("id IN (?)", req.Delete).Find(&deleteApis).Error
		if err != nil {
			return
		}
		cs := make([]models.SysRoleCasbin, 0)
		for _, api := range deleteApis {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: oldRole.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		_, err = my.BatchDeleteRoleCasbin(cs)
	}
	if len(req.Create) > 0 {
		createApis := make([]models.SysApi, 0)
		err = my.Q.Tx.Where("id IN (?)", req.Create).Find(&createApis).Error
		if err != nil {
			return
		}
		cs := make([]models.SysRoleCasbin, 0)
		for _, api := range createApis {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: oldRole.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		_, err = my.BatchCreateRoleCasbin(cs)

	}
	return
}

func (my MysqlService) DeleteRoleByIds(ids []uint) (err error) {
	var roles []models.SysRole
	err = my.Q.Tx.Preload("Users").Where("id IN (?)", ids).Find(&roles).Error
	if err != nil {
		return
	}
	newIds := make([]uint, 0)
	oldCasbins := make([]models.SysRoleCasbin, 0)
	for _, v := range roles {
		if len(v.Users) > 0 {
			return fmt.Errorf("role %s has %d associated users, please delete the user before deleting the role", v.Name, len(v.Users))
		}
		oldCasbins = append(oldCasbins, my.FindRoleCasbin(models.SysRoleCasbin{
			Keyword: v.Keyword,
		})...)
		newIds = append(newIds, v.Id)
	}
	if len(oldCasbins) > 0 {
		my.BatchDeleteRoleCasbin(oldCasbins)
	}
	if len(newIds) > 0 {
		err = my.Q.DeleteByIds(newIds, new(models.SysRole))
	}
	return
}
