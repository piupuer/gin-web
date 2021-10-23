package service

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/ms"
	"strings"
)

func (my MysqlService) FindOperationLog(req *request.OperationLogReq) ([]ms.SysOperationLog, error) {
	var err error
	list := make([]ms.SysOperationLog, 0)
	query := global.Mysql.
		Model(&ms.SysOperationLog{}).
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
	err = my.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}
