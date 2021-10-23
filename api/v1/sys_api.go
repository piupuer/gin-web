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

func FindApi(c *gin.Context) {
	var r request.ApiReq
	req.ShouldBind(c, &r)

	s := cache_service.New(c)
	list := s.FindApi(&r)
	resp.SuccessWithPageData(list, []response.ApiResp{}, r.Page)
}

func FindAllApiGroupByCategoryByRoleId(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	apis, ids, err := s.FindAllApiGroupByCategoryByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	resp.CheckErr(err)
	var rp response.ApiTreeWithAccessResp
	rp.AccessIds = ids
	utils.Struct2StructByJson(apis, &rp.List)
	resp.SuccessWithData(rp)
}

func CreateApi(c *gin.Context) {
	var r request.CreateApiReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())

	s := service.New(c)
	err := s.CreateApi(&r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateApiById(c *gin.Context) {
	var r request.UpdateApiReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.UpdateApiById(id, r)
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteApiByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.DeleteApiByIds(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}
