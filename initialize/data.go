package initialize

import (
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/utils"
)

// 初始化数据
func InitData() {
	// 1. 初始化角色
	creator := "系统自动创建"
	status := true
	visible := true
	roles := []models.SysRole{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:    "访客",
			Keyword: "guest",
			Desc:    "外来访问人员",
			Status:  &status,
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:    "测试",
			Keyword: "tester",
			Desc:    "系统测试工程师",
			Status:  &status,
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    "系统管理员",
			Status:  &status,
			Creator: creator,
		},
	}
	for _, role := range roles {
		oldRole := models.SysRole{}
		notFound := global.Mysql.Where("id = ?", role.Id).First(&oldRole).RecordNotFound()
		if notFound {
			global.Mysql.Create(&role)
		}
	}

	// 2. 初始化菜单
	noBreadcrumb := false
	menus := []models.SysMenu{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:       "dashboardRoot", // 对于想让子菜单显示在上层不显示的父级菜单不设置名字
			Title:      "首页根目录",
			Icon:       "dashboard",
			Path:       "/dashboard",
			Component:  "", // 如果包含子菜单, Component为空
			Sort:       0,
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb, // 面包屑不可见
			ParentId:   0,
			Creator:    creator,
			Roles:      roles,
		},
		{
			Model: models.Model{
				Id: 7,
			},
			Name:      "dashboard",
			Title:     "首页",
			Icon:      "dashboard",
			Path:      "index",
			Component: "/dashboard/index",
			Sort:      0,
			Status:    &status,
			Visible:   &visible,
			ParentId:  1,
			Creator:   creator,
			Roles:     roles,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:       "systemRoot",
			Title:      "系统设置根目录",
			Icon:       "component",
			Path:       "/system",
			Component:  "",
			Sort:       1,
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb, // 面包屑不可见
			ParentId:   0,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:       "testRoot",
			Title:      "测试用例根目录",
			Icon:       "bug",
			Path:       "/test",
			Component:  "",
			Sort:       2,
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   0,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
			},
		},
		{
			Model: models.Model{
				Id: 8,
			},
			Name:      "test",
			Title:     "测试用例",
			Icon:      "bug",
			Path:      "index",
			Component: "/test/index",
			Sort:      0,
			Status:    &status,
			Visible:   &visible,
			ParentId:  3,
			Creator:   creator,
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
			Icon:      "tree-table",
			Path:      "menu", // 子菜单不用全路径, 自动继承
			Component: "/system/menu",
			Sort:      0,
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
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
			Icon:      "peoples",
			Path:      "role",
			Component: "/system/role",
			Sort:      1,
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
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
			Icon:      "user",
			Path:      "user",
			Component: "/system/user",
			Sort:      2,
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[2],
			},
		},
	}
	for _, menu := range menus {
		oldMenu := models.SysMenu{}
		notFound := global.Mysql.Where("id = ?", menu.Id).First(&oldMenu).RecordNotFound()
		if notFound {
			global.Mysql.Create(&menu)
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
