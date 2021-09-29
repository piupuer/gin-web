package v1

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	var req request.OperationLogReq
	request.ShouldBind(c, &req)
	s := service.New(c)
	operationLogs, err := s.GetOperationLogs(&req)
	response.CheckErr(err)
	// 隐藏部分字段
	var respStruct []response.OperationLogResp
	utils.Struct2StructByJson(operationLogs, &respStruct)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 批量删除操作日志
func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.System.OperationLogAllowedToDelete {
		response.CheckErr("日志删除功能已被管理员关闭")
	}
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.DeleteByIds(req.GetUintIds(), new(models.SysOperationLog))
	response.CheckErr(err)
	response.Success()
}
