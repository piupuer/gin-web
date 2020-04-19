package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
)

// @Tags SysUser
// @Summary 获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUsers [post]
func GetUsers(c *gin.Context) {
	user, _ := c.Get("user")
	global.Log.Debug(fmt.Sprintf("当前登录用户: %v", user))
	users, err := service.GetUsers()
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithData(c, users)
}

