package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/utils"
	"strings"
)

// 获取所有角色
func (s *CommonService) GetRoles(req *request.RoleListRequestStruct) ([]models.SysRole, error) {
	var err error
	list := make([]models.SysRole, 0)
	db := global.Mysql
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = db.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 创建角色
func (s *CommonService) CreateRole(req *request.CreateRoleRequestStruct) (err error) {
	var role models.SysRole
	utils.Struct2StructByJson(req, &role)
	// 创建数据
	err = s.tx.Create(&role).Error
	return
}

// 更新角色
func (s *CommonService) UpdateRoleById(id uint, req gin.H) (err error) {
	var oldRole models.SysRole
	query := s.tx.Table(oldRole.TableName()).Where("id = ?", id).First(&oldRole)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldRole, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 更新角色的权限菜单
func (s *CommonService) UpdateRoleMenusById(id uint, req request.UpdateIncrementalIdsRequestStruct) (err error) {
	var oldRole models.SysRole
	query := s.tx.Model(&oldRole).Preload("Menus").Where("id = ?", id).First(&oldRole)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}
	// 获取当前菜单编号集合
	menuIds := make([]uint, 0)
	for _, menu := range oldRole.Menus {
		menuIds = append(menuIds, menu.Id)
	}
	// 获取菜单增量
	incremental := req.GetIncremental(menuIds)
	// 查询所有菜单
	var incrementalMenus []models.SysMenu
	err = s.tx.Where("id in (?)", incremental).Find(&incrementalMenus).Error
	if err != nil {
		return
	}
	// 替换菜单
	err = query.Association("Menus").Replace(&incrementalMenus).Error
	return
}

// 更新角色的权限接口
func (s *CommonService) UpdateRoleApisById(id uint, req request.UpdateIncrementalIdsRequestStruct) (err error) {
	var oldRole models.SysRole
	query := s.tx.Model(&oldRole).Where("id = ?", id).First(&oldRole)
	if query.RecordNotFound() {
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
func (s *CommonService) DeleteRoleByIds(ids []uint) (err error) {
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
		err = s.tx.Where("id IN (?)", newIds).Delete(models.SysRole{}).Error
	}
	return
}
