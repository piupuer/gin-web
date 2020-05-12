package global

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var (
	// 系统配置
	Conf Configuration
	// packr盒子用于打包配置文件到golang编译后的二进制程序中
	ConfBox packr.Box
	// zap日志
	Log *zap.SugaredLogger
	// mysql实例
	Mysql *gorm.DB
	// validation.v9校验器
	Validate *validator.Validate
	// validation.v9相关翻译器
	Translator ut.Translator
)

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
	if c != nil && Conf.System.Transaction {
		// 从context对象中读取事务对象
		txKey, _ := c.Get("tx")
		tx, _ = txKey.(*gorm.DB)
	}
	return tx
}
