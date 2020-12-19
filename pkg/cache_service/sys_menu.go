package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
)

// 获取权限菜单树
func (s *RedisService) GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetMenuTree(roleId)
	}
	tree := make([]models.SysMenu, 0)
	var role models.SysRole
	err := s.redis.Table(new(models.SysRole).TableName()).Preload("Menus").Where("id", "=", roleId).First(&role).Error
	menus := make([]models.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	// 生成菜单树
	tree = service.GenMenuTree(nil, role.Menus)
	return tree, nil
}

// 获取所有菜单
func (s *RedisService) GetMenus(currentRole models.SysRole) []models.SysMenu {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetMenus(currentRole)
	}
	tree := make([]models.SysMenu, 0)
	// 获取全部菜单
	menus := s.getAllMenu(currentRole)
	// 生成菜单树
	tree = service.GenMenuTree(nil, menus)
	return tree
}

// 根据权限编号获取全部菜单
func (s *RedisService) GetAllMenuByRoleId(currentRole models.SysRole, roleId uint) ([]models.SysMenu, []uint, error) {
	if !global.Conf.System.UseRedis {
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
	tree = service.GenMenuTree(nil, allMenu)
	// 获取id列表
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	// 只保留选中项目
	accessIds = models.GetCheckedMenuIds(accessIds, allMenu)
	return tree, accessIds, nil
}

// 获取权限菜单, 非菜单树
func (s *RedisService) getRoleMenus(roleId uint) []models.SysMenu {
	var role models.SysRole
	// 根据权限编号获取菜单
	err := s.redis.Table(new(models.SysRole).TableName()).Preload("Menus").Where("id", "=", roleId).First(&role).Error
	global.Log.Warn("[getRoleMenu]", err)
	return role.Menus
}

// 获取全部菜单, 非菜单树
func (s *RedisService) getAllMenu(currentRole models.SysRole) []models.SysMenu {
	menus := make([]models.SysMenu, 0)

	// 查询关系表
	relations := make([]models.RelationMenuRole, 0)
	menuIds := make([]uint, 0)
	query := s.redis.Table(new(models.RelationMenuRole).TableName())
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
		err = s.redis.Table(new(models.SysMenu).TableName()).Order("sort").Where("id", "in", menuIds).Find(&menus).Error
	} else {
		err = s.redis.Table(new(models.SysMenu).TableName()).Order("sort").Find(&menus).Error
	}

	global.Log.Warn("[getAllMenu]", err)
	return menus
}
