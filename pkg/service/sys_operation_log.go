package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

// 获取操作日志
func (s MysqlService) GetOperationLogs(req *request.OperationLogReq) ([]models.SysOperationLog, error) {
	var err error
	list := make([]models.SysOperationLog, 0)
	query := global.Mysql.
		Model(&models.SysOperationLog{}).
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		query = query.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	status := strings.TrimSpace(req.Status)
	if status != "" {
		query = query.Where("status LIKE ?", fmt.Sprintf("%%%s%%", status))
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}
