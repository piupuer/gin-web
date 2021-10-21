package models

import (
	"github.com/piupuer/go-helper/models"
	"time"
)

type SysOperationLog struct {
	models.M
	ApiDesc    string        `json:"apiDesc" gorm:"comment:'api description'"`
	Path       string        `json:"path" gorm:"comment:'url path'"`
	Method     string        `json:"method" gorm:"comment:'api method'"`
	Header     string        `json:"header" gorm:"type:blob;comment:'request header'"`
	Body       string        `json:"body" gorm:"type:blob;comment:'request body'"`
	Data       string        `json:"data" gorm:"type:blob;comment:'response data'"`
	Status     int           `json:"status" gorm:"comment:'response status'"`
	Username   string        `json:"username" gorm:"comment:'login username'"`
	RoleName   string        `json:"roleName" gorm:"comment:'login role name'"`
	Ip         string        `json:"ip" gorm:"comment:'IP'"`
	IpLocation string        `json:"ipLocation" gorm:"comment:'real location of the IP'"`
	Latency    time.Duration `json:"latency" gorm:"comment:'request time(ms)'"`
	UserAgent  string        `json:"userAgent" gorm:"comment:'browser user agent'"`
}
