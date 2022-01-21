package v1

import (
	"context"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
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
	var r req.Ids
	req.ShouldBind(c, &r)
	my := service.New(c)
	user := GetCurrentUser(c)
	err := my.DeleteLeaveByIds(r.Uints(), user)
	resp.CheckErr(err)
	resp.Success()
}

func LeaveTransition(ctx context.Context, logs ...resp.FsmApprovalLog) error {
	my := service.New(ctx)
	return my.LeaveTransition(logs...)
}

func GetLeaveFsmDetail(c *gin.Context, detail req.FsmSubmitterDetail) []resp.FsmSubmitterDetail {
	my := service.New(c)
	return my.GetLeaveFsmDetail(detail)
}

func UpdateLeaveFsmDetail(c *gin.Context, detail req.UpdateFsmSubmitterDetail) error {
	my := service.New(c)
	return my.UpdateLeaveFsmDetail(detail)
}
