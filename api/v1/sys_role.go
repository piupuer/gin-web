package v1

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
)

// @Tags SysRole
// @Summary 获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body true "分页获取角色列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/list [post]
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req request.RoleListRequestStruct
	_ = c.Bind(&req)
	menus, err := service.GetRoles(&req)
	if err != nil {
		response.Fail(c)
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.MenuListResponseStruct
	utils.Struct2StructByJson(menus, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(c, resp)
}
