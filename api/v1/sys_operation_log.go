package v1

import (
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindOperationLog(c *gin.Context) {
	var r request.OperationLogReq
	req.ShouldBind(c, &r)
	s := service.New(c)
	list, err := s.FindOperationLog(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.OperationLogResp{}, r.Page)
}

func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.Logs.OperationAllowedToDelete {
		resp.CheckErr("this feature has been turned off by the administrator")
	}
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.Uints(), new(ms.SysOperationLog))
	resp.CheckErr(err)
	resp.Success()
}
