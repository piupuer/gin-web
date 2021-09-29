package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取接口列表
func GetApis(c *gin.Context) {
	var req request.ApiRequestStruct
	request.ShouldBind(c, &req)

	s := cache_service.New(c)
	apis, err := s.GetApis(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.ApiListResponseStruct
	utils.Struct2StructByJson(apis, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 查询指定角色的接口(以分类分组)
func GetAllApiGroupByCategoryByRoleId(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	apis, ids, err := s.GetAllApiGroupByCategoryByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	response.CheckErr(err)
	var resp response.ApiTreeWithAccessResponseStruct
	resp.AccessIds = ids
	utils.Struct2StructByJson(apis, &resp.List)
	response.SuccessWithData(resp)
}

// 创建接口
func CreateApi(c *gin.Context) {
	user := GetCurrentUser(c)
	var req request.CreateApiRequestStruct
	request.ShouldBind(c, &req)
	request.Validate(c, req, req.FieldTrans())

	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.CreateApi(&req)
	response.CheckErr(err)
	response.Success()
}

// 更新接口
func UpdateApiById(c *gin.Context) {
	var req request.UpdateApiRequestStruct
	request.ShouldBind(c, &req)
	// 获取path中的apiId
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		response.CheckErr("接口编号不正确")
	}
	s := service.New(c)
	err := s.UpdateApiById(apiId, req)
	response.CheckErr(err)
	response.Success()
}

// 批量删除接口
func BatchDeleteApiByIds(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteApiByIds(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}
