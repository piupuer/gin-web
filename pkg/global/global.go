package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/packr/v2"
	"github.com/piupuer/go-helper/pkg/oss"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
)

var (
	// Mode project mode: development/staging/production
	// RuntimeRoot runtime root path prefix
	Mode           string
	RuntimeRoot    string
	Conf           Configuration
	ConfBox        *CustomConfBox
	Mysql          *gorm.DB
	Redis          redis.UniversalClient
	CasbinEnforcer *casbin.Enforcer
	Minio          *oss.MinioOss
)

// custom conf box
type CustomConfBox struct {
	// conf path
	ConfEnv  string
	PackrBox *packr.Box
	ViperIns *viper.Viper
}

// find config file by filename
func (c *CustomConfBox) Find(filename string) []byte {
	f := filename
	if c.ConfEnv != "" {
		f = c.ConfEnv + "/" + filename
	}
	// read from system
	bs, _ := ioutil.ReadFile(f)
	if len(bs) == 0 {
		// read from packr box
		bs, _ = c.PackrBox.Find(f)
	}
	return bs
}
