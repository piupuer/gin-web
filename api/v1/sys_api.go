package v1

import (
	"github.com/gin-gonic/gin"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
)

// 获取接口列表
func GetApis(c *gin.Context) {
	// 绑定参数
	var req request.ApiListRequestStruct
	_ = c.Bind(&req)
	apis, err := service.GetApis(&req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.ApiListResponseStruct
	utils.Struct2StructByJson(apis, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(c, resp)
}

// 创建接口
func CreateApi(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateApiRequestStruct
	_ = c.Bind(&req)
	// 参数校验
	err := global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	err = service.CreateApi(&req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 更新接口
func UpdateApiById(c *gin.Context) {
	// 绑定参数, 这里与创建接口用同一结构体即可
	var req request.CreateApiRequestStruct
	_ = c.Bind(&req)
	// 获取path中的apiId
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		response.FailWithMsg(c, "接口编号不正确")
		return
	}
	// 更新数据
	err := service.UpdateApiById(uint(apiId), &req)
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}

// 批量删除接口
func BatchDeleteApiByIds(c *gin.Context) {
	var req request.Req
	_ = c.Bind(&req)
	// 删除数据
	err := service.DeleteApiByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	response.Success(c)
}
