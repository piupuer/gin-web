package service

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/ms"
)

// get menu tree by role id
func (my MysqlService) GetMenuTree(roleId uint) ([]ms.SysMenu, error) {
	tree := make([]ms.SysMenu, 0)
	// query all menus
	allMenu := make([]ms.SysMenu, 0)
	my.Q.Tx.
		Model(&ms.SysMenu{}).
		Find(&allMenu)
	roleMenu := my.findMenuByRoleId(roleId)
	_, newMenus := addParentMenu(roleMenu, allMenu)

	tree = my.GenMenuTree(0, newMenus)
	return tree, nil
}

func (my MysqlService) FindMenu(currentRole models.SysRole) []ms.SysMenu {
	tree := make([]ms.SysMenu, 0)
	menus := my.findMenuByCurrentRole(currentRole)
	tree = my.GenMenuTree(0, menus)
	return tree
}

// generate menu tree
func (my MysqlService) GenMenuTree(parentId uint, roleMenus []ms.SysMenu) []ms.SysMenu {
	roleMenuIds := make([]uint, 0)
	allMenu := make([]ms.SysMenu, 0)
	my.Q.Tx.
		Model(&ms.SysMenu{}).
		Find(&allMenu)
	// add parent menu
	_, newRoleMenus := addParentMenu(roleMenus, allMenu)
	for _, menu := range newRoleMenus {
		if !utils.ContainsUint(roleMenuIds, menu.Id) {
			roleMenuIds = append(roleMenuIds, menu.Id)
		}
	}
	return genMenuTree(parentId, roleMenuIds, allMenu)
}

func genMenuTree(parentId uint, roleMenuIds []uint, allMenu []ms.SysMenu) []ms.SysMenu {
	tree := make([]ms.SysMenu, 0)
	for _, menu := range allMenu {
		if !utils.ContainsUint(roleMenuIds, menu.Id) {
			continue
		}
		if menu.ParentId == parentId {
			menu.Children = genMenuTree(menu.Id, roleMenuIds, allMenu)
			tree = append(tree, menu)
		}
	}
	return tree
}

func (my MysqlService) FindMenuByRoleId(currentRole models.SysRole, roleId uint) ([]ms.SysMenu, []uint, error) {
	tree := make([]ms.SysMenu, 0)
	accessIds := make([]uint, 0)
	allMenu := my.findMenuByCurrentRole(currentRole)
	roleMenus := my.findMenuByRoleId(roleId)
	tree = my.GenMenuTree(0, allMenu)
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	accessIds = FindCheckedMenuId(accessIds, allMenu)
	return tree, accessIds, nil
}

func FindCheckedMenuId(list []uint, allMenu []ms.SysMenu) []uint {
	checked := make([]uint, 0)
	for _, c := range list {
		children := FindChildrenId(c, allMenu)
		count := 0
		for _, child := range children {
			contains := false
			for _, v := range list {
				if v == child {
					contains = true
				}
			}
			if contains {
				count++
			}
		}
		if len(children) == count {
			// all checked
			checked = append(checked, c)
		}
	}
	return checked
}

// find children menu ids
func FindChildrenId(parentId uint, allMenu []ms.SysMenu) []uint {
	childrenIds := make([]uint, 0)
	for _, menu := range allMenu {
		if menu.ParentId == parentId {
			childrenIds = append(childrenIds, menu.Id)
		}
	}
	return childrenIds
}

func FindIncrementalMenu(req request.UpdateMenuIncrementalIdsReq, oldMenuIds []uint, allMenu []ms.SysMenu) []uint {
	createIds := FindCheckedMenuId(req.Create, allMenu)
	deleteIds := FindCheckedMenuId(req.Delete, allMenu)
	newList := make([]uint, 0)
	for _, oldItem := range oldMenuIds {
		// not in delete
		if !utils.Contains(deleteIds, oldItem) {
			newList = append(newList, oldItem)
		}
	}
	// need create
	return append(newList, createIds...)
}

func (my MysqlService) CreateMenu(currentRole models.SysRole, req *request.CreateMenuReq) (err error) {
	var menu ms.SysMenu
	utils.Struct2StructByJson(req, &menu)
	err = my.Q.Tx.Create(&menu).Error
	menuReq := request.UpdateMenuIncrementalIdsReq{
		Create: []uint{menu.Id},
	}
	err = my.UpdateMenuByRoleId(currentRole, currentRole.Id, menuReq)
	return
}

