package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取数据字典列表
func GetDicts(c *gin.Context) {
	var r request.DictReq
	req.ShouldBind(c, &r)

	s := service.New(c)
	list, err := s.GetDicts(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.DictResp{}, r.Page)
}

// 创建数据字典
func CreateDict(c *gin.Context) {
	var r request.CreateDictReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.CreateDict(&r)
	resp.CheckErr(err)
	resp.Success()
}

// 更新数据字典
func UpdateDictById(c *gin.Context) {
	var r request.UpdateDictReq
	req.ShouldBind(c, &r)

	// 获取path中的dictId
	dictId := utils.Str2Uint(c.Param("dictId"))
	if dictId == 0 {
		resp.CheckErr("数据字典编号不正确")
	}
	s := service.New(c)
	err := s.UpdateDictById(dictId, r)
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除数据字典
func BatchDeleteDictByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)

	s := service.New(c)
	err := s.DeleteDictByIds(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}

// 获取数据字典数据列表
func GetDictDatas(c *gin.Context) {
	var r request.DictDataReq
	req.ShouldBind(c, &r)

	s := service.New(c)
	list, err := s.GetDictDatas(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.DictDataResp{}, r.Page)
}

// 创建数据字典数据
func CreateDictData(c *gin.Context) {
	var r request.CreateDictDataReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.CreateDictData(&r)
	resp.CheckErr(err)
	resp.Success()
}

// 更新数据字典数据
func UpdateDictDataById(c *gin.Context) {
	var r request.UpdateDictDataReq
	req.ShouldBind(c, &r)

	// 获取path中的dictDataId
	dictDataId := utils.Str2Uint(c.Param("dictDataId"))
	if dictDataId == 0 {
		resp.CheckErr("数据字典数据编号不正确")
	}
	s := service.New(c)
	err := s.UpdateDictDataById(dictDataId, r)
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除数据字典数据
func BatchDeleteDictDataByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)

	s := service.New(c)
	err := s.DeleteDictDataByIds(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}
