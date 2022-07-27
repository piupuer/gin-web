package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
	"strings"
)

// GetUserInfo
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description GetUserInfo
// @Router /user/info [GET]
func GetUserInfo(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "GetUserInfo"))
	defer span.End()
	user := GetCurrentUser(c)
	oldCache, ok := CacheGetUserInfo(ctx, user.Id)
	if ok {
		resp.SuccessWithData(oldCache)
		return
	}

	var rp response.UserInfo
	utils.Struct2StructByJson(user, &rp)
	rp.Roles = []string{
		"admin",
	}
	rp.Keyword = user.Role.Keyword
	if user.Role.Sort != nil {
		rp.RoleSort = *user.Role.Sort
	} else {
		rp.RoleSort = 9999
	}
	CacheSetUserInfo(ctx, user.Id, rp)
	resp.SuccessWithData(rp)
}

// FindUser
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description FindUser
// @Param params query request.User true "params"
// @Router /user/list [GET]
func FindUser(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "FindUser"))
	defer span.End()
	var r request.User
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.CurrentRole = user.Role
	cs := cache_service.New(c)
	list := cs.FindUser(&r)
	resp.SuccessWithPageData(list, &[]response.User{}, r.Page)
}

// ChangePwd
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description ChangePwd
// @Param params body request.ChangePwd true "params"
// @Router /user/changePwd [PUT]
func ChangePwd(c *gin.Context) {
	var r request.ChangePwd
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	query := global.Mysql.Where("username = ?", user.Username).First(&user)
	err := query.Error
	resp.CheckErr(err)
	if ok := utils.ComparePwd(r.OldPassword, user.Password); !ok {
		resp.CheckErr("the original password is incorrect")
	}
	err = query.Update("password", utils.GenPwd(r.NewPassword)).Error
	resp.CheckErr(err)
	resp.Success()
}

func GetCurrentUser(c *gin.Context) models.SysUser {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "GetCurrentUser"))
	defer span.End()
	userId, exists := c.Get(constant.MiddlewareJwtUserCtxKey)
	var newUser models.SysUser
	if !exists {
		return newUser
	}
	uid := utils.Str2Uint(fmt.Sprintf("%d", userId))
	oldCache, ok := CacheGetUser(ctx, uid)
	if ok {
		return oldCache
	}
	my := service.New(c)
	newUser = my.GetUserById(uid)
	if newUser.Id > constant.Zero {
		// if user id exists, set to cache
		CacheSetUser(c, uid, newUser)
	}
	return newUser
}

func GetCurrentUserAndRole(c *gin.Context) ms.User {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "GetCurrentUserAndRole"))
	defer span.End()
	user := GetCurrentUser(c)
	cs := cache_service.New(c)
	var roleSort uint
	if user.Role.Sort != nil {
		roleSort = *user.Role.Sort
	}
	pathRoleId, err := req.UintIdWithErr(c)
	pathRoleKeyword := ""
	if err == nil {
		role := cs.GetRoleById(pathRoleId)
		pathRoleKeyword = role.Keyword
	}
	var u ms.User
	copier.Copy(&u, user)
	u.RoleId = user.RoleId
	u.RoleName = user.Role.Name
	u.RoleSort = roleSort
	u.RoleKeyword = user.Role.Keyword
	u.PathRoleId = pathRoleId
	u.PathRoleKeyword = pathRoleKeyword
	return u
}

func GetUserLoginStatus(c *gin.Context, r *req.UserStatus) (err error) {
	my := service.New(c)
	var u models.SysUser
	u, err = my.GetUserByUsername(r.Username)
	if err != nil {
		err = nil
		return
	}
	r.Locked = u.Locked
	r.LockExpire = u.LockExpire
	r.Wrong = u.Wrong
	return
}

// FindUserByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description FindUserByIds
// @Param ids path string true "ids"
// @Router /user/list/{ids} [GET]
func FindUserByIds(c *gin.Context) {
	ids := req.UintIds(c)
	cs := cache_service.New(c)
	list := cs.FindUserByIds(ids)
	resp.SuccessWithData(list)
}

func RouterFindUserByIds(c *gin.Context, userIds []uint) []ms.User {
	users := make([]models.SysUser, 0)
	global.Mysql.
		Model(&models.SysUser{}).
		Where("id IN (?)", userIds).
		Find(&users)
	newUsers := make([]ms.User, 0)
	utils.Struct2StructByJson(users, &newUsers)
	return newUsers
}

func RouterFindRoleByIds(c *gin.Context, roleIds []uint) []ms.Role {
	roles := make([]models.SysRole, 0)
	global.Mysql.
		Model(&models.SysRole{}).
		Where("id IN (?)", roleIds).
		Find(&roles)
	newRoles := make([]ms.Role, 0)
	utils.Struct2StructByJson(roles, &newRoles)
	return newRoles
}

// CreateUser
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description CreateUser
// @Param params body request.CreateUser true "params"
// @Router /user/create [POST]
func CreateUser(c *gin.Context) {
	var r request.CreateUser
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	my := service.New(c)
	// plaintext to ciphertext
	r.Password = utils.GenPwd(r.InitPassword)
	err := my.Q.Create(r, new(models.SysUser))
	resp.CheckErr(err)
	resp.Success()
}

// UpdateUserById
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description UpdateUserById
// @Param id path uint true "id"
// @Param params body request.UpdateUser true "params"
// @Router /user/update/{id} [PATCH]
func UpdateUserById(c *gin.Context) {
	var r request.UpdateUser
	req.ShouldBind(c, &r)
	id := req.UintId(c)

	// new password is not empty, update password
	if r.NewPassword != nil && strings.TrimSpace(*r.NewPassword) != "" {
		password := utils.GenPwd(*r.NewPassword)
		r.Password = &password
	}

	user := GetCurrentUser(c)
	if id == user.Id {
		if r.Status != nil && uint(*r.Status) == models.SysUserStatusDisabled {
			resp.CheckErr("cannot disable yourself")
		}
		if r.RoleId != nil && user.RoleId != *r.RoleId {
			if *user.Role.Sort != models.SysRoleSuperAdminSort {
				resp.CheckErr("cannot change your role")
			} else {
				resp.CheckErr("cannot change super admin's role")
			}
		}
	}

	if r.Locked != nil && *r.Locked == req.NullUint(constant.One) {
		var i int64 = 0
		var j = 0
		r.LockExpire = &i
		r.Wrong = &j
	}

	my := service.New(c)
	err := my.Q.UpdateById(id, r, new(models.SysUser))
	resp.CheckErr(err)
	CacheDeleteUserInfo(c, user.Id)
	CacheDeleteUser(c, user.Id)
	resp.Success()
}

// BatchDeleteUserByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags User
// @Description BatchDeleteUserByIds
// @Param ids body req.Ids true "ids"
// @Router /user/delete/batch [DELETE]
func BatchDeleteUserByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.Uints(), user.Id) {
		resp.CheckErr("cannot remove yourself")
	}

	my := service.New(c)
	err := my.Q.DeleteByIds(r.Uints(), new(models.SysUser))
	resp.CheckErr(err)
	CacheFlushUserInfo(c)
	CacheFlushUser(c)
	resp.Success()
}
