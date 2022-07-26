package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
)

// FindLeave
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Leave
// @Description FindLeave
// @Param params query request.Leave true "params"
// @Router /leave/list [GET]
func FindLeave(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "FindLeave"))
	defer span.End()
	var r request.Leave
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.UserId = user.Id
	my := service.New(c)
	list := my.FindLeave(&r)
	resp.SuccessWithPageData(list, &[]response.Leave{}, r.Page)
}

// CreateLeave
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Leave
// @Description CreateLeave
// @Param params body request.CreateLeave true "params"
// @Router /leave/create [POST]
func CreateLeave(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "CreateLeave"))
	defer span.End()
	user := GetCurrentUser(c)
	var r request.CreateLeave
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())
	r.User = user
	my := service.New(c)
	err := my.CreateLeave(&r)
	resp.CheckErr(err)
	resp.Success()
}

// UpdateLeaveById
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Leave
// @Description UpdateLeaveById
// @Param id path uint true "id"
// @Param params body request.UpdateLeave true "params"
// @Router /leave/update/{id} [PATCH]
func UpdateLeaveById(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "UpdateLeaveById"))
	defer span.End()
	var r request.UpdateLeave
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	my := service.New(c)
	user := GetCurrentUser(c)
	err := my.UpdateLeaveById(id, r, user)
	resp.CheckErr(err)
	resp.Success()
}

// BatchDeleteLeaveByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Leave
// @Description BatchDeleteLeaveByIds
// @Param ids body req.Ids true "ids"
// @Router /leave/delete/batch [DELETE]
func BatchDeleteLeaveByIds(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "BatchDeleteLeaveByIds"))
	defer span.End()
	var r req.Ids
	req.ShouldBind(c, &r)
	my := service.New(c)
	user := GetCurrentUser(c)
	err := my.DeleteLeaveByIds(r.Uints(), user)
	resp.CheckErr(err)
	resp.Success()
}
