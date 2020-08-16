package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
	"time"
)

// 查询指定用户未删除的消息
func (s *MysqlService) GetUnDeleteMessages(req request.MessageListRequestStruct) ([]response.MessageListResponseStruct, error) {
	sysMessageLogTableName := new(models.SysMessageLog).TableName()
	sysMessageTableName := new(models.SysMessage).TableName()
	sysUserTableName := new(models.SysUser).TableName()
	list := make([]response.MessageListResponseStruct, 0)
	// 自定义查询字段
	fields := []string{
		fmt.Sprintf("%s.id AS id", sysMessageLogTableName),
		fmt.Sprintf("%s.to_user_id AS to_user_id", sysMessageLogTableName),
		fmt.Sprintf("toUser.username AS to_username"),
		fmt.Sprintf("%s.status AS status", sysMessageLogTableName),
		fmt.Sprintf("%s.type AS type", sysMessageTableName),
		fmt.Sprintf("%s.title AS title", sysMessageTableName),
		fmt.Sprintf("%s.content AS content", sysMessageTableName),
		fmt.Sprintf("%s.created_at AS created_at", sysMessageTableName),
		fmt.Sprintf("%s.from_user_id AS from_user_id", sysMessageTableName),
		fmt.Sprintf("fromUser.username AS from_username"),
	}
	query := s.tx.
		Table(sysMessageLogTableName).
		Select(fields).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.message_id = %s.id", sysMessageTableName, sysMessageLogTableName, sysMessageTableName)).
		Joins(fmt.Sprintf("LEFT JOIN %s AS toUser ON %s.to_user_id = toUser.id", sysUserTableName, sysMessageLogTableName)).
		Joins(fmt.Sprintf("LEFT JOIN %s AS fromUser ON %s.from_user_id = fromUser.id", sysUserTableName, sysMessageTableName))

	// 添加条件
	query = query.Where(fmt.Sprintf("%s.to_user_id = ?", sysMessageLogTableName), req.ToUserId)
	title := strings.TrimSpace(req.Title)
	if title != "" {
		query = query.Where(fmt.Sprintf("%s.title LIKE ?", sysMessageTableName), fmt.Sprintf("%%%s%%", title))
	}
	content := strings.TrimSpace(req.Title)
	if content != "" {
		query = query.Where(fmt.Sprintf("%s.content LIKE ?", sysMessageTableName), fmt.Sprintf("%%%s%%", content))
	}
	if req.Type != nil {
		query = query.Where(fmt.Sprintf("%s.type = ?", sysMessageTableName), *req.Type)
	}
	if req.Status != nil {
		query = query.Where(fmt.Sprintf("%s.status = ?", sysMessageLogTableName), *req.Status)
	} else {
		// 未删除的
		query = query.Where(fmt.Sprintf("%s.status != ?", sysMessageLogTableName), models.SysMessageLogStatusDeleted)
	}

	// 多表联合查询不用Find用Scan
	// 查询条数
	err := query.Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Scan(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Scan(&list).Error
		}
	}
	return list, err
}

// 查询未读消息条数
func (s *MysqlService) GetUnReadMessageCount(userId uint) (uint, error) {
	var total uint
	err := s.tx.
		Table(new(models.SysMessageLog).TableName()).
		Where("to_user_id = ?", userId).
		Where("status = ?", models.SysMessageLogStatusUnRead).
		Count(&total).Error
	return total, err
}

// 更新为已读
func (s *MysqlService) BatchUpdateMessageRead(messageLogIds []uint) error {
	return s.BatchUpdateMessageStatus(messageLogIds, models.SysMessageLogStatusRead)
}

// 更新为已删除
func (s *MysqlService) BatchUpdateMessageDeleted(messageLogIds []uint) error {
	return s.BatchUpdateMessageStatus(messageLogIds, models.SysMessageLogStatusDeleted)
}

// 全标已读
func (s *MysqlService) UpdateAllMessageRead(userId uint) error {
	return s.UpdateAllMessageStatus(userId, models.SysMessageLogStatusRead)
}

// 全标删除
func (s *MysqlService) UpdateAllMessageDeleted(userId uint) error {
	return s.UpdateAllMessageStatus(userId, models.SysMessageLogStatusDeleted)
}

// 批量更新消息状态
func (s *MysqlService) BatchUpdateMessageStatus(messageLogIds []uint, status uint) error {
	return s.tx.
		Table(new(models.SysMessageLog).TableName()).
		// 已删除的消息不再标记
		Where("status != ?", models.SysMessageLogStatusDeleted).
		Where("id IN (?)", messageLogIds).
		Update("status", status).Error
}

