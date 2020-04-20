package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
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
	user := GetCurrentUser(c)
	global.Log.Debug(fmt.Sprintf("当前登录用户: %v", user))
	users, err := service.GetUsers()
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithData(c, users)
}

// @Tags SysUser
// @Summary 修改密码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body true "是否修改成功"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/changePwd [put]
func ChangePwd(c *gin.Context) {
	var msg string
	// 请求json绑定
	var req request.ChangePwdStruct
	_ = c.ShouldBindJSON(&req)
	// 获取当前用户
	user := GetCurrentUser(c)
	query := global.Mysql.Where("username = ?", user.Username).First(&user)
	// 查询用户
	err := query.Error
	if err != nil {
		msg = err.Error()
	} else {
		// 校验密码
		if ok := utils.ComparePwd(req.OldPassword, user.Password); !ok {
			msg = "原密码错误"
		} else {
			// 更新密码
			err = query.Update("password", utils.GenPwd(req.NewPassword)).Error
			if err != nil {
				msg = err.Error()
			}
		}
	}
	if msg != "" {
		response.FailWithMsg(c, msg)
		return
	}
	response.Success(c)
}

// 获取当前请求用户信息
func GetCurrentUser(c *gin.Context) models.SysUser {
	user, _ := c.Get("user")
	u, _ := user.(models.SysUser)
	return u
}
