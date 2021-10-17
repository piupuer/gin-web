package global

import (
	"gin-web/pkg/oss"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/packr/v2"
	"github.com/piupuer/go-helper/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
)

var (
	// 当前模式
	Mode string
	// 系统配置
	Conf Configuration
	// packr盒子用于打包配置文件到golang编译后的二进制程序中
	ConfBox *CustomConfBox
	// zap日志
	Log *logger.Logger
	// mysql实例
	Mysql *gorm.DB
	// redis实例
	Redis redis.UniversalClient
	// cabin实例
	CasbinEnforcer *casbin.Enforcer
	// minio实例
	Minio *oss.MinioOss
	// 运行时根目录
	RuntimeRoot string
)

// 自定义配置盒子
type CustomConfBox struct {
	// 配置文件路径环境变量
	ConfEnv string
	// packr盒子
	PackrBox *packr.Box
	// viper实例
	ViperIns *viper.Viper
}

// 查找指定配置
func (c *CustomConfBox) Find(filename string) []byte {
	f := filename
	if c.ConfEnv != "" {
		f = c.ConfEnv + "/" + filename
	}
	// 从文件系统中读取
	bs, _ := ioutil.ReadFile(f)
	if len(bs) == 0 {
		// 从packr box中读取
		bs, _ = c.PackrBox.Find(f)
	}
	return bs
}

// 获取事务对象
func GetTx(c *gin.Context) *gorm.DB {
	// 默认使用无事务的mysql
	tx := Mysql
	if c != nil {
		method := ""
		if c.Request != nil {
			method = c.Request.Method
		}
		if !(method == "OPTIONS" || method == "GET" || !Conf.System.Transaction) {
			// 从context对象中读取事务对象
			txKey, exists := c.Get("tx")
			if exists {
				if item, ok := txKey.(*gorm.DB); ok {
					tx = item
				}
			}
		}
	}
	return tx
}

// 获取携带request id的上下文
func RequestIdContext(requestId string) *gin.Context {
	if requestId == "" {
		uuid4 := uuid.NewV4()
		requestId = uuid4.String()
	}
	ctx := gin.Context{}
	ctx.Set(RequestIdContextKey, requestId)
	return &ctx
}
