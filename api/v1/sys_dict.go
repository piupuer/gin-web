package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取数据字典列表
func GetDicts(c *gin.Context) {
	var req request.DictRequestStruct
	request.ShouldBind(c, &req)

	s := service.New(c)
	dicts, err := s.GetDicts(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.DictResponseStruct
	utils.Struct2StructByJson(dicts, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建数据字典
func CreateDict(c *gin.Context) {
	var req request.CreateDictRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	s := service.New(c)
	err := s.CreateDict(&req)
	response.CheckErr(err)
	response.Success()
}

// 更新数据字典
func UpdateDictById(c *gin.Context) {
	var req request.UpdateDictRequestStruct
	request.ShouldBind(c, &req)

	// 获取path中的dictId
	dictId := utils.Str2Uint(c.Param("dictId"))
	if dictId == 0 {
		response.CheckErr("数据字典编号不正确")
	}
	s := service.New(c)
	err := s.UpdateDictById(dictId, req)
	response.CheckErr(err)
	response.Success()
}

// 批量删除数据字典
func BatchDeleteDictByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)

	s := service.New(c)
	err := s.DeleteDictByIds(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}

// 获取数据字典数据列表
func GetDictDatas(c *gin.Context) {
	var req request.DictDataRequestStruct
	request.ShouldBind(c, &req)

	s := service.New(c)
	dictDatas, err := s.GetDictDatas(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.DictDataResponseStruct
	utils.Struct2StructByJson(dictDatas, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建数据字典数据
func CreateDictData(c *gin.Context) {
	var req request.CreateDictDataRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())
	s := service.New(c)
	err := s.CreateDictData(&req)
	response.CheckErr(err)
	response.Success()
}

// 更新数据字典数据
func UpdateDictDataById(c *gin.Context) {
	var req request.UpdateDictDataRequestStruct
	request.ShouldBind(c, &req)

	// 获取path中的dictDataId
	dictDataId := utils.Str2Uint(c.Param("dictDataId"))
	if dictDataId == 0 {
		response.CheckErr("数据字典数据编号不正确")
	}
	s := service.New(c)
	err := s.UpdateDictDataById(dictDataId, req)
	response.CheckErr(err)
	response.Success()
}

// 批量删除数据字典数据
func BatchDeleteDictDataByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)

	s := service.New(c)
	err := s.DeleteDictDataByIds(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}
