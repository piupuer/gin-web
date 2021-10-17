package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取接口列表
func GetApis(c *gin.Context) {
	var r request.ApiReq
	req.ShouldBind(c, &r)

	s := cache_service.New(c)
	list, err := s.GetApis(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.ApiResp{}, r.Page)
}

// 查询指定角色的接口(以分类分组)
func GetAllApiGroupByCategoryByRoleId(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	apis, ids, err := s.GetAllApiGroupByCategoryByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	resp.CheckErr(err)
	var rp response.ApiTreeWithAccessResp
	rp.AccessIds = ids
	utils.Struct2StructByJson(apis, &rp.List)
	resp.SuccessWithData(rp)
}

// 创建接口
func CreateApi(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateApiReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())

	// 记录当前创建人信息
	r.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.CreateApi(&r)
	resp.CheckErr(err)
	resp.Success()
}

// 更新接口
func UpdateApiById(c *gin.Context) {
	var r request.UpdateApiReq
	req.ShouldBind(c, &r)
	// 获取path中的apiId
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		resp.CheckErr("接口编号不正确")
	}
	s := service.New(c)
	err := s.UpdateApiById(apiId, r)
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除接口
func BatchDeleteApiByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.DeleteApiByIds(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}
