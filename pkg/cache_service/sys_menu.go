package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"github.com/piupuer/go-helper/ms"
)

func (rd RedisService) FindMenu(currentRole models.SysRole) []ms.SysMenu {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindMenu(currentRole)
	}
	tree := make([]ms.SysMenu, 0)
	menus := rd.findMenuByCurrentRole(currentRole)
	tree = rd.mysql.GenMenuTree(0, menus)
	return tree
}

func (rd RedisService) FindMenuByRoleId(currentRole models.SysRole, roleId uint) ([]ms.SysMenu, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindMenuByRoleId(currentRole, roleId)
	}
	tree := make([]ms.SysMenu, 0)
	accessIds := make([]uint, 0)
	allMenu := rd.findMenuByCurrentRole(currentRole)
	roleMenus := rd.findMenuByRoleId(roleId)
	tree = rd.mysql.GenMenuTree(0, allMenu)
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	accessIds = service.FindCheckedMenuId(accessIds, allMenu)
	return tree, accessIds, nil
}

// find all menus by role id(not menu tree)
func (rd RedisService) findMenuByRoleId(roleId uint) []ms.SysMenu {
	// query current role menu relation
	relations := make([]ms.SysMenuRoleRelation, 0)
	menuIds := make([]uint, 0)
	rd.Q.
		Table("sys_menu_role_relation").
		Where("role_id", "=", roleId).
		Find(&relations)
	for _, relation := range relations {
		menuIds = append(menuIds, relation.MenuId)
	}
	roleMenu := make([]ms.SysMenu, 0)
	if len(menuIds) > 0 {
		rd.Q.
			Table("sys_menu").
			Where("id", "contains", menuIds).
			Order("sort").
			Find(&roleMenu)
	}
	return roleMenu
}

// find all menus by current role(not menu tree)
func (rd RedisService) findMenuByCurrentRole(currentRole models.SysRole) []ms.SysMenu {
	menus := make([]ms.SysMenu, 0)
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		// find menus by current role id
		menus = rd.findMenuByRoleId(currentRole.Id)
	} else {
		// super admin has all menus
		rd.Q.
			Table("sys_menu").
			Order("sort").
			Find(&menus)
	}
	return menus
}
