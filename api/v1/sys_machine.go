package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取机器列表
func GetMachines(c *gin.Context) {
	// 绑定参数
	var req request.MachineListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := cache_service.New(c)
	machines, err := s.GetMachines(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.MachineListResponseStruct
	utils.Struct2StructByJson(machines, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建机器
func CreateMachine(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateMachineRequestStruct
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
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	err = s.CreateMachine(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新机器
func UpdateMachineById(c *gin.Context) {
	// 绑定参数
	var req models.SysMachine
	var machineInfo request.CreateMachineRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 将部分参数转为pwd, 如果值不为空, 可能会用到
	utils.Struct2StructByJson(req, &machineInfo)
	// 获取path中的machineId
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		response.FailWithMsg("机器编号不正确")
		return
	}

	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateMachineById(machineId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 重连机器
func ConnectMachineById(c *gin.Context) {
	// 获取path中的machineId
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		response.FailWithMsg("机器编号不正确")
		return
	}

	// 创建服务
	s := service.New(c)
	// 更新数据
	err := s.ConnectMachine(machineId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除机器
func BatchDeleteMachineByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteMachineByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
