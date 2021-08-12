package initialize

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

var (
	creator      = "系统自动创建"
	status       = uint(1)
	visible      = uint(1)
	noBreadcrumb = uint(0)
)

// 初始化数据
func Data() {
	if !global.Conf.System.InitData {
		return
	}
	s := service.New(nil)
	db := global.Mysql.WithContext(s.RequestIdContext(requestId))
	// 1. 初始化角色
	newRoles := make([]models.SysRole, 0)
	roles := []models.SysRole{
		{
			Name:    "超级管理员",
			Keyword: "super",
			Desc:    "超级管理员",
		},
		{
			Name:    "访客",
			Keyword: "guest",
			Desc:    "外来访问人员",
		},
	}
	for i, role := range roles {
		sort := uint(i)
		id := uint(i + 1)
		roles[i].Id = id
		oldRole := models.SysRole{}
		err := db.Where("id = ?", id).First(&oldRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role.Id = id
			role.Creator = creator
			role.Status = &status
			if role.Sort == nil {
				role.Sort = &sort
			}
			newRoles = append(newRoles, role)
		}
	}
	if len(newRoles) > 0 {
		db.Create(&newRoles)
	}

	// 2. 初始化菜单

	menus := []models.SysMenu{
		{
			Name:  "dashboardRoot", // 对于想让子菜单显示在上层不显示的父级菜单不设置名字
			Title: "首页根目录",
			Icon:  "dashboard",
			Path:  "/dashboard",
			Roles: roles,
			Children: []models.SysMenu{
				{
					Name:      "dashboard",
					Title:     "首页",
					Icon:      "dashboard",
					Path:      "index",
					Component: "/dashboard/index",
					Roles:     roles,
				},
			},
		},
		{
			Name:  "systemRoot",
			Title: "系统设置根目录",
			Icon:  "component",
			Path:  "/system",
			Children: []models.SysMenu{
				{
					Name:      "menu",
					Title:     "菜单管理",
					Icon:      "tree-table",
					Path:      "menu", // 子菜单不用全路径, 自动继承
					Component: "/system/menu",
				},
				{
					Name:      "role",
					Title:     "角色管理",
					Icon:      "peoples",
					Path:      "role",
					Component: "/system/role",
				},
				{
					Name:      "user",
					Title:     "用户管理",
					Icon:      "user",
					Path:      "user",
					Component: "/system/user",
				},
				{
					Name:      "api",
					Title:     "接口管理",
					Icon:      "tree",
					Path:      "api",
					Component: "/system/api",
				},
				{
					Name:      "workflow",
					Title:     "工作流管理",
					Icon:      "example",
					Path:      "workflow",
					Component: "/system/workflow",
				},
				{
					Name:      "operation-log",
					Title:     "操作日志",
					Icon:      "example",
					Path:      "operation-log",
					Component: "/system/operation-log",
				},
				{
					Name:      "message-push",
					Title:     "消息推送",
					Icon:      "guide",
					Path:      "message-push",
					Component: "/system/message-push",
				},
				{
					Name:      "machine",
					Title:     "机器管理",
					Icon:      "guide",
					Path:      "machine",
					Component: "/system/machine",
				},
			},
		},
		{
			Name:  "testRoot",
			Title: "测试页面",
			Icon:  "bug",
			Path:  "/test",
			Children: []models.SysMenu{
				{
					Name:      "test",
					Title:     "测试用例",
					Icon:      "bug",
					Path:      "index",
					Component: "/test/index",
				},
				{
					Name:      "leave",
					Title:     "我的请假条",
					Icon:      "skill",
					Path:      "leave",
					Component: "/test/leave",
				},
				{
					Name:      "approving",
					Title:     "待审批列表",
					Icon:      "form",
					Path:      "approving",
					Component: "/test/approving",
				},
			},
		},
		{
			Name:  "uploader",
			Title: "上传组件",
			Icon:  "back-top",
			Path:  "/uploader",
			Children: []models.SysMenu{
				{
					Name:      "uploader1",
					Title:     "上传示例1",
					Icon:      "guide",
					Path:      "uploader1",
					Component: "/uploader/uploader1",
				},
				{
					Name:      "uploader2",
					Title:     "上传示例2",
					Icon:      "guide",
					Path:      "uploader2",
					Component: "/uploader/uploader2",
				},
			},
		},
	}
	menus = genMenu(0, menus, roles[0])
	createMenu(db, menus)

	// 3. 初始化用户
	// 默认头像
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := []models.SysUser{
		{
			Username:     "super",
			Password:     utils.GenPwd("123456"),
			Mobile:       "19999999999",
			Avatar:       avatar,
			Nickname:     "超级管理员",
			Introduction: "我是超管我怕谁？",
		},
		{
			Username:     "guest",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13999999999",
			Avatar:       avatar,
			Nickname:     "访客",
			Introduction: "这个人很懒, 什么也没留下",
		},
	}
	newUsers := make([]models.SysUser, 0)
	for i, user := range users {
		id := uint(i + 1)
		oldUser := models.SysUser{}
		err := db.Where("id = ?", id).First(&oldUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.Id = id
			user.Creator = creator
			if user.RoleId == 0 {
				user.RoleId = id
			}
			newUsers = append(newUsers, user)
		}
	}
	if len(newUsers) > 0 {
		db.Create(&newUsers)
	}

	// 4. 初始化接口
	apis := []models.SysApi{
		{
			Method:   "POST",
			Path:     "/v1/base/login",
			Category: "base",
			Desc:     "用户登录",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/logout",
			Category: "base",
			Desc:     "用户登出",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/refreshToken",
			Category: "base",
			Desc:     "刷新JWT令牌",
		},
		{
			Method:   "GET",
			Path:     "/v1/base/idempotenceToken",
			Category: "base",
			Desc:     "获取幂等性token",
		},
		{
			Method:   "POST",
			Path:     "/v1/user/info",
			Category: "user",
			Desc:     "获取当前登录用户信息",
		},
		{
			Method:   "GET",
			Path:     "/v1/user/list",
			Category: "user",
			Desc:     "获取用户列表",
		},
		{
			Method:   "PUT",
			Path:     "/v1/user/changePwd",
			Category: "user",
			Desc:     "修改用户登录密码",
		},
		{
			Method:   "POST",
			Path:     "/v1/user/create",
			Category: "user",
			Desc:     "创建用户",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/user/update/:userId",
			Category: "user",
			Desc:     "更新用户",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/user/delete/batch",
			Category: "user",
			Desc:     "批量删除用户",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/tree",
			Category: "menu",
			Desc:     "获取权限菜单",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/list",
			Category: "menu",
			Desc:     "获取菜单列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/menu/create",
			Category: "menu",
			Desc:     "创建菜单",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/menu/update/:menuId",
			Category: "menu",
			Desc:     "更新菜单",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/menu/delete/batch",
			Category: "menu",
			Desc:     "批量删除菜单",
		},
		{
			Method:   "GET",
			Path:     "/v1/role/list",
			Category: "role",
			Desc:     "获取角色列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "role",
			Desc:     "创建角色",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/update/:roleId",
			Category: "role",
			Desc:     "更新角色",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/role/delete/batch",
			Category: "role",
			Desc:     "批量删除角色",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/list",
			Category: "api",
			Desc:     "获取接口列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/api/create",
			Category: "api",
			Desc:     "创建接口",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/api/update/:roleId",
			Category: "api",
			Desc:     "更新接口",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/api/delete/batch",
			Category: "api",
			Desc:     "批量删除接口",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/all/:roleId",
			Category: "menu",
			Desc:     "查询指定角色的菜单树",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/all/category/:roleId",
			Category: "api",
			Desc:     "查询指定角色的接口(以分类分组)",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/menus/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限菜单",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/apis/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限接口",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/list",
			Category: "workflow",
			Desc:     "获取工作流列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/workflow/create",
			Category: "workflow",
			Desc:     "创建工作流",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/update/:roleId",
			Category: "workflow",
			Desc:     "更新工作流",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/workflow/delete/batch",
			Category: "workflow",
			Desc:     "批量删除工作流",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/line/list",
			Category: "workflow",
			Desc:     "获取流水线列表",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/line/update",
			Category: "workflow",
			Desc:     "更新流水线",
		},
		{
			Method:   "GET",
			Path:     "/v1/leave/list",
			Category: "leave",
			Desc:     "获取请假条列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/leave/create",
			Category: "leave",
			Desc:     "创建请假条",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/leave/update/:leaveId",
			Category: "leave",
			Desc:     "更新请假条",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/leave/delete/batch",
			Category: "leave",
			Desc:     "批量删除请假条",
		},
		{
			Method:   "GET",
			Path:     "/v1/leave/approval/list/:leaveId",
			Category: "leave",
			Desc:     "获取请假条审批记录",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/approving/list",
			Category: "workflow",
			Desc:     "获取当前登录用户待审批记录",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/log/approval",
			Category: "workflow",
			Desc:     "审批工作流日志",
		},
		{
			Method:   "GET",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "获取文件块信息以及上传完成部分",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "上传文件(分片)",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/merge",
			Category: "upload",
			Desc:     "合并分片文件",
		},
		{
			Method:   "GET",
			Path:     "/v1/operation/log/list",
			Category: "operation-log",
			Desc:     "获取操作日志列表",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/operation/log/delete/batch",
			Category: "operation-log",
			Desc:     "批量删除操作日志",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/unzip",
			Category: "upload",
			Desc:     "解压ZIP文件",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/all",
			Category: "message",
			Desc:     "获取全部消息",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/unRead/count",
			Category: "message",
			Desc:     "获取未读消息条数",
		},
		{
			Method:   "POST",
			Path:     "/v1/message/push",
			Category: "message",
			Desc:     "发送新消息",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/read/batch",
			Category: "message",
			Desc:     "批量标为已读",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/deleted/batch",
			Category: "message",
			Desc:     "批量标为删除",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/read/all",
			Category: "message",
			Desc:     "全部标为已读",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/deleted/all",
			Category: "message",
			Desc:     "全部标为删除",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/ws",
			Category: "message",
			Desc:     "消息中心长连接",
		},
		{
			Method:   "GET",
			Path:     "/v1/machine/shell/ws",
			Category: "machine",
			Desc:     "机器终端shell长连接",
		},
		{
			Method:   "GET",
			Path:     "/v1/machine/list",
			Category: "machine",
			Desc:     "获取机器列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/machine/create",
			Category: "machine",
			Desc:     "创建机器",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/machine/update/:machineId",
			Category: "machine",
			Desc:     "更新机器",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/machine/connect/:machineId",
			Category: "machine",
			Desc:     "连接或刷新机器状态",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/machine/delete/batch",
			Category: "machine",
			Desc:     "批量删除机器",
		},
		{
			Method:   "GET",
			Path:     "/v1/dict/list",
			Category: "dict",
			Desc:     "获取数据字典列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/dict/create",
			Category: "dict",
			Desc:     "创建数据字典",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/dict/update/:dictId",
			Category: "dict",
			Desc:     "更新数据字典",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/dict/delete/batch",
			Category: "dict",
			Desc:     "批量删除数据字典",
		},
		{
			Method:   "GET",
			Path:     "/v1/dict/data/list",
			Category: "dict",
			Desc:     "获取数据字典数据列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/dict/data/create",
			Category: "dict",
			Desc:     "创建数据字典数据",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/dict/data/update/:dictDataId",
			Category: "dict",
			Desc:     "更新数据字典数据",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/dict/data/delete/batch",
			Category: "dict",
			Desc:     "批量删除数据字典数据",
		},
	}
	newApis := make([]models.SysApi, 0)
	newRoleCasbins := make([]models.SysRoleCasbin, 0)
	for i, api := range apis {
		id := uint(i + 1)
		oldApi := models.SysApi{}
		err := db.Where("id = ?", id).First(&oldApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.Id = id
			api.Creator = creator
			newApis = append(newApis, api)
			// 超级管理员拥有所有API权限
			newRoleCasbins = append(newRoleCasbins, models.SysRoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/base/idempotenceToken",
				"/user/info",
				"/menu/tree",
				"/message/all",
				"/message/unRead/count",
				"/message/read/batch",
				"/message/deleted/batch",
				"/message/read/all",
				"/message/ws",
			}
			p := strings.TrimPrefix(api.Path, "/"+global.Conf.System.ApiVersion)
			if utils.Contains(basePaths, p) {
				// 非超级管理员有基础权限
				for i := 1; i < len(roles); i++ {
					newRoleCasbins = append(newRoleCasbins, models.SysRoleCasbin{
						Keyword: roles[i].Keyword,
						Path:    api.Path,
						Method:  api.Method,
					})
				}
			}
		}
	}
	if len(newApis) > 0 {
		db.Create(&newApis)
	}
	if len(newRoleCasbins) > 0 {
		s := service.New(nil)
		s.CreateRoleCasbins(newRoleCasbins)
	}
}

