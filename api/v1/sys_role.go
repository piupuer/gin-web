package v1

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
)

func FindRole(c *gin.Context) {
	var r request.RoleReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	// bind current user role sort(low level cannot view high level)
	r.CurrentRoleSort = *user.Role.Sort

	s := service.New(c)
	list := s.FindRole(&r)
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

func FindRoleKeywordByRoleIds(c *gin.Context, roleIds []uint) []string {
	s := cache_service.New(c)
	roles := s.FindRoleByIds(roleIds)
	keywords := make([]string, 0)
	for _, role := range roles {
		keywords = append(keywords, role.Keyword)
	}
	return keywords
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
