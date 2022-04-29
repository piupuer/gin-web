package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/oss"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

var (
	// Mode project mode: development/staging/production
	// RuntimeRoot runtime root path prefix
	Mode           string
	RuntimeRoot    string
	Conf           Configuration
	ConfBox        ms.ConfBox
	Tracer         *trace.TracerProvider
	Mysql          *gorm.DB
	Redis          redis.UniversalClient
	CasbinEnforcer *casbin.Enforcer
	Minio          *oss.MinioOss
)
