package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindDict(c *gin.Context) {
	var r request.DictReq
	req.ShouldBind(c, &r)

	s := service.New(c)
	list, err := s.FindDict(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.DictResp{}, r.Page)
}

func CreateDict(c *gin.Context) {
	var r request.CreateDictReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.CreateDict(&r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateDictById(c *gin.Context) {
	var r request.UpdateDictReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.UpdateDictById(id, r)
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteDictByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)

	s := service.New(c)
	err := s.DeleteDictByIds(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}

func FindDictData(c *gin.Context) {
	var r request.DictDataReq
	req.ShouldBind(c, &r)

	s := service.New(c)
	list, err := s.FindDictData(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.DictDataResp{}, r.Page)
}

func CreateDictData(c *gin.Context) {
	var r request.CreateDictDataReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.CreateDictData(&r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateDictDataById(c *gin.Context) {
	var r request.UpdateDictDataReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.UpdateDictDataById(id, r)
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteDictDataByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)

	s := service.New(c)
	err := s.DeleteDictDataByIds(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}
