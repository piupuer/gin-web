package initialize

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// 初始化数据
func InitData() {
	// 1. 初始化角色
	creator := "系统自动创建"
	status := uint(1)
	visible := uint(1)
	edit := uint(1)
	sorts := []uint{
		0, 1, 2, 3, 4, 5, 6, 7,
	}
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
			Sort:    &sorts[3],
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
			Sort:    &sorts[2],
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:    "超级管理员",
			Keyword: "super",
			Desc:    "超级管理员",
			Status:  &status,
			Creator: creator,
			Sort:    &sorts[0],
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Name:    "系统管理员",
			Keyword: "admin",
			Desc:    "系统管理员",
			Status:  &status,
			Creator: creator,
			Sort:    &sorts[1],
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Name:    "请假条提交员",
			Keyword: "leave",
			Desc:    "请假条提交员",
			Status:  &status,
			Creator: creator,
			Sort:    &sorts[2],
		},
	}
	for _, role := range roles {
		oldRole := models.SysRole{}
		err := global.Mysql.Where("id = ?", role.Id).First(&oldRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&role)
		}
	}

	// 2. 初始化菜单
	noBreadcrumb := uint(0)
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
			Sort:       &sorts[0],
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
			Sort:      &sorts[0],
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
			Sort:       &sorts[1],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb, // 面包屑不可见
			ParentId:   0,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[2],
				roles[3],
			},
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:       "testRoot",
			Title:      "测试页面",
			Icon:       "bug",
			Path:       "/test",
			Component:  "",
			Sort:       &sorts[2],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   0,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
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
			Sort:      &sorts[0],
			Status:    &status,
			Visible:   &visible,
			ParentId:  3,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
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
			Sort:      &sorts[0],
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
			Sort:      &sorts[1],
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[2],
				roles[3],
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
			Sort:      &sorts[2],
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[2],
				roles[3],
			},
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Name:      "api",
			Title:     "接口管理",
			Icon:      "tree",
			Path:      "api",
			Component: "/system/api",
			Sort:      &sorts[3],
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
				Id: 10,
			},
			Name:      "workflow",
			Title:     "工作流管理",
			Icon:      "example",
			Path:      "workflow",
			Component: "/system/workflow",
			Sort:      &sorts[4],
			Status:    &status,
			Visible:   &visible,
			ParentId:  2,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[2],
				roles[3],
			},
		},
		{
			Model: models.Model{
				Id: 15,
			},
			Name:      "operation-log",
			Title:     "操作日志",
			Icon:      "example",
			Path:      "operation-log",
			Component: "/system/operation-log",
			Sort:      &sorts[5],
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
				Id: 11,
			},
			Name:      "leave",
			Title:     "我的请假条",
			Icon:      "skill",
			Path:      "leave",
			Component: "/test/leave",
			Sort:      &sorts[2],
			Status:    &status,
			Visible:   &visible,
			ParentId:  3,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 12,
			},
			Name:      "approving",
			Title:     "待审批列表",
			Icon:      "form",
			Path:      "approving",
			Component: "/test/approving",
			Sort:      &sorts[3],
			Status:    &status,
			Visible:   &visible,
			ParentId:  3,
			Creator:   creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 13,
			},
			Name:       "uploader",
			Title:      "上传组件",
			Icon:       "back-top",
			Path:       "/uploader",
			Component:  "",
			Sort:       &sorts[3],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   0,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 14,
			},
			Name:       "uploader1",
			Title:      "上传示例1",
			Icon:       "guide",
			Path:       "uploader1",
			Component:  "/uploader/uploader1",
			Sort:       &sorts[0],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   13,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 16,
			},
			Name:       "uploader2",
			Title:      "上传示例2",
			Icon:       "guide",
			Path:       "uploader2",
			Component:  "/uploader/uploader2",
			Sort:       &sorts[1],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   13,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[1],
				roles[2],
				roles[3],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 17,
			},
			Name:       "message-push",
			Title:      "消息推送",
			Icon:       "guide",
			Path:       "message-push",
			Component:  "/system/message-push",
			Sort:       &sorts[6],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   2,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[2],
				roles[3],
				roles[4],
			},
		},
		{
			Model: models.Model{
				Id: 18,
			},
			Name:       "machine",
			Title:      "机器管理",
			Icon:       "guide",
			Path:       "machine",
			Component:  "/system/machine",
			Sort:       &sorts[7],
			Status:     &status,
			Visible:    &visible,
			Breadcrumb: &noBreadcrumb,
			ParentId:   2,
			Creator:    creator,
			Roles: []models.SysRole{
				roles[2],
			},
		},
	}
	for _, menu := range menus {
		oldMenu := models.SysMenu{}
		err := global.Mysql.Where("id = ?", menu.Id).First(&oldMenu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&menu)
		}
	}

	// 3. 初始化用户
	// 默认头像
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := [5]models.SysUser{
		{
			Model: models.Model{
				Id: 1,
			},
			Username:     "super",
			Password:     utils.GenPwd("123456"),
			Mobile:       "19999999999",
			Avatar:       avatar,
			Nickname:     "超级管理员",
			Introduction: "我是超管我怕谁？",
			RoleId:       3,
			Creator:      creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Username:     "admin",
			Password:     utils.GenPwd("123456"),
			Mobile:       "18888888888",
			Avatar:       avatar,
			Nickname:     "系统管理员",
			Introduction: "妖怪, 哪里跑",
			RoleId:       4,
			Creator:      creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Username:     "zhangsan",
			Password:     utils.GenPwd("123456"),
			Mobile:       "15888888888",
			Avatar:       avatar,
			Nickname:     "张三",
			Introduction: "今天是个好日子",
			RoleId:       2,
			Creator:      creator,
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Username:     "lisi",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13888888888",
			Avatar:       avatar,
			Nickname:     "李四",
			Introduction: "这个人很懒, 什么也没留下",
			RoleId:       1,
			Creator:      creator,
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Username:     "wangwu",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13999999999",
			Avatar:       avatar,
			Nickname:     "王武",
			Introduction: "这个人很懒, 什么也没留下",
			RoleId:       5,
			Creator:      creator,
		},
	}
	for _, user := range users {
		oldUser := models.SysUser{}
		err := global.Mysql.Where("username = ?", user.Username).First(&oldUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&user)
		}
	}

	// 4. 初始化接口
	apis := []models.SysApi{
		{
			Model: models.Model{
				Id: 1,
			},
			Method:   "POST",
			Path:     "/v1/base/login",
			Category: "base",
			Desc:     "用户登录",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Method:   "POST",
			Path:     "/v1/base/logout",
			Category: "base",
			Desc:     "用户登出",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Method:   "POST",
			Path:     "/v1/base/refresh_token",
			Category: "base",
			Desc:     "刷新JWT令牌",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Method:   "POST",
			Path:     "/v1/user/info",
			Category: "user",
			Desc:     "获取当前登录用户信息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Method:   "GET",
			Path:     "/v1/user/list",
			Category: "user",
			Desc:     "获取用户列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 6,
			},
			Method:   "PUT",
			Path:     "/v1/user/changePwd",
			Category: "user",
			Desc:     "修改用户登录密码",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 7,
			},
			Method:   "POST",
			Path:     "/v1/user/create",
			Category: "user",
			Desc:     "创建用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 8,
			},
			Method:   "PATCH",
			Path:     "/v1/user/update/:userId",
			Category: "user",
			Desc:     "更新用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Method:   "DELETE",
			Path:     "/v1/user/delete/batch",
			Category: "user",
			Desc:     "批量删除用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 10,
			},
			Method:   "GET",
			Path:     "/v1/menu/tree",
			Category: "menu",
			Desc:     "获取权限菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 11,
			},
			Method:   "GET",
			Path:     "/v1/menu/list",
			Category: "menu",
			Desc:     "获取菜单列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 12,
			},
			Method:   "POST",
			Path:     "/v1/menu/create",
			Category: "menu",
			Desc:     "创建菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 13,
			},
			Method:   "PATCH",
			Path:     "/v1/menu/update/:menuId",
			Category: "menu",
			Desc:     "更新菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 14,
			},
			Method:   "DELETE",
			Path:     "/v1/menu/delete/batch",
			Category: "menu",
			Desc:     "批量删除菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 15,
			},
			Method:   "GET",
			Path:     "/v1/role/list",
			Category: "role",
			Desc:     "获取角色列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 16,
			},
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "role",
			Desc:     "创建角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 17,
			},
			Method:   "PATCH",
			Path:     "/v1/role/update/:roleId",
			Category: "role",
			Desc:     "更新角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 18,
			},
			Method:   "DELETE",
			Path:     "/v1/role/delete/batch",
			Category: "role",
			Desc:     "批量删除角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 19,
			},
			Method:   "GET",
			Path:     "/v1/api/list",
			Category: "api",
			Desc:     "获取接口列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 20,
			},
			Method:   "POST",
			Path:     "/v1/api/create",
			Category: "api",
			Desc:     "创建接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 21,
			},
			Method:   "PATCH",
			Path:     "/v1/api/update/:roleId",
			Category: "api",
			Desc:     "更新接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 22,
			},
			Method:   "DELETE",
			Path:     "/v1/api/delete/batch",
			Category: "api",
			Desc:     "批量删除接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 23,
			},
			Method:   "GET",
			Path:     "/v1/menu/all/:roleId",
			Category: "menu",
			Desc:     "查询指定角色的菜单树",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 24,
			},
			Method:   "GET",
			Path:     "/v1/api/all/category/:roleId",
			Category: "api",
			Desc:     "查询指定角色的接口(以分类分组)",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 25,
			},
			Method:   "PATCH",
			Path:     "/v1/role/menus/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 26,
			},
			Method:   "PATCH",
			Path:     "/v1/role/apis/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 27,
			},
			Method:   "GET",
			Path:     "/v1/workflow/list",
			Category: "workflow",
			Desc:     "获取工作流列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 28,
			},
			Method:   "POST",
			Path:     "/v1/workflow/create",
			Category: "workflow",
			Desc:     "创建工作流",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 29,
			},
			Method:   "PATCH",
			Path:     "/v1/workflow/update/:roleId",
			Category: "workflow",
			Desc:     "更新工作流",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 30,
			},
			Method:   "DELETE",
			Path:     "/v1/workflow/delete/batch",
			Category: "workflow",
			Desc:     "批量删除工作流",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 31,
			},
			Method:   "GET",
			Path:     "/v1/workflow/line/list",
			Category: "workflow",
			Desc:     "获取流水线列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 32,
			},
			Method:   "PATCH",
			Path:     "/v1/workflow/line/update",
			Category: "workflow",
			Desc:     "更新流水线",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 33,
			},
			Method:   "GET",
			Path:     "/v1/leave/list",
			Category: "leave",
			Desc:     "获取请假条列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 34,
			},
			Method:   "POST",
			Path:     "/v1/leave/create",
			Category: "leave",
			Desc:     "创建请假条",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 35,
			},
			Method:   "PATCH",
			Path:     "/v1/leave/update/:leaveId",
			Category: "leave",
			Desc:     "更新请假条",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 36,
			},
			Method:   "DELETE",
			Path:     "/v1/leave/delete/batch",
			Category: "leave",
			Desc:     "批量删除请假条",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 37,
			},
			Method:   "GET",
			Path:     "/v1/leave/approval/list/:leaveId",
			Category: "leave",
			Desc:     "获取请假条审批记录",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 38,
			},
			Method:   "GET",
			Path:     "/v1/workflow/approving/list",
			Category: "workflow",
			Desc:     "获取当前登录用户待审批记录",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 39,
			},
			Method:   "PATCH",
			Path:     "/v1/workflow/log/approval",
			Category: "workflow",
			Desc:     "审批工作流日志",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 40,
			},
			Method:   "GET",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "获取文件块信息以及上传完成部分",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 41,
			},
			Method:   "POST",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "上传文件(分片)",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 42,
			},
			Method:   "POST",
			Path:     "/v1/upload/merge",
			Category: "upload",
			Desc:     "合并分片文件",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 43,
			},
			Method:   "GET",
			Path:     "/v1/operation/log/list",
			Category: "operation-log",
			Desc:     "获取操作日志列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 44,
			},
			Method:   "DELETE",
			Path:     "/v1/operation/log/delete/batch",
			Category: "operation-log",
			Desc:     "批量删除操作日志",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 45,
			},
			Method:   "POST",
			Path:     "/v1/upload/unzip",
			Category: "upload",
			Desc:     "解压ZIP文件",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 46,
			},
			Method:   "GET",
			Path:     "/v1/message/all",
			Category: "message",
			Desc:     "获取全部消息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 47,
			},
			Method:   "GET",
			Path:     "/v1/message/unRead/count",
			Category: "message",
			Desc:     "获取未读消息条数",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 48,
			},
			Method:   "POST",
			Path:     "/v1/message/push",
			Category: "message",
			Desc:     "发送新消息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 49,
			},
			Method:   "PATCH",
			Path:     "/v1/message/read/batch",
			Category: "message",
			Desc:     "批量标为已读",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 50,
			},
			Method:   "PATCH",
			Path:     "/v1/message/deleted/batch",
			Category: "message",
			Desc:     "批量标为删除",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 51,
			},
			Method:   "PATCH",
			Path:     "/v1/message/read/all",
			Category: "message",
			Desc:     "全部标为已读",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 52,
			},
			Method:   "PATCH",
			Path:     "/v1/message/deleted/all",
			Category: "message",
			Desc:     "全部标为删除",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 53,
			},
			Method:   "GET",
			Path:     "/v1/message/ws",
			Category: "message",
			Desc:     "消息中心长连接",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 54,
			},
			Method:   "GET",
			Path:     "/v1/machine/shell/ws",
			Category: "machine",
			Desc:     "机器终端shell长连接",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 55,
			},
			Method:   "GET",
			Path:     "/v1/machine/list",
			Category: "machine",
			Desc:     "获取机器列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 56,
			},
			Method:   "POST",
			Path:     "/v1/machine/create",
			Category: "machine",
			Desc:     "创建机器",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 57,
			},
			Method:   "PATCH",
			Path:     "/v1/machine/update/:machineId",
			Category: "machine",
			Desc:     "更新机器",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 58,
			},
			Method:   "PATCH",
			Path:     "/v1/machine/connect/:machineId",
			Category: "machine",
			Desc:     "连接或刷新机器状态",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 59,
			},
			Method:   "DELETE",
			Path:     "/v1/machine/delete/batch",
			Category: "machine",
			Desc:     "批量删除机器",
			Creator:  creator,
		},
	}
	for _, api := range apis {
		oldApi := models.SysApi{}
		err := global.Mysql.Where("id = ?", api.Id).First(&oldApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&api)
			// 创建服务
			s := service.New(nil)
			// 超级管理员拥有所有API权限role[2]
			s.CreateRoleCasbin(models.SysRoleCasbin{
				Keyword: roles[2].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
			// 管理员没有菜单管理/接口管理/操作日志/机器管理权限
			if !utils.Contains([]uint{
				11,
				12,
				13,
				14,
				19,
				20,
				21,
				22,
				43,
				44,
				54,
				55,
				56,
				57,
				58,
				59,
			}, api.Id) {
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[3].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}

			// 测试/请假条测试员: 拥有请假条与审批请假条权限
			if api.Id >= 33 && api.Id <= 39 {
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[4].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
			// 消息推送权限只给管理员和测试员
			if api.Id == 48 {
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
			// 其他人暂时只有登录/获取用户信息的权限
			if api.Id < 5 || api.Id == 10 || (api.Id > 45 && api.Id != 48 && api.Id < 54) {
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[0].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				s.CreateRoleCasbin(models.SysRoleCasbin{
					Keyword: roles[4].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	// 5. 初始化工作流
	workflows := []models.SysWorkflow{
		{
			Model: models.Model{
				Id: 1,
			},
			Uuid:     uuid.NewV4().String(),
			Category: models.SysWorkflowTargetCategoryLeave,
			Name:     "请假审批流程",
			Desc:     "用于员工请假",
			Creator:  creator,
		},
	}
	for _, workflow := range workflows {
		oldWorkflow := models.SysWorkflow{}
		err := global.Mysql.Where("id = ?", workflow.Id).First(&oldWorkflow).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&workflow)
			// 创建服务
			s := service.New(nil)
			// 设置张三和管理员审批
			req := request.UpdateWorkflowLineIncrementalRequestStruct{
				FlowId: 1,
				Create: []request.UpdateWorkflowLineRequestStruct{
					{
						Id:     1,
						FlowId: 1,
						UserIds: []uint{
							3,
						},
						Edit: &edit,
						Name: "主管审批",
					},
					{
						Id:     2,
						FlowId: 1,
						RoleId: &roles[3].Id,
						Edit:   &edit,
						Name:   "总经理审批",
					},
				},
			}
			s.UpdateWorkflowLineByIncremental(&req)
		}
	}
}
