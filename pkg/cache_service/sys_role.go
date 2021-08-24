package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

// 根据当前角色顺序获取角色编号集合(主要功能是针对不同角色用户登录系统隐藏特定菜单)
func (s RedisService) GetRoleIdsBySort(currentRoleSort uint) ([]uint, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetRoleIdsBySort(currentRoleSort)
	}
	roles := make([]models.SysRole, 0)
	roleIds := make([]uint, 0)
	err := s.redis.Table(new(models.SysRole).TableName()).Where("sort", ">=", currentRoleSort).Find(&roles).Error
	if err != nil {
		return roleIds, err
	}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds, nil
}

// 获取所有角色
func (s RedisService) GetRoles(req *request.RoleRequestStruct) ([]models.SysRole, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetRoles(req)
	}
	var err error
	list := make([]models.SysRole, 0)
	query := s.redis.
		Table(new(models.SysRole).TableName()).
		Order("created_at DESC").
		Where("sort", ">=", req.CurrentRoleSort)
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name", "contains", name)
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		query = query.Where("keyword", "contains", keyword)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status", "=", 1)
		} else {
			query = query.Where("status", "=", 0)
		}
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}
