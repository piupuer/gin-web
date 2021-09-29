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

// 获取机器列表
func GetMachines(c *gin.Context) {
	var req request.MachineRequestStruct
	request.ShouldBind(c, &req)
	s := cache_service.New(c)
	machines, err := s.GetMachines(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.MachineListResponseStruct
	utils.Struct2StructByJson(machines, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建机器
func CreateMachine(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateMachineRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.Create(req, new(models.SysMachine))
	response.CheckErr(err)
	response.Success()
}

// 更新机器
func UpdateMachineById(c *gin.Context) {
	var req request.UpdateMachineRequestStruct
	request.ShouldBind(c, &req)
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		response.CheckErr("机器编号不正确")
	}

	s := service.New(c)
	err := s.UpdateById(machineId, req, new(models.SysMachine))
	response.CheckErr(err)
	response.Success()
}

// 重连机器
func ConnectMachineById(c *gin.Context) {
	// 获取path中的machineId
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		response.CheckErr("机器编号不正确")
		return
	}

	s := service.New(c)
	err := s.ConnectMachine(machineId)
	response.CheckErr(err)
	response.Success()
}

// 批量删除机器
func BatchDeleteMachineByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysMachine))
	response.CheckErr(err)
	response.Success()
}
