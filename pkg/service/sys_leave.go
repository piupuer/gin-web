package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// 获取所有请假
func (s *MysqlService) GetLeaves(req *request.LeaveListRequestStruct) ([]models.SysLeave, error) {
	var err error
	list := make([]models.SysLeave, 0)
	query := s.tx
	desc := strings.TrimSpace(req.Desc)
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if desc != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", desc))
	}
	// 查询条数
	err = query.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 创建请假
func (s *MysqlService) CreateLeave(req *request.CreateLeaveRequestStruct) (err error) {
	var leave models.SysLeave
	utils.Struct2StructByJson(req, &leave)
	// 创建数据
	err = s.tx.Create(&leave).Error
	return
}

// 更新请假
func (s *MysqlService) UpdateLeaveById(id uint, req gin.H) (err error) {
	var leave models.SysLeave
	query := s.tx.Table(leave.TableName()).Where("id = ?", id).First(&leave)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(leave, req, &m)

	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除请假
func (s *MysqlService) DeleteLeaveByIds(ids []uint) (err error) {
	var list []models.SysLeave
	query := s.tx.Where("id IN (?)", ids).Find(&list)
	if query.Error != nil {
		return
	}
	return query.Delete(models.SysLeave{}).Error
}
