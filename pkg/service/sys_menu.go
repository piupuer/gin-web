package service

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"gorm.io/gorm"
)

// 获取权限菜单树
func (s *MysqlService) GetMenuTree(roleId uint) ([]models.SysMenu, error) {
	tree := make([]models.SysMenu, 0)
	var role models.SysRole
	err := s.tx.Table(new(models.SysRole).TableName()).Preload("Menus").Where("id = ?", roleId).Find(&role).Error
	menus := make([]models.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	// 生成菜单树
	tree = GenMenuTree(nil, role.Menus)
	return tree, nil
}

// 获取所有菜单
func (s *MysqlService) GetMenus() []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	menus := s.getAllMenu()
	// 生成菜单树
	tree = GenMenuTree(nil, menus)
	return tree
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
	// 查询全部菜单
	allMenu := s.getAllMenu()
	// 查询角色拥有菜单
	roleMenus := s.getRoleMenus(roleId)
	// 生成菜单树
	tree = GenMenuTree(nil, allMenu)
	// 获取id列表
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	// 只保留选中项目
	accessIds = models.GetCheckedMenuIds(accessIds, allMenu)
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
func (s *MysqlService) UpdateMenuById(id uint, req map[string]interface{}) (err error) {
	var oldMenu models.SysMenu
	query := s.tx.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(map[string]interface{}, 0)
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

// 获取权限菜单, 非菜单树
func (s *MysqlService) getRoleMenus(roleId uint) []models.SysMenu {
	var role models.SysRole
	// 根据权限编号获取菜单
	err := s.tx.Preload("Menus").Where("id = ?", roleId).First(&role).Error
	global.Log.Warn("[getRoleMenu]", err)
	return role.Menus
}

// 获取全部菜单, 非菜单树
func (s *MysqlService) getAllMenu() []models.SysMenu {
	menus := make([]models.SysMenu, 0)
	// 查询所有菜单
	err := s.tx.Order("sort").Find(&menus).Error
	global.Log.Warn("[getAllMenu]", err)
	return menus
}
