package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := CacheGetMenuTree(c, user.Id)
	if ok {
		resp.SuccessWithData(oldCache)
		return
	}

	s := service.New(c)
	menus, err := s.GetMenuTree(user.RoleId)
	resp.CheckErr(err)
	var rp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &rp)
	CacheSetMenuTree(c, user.Id, rp)
	resp.SuccessWithData(rp)
}

func FindMenuByRoleId(c *gin.Context) {
	id := req.UintId(c)
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	menus, ids, err := s.FindMenuByRoleId(user.Role, id)
	resp.CheckErr(err)
	var rp response.MenuTreeWithAccessResp
	rp.AccessIds = ids
	utils.Struct2StructByJson(menus, &rp.List)
	resp.SuccessWithData(rp)
}

func FindMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	menus := s.FindMenu(user.Role)
	var rp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &rp)
	resp.SuccessWithData(rp)
}

func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateMenuReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	s := service.New(c)
	err := s.CreateMenu(user.Role, &r)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateMenuById(c *gin.Context) {
	var r request.UpdateMenuReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(ms.SysMenu))
	resp.CheckErr(err)
	resp.Success()
}

func BatchDeleteMenuByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.Uints(), new(ms.SysMenu))
	resp.CheckErr(err)
	resp.Success()
}
