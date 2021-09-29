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
	var req request.MessageRequestStruct
	request.ShouldBind(c, &req)
	user := GetCurrentUser(c)
	req.ToUserId = user.Id
	s := cache_service.New(c)

	messages, err := s.GetUnDeleteMessages(&req)
	response.CheckErr(err)
	// 返回分页数据
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = messages
	response.SuccessWithData(resp)
}

// 未读消息条数
func GetUnReadMessageCount(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)

	total, err := s.GetUnReadMessageCount(user.Id)
	response.CheckErr(err)
	response.SuccessWithData(total)
}

// 推送消息
func PushMessage(c *gin.Context) {
	var req request.PushMessageRequestStruct
	request.ShouldBind(c, &req)
	user := GetCurrentUser(c)
	s := service.New(c)
	req.FromUserId = user.Id
	err := s.CreateMessage(&req)
	response.CheckErr(err)
	response.Success()
}

// 批量更新为已读
func BatchUpdateMessageRead(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.BatchUpdateMessageRead(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}

// 批量更新为删除
func BatchUpdateMessageDeleted(c *gin.Context) {
	var req request.Req
	request.ShouldBind(c, &req)
	s := service.New(c)
	err := s.BatchUpdateMessageDeleted(req.GetUintIds())
	response.CheckErr(err)
	response.Success()
}

// 全部更新为已读
func UpdateAllMessageRead(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageRead(user.Id)
	response.CheckErr(err)
	response.Success()
}

// 全部更新为删除
func UpdateAllMessageDeleted(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageDeleted(user.Id)
	response.CheckErr(err)
	response.Success()
}
