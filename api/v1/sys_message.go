package v1

import (
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FindMessage(c *gin.Context) {
	var r request.MessageReq
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	r.ToUserId = user.Id
	s := cache_service.New(c)

	list, err := s.FindUnDeleteMessage(&r)
	resp.CheckErr(err)
	resp.SuccessWithPageData(list, []response.MessageResp{}, r.Page)
}

func GetUnReadMessageCount(c *gin.Context) {
	user := GetCurrentUser(c)
	s := cache_service.New(c)

	total, err := s.GetUnReadMessageCount(user.Id)
	resp.CheckErr(err)
	resp.SuccessWithData(total)
}

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

func BatchUpdateMessageRead(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.BatchUpdateMessageRead(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}

func BatchUpdateMessageDeleted(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	s := service.New(c)
	err := s.BatchUpdateMessageDeleted(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}

func UpdateAllMessageRead(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageRead(user.Id)
	resp.CheckErr(err)
	resp.Success()
}

func UpdateAllMessageDeleted(c *gin.Context) {
	s := service.New(c)
	user := GetCurrentUser(c)
	err := s.UpdateAllMessageDeleted(user.Id)
	resp.CheckErr(err)
	resp.Success()
}
