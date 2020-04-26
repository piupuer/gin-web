package initialize

import (
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/utils"
)

// 初始化数据
func InitData() {
	// 1. 初始化角色
	roles := []models.SysRole{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:    "访客",
			Keyword: "guest",
			Desc:    "外来访问人员",
			Status:  true,
			Creator: "系统自动创建",
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:    "测试",
			Keyword: "tester",
			Desc:    "系统测试工程师",
			Status:  true,
			Creator: "系统自动创建",
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    "系统管理员",
			Status:  true,
			Creator: "系统自动创建",
		},
	}
	for _, role := range roles {
		oldRole := models.SysRole{}
		notFound := global.Mysql.Where("keyword = ?", role.Keyword).First(&oldRole).RecordNotFound()
		if notFound {
			global.Mysql.Create(&role)
		} else {
			_ = global.Mysql.Table(role.TableName()).Where("keyword = ?", role.Keyword).Update(&role).Error
		}
	}

	// 2. 初始化菜单
	menus := []models.SysMenu{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:      "dashboard",
			Title:     "主页",
			Icon:      "",
			Path:      "/dashboard",
			Component: "/dashboard",
			Sort:      0,
			Status:    true,
			Visible:   true,
			ParentId:  0,
			Roles:     roles,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:      "system",
			Title:     "系统设置",
			Icon:      "",
			Path:      "/system",
			Component: "/system",
			Sort:      1,
			Status:    true,
			Visible:   true,
			ParentId:  0,
			Roles: []models.SysRole{
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:      "example",
			Title:     "示例菜单",
			Icon:      "",
			Path:      "/example",
			Component: "/example",
			Sort:      2,
			Status:    true,
			Visible:   true,
			ParentId:  0,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Name:      "menu",
			Title:     "菜单管理",
			Icon:      "",
			Path:      "/system/menu",
			Component: "/system/user",
			Sort:      0,
			Status:    true,
			Visible:   true,
			ParentId:  2,
			Roles: []models.SysRole{
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Name:      "role",
			Title:     "角色管理",
			Icon:      "",
			Path:      "/system/role",
			Component: "/system/role",
			Sort:      1,
			Status:    true,
			Visible:   true,
			ParentId:  2,
			Roles: []models.SysRole{
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 6,
			},
			Name:      "user",
			Title:     "用户管理",
			Icon:      "",
			Path:      "/system/user",
			Component: "/system/user",
			Sort:      2,
			Status:    true,
			Visible:   true,
			ParentId:  2,
			Roles: []models.SysRole{
				roles[2],
			},
		},
	}
	for _, menu := range menus {
		oldMenu := models.SysMenu{}
		notFound := global.Mysql.Where("path = ?", menu.Path).First(&oldMenu).RecordNotFound()
		if notFound {
			global.Mysql.Create(&menu)
		} else {
			// 角色不更新
			menu.Roles = nil
			_ = global.Mysql.Table(menu.TableName()).Where("path = ?", menu.Path).Update(&menu).Error
		}
	}

	// 3. 初始化用户
	// 默认头像
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := [3]models.SysUser{
		{
			Username:     "admin",
			Password:     utils.GenPwd("123456"),
			Mobile:       "18888888888",
			Avatar:       avatar,
			Nickname:     "管理员",
			Introduction: "妖怪, 哪里跑",
			RoleId:       3,
		},
		{
			Username:     "zhangsan",
			Password:     utils.GenPwd("123456"),
			Mobile:       "15888888888",
			Avatar:       avatar,
			Nickname:     "张三",
			Introduction: "今天是个好日子",
			RoleId:       2,
		},
		{
			Username:     "lisi",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13888888888",
			Avatar:       avatar,
			Nickname:     "李四",
			Introduction: "这个人很懒, 什么也没留下",
			RoleId:       1,
		},
	}
	for _, user := range users {
		oldUser := models.SysUser{}
		notFound := global.Mysql.Where("username = ?", user.Username).First(&oldUser).RecordNotFound()
		if notFound {
			global.Mysql.Create(&user)
		}
	}
}
