package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindRole(c *gin.Context) {
	var r request.RoleReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	// bind current user role sort(low level cannot view high level)
	r.CurrentRoleSort = *user.Role.Sort

	s := cache_service.New(c)
	list, err := s.FindRole(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.RoleResp{}, r.Page)
}

func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateRoleReq
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())

	if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
		resp.CheckErr("sort must >= %d", *user.Role.Sort)
	}

	s := service.New(c)
	err := s.Q.Create(r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

func UpdateRoleById(c *gin.Context) {
	var r request.UpdateRoleReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	if r.Sort != nil {
		user := GetCurrentUser(c)
		if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
			resp.CheckErr("sort must >= %d", *user.Role.Sort)
		}
	}

	user := GetCurrentUser(c)
	if r.Status != nil && uint(*r.Status) == models.SysRoleStatusDisabled && id == user.RoleId {
		resp.CheckErr("cannot disable your role")
	}

	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

func UpdateRoleMenusById(c *gin.Context) {
	var r request.UpdateMenuIncrementalIdsReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	user := GetCurrentUser(c)
	if user.RoleId == id {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(r.Delete) > 0 {
			resp.CheckErr("cannot remove super admin privileges")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			resp.CheckErr("cannot change your permissions")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleMenusById(user.Role, id, r)
	resp.CheckErr(err)
	CacheFlushMenuTree(c)
	resp.Success()
}

func UpdateRoleApisById(c *gin.Context) {
	var r request.UpdateMenuIncrementalIdsReq
	req.ShouldBind(c, &r)
	id := req.UintId(c)

	user := GetCurrentUser(c)

	if user.RoleId == id {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(r.Delete) > 0 {
			resp.CheckErr("cannot remove super admin privileges")
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			resp.CheckErr("cannot change your permissions")
		}
	}

	s := service.New(c)
	err := s.UpdateRoleApisById(id, r)
	resp.CheckErr(err)
	CacheFlushMenuTree(c)
	resp.Success()
}

func BatchDeleteRoleByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.Uints(), user.RoleId) {
		resp.CheckErr("cannot delete your role")
	}

	s := service.New(c)
	err := s.DeleteRoleByIds(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}
