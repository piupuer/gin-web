package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	// 绑定参数
	var req request.OperationLogListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := cache_service.New(c)
	operationLogs, err := s.GetOperationLogs(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.OperationLogListResponseStruct
	utils.Struct2StructByJson(operationLogs, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// 批量删除操作日志
func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.System.OperationLogAllowedToDelete {
		response.FailWithMsg("日志删除功能已被管理员关闭")
		return
	}
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysOperationLog))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
