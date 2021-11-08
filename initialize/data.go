package initialize

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/utils"
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
			Name:    "Super Admin Role",
			Keyword: "super",
			Desc:    "Super administrator role",
		},
		{
			Name:    "Guest Role",
			Keyword: "guest",
			Desc:    "External visitor role",
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
	menus := []ms.SysMenu{
		{
			Name:  "dashboardRoot",
			Title: "Dashboard Root",
			Icon:  "dashboard",
			Path:  "/dashboard",
			Children: []ms.SysMenu{
				{
					Name:      "dashboard",
					Title:     "Index",
					Icon:      "dashboard",
					Path:      "index",
					Component: "/dashboard/index",
				},
			},
			RoleIds: []uint{
				roles[1].Id,
			},
		},
		{
			Name:  "systemRoot",
			Title: "System Root",
			Icon:  "system",
			Path:  "/system",
			Children: []ms.SysMenu{
				{
					Name:      "menu",
					Title:     "Menus",
					Icon:      "menu",
					Path:      "menu",
					Component: "/system/menu",
				},
				{
					Name:      "role",
					Title:     "Roles",
					Icon:      "role",
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
					Icon:      "api",
					Path:      "api",
					Component: "/system/api",
				},
				{
					Name:      "operationLog",
					Title:     "Operation Logs",
					Icon:      "log",
					Path:      "operation-log",
					Component: "/system/operation-log",
				},
				{
					Name:      "messagePush",
					Title:     "Message Push",
					Icon:      "push",
					Path:      "message-push",
					Component: "/system/message-push",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
				{
					Name:      "machine",
					Title:     "Machines",
					Icon:      "machine",
					Path:      "machine",
					Component: "/system/machine",
				},
			},
		},
		{
			Name:  "fsmRoot",
			Title: "Fsm",
			Icon:  "fsm",
			Path:  "/fsm",
			Children: []ms.SysMenu{
				{
					Name:      "fsm",
					Title:     "Machines",
					Icon:      "fsm-machine",
					Path:      "index",
					Component: "/fsm/machine",
				},
				{
					Name:      "leave",
					Title:     "My Leave",
					Icon:      "leave",
					Path:      "leave",
					Component: "/fsm/leave",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
				{
					Name:      "approving",
					Title:     "Approving",
					Icon:      "approve",
					Path:      "approving",
					Component: "/fsm/approving",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
			},
		},
		{
			Name:  "uploader",
			Title: "Uploader",
			Icon:  "upload",
			Path:  "/uploader",
			Children: []ms.SysMenu{
				{
					Name:      "uploader1",
					Title:     "Uploader1",
					Icon:      "image",
					Path:      "uploader1",
					Component: "/uploader/uploader1",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
				{
					Name:      "uploader2",
					Title:     "Uploader2",
					Icon:      "zip",
					Path:      "uploader2",
					Component: "/uploader/uploader2",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
			},
		},
		{
			Name:  "testRoot",
			Title: "Tests",
			Icon:  "debug",
			Path:  "/test",
			Children: []ms.SysMenu{
				{
					Name:      "test",
					Title:     "Test Case",
					Icon:      "test",
					Path:      "index",
					Component: "/test/index",
					RoleIds: []uint{
						roles[1].Id,
					},
				},
			},
		},
	}
	menus = genMenu(0, menus, roles[0])
	relations := createMenu(db, menus)
	if len(relations) > 0 {
		db.Create(relations)
	}

	// 3. init users
	// default avatar image
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := []models.SysUser{
		{
			Username:     "super",
			Password:     utils.GenPwd("123456"),
			Mobile:       "19999999999",
			Avatar:       avatar,
			Nickname:     "Super Admin",
			Introduction: "I'm super. Who am I afraid of ?",
		},
		{
			Username:     "guest",
			Password:     utils.GenPwd("123456"),
			Mobile:       "13999999999",
			Avatar:       avatar,
			Nickname:     "Guest",
			Introduction: "The man was lazy and left nothing",
		},
		{
			Username:     "leave",
			Password:     utils.GenPwd("123456"),
			Mobile:       "15999999999",
			Avatar:       avatar,
			Nickname:     "Leave Tester",
			Introduction: "The man was lazy and left nothing",
			RoleId:       roles[1].Id,
		},
	}
	newUsers := make([]models.SysUser, 0)
	for i, user := range users {
		id := uint(i + 1)
		users[i].Id = id
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
	apis := []ms.SysApi{
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
			Method:   "GET",
			Path:     "/v1/user/list/:ids",
			Category: "user",
			Desc:     "find users by ids",
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
			Path:     "/v1/user/update/:id",
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
			Path:     "/v1/menu/all/:id",
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
			Path:     "/v1/menu/update/:id",
			Category: "menu",
			Desc:     "update menu",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/menu/role/update/:id",
			Category: "menu",
			Desc:     "update role menus",
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
			Method:   "GET",
			Path:     "/v1/role/list/:ids",
			Category: "role",
			Desc:     "find roles by ids",
		},
		{
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "role",
			Desc:     "create role",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/update/:id",
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
			Path:     "/v1/api/update/:id",
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
			Path:     "/v1/api/all/category/:id",
			Category: "api",
			Desc:     "get all api by role id",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/api/role/update/:id",
			Category: "api",
			Desc:     "update role apis",
		},
		{
			Method:   "GET",
			Path:     "/v1/fsm/list",
			Category: "fsm",
			Desc:     "find fsm machines",
		},
		{
			Method:   "POST",
			Path:     "/v1/fsm/create",
			Category: "fsm",
			Desc:     "create fsm machine",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/fsm/update/:id",
			Category: "fsm",
			Desc:     "update fsm machine",
		},
		{
			Method:   "GET",
			Path:     "/v1/fsm/approving/list",
			Category: "fsm",
			Desc:     "find fsm pending approve logs",
		},
		{
			Method:   "GET",
			Path:     "/v1/fsm/log/track",
			Category: "fsm",
			Desc:     "find fsm log history track",
		},
		{
			Method:   "GET",
			Path:     "/v1/fsm/submitter/detail",
			Category: "fsm",
			Desc:     "get submitter fsm log detail",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/fsm/submitter/detail",
			Category: "fsm",
			Desc:     "update submitter fsm log detail",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/fsm/approve",
			Category: "fsm",
			Desc:     "approved/refused fsm log",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/fsm/cancel",
			Category: "fsm",
			Desc:     "cancelled fsm log",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/fsm/delete/batch",
			Category: "fsm",
			Desc:     "batch delete fsm log",
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
			Path:     "/v1/leave/update/:id",
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
			Path:     "/v1/machine/update/:id",
			Category: "machine",
			Desc:     "update machine",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/machine/connect/:id",
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
			Path:     "/v1/dict/update/:id",
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
			Path:     "/v1/dict/data/update/:id",
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
	newApis := make([]ms.SysApi, 0)
	newRoleCasbins := make([]ms.SysRoleCasbin, 0)
	for i, api := range apis {
		id := uint(i + 1)
		oldApi := ms.SysApi{}
		err := db.Where("id = ?", id).First(&oldApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.Id = id
			newApis = append(newApis, api)
			// super has all api permission
			newRoleCasbins = append(newRoleCasbins, ms.SysRoleCasbin{
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
				"/leave/list",
				"/leave/create",
				"/leave/update/:id",
				"/leave/delete/batch",
				"/fsm/approving/list",
				"/fsm/log/track",
				"/fsm/submitter/detail",
				"/fsm/approve",
				"/fsm/cancel",
			}
			p := strings.TrimPrefix(api.Path, "/"+global.Conf.System.ApiVersion)
			if utils.Contains(basePaths, p) {
				// basic permission
				for i := 1; i < len(roles); i++ {
					newRoleCasbins = append(newRoleCasbins, ms.SysRoleCasbin{
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
		s.Q.BatchCreateRoleCasbin(newRoleCasbins)
	}

	// 5. init leave fsm machine
	// auto migrate fsm
	f := fsm.New(global.Mysql, fsm.WithCtx(ctx))
	f.CreateMachine(req.FsmCreateMachine{
		Category:            req.NullUint(global.FsmCategoryLeave),
		Name:                "Leave approval workflow",
		SubmitterName:       "Leave Submitter",
		SubmitterEditFields: "Desc,StartTime,EndTime",
		Levels: []req.FsmCreateEvent{
			{
				Name:       "First Level",
				Edit:       req.NullUint(constant.One),
				EditFields: "StartTime,EndTime",
				Refuse:     req.NullUint(constant.One),
				Users:      req.IdsStr(fmt.Sprintf("%d", users[1].Id)),
			},
			{
				Name:       "Second Level",
				Edit:       req.NullUint(constant.One),
				EditFields: "StartTime,EndTime",
				Refuse:     req.NullUint(constant.One),
				Roles:      req.IdsStr(fmt.Sprintf("%d", roles[0].Id)),
			},
		},
	})
}

var menuTotal = 0

func genMenu(parentId uint, menus []ms.SysMenu, superRole models.SysRole) []ms.SysMenu {
	newMenus := make([]ms.SysMenu, len(menus))
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
		// add super role
		if !utils.ContainsUint(menu.RoleIds, superRole.Id) {
			menu.RoleIds = append(menu.RoleIds, superRole.Id)
		}
		newMenus[i] = menu
	}
	return newMenus
}

func createMenu(db *gorm.DB, menus []ms.SysMenu) []ms.SysMenuRoleRelation {
	relations := make([]ms.SysMenuRoleRelation, 0)
	for _, menu := range menus {
		oldMenu := ms.SysMenu{}
		err := db.Where("id = ?", menu.Id).First(&oldMenu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&menu)
			for _, id := range menu.RoleIds {
				relations = append(relations, ms.SysMenuRoleRelation{
					MenuId: menu.Id,
					RoleId: id,
				})
			}
		}
		if len(menu.Children) > 0 {
			childrenRelations := createMenu(db, menu.Children)
			if len(childrenRelations) > 0 {
				relations = append(relations, childrenRelations...)
			}
		}
	}
	return relations
}
