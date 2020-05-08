package v1

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
)

// 查询当前用户菜单树
func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)
	var u models.SysUser
	// GetCurrentUser可能不是最新的数据, 查询最新的用户role
	err := global.Mysql.Where("id = ?", user.Id).First(&u).Error
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	menus, err := service.GetMenuTree(u.RoleId)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(c, resp)
}

// 查询指定角色的菜单树
func GetAllMenuByRoleId(c *gin.Context) {
	// 绑定参数
	menus, ids, err := service.GetAllMenuByRoleId(utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	var resp response.MenuTreeWithAccessResponseStruct
	resp.AccessIds = ids
	utils.Struct2StructByJson(menus, &resp.List)
	response.SuccessWithData(c, resp)
}

// 查询所有菜单
func GetMenus(c *gin.Context) {
	menus, err := service.GetMenus()
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(c, resp)
}

// 创建菜单
func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateMenuRequestStruct
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	err = service.CreateMenu(&req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 更新菜单
func UpdateMenuById(c *gin.Context) {
	// 绑定参数
	var req request.CreateMenuRequestStruct
	_ = c.Bind(&req)
	// 获取path中的menuId
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		response.FailWithMsg(c, "菜单编号不正确")
		return
	}
	// 更新数据
	err := service.UpdateMenuById(menuId, &req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 批量删除菜单
func BatchDeleteMenuByIds(c *gin.Context) {
	var req request.Req
	_ = c.Bind(&req)
	// 删除数据
	err := service.DeleteMenuByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}
