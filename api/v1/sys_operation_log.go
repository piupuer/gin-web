package v1

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func GetOperationLogs(c *gin.Context) {
	var r request.OperationLogReq
	req.ShouldBind(c, &r)
	s := service.New(c)
	list, err := s.GetOperationLogs(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.OperationLogResp{}, r.Page)
}

func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.Logs.OperationAllowedToDelete {
		resp.CheckErr("this feature has been turned off by the administrator")
	}
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.GetUintIds(), new(models.SysOperationLog))
	resp.CheckErr(err)
	resp.Success()
}
