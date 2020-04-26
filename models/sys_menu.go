package models

// 系统菜单表
type SysMenu struct {
	Model
	Name       string    `gorm:"comment:'菜单名称'" json:"name"`
	Title      string    `gorm:"comment:'菜单标题'" json:"title"`
	Icon       string    `gorm:"comment:'菜单图标'" json:"icon"`
	Path       string    `gorm:"unique;comment:'菜单访问路径'" json:"path"`
	Component  string    `gorm:"comment:'前端组件路径'" json:"component"`
	Permission string    `gorm:"comment:'权限标识'" json:"permission"`
	Sort       int       `gorm:"type:int(3);comment:'菜单顺序(同级菜单, 从0开始, 越小显示越靠前)'" json:"sort"`
	Status     bool      `gorm:"type:tinyint;comment:'菜单状态(正常/禁用)'" json:"status"`
	Visible    bool      `gorm:"type:tinyint;comment:'菜单可见性(可见/隐藏)'" json:"visible"`
	ParentId   uint      `gorm:"default:0;comment:'父菜单编号'" json:"parent_id"`
	Children   []SysMenu `gorm:"-" json:"children"`                          // 子菜单集合
	Roles      []SysRole `gorm:"many2many:relation_role_menu;" json:"roles"` // 角色菜单多对多关系
}

func (m SysMenu) TableName() string {
	return m.Model.TableName("sys_menu")
}
