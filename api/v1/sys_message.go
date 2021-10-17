package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// 获取全部消息
func GetAllMessages(c *gin.Context) {
	var r request.MessageReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.ToUserId = user.Id
	s := cache_service.New(c)

	messages, err := s.GetUnDeleteMessages(&r)
	resp.CheckErr(err)
	// 返回分页数据
	var rp resp.PageData
	rp.Page = r.Page
	rp.List = messages
	resp.SuccessWithData(rp)
}

// 未读消息条数
func GetUnReadMessageCount(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)

	total, err := s.GetUnReadMessageCount(user.Id)
	resp.CheckErr(err)
	resp.SuccessWithData(total)
}

// 推送消息
func PushMessage(c *gin.Context) {
	var r request.PushMessageReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	s := service.New(c)
	r.FromUserId = user.Id
	err := s.CreateMessage(&r)
	resp.CheckErr(err)
	resp.Success()
}

// 批量更新为已读
func BatchUpdateMessageRead(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.BatchUpdateMessageRead(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}

// 批量更新为删除
func BatchUpdateMessageDeleted(c *gin.Context) {
	var r request.Req
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.BatchUpdateMessageDeleted(r.GetUintIds())
	resp.CheckErr(err)
	resp.Success()
}

// 全部更新为已读
func UpdateAllMessageRead(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageRead(user.Id)
	resp.CheckErr(err)
	resp.Success()
}

// 全部更新为删除
func UpdateAllMessageDeleted(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageDeleted(user.Id)
	resp.CheckErr(err)
	resp.Success()
}
