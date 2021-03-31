package v1

import (
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取数据字典列表
func GetDicts(c *gin.Context) {
	// 绑定参数
	var req request.DictRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	dicts, err := s.GetDicts(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.DictResponseStruct
	utils.Struct2StructByJson(dicts, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建数据字典
func CreateDict(c *gin.Context) {
	// 绑定参数
	var req request.CreateDictRequestStruct
	err := c.Bind(&req)
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
	// 创建服务
	s := service.New(c)
	err = s.CreateDict(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新数据字典
func UpdateDictById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateDictRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的dictId
	dictId := utils.Str2Uint(c.Param("dictId"))
	if dictId == 0 {
		response.FailWithMsg("数据字典编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateDictById(dictId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除数据字典
func BatchDeleteDictByIds(c *gin.Context) {
	var req request.Req
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteDictByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 获取数据字典数据列表
func GetDictDatas(c *gin.Context) {
	// 绑定参数
	var req request.DictDataRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	dictDatas, err := s.GetDictDatas(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.DictDataResponseStruct
	utils.Struct2StructByJson(dictDatas, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 创建数据字典数据
func CreateDictData(c *gin.Context) {
	// 绑定参数
	var req request.CreateDictDataRequestStruct
	err := c.Bind(&req)
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
	// 创建服务
	s := service.New(c)
	err = s.CreateDictData(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新数据字典数据
func UpdateDictDataById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateDictDataRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的dictDataId
	dictDataId := utils.Str2Uint(c.Param("dictDataId"))
	if dictDataId == 0 {
		response.FailWithMsg("数据字典数据编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateDictDataById(dictDataId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除数据字典数据
func BatchDeleteDictDataByIds(c *gin.Context) {
	var req request.Req
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteDictDataByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
