package service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取权限菜单树
func (s *MysqlService) GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	query := s.tx.First(&models.SysRole{
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
	tree = GenMenuTree(nil, menus)
	return tree, nil
}

// 获取所有菜单
func (s *MysqlService) GetMenus() ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	menus := make([]models.SysMenu, 0)
	// 查询所有菜单
	err := s.tx.Order("sort").Find(&menus).Error
	// 生成菜单树
	tree = GenMenuTree(nil, menus)
	return tree, err
}

// 生成菜单树
func GenMenuTree(parent *models.SysMenu, menus []models.SysMenu) []models.SysMenu {
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
			menu.Children = GenMenuTree(&menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	return tree
}

// 根据权限编号获取全部菜单
func (s *MysqlService) GetAllMenuByRoleId(roleId uint) ([]models.SysMenu, []uint, error) {
	// 菜单树
	tree := make([]models.SysMenu, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allMenu := make([]models.SysMenu, 0)
	// 查询全部菜单
	err := s.tx.Order("sort").Find(&allMenu).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 查询角色拥有的全部菜单
	var role models.SysRole
	err = s.tx.Preload("Menus").Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 生成菜单树
	tree = GenMenuTree(nil, allMenu)
	// 获取id列表
	for _, menu := range role.Menus {
		accessIds = append(accessIds, menu.Id)
	}
	return tree, accessIds, nil
}

// 创建菜单
func (s *MysqlService) CreateMenu(req *request.CreateMenuRequestStruct) (err error) {
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = s.tx.Create(&menu).Error
	return
}

// 更新菜单
func (s *MysqlService) UpdateMenuById(id uint, req gin.H) (err error) {
	var oldMenu models.SysMenu
	query := s.tx.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
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
func (s *MysqlService) DeleteMenuByIds(ids []uint) (err error) {
	// 执行删除
	return s.tx.Where("id IN (?)", ids).Delete(models.SysMenu{}).Error
}