var menuTotal = 0

// 生成菜单
func genMenu(parentId uint, menus []models.SysMenu, superRole models.SysRole) []models.SysMenu {
	newMenus := make([]models.SysMenu, len(menus))
	// sort
	for i, menu := range menus {
		sort := uint(i)
		menu.Sort = &sort
		menu.Creator = creator
		menu.Status = &status
		menu.Visible = &visible
		newMenus[i] = menu
	}
	// id
	for i, menu := range newMenus {
		menuTotal++
		menu.Id = uint(menuTotal)
		newMenus[i] = menu
	}
	// children
	for i, menu := range newMenus {
		menu.Children = genMenu(menu.Id, menu.Children, superRole)
		newMenus[i] = menu
	}
	// parentId
	for i, menu := range newMenus {
		if parentId > 0 {
			menu.ParentId = parentId
		} else {
			menu.Component = ""             // 如果包含子菜单, Component为空
			menu.Breadcrumb = &noBreadcrumb // 面包屑不可见
		}
		if menu.Roles == nil {
			menu.Roles = []models.SysRole{
				superRole,
			}
		}
		newMenus[i] = menu
	}
	return newMenus
}

// 创建菜单
func createMenu(db *gorm.DB, menus []models.SysMenu) {
	for _, menu := range menus {
		oldMenu := models.SysMenu{}
		err := db.Where("id = ?", menu.Id).First(&oldMenu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&menu)
		}
		if len(menu.Children) > 0 {
			createMenu(db, menu.Children)
		}
	}
}
