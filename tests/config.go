package tests

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configBoxName = "gin-conf-test-box"
	configType = "yml"
	// 配置文件目录, packr.Box基于当前包目录, 文件名需要写完整, 即使viper可以自动获取
	developmentConfig = "config.dev.yml"
	stagingConfig     = "config.stage.yml"
	productionConfig  = "config.prod.yml"
)

// 初始化配置文件
func InitConfig() {
	if os.Getenv("TEST_CONF") == "" {
		panic("[单元测试]请检查环境变量TEST_CONF")
	}
	// 使用packr将配置文件打包到二进制文件中, 如果以docker镜像方式运行将会非常舒服
	global.ConfBox = packr.New(configBoxName, os.Getenv("TEST_CONF"))
	// 获取实例(可创建多实例读取多个配置文件, 这里不做演示)
	v := viper.New()

	// 读取开发环境配置作为默认配置项
	readConfig(v, developmentConfig)
	// 将default中的配置全部以默认配置写入
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	// 读取当前go运行环境变量
	env := strings.ToLower(os.Getenv("GIN_WEB_MODE"))
	configName := ""
	if env == "staging" {
		configName = stagingConfig
	} else if env == "production" {
		configName = productionConfig
	}
	if configName != "" {
		// 读取不同环境中的差异部分
		readConfig(v, configName)
	}
	// 转换为结构体
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v", err))
	}
	fmt.Println("[单元测试]初始化配置文件完成")
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config, err := global.ConfBox.Find(configFile)
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v", err))
	}
	// 加载配置
	if err = v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v", err))
	}
}
