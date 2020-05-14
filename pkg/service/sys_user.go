package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// 登录校验
func (s *CommonService) LoginCheck(user *models.SysUser) (*models.SysUser, error) {
	var u models.SysUser
	// 查询用户及其角色
	err := s.tx.Preload("Role").Where("username = ?", user.Username).First(&u).Error
	if err != nil {
		return nil, err
	}
	// 校验密码
	if ok := utils.ComparePwd(user.Password, u.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	return &u, err
}

// 获取用户
func (s *CommonService) GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	var err error
	list := make([]models.SysUser, 0)
	db := global.Mysql
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = db.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 获取单个用户
func (s *CommonService) GetUserById(id uint) (models.SysUser, error) {
	var user models.SysUser
	var err error
	err = s.tx.Preload("Role").Where("id = ?", id).First(&user).Error
	return user, err
}

// 创建用户
func (s *CommonService) CreateUser(req *request.CreateUserRequestStruct) (err error) {
	var user models.SysUser
	utils.Struct2StructByJson(req, &user)
	// 将初始密码转为密文
	user.Password = utils.GenPwd(req.InitPassword)
	// 创建数据
	err = s.tx.Create(&user).Error
	return
}

// 更新用户
func (s *CommonService) UpdateUserById(id uint, newPassword string, req gin.H) (err error) {
	var oldUser models.SysUser
	query := s.tx.Table(oldUser.TableName()).Where("id = ?", id).First(&oldUser)
	if query.RecordNotFound() {
		return errors.New("记录不存在")
	}

	password := ""
	// 填写了新密码
	if strings.TrimSpace(newPassword) != "" {
		password = utils.GenPwd(newPassword)
	}
	// 比对增量字段
	m := make(gin.H, 0)
	utils.CompareDifferenceStructByJson(oldUser, req, &m)

	if password != "" {
		// 更新密码以及其他指定列
		err = query.Update("password", password).Updates(m).Error
	} else {
		// 更新指定列
		err = query.Updates(m).Error
	}
	return
}

// 批量删除用户
func (s *CommonService) DeleteUserByIds(ids []uint) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(models.SysUser{}).Error
}
