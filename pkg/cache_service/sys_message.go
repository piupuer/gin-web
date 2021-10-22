package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/pkg/resp"
	"strings"
)

// find status!=models.SysMessageLogStatusDeleted messages
func (rd RedisService) FindUnDeleteMessage(req *request.MessageReq) ([]response.MessageResp, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindUnDeleteMessage(req)
	}
	currentUserAllLogs := make([]models.SysMessageLog, 0)
	err := rd.Q.
		Table("sys_message_log").
		Preload("Message").
		Preload("Message.FromUser").
		Preload("ToUser").
		Where("to_user_id", "=", req.ToUserId).
		// un delete
		Where("status", "!=", models.SysMessageLogStatusDeleted).
		Find(&currentUserAllLogs).Error
	if err != nil {
		return nil, err
	}

	messageLogs := make([]models.SysMessageLog, 0)
	// all log json
	query := rd.Q.
		FromString(utils.Struct2Json(currentUserAllLogs)).
		Order("created_at DESC")
	title := strings.TrimSpace(req.Title)
	if title != "" {
		query = query.Where("message.title", "contains", title)
	}
	content := strings.TrimSpace(req.Content)
	if content != "" {
		query = query.Where("message.content", "contains", content)
	}
	if req.Type != nil {
		query = query.Where("type", "=", *req.Type)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	err = rd.Q.FindWithPage(query, &req.Page, &messageLogs)
	if err != nil {
		return nil, err
	}
	// convert to MessageResp
	list := make([]response.MessageResp, 0)
	for _, log := range messageLogs {
		res := response.MessageResp{
			Base: resp.Base{
				Id:        log.Id,
				CreatedAt: log.CreatedAt,
				UpdatedAt: log.UpdatedAt,
			},
			Status:       log.Status,
			ToUserId:     log.ToUserId,
			ToUsername:   log.ToUser.Username,
			Type:         log.Message.Type,
			Title:        log.Message.Title,
			Content:      log.Message.Content,
			FromUserId:   log.Message.FromUserId,
			FromUsername: log.Message.FromUser.Username,
		}
		list = append(list, res)
	}

	return list, err
}

// un read total count
func (rd RedisService) GetUnReadMessageCount(userId uint) (int64, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.GetUnReadMessageCount(userId)
	}
	var total int64
	err := rd.Q.
		Table("sys_message_log").
		Where("to_user_id", "=", userId).
		Where("status", "=", models.SysMessageLogStatusUnRead).
		Count(&total).Error
	return total, err
}
