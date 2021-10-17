package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"time"
)

var (
	// 定期缓存, 避免每次频繁查询数据库
	menuTreeCache = cache.New(24*time.Hour, 48*time.Hour)
)

// 查询当前用户菜单树
func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := menuTreeCache.Get(fmt.Sprintf("%d", user.Id))
	if ok {
		rp, _ := oldCache.([]response.MenuTreeResp)
		resp.SuccessWithData(rp)
		return
	}

	s := service.New(c)
	menus, err := s.GetMenuTree(user.RoleId)
	resp.CheckErr(err)
	var rp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &rp)
	// 写入缓存
	menuTreeCache.Set(fmt.Sprintf("%d", user.Id), rp, cache.DefaultExpiration)
	resp.SuccessWithData(rp)
}

// 查询指定角色的菜单树
func GetAllMenuByRoleId(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	menus, ids, err := s.GetAllMenuByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	resp.CheckErr(err)
	var rp response.MenuTreeWithAccessResp
	rp.AccessIds = ids
	utils.Struct2StructByJson(menus, &rp.List)
	resp.SuccessWithData(rp)
}

// 查询所有菜单
func GetMenus(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	s := cache_service.New(c)
	menus := s.GetMenus(user.Role)
	var rp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &rp)
	resp.SuccessWithData(rp)
}

// 创建菜单
func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateMenuReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	// 记录当前创建人信息
	r.Creator = user.Nickname + user.Username
	s := service.New(c)
	err := s.CreateMenu(user.Role, &r)
	resp.CheckErr(err)
	resp.Success()
}

// 更新菜单
func UpdateMenuById(c *gin.Context) {
	var r request.UpdateMenuReq
	req.ShouldBind(c, &r)
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		resp.CheckErr("菜单编号不正确")
	}
	s := service.New(c)
	err := s.UpdateById(menuId, r, new(models.SysMenu))
	resp.CheckErr(err)
	resp.Success()
}

// 批量删除菜单
func BatchDeleteMenuByIds(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.DeleteByIds(r.GetUintIds(), new(models.SysMenu))
	resp.CheckErr(err)
	resp.Success()
}