func (my MysqlService) UpdateMenuByRoleId(currentRole models.SysRole, targetRoleId uint, req request.UpdateMenuIncrementalIdsReq) (err error) {
	allMenu := my.FindMenu(currentRole)
	roleMenus := my.findMenuByRoleId(targetRoleId)
	menuIds := make([]uint, 0)
	for _, menu := range roleMenus {
		menuIds = append(menuIds, menu.Id)
	}
	incremental := FindIncrementalMenu(req, menuIds, allMenu)
	incrementalMenus := make([]ms.SysMenu, 0)
	my.Q.Tx.
		Model(&ms.SysMenu{}).
		Where("id in (?)", incremental).
		Find(&incrementalMenus)
	newRelations := make([]ms.SysMenuRoleRelation, 0)
	for _, menu := range incrementalMenus {
		newRelations = append(newRelations, ms.SysMenuRoleRelation{
			MenuId: menu.Id,
			RoleId: targetRoleId,
		})
	}
	my.Q.Tx.
		Where("role_id = ?", targetRoleId).
		Delete(&ms.SysMenuRoleRelation{})
	my.Q.Tx.
		Model(&ms.SysMenuRoleRelation{}).
		Create(&newRelations)
	return
}

// find all menus by role id(not menu tree)
func (my MysqlService) findMenuByRoleId(roleId uint) []ms.SysMenu {
	// query current role menu relation
	menuIds := make([]uint, 0)
	my.Q.Tx.
		Model(&ms.SysMenuRoleRelation{}).
		Where("role_id = ?", roleId).
		Pluck("menu_id", &menuIds)
	roleMenu := make([]ms.SysMenu, 0)
	if len(menuIds) > 0 {
		my.Q.Tx.
			Model(&ms.SysMenu{}).
			Where("id IN (?)", menuIds).
			Order("sort").
			Find(&roleMenu)
	}
	return roleMenu
}

// find all menus by current role(not menu tree)
func (my MysqlService) findMenuByCurrentRole(currentRole models.SysRole) []ms.SysMenu {
	menus := make([]ms.SysMenu, 0)
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		// find menus by current role id
		menus = my.findMenuByRoleId(currentRole.Id)
	} else {
		// super admin has all menus
		my.Q.Tx.
			Order("sort").
			Find(&menus)
	}
	return menus
}

func addParentMenu(menus, all []ms.SysMenu) ([]uint, []ms.SysMenu) {
	parentIds := make([]uint, 0)
	menuIds := make([]uint, 0)
	for _, menu := range menus {
		if menu.ParentId > 0 {
			parentIds = append(parentIds, menu.ParentId)
			// find parent menu
			parentMenuIds := findParentMenuId(menu.ParentId, all)
			if len(parentMenuIds) > 0 {
				parentIds = append(parentIds, parentMenuIds...)
			}
		}
		menuIds = append(menuIds, menu.Id)
	}
	// merge parent menu
	if len(parentIds) > 0 {
		menuIds = append(menuIds, parentIds...)
	}
	newMenuIds := make([]uint, 0)
	newMenus := make([]ms.SysMenu, 0)
	for _, menu := range all {
		for _, id := range menuIds {
			if id == menu.Id && !utils.ContainsUint(newMenuIds, id) {
				newMenus = append(newMenus, menu)
				newMenuIds = append(newMenuIds, id)
			}
		}
	}
	return newMenuIds, newMenus
}

// find parent menu ids
func findParentMenuId(menuId uint, all []ms.SysMenu) []uint {
	var currentMenu ms.SysMenu
	parentIds := make([]uint, 0)
	for _, menu := range all {
		if menuId == menu.Id {
			currentMenu = menu
			break
		}
	}
	if currentMenu.ParentId == 0 {
		return parentIds
	}
	parentIds = append(parentIds, currentMenu.ParentId)
	newParentIds := findParentMenuId(currentMenu.ParentId, all)
	if len(newParentIds) > 0 {
		parentIds = append(parentIds, newParentIds...)
	}
	return parentIds
}
