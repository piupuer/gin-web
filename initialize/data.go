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
	status       = uint(1)
	visible      = uint(1)
	noBreadcrumb = uint(0)
)

func Data() {
	if !global.Conf.Mysql.InitData {
		return
	}
	db := global.Mysql.WithContext(ctx)
	// 1. init roles
	newRoles := make([]models.SysRole, 0)
	roles := []models.SysRole{
		{
			Name:    "Super Admin",
			Keyword: "super",
			Desc:    "Super Admin",
		},
		{
			Name:    "Guest",
			Keyword: "guest",
			Desc:    "foreign visitors",
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

	// 2. init menus
	menus := []models.SysMenu{
		{
			Name:  "dashboardRoot", 
			Title: "Dashboard Root",
			Icon:  "dashboard",
			Path:  "/dashboard",
			Roles: roles,
			Children: []models.SysMenu{
				{
					Name:      "dashboard",
					Title:     "Index",
					Icon:      "dashboard",
					Path:      "index",
					Component: "/dashboard/index",
					Roles:     roles,
				},
			},
		},
		{
			Name:  "systemRoot",
			Title: "System Root",
			Icon:  "component",
			Path:  "/system",
			Children: []models.SysMenu{
				{
					Name:      "menu",
					Title:     "Menus",
					Icon:      "tree-table",
					Path:      "menu", 
					Component: "/system/menu",
				},
				{
					Name:      "role",
					Title:     "Roles",
					Icon:      "peoples",
					Path:      "role",
					Component: "/system/role",
				},
				{
					Name:      "user",
					Title:     "Users",
					Icon:      "user",
					Path:      "user",
					Component: "/system/user",
				},
				{
					Name:      "api",
					Title:     "Apis",
					Icon:      "tree",
					Path:      "api",
					Component: "/system/api",
				},
				{
					Name:      "workflow",
					Title:     "Workflows",
					Icon:      "example",
					Path:      "workflow",
					Component: "/system/workflow",
				},
				{
					Name:      "operation-log",
					Title:     "Operation Logs",
					Icon:      "example",
					Path:      "operation-log",
					Component: "/system/operation-log",
				},
				{
					Name:      "message-push",
					Title:     "Message Push",
					Icon:      "guide",
					Path:      "message-push",
					Component: "/system/message-push",
				},
				{
					Name:      "machine",
					Title:     "Machines",
					Icon:      "guide",
					Path:      "machine",
					Component: "/system/machine",
				},
			},
		},
		{
			Name:  "testRoot",
			Title: "Tests",
			Icon:  "bug",
			Path:  "/test",
			Children: []models.SysMenu{
				{
					Name:      "test",
					Title:     "Test Case",
					Icon:      "bug",
					Path:      "index",
					Component: "/test/index",
				},
				{
					Name:      "leave",
					Title:     "My Leave",
					Icon:      "skill",
					Path:      "leave",
					Component: "/test/leave",
				},
				{
					Name:      "approving",
					Title:     "Approving",
					Icon:      "form",
					Path:      "approving",
					Component: "/test/approving",
				},
			},
		},
		{
			Name:  "uploader",
			Title: "Uploader",
			Icon:  "back-top",
			Path:  "/uploader",
			Children: []models.SysMenu{
				{
					Name:      "uploader1",
					Title:     "Uploader1",
					Icon:      "guide",
					Path:      "uploader1",
					Component: "/uploader/uploader1",
				},
				{
					Name:      "uploader2",
					Title:     "Uploader2",
					Icon:      "guide",
					Path:      "uploader2",
					Component: "/uploader/uploader2",
				},
			},
		},
	}
	menus = genMenu(0, menus, roles[0])
	createMenu(db, menus)

	// 3. init users
	// default avatar image
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := []models.SysUser{
		{
			Username:     "super",
			Password:     utils.GenPwd("123456"),
			Mobile:       "19999999999",
			Avatar:       avatar,
			Nickname:     "super admin",
			Introduction: "I'm super. Who am I afraid of ?",
		},
		{
			Username:     "guest",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13999999999",
			Avatar:       avatar,
			Nickname:     "guest",
			Introduction: "The man was lazy and left nothing",
		},
	}
	newUsers := make([]models.SysUser, 0)
	for i, user := range users {
		id := uint(i + 1)
		oldUser := models.SysUser{}
		err := db.Where("id = ?", id).First(&oldUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.Id = id
			if user.RoleId == 0 {
				user.RoleId = id
			}
			newUsers = append(newUsers, user)
		}
	}
	if len(newUsers) > 0 {
		db.Create(&newUsers)
	}

	// 4. init apis
	apis := []models.SysApi{
		{
			Method:   "POST",
			Path:     "/v1/base/login",
			Category: "base",
			Desc:     "login",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/logout",
			Category: "base",
			Desc:     "logout",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/refreshToken",
			Category: "base",
			Desc:     "refresh token",
		},
		{
			Method:   "GET",
			Path:     "/v1/base/idempotenceToken",
			Category: "base",
			Desc:     "get idempotence token",
		},
		{
			Method:   "POST",
			Path:     "/v1/user/info",
			Category: "user",
			Desc:     "get current login user info",
		},
		{
			Method:   "GET",
			Path:     "/v1/user/list",
			Category: "user",
			Desc:     "find users",
		},
		{
			Method:   "PUT",
			Path:     "/v1/user/changePwd",
			Category: "user",
			Desc:     "change user password",
		},
		{
			Method:   "POST",
			Path:     "/v1/user/create",
			Category: "user",
			Desc:     "create user",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/user/update/:userId",
			Category: "user",
			Desc:     "update user",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/user/delete/batch",
			Category: "user",
			Desc:     "batch delete users",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/tree",
			Category: "menu",
			Desc:     "get menu tree",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/list",
			Category: "menu",
			Desc:     "find menus",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/all/:roleId",
			Category: "menu",
			Desc:     "get all menu by role id",
		},
		{
			Method:   "POST",
			Path:     "/v1/menu/create",
			Category: "menu",
			Desc:     "create menu",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/menu/update/:menuId",
			Category: "menu",
			Desc:     "update menu",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/menu/delete/batch",
			Category: "menu",
			Desc:     "batch delete menu",
		},
		{
			Method:   "GET",
			Path:     "/v1/role/list",
			Category: "role",
			Desc:     "find roles",
		},
		{
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "role",
			Desc:     "create role",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/update/:roleId",
			Category: "role",
			Desc:     "update role",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/role/delete/batch",
			Category: "role",
			Desc:     "batch delete role",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/list",
			Category: "api",
			Desc:     "find apis",
		},
		{
			Method:   "POST",
			Path:     "/v1/api/create",
			Category: "api",
			Desc:     "create api",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/api/update/:roleId",
			Category: "api",
			Desc:     "update api",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/api/delete/batch",
			Category: "api",
			Desc:     "batch delete api",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/all/category/:roleId",
			Category: "api",
			Desc:     "get all api by role id",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/menus/update/:roleId",
			Category: "role",
			Desc:     "update role menus",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/apis/update/:roleId",
			Category: "role",
			Desc:     "update role apis",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/list",
			Category: "workflow",
			Desc:     "find workflows",
		},
		{
			Method:   "POST",
			Path:     "/v1/workflow/create",
			Category: "workflow",
			Desc:     "create workflow",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/update/:roleId",
			Category: "workflow",
			Desc:     "update workflow",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/workflow/delete/batch",
			Category: "workflow",
			Desc:     "batch delete workflow",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/line/list",
			Category: "workflow",
			Desc:     "find workflow lines",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/line/update",
			Category: "workflow",
			Desc:     "update workflow line",
		},
		{
			Method:   "GET",
			Path:     "/v1/leave/list",
			Category: "leave",
			Desc:     "find leaves",
		},
		{
			Method:   "POST",
			Path:     "/v1/leave/create",
			Category: "leave",
			Desc:     "create leave",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/leave/update/:leaveId",
			Category: "leave",
			Desc:     "update leave",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/leave/delete/batch",
			Category: "leave",
			Desc:     "batch delete leave",
		},
		{
			Method:   "GET",
			Path:     "/v1/leave/approval/list/:leaveId",
			Category: "leave",
			Desc:     "find leave approval logs",
		},
		{
			Method:   "GET",
			Path:     "/v1/workflow/approving/list",
			Category: "workflow",
			Desc:     "find current user approvings",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/workflow/log/approval",
			Category: "workflow",
			Desc:     "approve workflow log",
		},
		{
			Method:   "GET",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "get uploaded file info",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "upload file",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/merge",
			Category: "upload",
			Desc:     "merge file",
		},
		{
			Method:   "GET",
			Path:     "/v1/operation/log/list",
			Category: "operation-log",
			Desc:     "find operation logs",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/operation/log/delete/batch",
			Category: "operation-log",
			Desc:     "batch delete operation log",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/unzip",
			Category: "upload",
			Desc:     "unzip",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/all",
			Category: "message",
			Desc:     "find messages",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/unRead/count",
			Category: "message",
			Desc:     "get unread message count",
		},
		{
			Method:   "POST",
			Path:     "/v1/message/push",
			Category: "message",
			Desc:     "push new message",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/read/batch",
			Category: "message",
			Desc:     "batch marked as read",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/deleted/batch",
			Category: "message",
			Desc:     "batch marked as deleted",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/read/all",
			Category: "message",
			Desc:     "all marked as read",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/message/deleted/all",
			Category: "message",
			Desc:     "batch marked as deleted",
		},
		{
			Method:   "GET",
			Path:     "/v1/message/ws",
			Category: "message",
			Desc:     "message websocket",
		},
		{
			Method:   "GET",
			Path:     "/v1/machine/shell/ws",
			Category: "machine",
			Desc:     "machine shell websocket",
		},
		{
			Method:   "GET",
			Path:     "/v1/machine/list",
			Category: "machine",
			Desc:     "find machines",
		},
		{
			Method:   "POST",
			Path:     "/v1/machine/create",
			Category: "machine",
			Desc:     "create machine",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/machine/update/:machineId",
			Category: "machine",
			Desc:     "update machine",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/machine/connect/:machineId",
			Category: "machine",
			Desc:     "connect or refresh machine status",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/machine/delete/batch",
			Category: "machine",
			Desc:     "batch delete machine",
		},
		{
			Method:   "GET",
			Path:     "/v1/dict/list",
			Category: "dict",
			Desc:     "find dicts",
		},
		{
			Method:   "POST",
			Path:     "/v1/dict/create",
			Category: "dict",
			Desc:     "create dict",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/dict/update/:dictId",
			Category: "dict",
			Desc:     "update dict",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/dict/delete/batch",
			Category: "dict",
			Desc:     "batch delete dict",
		},
		{
			Method:   "GET",
			Path:     "/v1/dict/data/list",
			Category: "dict",
			Desc:     "find dict datas",
		},
		{
			Method:   "POST",
			Path:     "/v1/dict/data/create",
			Category: "dict",
			Desc:     "create dict data",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/dict/data/update/:dictDataId",
			Category: "dict",
			Desc:     "update dict data",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/dict/data/delete/batch",
			Category: "dict",
			Desc:     "batch delete dict data",
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
			newApis = append(newApis, api)
			// super has all api permission
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
				// basic permission
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

func genMenu(parentId uint, menus []models.SysMenu, superRole models.SysRole) []models.SysMenu {
	newMenus := make([]models.SysMenu, len(menus))
	// sort
	for i, menu := range menus {
		sort := uint(i)
		menu.Sort = &sort
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
			// The component of the submenu is empty
			menu.Component = ""
			menu.Breadcrumb = &noBreadcrumb
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
