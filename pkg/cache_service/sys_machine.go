package cache_service

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"strings"
)

func (s RedisService) GetMachines(req *request.MachineReq) ([]models.SysMachine, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		// 不使用redis
		return s.mysql.GetMachines(req)
	}
	var err error
	list := make([]models.SysMachine, 0)
	query := s.Q.
		Table("sys_machine").
		Order("created_at DESC")
	host := strings.TrimSpace(req.Host)
	if host != "" {
		query = query.Where("host", "contains", host)
	}
	loginName := strings.TrimSpace(req.LoginName)
	if loginName != "" {
		query = query.Where("login_name", "contains", loginName)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	// 查询列表
	err = s.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}
