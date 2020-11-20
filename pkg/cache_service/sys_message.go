package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

// 查询指定用户未删除的消息
func (s *RedisService) GetUnDeleteMessages(req *request.MessageListRequestStruct) ([]response.MessageListResponseStruct, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetUnDeleteMessages(req)
	}
	// 取出当前用户的所有消息
	currentUserAllLogs := make([]models.SysMessageLog, 0)
	err := s.redis.
		Table(new(models.SysMessageLog).TableName()).
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
	query := s.redis.FromString(utils.Struct2Json(currentUserAllLogs))

	title := strings.TrimSpace(req.Title)
	if title != "" {
		query = query.Where("message.title", "contains", title)
	}
	content := strings.TrimSpace(req.Content)
	if content != "" {
		query = query.Where("message.content", "contains", content)
	}
	if req.Type != nil {
		query = query.Where("message.type", "=", *req.Type)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	err = query.Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Find(&messageLogs).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&messageLogs).Error
		}
	}
	// 将数据转为响应格式
	list := make([]response.MessageListResponseStruct, 0)

	for _, log := range messageLogs {
		res := response.MessageListResponseStruct{
			Id:           log.Id,
			Status:       log.Status,
			ToUserId:     log.ToUserId,
			ToUsername:   log.ToUser.Username,
			Type:         log.Message.Type,
			Title:        log.Message.Title,
			Content:      log.Message.Content,
			CreatedAt:    log.Message.CreatedAt,
			FromUserId:   log.Message.FromUserId,
			FromUsername: log.Message.FromUser.Username,
		}
		list = append(list, res)
	}

	return list, err
}

// 查询未读消息条数
func (s *RedisService) GetUnReadMessageCount(userId uint) (int64, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetUnReadMessageCount(userId)
	}
	var total int64
	err := s.redis.
		Table(new(models.SysMessageLog).TableName()).
		Where("to_user_id", "=", userId).
		Where("status", "=", models.SysMessageLogStatusUnRead).
		Count(&total).Error
	return total, err
}
