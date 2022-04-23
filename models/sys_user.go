package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/ms"
)

const (
	SysUserStatusDisabled    uint   = 0
	SysUserStatusEnable      uint   = 1
	SysUserStatusDisabledStr string = "disabled"
	SysUserStatusEnableStr   string = "enable"
)

var SysUserStatusConst = map[uint]string{
	SysUserStatusDisabled: SysUserStatusDisabledStr,
	SysUserStatusEnable:   SysUserStatusEnableStr,
}

type SysUser struct {
	ms.M
	Username     string          `gorm:"index:idx_username,unique;comment:user login name" json:"username"`
	Password     string          `gorm:"comment:password" json:"password"`
	Mobile       string          `gorm:"comment:mobile number" json:"mobile"`
	Avatar       string          `gorm:"comment:avatar url" json:"avatar"`
	Nickname     string          `gorm:"comment:nickname" json:"nickname"`
	Introduction string          `gorm:"comment:introduction" json:"introduction"`
	Status       *uint           `gorm:"type:tinyint(1);default:1;comment:status(0: disabled, 1: enable)" json:"status"`
	RoleId       uint            `gorm:"comment:role id" json:"roleId"`
	Role         SysRole         `gorm:"foreignKey:RoleId" json:"role"`
	LastLogin    carbon.DateTime `gorm:"comment:last login time" json:"lastLogin"`
	Locked       uint            `gorm:"type:tinyint(1);default:0;comment:locked(0: unlock, 1: locked)" json:"locked"`
	LockExpire   int64           `gorm:"comment:lock expiration time" json:"lockExpire"`
	Wrong        int             `gorm:"comment:type wrong password count" json:"wrong"`
}
