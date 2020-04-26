package service

import (
	"errors"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
)

// 获取菜单树
func GetMenuTree(roleId uint) (tree []models.SysMenu, err error) {
	query := global.Mysql.First(&models.SysRole{
		Model: models.Model{
			Id: roleId,
		},
	})
	if query.RecordNotFound() {
		return nil, errors.New("菜单为空")
	}
	menus := make([]models.SysMenu, 0)
	// 查询当前role关联的所有菜单
	query.Association("menus").Find(&menus)
	// 生成菜单树
	tree = genMenuTree(nil, menus)
	return
}

// 生成菜单树
func genMenuTree(parent *models.SysMenu, menus []models.SysMenu) (tree []models.SysMenu) {
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}

	tree = make([]models.SysMenu, 0)
	for _, menu := range menus {
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = genMenuTree(&menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	return
}
