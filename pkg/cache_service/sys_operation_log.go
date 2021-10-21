package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strconv"
	"strings"
)

// 获取所有操作日志
func (s RedisService) GetOperationLogs(req *request.OperationLogReq) ([]models.SysOperationLog, error) {
	if !global.Conf.System.UseRedis || !global.Conf.System.UseRedisService {
		// 不使用redis
		return s.mysql.GetOperationLogs(req)
	}
	var err error
	list := make([]models.SysOperationLog, 0)
	query := s.Q.
		Table("sys_operation_log").
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method", "contains", method)
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path", "contains", path)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		query = query.Where("username", "contains", username)
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		query = query.Where("ip", "contains", ip)
	}
	status := strings.TrimSpace(req.Status)
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status", "contains", s)
	}
	// 查询列表
	err = s.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}
