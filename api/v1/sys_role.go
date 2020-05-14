package v1

import (
	"fmt"
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
	_ = c.Bind(&req)
	// 创建服务
	s := service.New(c)
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
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
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
	var req gin.H
	_ = c.Bind(&req)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err := s.UpdateRoleById(roleId, req)
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
	err := c.Bind(&req)
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
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateRoleMenusById(roleId, req)
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
	err := c.Bind(&req)
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
	_ = c.Bind(&req)
	// 创建服务
	s := service.New(c)
	// 删除数据
	err := s.DeleteRoleByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
