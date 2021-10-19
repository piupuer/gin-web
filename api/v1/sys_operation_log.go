package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	var r request.OperationLogReq
	req.ShouldBind(c, &r)
	s := cache_service.New(c)
	list, err := s.GetOperationLogs(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.OperationLogResp{}, r.Page)
}

// 批量删除操作日志
func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.System.OperationLogAllowedToDelete {
		resp.CheckErr("日志删除功能已被管理员关闭")
	}
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.Q.DeleteByIds(r.GetUintIds(), new(models.SysOperationLog))
	resp.CheckErr(err)
	resp.Success()
}
