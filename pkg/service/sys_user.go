package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
	"strings"
)

func (my MysqlService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	var u models.SysUser
	err := my.Q.Tx.Preload("Role").Where("username = ?", user.Username).First(&u).Error
	if err != nil {
		return nil, fmt.Errorf(resp.LoginCheckErrorMsg)
	}
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, fmt.Errorf(resp.LoginCheckErrorMsg)
	}
	return &u, err
}

func (my MysqlService) FindUser(req *request.UserReq) []models.SysUser {
	list := make([]models.SysUser, 0)
	query := my.Q.Tx.
		Model(&models.SysUser{}).
		Order("created_at DESC")
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds := my.FindRoleIdBySort(*req.CurrentRole.Sort)
		query = query.Where("role_id IN (?)", roleIds)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		query = query.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		query = query.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		query = query.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status = ?", 1)
		} else {
			query = query.Where("status = ?", 0)
		}
	}
	my.Q.FindWithPage(query, &req.Page, &list)
	return list
}

func (my MysqlService) GetUserById(id uint) (models.SysUser, error) {
	var user models.SysUser
	var err error
	err = my.Q.Tx.Preload("Role").
		Where("id = ?", id).
		Where("status = ?", models.SysUserStatusEnable).
		First(&user).Error
	return user, err
}
