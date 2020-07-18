package service

import (
	"gin-web/models"
)

// 批量删除操作记录
func (s *MysqlService) DeleteOperationLogByIds(ids []uint) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(models.SysOperationLog{}).Error
}
