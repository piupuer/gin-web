package v1

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
)

// @Tags GetMenuTree
// @Summary 查询当前用户菜单树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body true "查询当前用户菜单树"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/info [post]
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
