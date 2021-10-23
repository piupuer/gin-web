package cache_service

import (
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/ms"
	"strings"
)

func (rd RedisService) FindMachine(req *request.MachineReq) ([]ms.SysMachine, error) {
	if !global.Conf.Redis.Enable || !global.Conf.Redis.EnableService {
		return rd.mysql.FindMachine(req)
	}
	var err error
	list := make([]ms.SysMachine, 0)
	query := rd.Q.
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
	if req.Status != nil {
		query = query.Where("status", "=", *req.Status)
	}
	err = rd.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}
