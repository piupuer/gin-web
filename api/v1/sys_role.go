package v1

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req request.RoleListRequestStruct
	_ = c.Bind(&req)
	roles, err := service.GetRoles(&req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
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
	response.SuccessWithData(c, resp)
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
		response.FailWithMsg(c, err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	err = service.CreateRole(&req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	// 绑定参数, 这里与创建角色用同一结构体即可
	var req request.CreateRoleRequestStruct
	_ = c.Bind(&req)
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg(c, "角色编号不正确")
		return
	}
	// 更新数据
	err := service.UpdateRoleById(roleId, &req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request.Req
	_ = c.Bind(&req)
	// 删除数据
	err := service.DeleteRoleByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}
