package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
)

// 获取所有菜单
func (s RedisService) GetMenus(currentRole models.SysRole) []models.SysMenu {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		// 不使用redis
		return s.mysql.GetMenus(currentRole)
	}
	tree := make([]models.SysMenu, 0)
	// 获取全部菜单
	menus := s.getAllMenu(currentRole)
	// 生成菜单树
	tree = s.mysql.GenMenuTree(0, menus)
	return tree
}

// 根据权限编号获取全部菜单
func (s RedisService) GetAllMenuByRoleId(currentRole models.SysRole, roleId uint) ([]models.SysMenu, []uint, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		// 不使用redis
		return s.mysql.GetAllMenuByRoleId(currentRole, roleId)
	}
	// 菜单树
	tree := make([]models.SysMenu, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	// 查询全部菜单
	allMenu := s.getAllMenu(currentRole)
	// 查询角色拥有菜单
	roleMenus := s.getRoleMenus(roleId)
	// 生成菜单树
	tree = s.mysql.GenMenuTree(0, allMenu)
	// 获取id列表
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	// 只保留选中项目
	accessIds = models.FindCheckedMenuId(accessIds, allMenu)
	return tree, accessIds, nil
}

// 获取权限菜单, 非菜单树
func (s RedisService) getRoleMenus(roleId uint) []models.SysMenu {
	var role models.SysRole
	// 根据权限编号获取菜单
	err := s.Q.
		Table("sys_role").
		Preload("Menus").
		Where("id", "=", roleId).
		First(&role).Error
	if err != nil {
		global.Log.Warn(s.Q.Ctx, "[getRoleMenu]", err)
	}
	return role.Menus
}

// 获取全部菜单, 非菜单树
func (s RedisService) getAllMenu(currentRole models.SysRole) []models.SysMenu {
	menus := make([]models.SysMenu, 0)

	// 查询关系表
	relations := make([]models.SysRoleMenuRelation, 0)
	menuIds := make([]uint, 0)
	query := s.Q.Table("sys_role_menu_relation")
	var err error
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		query = query.Where("sys_role_id", "=", currentRole.Id)
		err = query.Find(&relations).Error
		if err != nil {
			return menus
		}
		for _, relation := range relations {
			menuIds = append(menuIds, relation.SysMenuId)
		}
		// 查询所有菜单
		err = s.Q.
			Table("sys_menu").
			Where("id", "in", menuIds).
			Order("sort").
			Find(&menus).Error
	} else {
		err = s.Q.
			Table("sys_menu").
			Order("sort").
			Find(&menus).Error
	}

	if err != nil {
		global.Log.Warn(s.Q.Ctx, "[getAllMenu]", err)
	}
	return menus
}
