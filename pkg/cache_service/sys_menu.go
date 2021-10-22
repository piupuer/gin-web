package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
)

func (rd RedisService) FindMenu(currentRole models.SysRole) []models.SysMenu {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindMenu(currentRole)
	}
	tree := make([]models.SysMenu, 0)
	menus := rd.findMenuByCurrentRole(currentRole)
	tree = rd.mysql.GenMenuTree(0, menus)
	return tree
}

func (rd RedisService) FindMenuByRoleId(currentRole models.SysRole, roleId uint) ([]models.SysMenu, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindMenuByRoleId(currentRole, roleId)
	}
	tree := make([]models.SysMenu, 0)
	accessIds := make([]uint, 0)
	allMenu := rd.findMenuByCurrentRole(currentRole)
	roleMenus := rd.findMenuByRoleId(roleId)
	tree = rd.mysql.GenMenuTree(0, allMenu)
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	accessIds = models.FindCheckedMenuId(accessIds, allMenu)
	return tree, accessIds, nil
}

// find all menus by role id(not menu tree)
func (rd RedisService) findMenuByRoleId(roleId uint) []models.SysMenu {
	var role models.SysRole
	err := rd.Q.
		Table("sys_role").
		Preload("Menus").
		Where("id", "=", roleId).
		First(&role).Error
	if err != nil {
		global.Log.Warn(rd.Q.Ctx, "%v", err)
	}
	return role.Menus
}

// find all menus by current role(not menu tree)
func (rd RedisService) findMenuByCurrentRole(currentRole models.SysRole) []models.SysMenu {
	menus := make([]models.SysMenu, 0)

	relations := make([]models.SysRoleMenuRelation, 0)
	menuIds := make([]uint, 0)
	query := rd.Q.Table("sys_role_menu_relation")
	var err error
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		// find menus by current role id
		query = query.Where("sys_role_id", "=", currentRole.Id)
		err = query.Find(&relations).Error
		if err != nil {
			return menus
		}
		for _, relation := range relations {
			menuIds = append(menuIds, relation.SysMenuId)
		}
		err = rd.Q.
			Table("sys_menu").
			Where("id", "in", menuIds).
			Order("sort").
			Find(&menus).Error
	} else {
		// super admin has all menus
		err = rd.Q.
			Table("sys_menu").
			Order("sort").
			Find(&menus).Error
	}

	if err != nil {
		global.Log.Warn(rd.Q.Ctx, "%v", err)
	}
	return menus
}
