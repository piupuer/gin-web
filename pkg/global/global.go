package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/oss"
	"gorm.io/gorm"
)

var (
	// Mode project mode: development/staging/production
	// RuntimeRoot runtime root path prefix
	Mode           string
	RuntimeRoot    string
	Conf           Configuration
	ConfBox        ms.ConfBox
	Mysql          *gorm.DB
	Redis          redis.UniversalClient
	CasbinEnforcer *casbin.Enforcer
	Minio          *oss.MinioOss
)
