package models

import (
	"github.com/piupuer/go-helper/models"
	"time"
)

type SysOperationLog struct {
	models.M
	ApiDesc    string        `gorm:"comment:'api description'" json:"apiDesc"`
	Path       string        `gorm:"comment:'url path'" json:"path"`
	Method     string        `gorm:"comment:'api method'" json:"method"`
	Header     string        `gorm:"type:blob;comment:'request header'" json:"header"`
	Body       string        `gorm:"type:blob;comment:'request body'" json:"body"`
	Data       string        `gorm:"type:blob;comment:'response data'" json:"data"`
	Status     int           `gorm:"comment:'response status'" json:"status"`
	Username   string        `gorm:"comment:'login username'" json:"username"`
	RoleName   string        `gorm:"comment:'login role name'" json:"roleName"`
	Ip         string        `gorm:"comment:'IP'" json:"ip"`
	IpLocation string        `gorm:"comment:'real location of the IP'" json:"ipLocation"`
	Latency    time.Duration `gorm:"comment:'request time(ms)'" json:"latency"`
	UserAgent  string        `gorm:"comment:'browser user agent'" json:"userAgent"`
}
