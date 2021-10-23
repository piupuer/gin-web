package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/golang-module/carbon"
	"github.com/piupuer/go-helper/ms"
	"strings"
	"time"
)

// find status!=ms.SysMessageLogStatusDeleted messages
func (my MysqlService) FindUnDeleteMessage(req *request.MessageReq) []response.MessageResp {
	sysMessageLogTableName := global.Mysql.NamingStrategy.TableName("sys_message_log")
	sysMessageTableName := global.Mysql.NamingStrategy.TableName("sys_message")
	sysUserTableName := global.Mysql.NamingStrategy.TableName("sys_user")
	list := make([]response.MessageResp, 0)
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
	query := my.Q.Tx.
		Model(&ms.SysMessageLog{}).
		Select(fields).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.message_id = %s.id", sysMessageTableName, sysMessageLogTableName, sysMessageTableName)).
		Joins(fmt.Sprintf("LEFT JOIN %s AS toUser ON %s.to_user_id = toUser.id", sysUserTableName, sysMessageLogTableName)).
		Joins(fmt.Sprintf("LEFT JOIN %s AS fromUser ON %s.from_user_id = fromUser.id", sysUserTableName, sysMessageTableName))

	query = query.
		Order(fmt.Sprintf("%s.created_at DESC", sysMessageLogTableName)).
		Where(fmt.Sprintf("%s.to_user_id = ?", sysMessageLogTableName), req.ToUserId)
	title := strings.TrimSpace(req.Title)
	if title != "" {
		query = query.Where(fmt.Sprintf("%s.title LIKE ?", sysMessageTableName), fmt.Sprintf("%%%s%%", title))
	}
	content := strings.TrimSpace(req.Title)
	if content != "" {
		query = query.Where(fmt.Sprintf("%s.content LIKE ?", sysMessageTableName), fmt.Sprintf("%%%s%%", content))
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		query = query.Where(fmt.Sprintf("%s.status = ?", sysMessageLogTableName), *req.Status)
	} else {
		// un delete
		query = query.Where(fmt.Sprintf("%s.status != ?", sysMessageLogTableName), ms.SysMessageLogStatusDeleted)
	}

	// multi tables use ScanWithPage not FindWithPage
	my.Q.ScanWithPage(query, &req.Page, &list)
	return list
}

func (my MysqlService) GetUnReadMessageCount(userId uint) (int64, error) {
	var total int64
	err := my.Q.Tx.
		Model(&ms.SysMessageLog{}).
		Where("to_user_id = ?", userId).
		Where("status = ?", ms.SysMessageLogStatusUnRead).
		Count(&total).Error
	return total, err
}

func (my MysqlService) BatchUpdateMessageRead(messageLogIds []uint) error {
	return my.BatchUpdateMessageStatus(messageLogIds, ms.SysMessageLogStatusRead)
}

func (my MysqlService) BatchUpdateMessageDeleted(messageLogIds []uint) error {
	return my.BatchUpdateMessageStatus(messageLogIds, ms.SysMessageLogStatusDeleted)
}

func (my MysqlService) UpdateAllMessageRead(userId uint) error {
	return my.UpdateAllMessageStatus(userId, ms.SysMessageLogStatusRead)
}

func (my MysqlService) UpdateAllMessageDeleted(userId uint) error {
	return my.UpdateAllMessageStatus(userId, ms.SysMessageLogStatusDeleted)
}

func (my MysqlService) BatchUpdateMessageStatus(messageLogIds []uint, status uint) error {
	return my.Q.Tx.
		Model(&ms.SysMessageLog{}).
		Where("status != ?", ms.SysMessageLogStatusDeleted).
		Where("id IN (?)", messageLogIds).
		Update("status", status).Error
}

func (my MysqlService) UpdateAllMessageStatus(userId uint, status uint) error {
	var log ms.SysMessageLog
	log.ToUserId = userId
	return my.Q.Tx.
		Model(&log).
		Where("status != ?", ms.SysMessageLogStatusDeleted).
		Where(&log).
		Update("status", status).Error
}

