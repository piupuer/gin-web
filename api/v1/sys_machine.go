package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindMachine(c *gin.Context) {
	var r request.MachineReq
	req.ShouldBind(c, &r)
	s := cache_service.New(c)
	list, err := s.FindMachine(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.MachineResp{}, r.Page)
}

func CreateMachine(c *gin.Context) {
	var r request.CreateMachineReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.Q.Create(r, new(ms.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}

func UpdateMachineById(c *gin.Context) {
	var r request.UpdateMachineReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(ms.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}

func ConnectMachineById(c *gin.Context) {
	id := req.UintId(c)
	s := service.New(c)
	err := s.ConnectMachine(id)
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteMachineByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.Uints(), new(ms.SysMachine))
	resp.CheckErr(err)
	resp.Success()
}
