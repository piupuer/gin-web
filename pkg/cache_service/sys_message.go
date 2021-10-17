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

// 查询指定用户未删除的消息
func (s RedisService) GetUnDeleteMessages(req *request.MessageReq) ([]response.MessageResp, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetUnDeleteMessages(req)
	}
	// 取出当前用户的所有消息
	currentUserAllLogs := make([]models.SysMessageLog, 0)
	err := s.redis.
		Table("sys_message_log").
		Preload("Message").
		Preload("Message.FromUser").
		Preload("ToUser").
		Where("to_user_id", "=", req.ToUserId).
		// 未删除的
		Where("status", "!=", models.SysMessageLogStatusDeleted).
		Find(&currentUserAllLogs).Error
	if err != nil {
		return nil, err
	}

	messageLogs := make([]models.SysMessageLog, 0)
	// 转为json, 再匹配前端其他条件
	query := s.redis.
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
	// 查询列表
	err = s.Find(query, &req.Page, &messageLogs)
	if err != nil {
		return nil, err
	}
	// 将数据转为响应格式
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

// 查询未读消息条数
func (s RedisService) GetUnReadMessageCount(userId uint) (int64, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetUnReadMessageCount(userId)
	}
	var total int64
	err := s.redis.
		Table("sys_message_log").
		Where("to_user_id", "=", userId).
		Where("status", "=", models.SysMessageLogStatusUnRead).
		Count(&total).Error
	return total, err
}
