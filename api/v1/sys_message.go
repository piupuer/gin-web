package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
)

// 获取全部消息
func GetAllMessages(c *gin.Context) {
	// 绑定参数
	var req request.MessageListRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 设置接收人为当前登录用户
	user := GetCurrentUser(c)
	req.ToUserId = user.Id
	// 创建服务
	s := cache_service.New(c)

	messages, err := s.GetUnDeleteMessages(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.List = messages
	response.SuccessWithData(resp)
}

// 未读消息条数
func GetUnReadMessageCount(c *gin.Context) {
	user := GetCurrentUser(c)
	// 创建服务
	s := cache_service.New(c)

	total, err := s.GetUnReadMessageCount(user.Id)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(total)
}

// 批量更新为已读
func BatchUpdateMessageRead(c *gin.Context) {
	var req request.Req
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 更新
	err = s.BatchUpdateMessageRead(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量更新为删除
func BatchUpdateMessageDeleted(c *gin.Context) {
	var req request.Req
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 更新
	err = s.BatchUpdateMessageDeleted(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 全部更新为已读
func UpdateAllMessageRead(c *gin.Context) {
	// 创建服务
	s := service.New(c)
	user := GetCurrentUser(c)
	// 更新
	err := s.UpdateAllMessageRead(user.Id)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 全部更新为删除
func UpdateAllMessageDeleted(c *gin.Context) {
	// 创建服务
	s := service.New(c)
	user := GetCurrentUser(c)
	// 更新
	err := s.UpdateAllMessageDeleted(user.Id)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