// 更新某个用户的全部消息状态
func (s *MysqlService) UpdateAllMessageStatus(userId uint, status uint) error {
	var log models.SysMessageLog
	log.ToUserId = userId
	return s.tx.
		Table(log.TableName()).
		// 已删除的消息不再标记
		Where("status != ?", models.SysMessageLogStatusDeleted).
		Where(&log).
		Update("status", status).Error
}

// 同步消息(某个用户登录时, 将消息关联到log)
func (s *MysqlService) SyncMessageByUserIds(userIds []uint) error {
	// 查询用户
	users := make([]models.SysUser, 0)
	err := s.tx.Where("id IN (?)", userIds).Find(&users).Error
	if err != nil {
		return err
	}
	for _, user := range users {
		messages := make([]models.SysMessage, 0)
		err = s.tx.
			// 用户注册时间早于消息创建时间
			Where("created_at > ?", user.CreatedAt).
			// 消息有效期晚于当前时间
			Where("expired_at > ?", time.Now()).
			// 一对多需要角色一致, 系统消息不需要
			Where("(type = ? AND role_id = ?) OR type = ?", models.SysMessageTypeOneToMany, user.RoleId, models.SysMessageTypeSystem).
			Find(&messages).Error
		if err != nil {
			return err
		}
		messageIds := make([]uint, 0)
		for _, message := range messages {
			messageIds = append(messageIds, message.Id)
		}
		// 判断消息是否已经关联
		logs := make([]models.SysMessageLog, 0)
		err := s.tx.
			Where("to_user_id = ?", user.Id).
			Where("message_id IN (?)", messageIds).
			Find(&logs).Error
		if err != nil {
			return err
		}
		// 已关联的旧消息
		oldMessageIds := make([]uint, 0)
		for _, log := range logs {
			if !utils.ContainsUint(oldMessageIds, log.MessageId) {
				oldMessageIds = append(oldMessageIds, log.MessageId)
			}
		}
		for _, messageId := range messageIds {
			// 当前消息id不在旧的列表中
			if !utils.ContainsUint(oldMessageIds, messageId) {
				// 新增消息log
				s.tx.Create(&models.SysMessageLog{
					ToUserId:  user.Id,
					MessageId: messageId,
				})
			}
		}
	}
	return nil
}

// 创建一对一的消息(这里支持多个接收者, 但消息类型为一对一, 适合于发送给少量人群的批量操作)
func (s *MysqlService) BatchCreateOneToOneMessage(message models.SysMessage, toIds []uint) error {
	// 强制修改类型为一对一
	message.Type = models.SysMessageTypeOneToOne

	// 设置默认过期时间
	if message.ExpiredAt == nil {
		now := time.Now()
		d, _ := time.ParseDuration("720h")
		expired := now.Add(d)
		message.ExpiredAt = &models.LocalTime{
			Time: expired,
		}
	}

	// 创建消息内容
	err := s.tx.Create(&message).Error
	if err != nil {
		return err
	}
	// 记录接收人
	for _, id := range toIds {
		var log models.SysMessageLog
		log.MessageId = message.Id
		log.ToUserId = id
		err = s.tx.Create(&log).Error
		if err != nil {
			return err
		}
	}

	return err
}

// 创建一对多的消息(这里支持多个接收角色, 适合于发送给指定角色的批量操作)
func (s *MysqlService) BatchCreateOneToManyMessage(message models.SysMessage, toRoleIds []uint) error {
	// 强制修改类型为一对多
	message.Type = models.SysMessageTypeOneToMany

	// 设置默认过期时间
	if message.ExpiredAt == nil {
		now := time.Now()
		d, _ := time.ParseDuration("720h")
		expired := now.Add(d)
		message.ExpiredAt = &models.LocalTime{
			Time: expired,
		}
	}

	// 记录接收人
	for _, id := range toRoleIds {
		// 清理消息编号
		message.Id = 0
		// 记录角色编号
		message.RoleId = id
		// 创建消息内容
		err := s.tx.Create(&message).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// 创建系统消息(适合于发送全部用户系统公告)
func (s *MysqlService) CreateSystemMessage(message models.SysMessage) error {
	// 强制修改类型为系统
	message.Type = models.SysMessageTypeSystem

	// 设置默认过期时间
	if message.ExpiredAt == nil {
		now := time.Now()
		d, _ := time.ParseDuration("720h")
		expired := now.Add(d)
		message.ExpiredAt = &models.LocalTime{
			Time: expired,
		}
	}

	// 创建消息内容
	return s.tx.Create(&message).Error
}
