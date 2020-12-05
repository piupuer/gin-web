package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req request.RoleListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	req.CurrentRoleSort = *user.Role.Sort

	// 创建服务
	s := cache_service.New(c)
	roles, err := s.GetRoles(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.RoleListResponseStruct
	utils.Struct2StructByJson(roles, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建角色
func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateRoleRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 参数校验
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	if *user.Role.Sort > *req.Sort {
		response.FailWithMsg(fmt.Sprintf("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort))
		return
	}

	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	err = s.CreateRole(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	// 绑定参数
	var req models.SysRole
	var roleInfo request.CreateRoleRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	utils.Struct2StructByJson(req, &roleInfo)

	if roleInfo.Sort != nil {
		// 绑定当前用户角色排序(隐藏特定用户)
		user := GetCurrentUser(c)
		if *user.Role.Sort > *roleInfo.Sort {
			response.FailWithMsg(fmt.Sprintf("角色排序不允许比当前登录账号序号(%d)小", *user.Role.Sort))
			return
		}
	}

	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}

	user := GetCurrentUser(c)
	if roleInfo.Status != nil && *roleInfo.Status == models.SysRoleStatusDisabled && roleId == user.RoleId {
		response.FailWithMsg("不能禁用自己所在的角色")
		return
	}

	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateRoleById(roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新角色的权限菜单
func UpdateRoleMenusById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateIncrementalIdsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.FailWithMsg("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
			return
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.FailWithMsg("无法更改自己的权限, 如需更改请联系上级领导")
			return
		}
	}

	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateRoleMenusById(user.Role, roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新角色的权限接口
func UpdateRoleApisById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateIncrementalIdsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}

	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.FailWithMsg("无法移除超级管理员的权限, 如有疑问请联系网站开发者")
			return
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.FailWithMsg("无法更改自己的权限, 如需更改请联系上级领导")
			return
		}
	}

	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateRoleApisById(roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	user := GetCurrentUser(c)
	if utils.ContainsUint(req.GetUintIds(), user.RoleId) {
		response.FailWithMsg("不能删除自己所在的角色")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteRoleByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
