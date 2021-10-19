package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取机器列表
func GetMachines(c *gin.Context) {
	var r request.MachineReq
	req.ShouldBind(c, &r)
	s := cache_service.New(c)
	list, err := s.GetMachines(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.MachineResp{}, r.Page)
}

// 创建机器
func CreateMachine(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateMachineReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	// 记录当前创建人信息
	r.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.Q.Create(r, new(models.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}

// 更新机器
func UpdateMachineById(c *gin.Context) {
	var r request.UpdateMachineReq
	req.ShouldBind(c, &r)
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		resp.CheckErr("机器编号不正确")
	}

	s := service.New(c)
	err := s.Q.UpdateById(machineId, r, new(models.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}

// 重连机器
func ConnectMachineById(c *gin.Context) {
	// 获取path中的machineId
	machineId := utils.Str2Uint(c.Param("machineId"))
	if machineId == 0 {
		resp.CheckErr("机器编号不正确")
		return
	}

	s := service.New(c)
	err := s.ConnectMachine(machineId)
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除机器
func BatchDeleteMachineByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.GetUintIds(), new(models.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}
