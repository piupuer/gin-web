package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
)

// get menu tree by role id
func (my MysqlService) GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	menus := make([]models.SysMenu, 0)
	// query all menus
	allMenu := make([]models.SysMenu, 0)
	err := my.Q.Tx.
		Model(&models.SysMenu{}).
		Find(&allMenu).Error
	if err != nil {
		return menus, err
	}
	// query current role's menus
	var role models.SysRole
	err = my.Q.Tx.
		Model(&models.SysRole{}).
		Preload("Menus").
		Where("id = ?", roleId).
		First(&role).Error
	if err != nil {
		return menus, err
	}
	_, newMenus := addParentMenu(role.Menus, allMenu)

	tree = my.GenMenuTree(0, newMenus)
	return tree, nil
}

func (my MysqlService) FindMenu(currentRole models.SysRole) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	menus := my.findMenuByCurrentRole(currentRole)
	tree = my.GenMenuTree(0, menus)
	return tree
}

// generate menu tree
func (my MysqlService) GenMenuTree(parentId uint, roleMenus []models.SysMenu) []models.SysMenu {
	roleMenuIds := make([]uint, 0)
	allMenu := make([]models.SysMenu, 0)
	err := my.Q.Tx.
		Model(&models.SysMenu{}).
		Find(&allMenu).Error
	if err != nil {
		return roleMenus
	}
	// add parent menu
	_, newRoleMenus := addParentMenu(roleMenus, allMenu)
	for _, menu := range newRoleMenus {
		if !utils.ContainsUint(roleMenuIds, menu.Id) {
			roleMenuIds = append(roleMenuIds, menu.Id)
		}
	}
	return genMenuTree(parentId, roleMenuIds, allMenu)
}

func genMenuTree(parentId uint, roleMenuIds []uint, allMenu []models.SysMenu) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
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

func (my MysqlService) FindMenuByRoleId(currentRole models.SysRole, roleId uint) ([]models.SysMenu, []uint, error) {
	tree := make([]models.SysMenu, 0)
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

func FindCheckedMenuId(list []uint, allMenu []models.SysMenu) []uint {
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
func FindChildrenId(parentId uint, allMenu []models.SysMenu) []uint {
	childrenIds := make([]uint, 0)
	for _, menu := range allMenu {
		if menu.ParentId == parentId {
			childrenIds = append(childrenIds, menu.Id)
		}
	}
	return childrenIds
}

func FindIncremental(req request.UpdateMenuIncrementalIdsReq, oldMenuIds []uint, allMenu []models.SysMenu) []uint {
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
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	err = my.Q.Tx.Create(&menu).Error
	menuReq := request.UpdateMenuIncrementalIdsReq{
		Create: []uint{menu.Id},
	}
	err = my.UpdateRoleMenuById(currentRole, currentRole.Id, menuReq)
	return
}

// find all menus by role id(not menu tree)
func (my MysqlService) findMenuByRoleId(roleId uint) []models.SysMenu {
	var role models.SysRole
	err := my.Q.Tx.
		Preload("Menus").
		Where("id = ?", roleId).
		First(&role).Error
	if err != nil {
		global.Log.Warn(my.Q.Ctx, "%v", err)
	}
	return role.Menus
}

// find all menus by current role(not menu tree)
func (my MysqlService) findMenuByCurrentRole(currentRole models.SysRole) []models.SysMenu {
	menus := make([]models.SysMenu, 0)
	relations := make([]models.SysRoleMenuRelation, 0)
	menuIds := make([]uint, 0)
	query := my.Q.Tx.Model(&models.SysRoleMenuRelation{})
	var err error
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		// find menus by current role id
		query = query.Where("sys_role_id = ?", currentRole.Id)
		err = query.Find(&relations).Error
		if err != nil {
			return menus
		}
		for _, relation := range relations {
			menuIds = append(menuIds, relation.SysMenuId)
		}
		err = my.Q.Tx.
			Where("id IN (?)", menuIds).
			Order("sort").
			Find(&menus).Error
	} else {
		// super admin has all menus
		err = my.Q.Tx.
			Order("sort").
			Find(&menus).Error
	}
	if err != nil {
		global.Log.Warn(my.Q.Ctx, "%v", err)
	}
	return menus
}

func addParentMenu(menus, all []models.SysMenu) ([]uint, []models.SysMenu) {
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
	newMenus := make([]models.SysMenu, 0)
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
func findParentMenuId(menuId uint, all []models.SysMenu) []uint {
	var currentMenu models.SysMenu
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

// check whether has children menu
func hasChildrenMenu(menuId uint, all []models.SysMenu) bool {
	var currentMenu models.SysMenu
	for _, menu := range all {
		if menuId == menu.Id {
			currentMenu = menu
			break
		}
	}
	for _, menu := range all {
		if menu.ParentId == currentMenu.Id {
			return true
		}
	}
	return false
}
