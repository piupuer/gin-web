package global

import (
	"errors"
	"gin-web/pkg/oss"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

var (
	// 系统配置
	Conf Configuration
	// packr盒子用于打包配置文件到golang编译后的二进制程序中
	ConfBox *CustomConfBox
	// zap日志
	Log *zap.SugaredLogger
	// mysql实例
	Mysql *gorm.DB
	// redis实例
	Redis *redis.Client
	// minio实例
	Minio *oss.MinioOss
	// validation.v9校验器
	Validate *validator.Validate
	// validation.v9相关翻译器
	Translator ut.Translator
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
func (c *CustomConfBox) Find(filename string) ([]byte, error) {
	if c.ConfEnv == "" {
		// 从packr box中读取
		return c.PackrBox.Find(filename)
	} else {
		// 从文件系统中读取
		return ioutil.ReadFile(c.ConfEnv + "/" + filename)
	}
}

// 只返回一个错误即可
func NewValidatorError(err error, custom map[string]string) (e error) {
	if err == nil {
		return
	}
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		tranStr := e.Translate(Translator)
		// 判断错误字段是否在自定义集合中，如果在，则替换错误信息中的字段
		if v, ok := custom[e.Field()]; ok {
			return errors.New(strings.Replace(tranStr, e.Field(), v, 1))
		} else {
			return errors.New(tranStr)
		}
	}
	return
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