func (my MysqlService) SyncMessageByUserIds(userIds []uint) error {
	users := make([]models.SysUser, 0)
	my.Q.Tx.
		Where("id IN (?)", userIds).
		Find(&users)
	for _, user := range users {
		messages := make([]ms.SysMessage, 0)
		my.Q.Tx.
			// > user register time
			Where("created_at > ?", user.CreatedAt).
			// expire < now
			Where("expired_at > ?", time.Now()).
			// one2many requires consistent roles, system is not required
			Where("(type = ? AND role_id = ?) OR type = ?", ms.SysMessageTypeOneToMany, user.RoleId, ms.SysMessageTypeSystem).
			Find(&messages)
		messageIds := make([]uint, 0)
		for _, message := range messages {
			messageIds = append(messageIds, message.Id)
		}
		// check whether is synced
		logs := make([]ms.SysMessageLog, 0)
		my.Q.Tx.
			Where("to_user_id = ?", user.Id).
			Where("message_id IN (?)", messageIds).
			Find(&logs)
		// old messages
		oldMessageIds := make([]uint, 0)
		for _, log := range logs {
			if !utils.ContainsUint(oldMessageIds, log.MessageId) {
				oldMessageIds = append(oldMessageIds, log.MessageId)
			}
		}
		for _, messageId := range messageIds {
			if !utils.ContainsUint(oldMessageIds, messageId) {
				// need create
				my.Q.Tx.Create(&ms.SysMessageLog{
					ToUserId:  user.Id,
					MessageId: messageId,
				})
			}
		}
	}
	return nil
}

func (my MysqlService) CreateMessage(req *request.PushMessageReq) error {
	if req.Type != nil {
		message := ms.SysMessage{
			FromUserId: req.FromUserId,
			Title:      req.Title,
			Content:    req.Content,
			Type:       uint(*req.Type),
		}
		switch uint(*req.Type) {
		case ms.SysMessageTypeOneToOne:
			if len(req.ToUserIds) == 0 {
				return fmt.Errorf("to user is empty")
			}
			return my.BatchCreateOneToOneMessage(message, req.ToUserIds)
		case ms.SysMessageTypeOneToMany:
			if len(req.ToRoleIds) == 0 {
				return fmt.Errorf("to role is empty")
			}
			return my.BatchCreateOneToManyMessage(message, req.ToRoleIds)
		case ms.SysMessageTypeSystem:
			return my.CreateSystemMessage(message)
		}
	}
	return fmt.Errorf("message type is illegal")
}

// one2one message
func (my MysqlService) BatchCreateOneToOneMessage(message ms.SysMessage, toIds []uint) error {
	message.Type = ms.SysMessageTypeOneToOne

	// default expire
	if message.ExpiredAt == nil {
		message.ExpiredAt = &carbon.ToDateTimeString{
			Carbon: carbon.Now().AddDays(30),
		}
	}

	err := my.Q.Tx.Create(&message).Error
	if err != nil {
		return err
	}
	// save ToUsers
	for _, id := range toIds {
		var log ms.SysMessageLog
		log.MessageId = message.Id
		log.ToUserId = id
		err = my.Q.Tx.Create(&log).Error
		if err != nil {
			return err
		}
	}

	return err
}

// one2many message
func (my MysqlService) BatchCreateOneToManyMessage(message ms.SysMessage, toRoleIds []uint) error {
	message.Type = ms.SysMessageTypeOneToMany

	if message.ExpiredAt == nil {
		message.ExpiredAt = &carbon.ToDateTimeString{
			Carbon: carbon.Now().AddDays(30),
		}
	}

	// save ToRoles
	for _, id := range toRoleIds {
		message.Id = 0
		message.RoleId = id
		err := my.Q.Tx.Create(&message).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// one2all message
func (my MysqlService) CreateSystemMessage(message ms.SysMessage) error {
	message.Type = ms.SysMessageTypeSystem

	if message.ExpiredAt == nil {
		message.ExpiredAt = &carbon.ToDateTimeString{
			Carbon: carbon.Now().AddDays(30),
		}
	}

	return my.Q.Tx.Create(&message).Error
}
