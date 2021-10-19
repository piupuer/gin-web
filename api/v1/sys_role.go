package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	var r request.RoleReq
	req.ShouldBind(c, &r)
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	r.CurrentRoleSort = *user.Role.Sort

	s := cache_service.New(c)
	list, err := s.GetRoles(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.RoleResp{}, r.Page)
}

// 创建角色
func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateRoleReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())

	if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
		resp.CheckErr("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort)
	}

	// 记录当前创建人信息
	r.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.Q.Create(r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	var r request.UpdateRoleReq
	req.ShouldBind(c, &r)
	if r.Sort != nil {
		// 绑定当前用户角色排序(隐藏特定用户)
		user := GetCurrentUser(c)
		if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
			resp.CheckErr("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort)
		}
	}

	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		resp.CheckErr("角色编号不正确")
	}

	user := GetCurrentUser(c)
	if r.Status != nil && uint(*r.Status) == models.SysRoleStatusDisabled && roleId == user.RoleId {
		resp.CheckErr("不能禁用自己所在的角色")
	}

	s := service.New(c)
	err := s.Q.UpdateById(roleId, r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

// 更新角色的权限菜单
func UpdateRoleMenusById(c *gin.Context) {
	var r request.UpdateIncrementalIdsRequestStruct
	req.ShouldBind(c, &r)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		resp.CheckErr("角色编号不正确")
	}
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(r.Delete) > 0 {
			resp.CheckErr("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			resp.CheckErr("无法更改自己的权限, 如需更改请联系上级领导")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleMenusById(user.Role, roleId, r)
	resp.CheckErr(err)
	// 清理菜单树缓存
	menuTreeCache.Flush()
	resp.Success()
}

// 更新角色的权限接口
func UpdateRoleApisById(c *gin.Context) {
	var r request.UpdateIncrementalIdsRequestStruct
	req.ShouldBind(c, &r)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		resp.CheckErr("角色编号不正确")
	}

	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(r.Delete) > 0 {
			resp.CheckErr("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			resp.CheckErr("无法更改自己的权限, 如需更改请联系上级领导")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleApisById(roleId, r)
	resp.CheckErr(err)
	// 清理菜单树缓存
	menuTreeCache.Flush()
	resp.Success()
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.GetUintIds(), user.RoleId) {
		resp.CheckErr("不能删除自己所在的角色")
	}

	s := service.New(c)
	err := s.DeleteRoleByIds(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}
