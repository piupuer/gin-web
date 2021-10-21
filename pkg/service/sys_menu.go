package service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
)

// 获取权限菜单树
func (my MysqlService) GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	menus := make([]models.SysMenu, 0)
	// 查询全部菜单
	allMenu := make([]models.SysMenu, 0)
	err := my.Q.Tx.
		Model(&models.SysMenu{}).
		Find(&allMenu).Error
	if err != nil {
		return menus, err
	}
	// 查询当前权限
	var role models.SysRole
	err = my.Q.Tx.
		Model(&models.SysRole{}).
		Preload("Menus").
		Where("id = ?", roleId).
		First(&role).Error
	if err != nil {
		return menus, err
	}
	// 子级菜单需要将全部父级菜单加上
	_, newMenus := addParentMenu(role.Menus, allMenu)

	// 生成菜单树
	tree = my.GenMenuTree(0, newMenus)
	return tree, nil
}

// 获取所有菜单
func (my MysqlService) GetMenus(currentRole models.SysRole) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	menus := my.getAllMenu(currentRole)
	// 生成菜单树
	tree = my.GenMenuTree(0, menus)
	return tree
}

// 生成菜单树
// parentId: 父菜单编号
// roleMenus: 有权限的菜单列表
func (my MysqlService) GenMenuTree(parentId uint, roleMenus []models.SysMenu) []models.SysMenu {
	roleMenuIds := make([]uint, 0)
	// 查询全部菜单
	allMenu := make([]models.SysMenu, 0)
	err := my.Q.Tx.
		Model(&models.SysMenu{}).
		Find(&allMenu).Error
	if err != nil {
		return roleMenus
	}
	// 加上父级菜单
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
			// 无权限
			continue
		}
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = genMenuTree(menu.Id, roleMenuIds, allMenu)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	return tree
}

// 根据权限编号获取全部菜单
func (my MysqlService) GetAllMenuByRoleId(currentRole models.SysRole, roleId uint) ([]models.SysMenu, []uint, error) {
	// 菜单树
	tree := make([]models.SysMenu, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	// 查询全部菜单
	allMenu := my.getAllMenu(currentRole)
	// 查询角色拥有菜单
	roleMenus := my.getRoleMenus(roleId)
	// 生成菜单树
	tree = my.GenMenuTree(0, allMenu)
	// 获取id列表
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	// 只保留选中项目
	accessIds = models.FindCheckedMenuId(accessIds, allMenu)
	return tree, accessIds, nil
}

// 创建菜单
func (my MysqlService) CreateMenu(currentRole models.SysRole, req *request.CreateMenuReq) (err error) {
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = my.Q.Tx.Create(&menu).Error
	// 自己创建的菜单需绑定权限
	menuReq := request.UpdateIncrementalIdsRequestStruct{
		Create: []uint{menu.Id},
	}
	err = my.UpdateRoleMenusById(currentRole, currentRole.Id, menuReq)
	return
}

// 获取权限菜单, 非菜单树
func (my MysqlService) getRoleMenus(roleId uint) []models.SysMenu {
	var role models.SysRole
	// 根据权限编号获取菜单
	err := my.Q.Tx.
		Preload("Menus").
		Where("id = ?", roleId).
		First(&role).Error
	if err != nil {
		global.Log.Warn(my.Q.Ctx, "[getRoleMenu]", err)
	}
	return role.Menus
}

// 获取全部菜单, 非菜单树
func (my MysqlService) getAllMenu(currentRole models.SysRole) []models.SysMenu {
	menus := make([]models.SysMenu, 0)
	// 查询关系表
	relations := make([]models.SysRoleMenuRelation, 0)
	menuIds := make([]uint, 0)
	query := my.Q.Tx.Model(&models.SysRoleMenuRelation{})
	var err error
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		query = query.Where("sys_role_id = ?", currentRole.Id)
		err = query.Find(&relations).Error
		if err != nil {
			return menus
		}
		for _, relation := range relations {
			menuIds = append(menuIds, relation.SysMenuId)
		}
		// 查询所有菜单
		err = my.Q.Tx.
			Where("id IN (?)", menuIds).
			Order("sort").
			Find(&menus).Error
	} else {
		err = my.Q.Tx.
			Order("sort").
			Find(&menus).Error
	}
	if err != nil {
		global.Log.Warn(my.Q.Ctx, "[getAllMenu]", err)
	}
	return menus
}

// 给定菜单, 如果父菜单不存在则加入列表
func addParentMenu(menus, all []models.SysMenu) ([]uint, []models.SysMenu) {
	parentIds := make([]uint, 0)
	menuIds := make([]uint, 0)
	for _, menu := range menus {
		if menu.ParentId > 0 {
			parentIds = append(parentIds, menu.ParentId)
			// 向上获取父级菜单
			parentMenuIds := getParentMenuIds(menu.ParentId, all)
			if len(parentMenuIds) > 0 {
				parentIds = append(parentIds, parentMenuIds...)
			}
		}
		menuIds = append(menuIds, menu.Id)
	}
	// 合并父级菜单
	if len(parentIds) > 0 {
		menuIds = append(menuIds, parentIds...)
	}
	newMenuIds := make([]uint, 0)
	newMenus := make([]models.SysMenu, 0)
	for _, menu := range all {
		for _, id := range menuIds {
			// 保证id一致且不重复
			if id == menu.Id && !utils.ContainsUint(newMenuIds, id) {
				newMenus = append(newMenus, menu)
				newMenuIds = append(newMenuIds, id)
			}
		}
	}
	return newMenuIds, newMenus
}

// 获取全部父级菜单id
func getParentMenuIds(menuId uint, all []models.SysMenu) []uint {
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
	// 继续向上寻找
	newParentIds := getParentMenuIds(currentMenu.ParentId, all)
	if len(newParentIds) > 0 {
		parentIds = append(parentIds, newParentIds...)
	}
	return parentIds
}

// 是否包含子菜单
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
