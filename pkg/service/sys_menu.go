package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/utils"
)

// 获取权限菜单树
func GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	query := global.Mysql.First(&models.SysRole{
		Model: models.Model{
			Id: roleId,
		},
	})
	if query.RecordNotFound() {
		return tree, errors.New("菜单为空")
	}
	menus := make([]models.SysMenu, 0)
	// 查询当前role关联的所有菜单
	query.Where("status = ?", 1).Order("sort").Association("menus").Find(&menus)
	// 生成菜单树
	tree = genMenuTree(nil, menus)
	return tree, nil
}

// 获取所有菜单
func GetMenus() ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	menus := make([]models.SysMenu, 0)
	// 查询所有菜单
	err := global.Mysql.Order("sort").Find(&menus).Error
	// 生成菜单树
	tree = genMenuTree(nil, menus)
	return tree, err
}

// 生成菜单树
func genMenuTree(parent *models.SysMenu, menus []models.SysMenu) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}

	for _, menu := range menus {
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = genMenuTree(&menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	return tree
}

// 根据权限编号获取全部菜单
func GetAllMenuByRoleId(roleId uint) ([]models.SysMenu, []uint, error) {
	// 菜单树
	tree := make([]models.SysMenu, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allMenu := make([]models.SysMenu, 0)
	// 查询全部菜单
	err := global.Mysql.Order("sort").Find(&allMenu).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 查询角色拥有的全部菜单
	var role models.SysRole
	err = global.Mysql.Preload("Menus").Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 生成菜单树
	tree = genMenuTree(nil, allMenu)
	// 获取id列表
	for _, menu := range role.Menus {
		accessIds = append(accessIds, menu.Id)
	}
	return tree, accessIds, nil
}

// 创建菜单
func CreateMenu(req *request.CreateMenuRequestStruct) (err error) {
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = global.Mysql.Create(&menu).Error
	return
}

// 更新菜单
func UpdateMenuById(id uint, req gin.H) (err error) {
	var oldMenu models.SysMenu
	query := global.Mysql.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldMenu, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除菜单
func DeleteMenuByIds(ids []uint) (err error) {
	// 执行删除
	return global.Mysql.Where("id IN (?)", ids).Delete(models.SysMenu{}).Error
}
