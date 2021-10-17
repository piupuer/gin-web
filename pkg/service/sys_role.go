package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gorm.io/gorm"
	"strings"
)

// 根据当前角色顺序获取角色编号集合(主要功能是针对不同角色用户登录系统隐藏特定菜单)
func (s MysqlService) GetRoleIdsBySort(currentRoleSort uint) ([]uint, error) {
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	err := s.tx.Model(&models.SysRole{}).Where("sort >= ?", currentRoleSort).Find(&roles).Error
	if err != nil {
		return roleIds, err
	}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds, nil
}

// 获取所有角色
func (s MysqlService) GetRoles(req *request.RoleReq) ([]models.SysRole, error) {
	var err error
	list := make([]models.SysRole, 0)
	query := s.tx.
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
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status = ?", 1)
		} else {
			query = query.Where("status = ?", 0)
		}
	}
	// 查询列表
	err = s.Find(query, &req.Page, &list)
	return list, err
}

// 更新角色的权限菜单
func (s MysqlService) UpdateRoleMenusById(currentRole models.SysRole, id uint, req request.UpdateIncrementalIdsRequestStruct) (err error) {
	// 查询全部菜单
	allMenu := s.getAllMenu(currentRole)
	// 查询角色拥有菜单
	roleMenus := s.getRoleMenus(id)
	// 获取当前菜单编号集合
	menuIds := make([]uint, 0)
	for _, menu := range roleMenus {
		menuIds = append(menuIds, menu.Id)
	}
	// 获取菜单增量
	incremental := req.GetIncremental(menuIds, allMenu)
	// 查询所有菜单
	incrementalMenus := make([]models.SysMenu, 0)
	err = s.tx.
		Model(&models.SysMenu{}).
		Where("id in (?)", incremental).
		Find(&incrementalMenus).Error
	if err != nil {
		return
	}
	// 去除菜单增量中包含子菜单的父菜单
	newIncrementalMenus := make([]models.SysMenu, 0)
	for _, menu := range incrementalMenus {
		if !hasChildrenMenu(menu.Id, allMenu) {
			newIncrementalMenus = append(newIncrementalMenus, menu)
		}
	}
	// 查询role
	var role models.SysRole
	err = s.tx.Where("id = ?", id).First(&role).Error
	if err != nil {
		return
	}
	// 替换菜单
	err = s.tx.Model(&role).Association("Menus").Replace(&incrementalMenus)
	return
}

// 更新角色的权限接口
func (s MysqlService) UpdateRoleApisById(id uint, req request.UpdateIncrementalIdsRequestStruct) (err error) {
	var oldRole models.SysRole
	query := s.tx.Model(&oldRole).Where("id = ?", id).First(&oldRole)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	if len(req.Delete) > 0 {
		// 查询需要删除的api
		deleteApis := make([]models.SysApi, 0)
		err = s.tx.Where("id IN (?)", req.Delete).Find(&deleteApis).Error
		if err != nil {
			return
		}
		// 构建casbin规则
		cs := make([]models.SysRoleCasbin, 0)
		for _, api := range deleteApis {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: oldRole.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		// 批量删除
		_, err = s.BatchDeleteRoleCasbins(cs)
	}
	if len(req.Create) > 0 {
		// 查询需要新增的api
		createApis := make([]models.SysApi, 0)
		err = s.tx.Where("id IN (?)", req.Create).Find(&createApis).Error
		if err != nil {
			return
		}
		// 构建casbin规则
		cs := make([]models.SysRoleCasbin, 0)
		for _, api := range createApis {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: oldRole.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		// 批量创建
		_, err = s.BatchCreateRoleCasbins(cs)

	}
	return
}

// 批量删除角色
func (s MysqlService) DeleteRoleByIds(ids []uint) (err error) {
	var roles []models.SysRole
	// 查询符合条件的角色, 以及关联的用户
	err = s.tx.Preload("Users").Where("id IN (?)", ids).Find(&roles).Error
	if err != nil {
		return
	}
	newIds := make([]uint, 0)
	oldCasbins := make([]models.SysRoleCasbin, 0)
	for _, v := range roles {
		if len(v.Users) > 0 {
			return errors.New(fmt.Sprintf("角色[%s]仍有%d位关联用户, 请先删除用户再删除角色", v.Name, len(v.Users)))
		}
		oldCasbins = append(oldCasbins, s.GetRoleCasbins(models.SysRoleCasbin{
			Keyword: v.Keyword,
		})...)
		newIds = append(newIds, v.Id)
	}
	if len(oldCasbins) > 0 {
		// 删除关联的casbin
		s.BatchDeleteRoleCasbins(oldCasbins)
	}
	if len(newIds) > 0 {
		// 执行删除
		err = s.DeleteByIds(newIds, new(models.SysRole))
	}
	return
}
