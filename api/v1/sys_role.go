package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	var req request.RoleRequestStruct
	request.ShouldBind(c, &req)
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	req.CurrentRoleSort = *user.Role.Sort

	s := cache_service.New(c)
	roles, err := s.GetRoles(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.RoleListResponseStruct
	utils.Struct2StructByJson(roles, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建角色
func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateRoleRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())

	if req.Sort != nil && *user.Role.Sort > uint(*req.Sort) {
		response.CheckErr("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort)
	}

	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.Create(req, new(models.SysRole))
	response.CheckErr(err)
	response.Success()
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	var req request.UpdateRoleRequestStruct
	request.ShouldBind(c, &req)
	if req.Sort != nil {
		// 绑定当前用户角色排序(隐藏特定用户)
		user := GetCurrentUser(c)
		if req.Sort != nil && *user.Role.Sort > uint(*req.Sort) {
			response.CheckErr("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort)
		}
	}

	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.CheckErr("角色编号不正确")
	}

	user := GetCurrentUser(c)
	if req.Status != nil && uint(*req.Status) == models.SysRoleStatusDisabled && roleId == user.RoleId {
		response.CheckErr("不能禁用自己所在的角色")
	}

	s := service.New(c)
	err := s.UpdateById(roleId, req, new(models.SysRole))
	response.CheckErr(err)
	response.Success()
}

// 更新角色的权限菜单
func UpdateRoleMenusById(c *gin.Context) {
	var req request.UpdateIncrementalIdsRequestStruct
	request.ShouldBind(c, &req)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.CheckErr("角色编号不正确")
	}
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.CheckErr("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.CheckErr("无法更改自己的权限, 如需更改请联系上级领导")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleMenusById(user.Role, roleId, req)
	response.CheckErr(err)
	// 清理菜单树缓存
	menuTreeCache.Flush()
	response.Success()
}

// 更新角色的权限接口
func UpdateRoleApisById(c *gin.Context) {
	var req request.UpdateIncrementalIdsRequestStruct
	request.ShouldBind(c, &req)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.CheckErr("角色编号不正确")
	}

	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.CheckErr("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.CheckErr("无法更改自己的权限, 如需更改请联系上级领导")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleApisById(roleId, req)
	response.CheckErr(err)
	// 清理菜单树缓存
	menuTreeCache.Flush()
	response.Success()
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	user := GetCurrentUser(c)
	if utils.ContainsUint(req.GetUintIds(), user.RoleId) {
		response.CheckErr("不能删除自己所在的角色")
	}

	s := service.New(c)
	err := s.DeleteRoleByIds(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}
