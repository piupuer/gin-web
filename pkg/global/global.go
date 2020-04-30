package global

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var (
	Log        *zap.SugaredLogger
	Mysql      *gorm.DB
	Validate   *validator.Validate
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
