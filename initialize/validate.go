package initialize

import (
	"gin-web/pkg/global"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

// 初始化校验器
func Validate() {
	// 实例化需要转换的语言, 中文
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()

	// 注册转换的语言为默认语言
	_ = translations.RegisterDefaultTranslations(validate, trans)

	global.Validate = validate
	global.Translator = trans
	global.Log.Debug("初始化validator.v9校验器完成")
}
