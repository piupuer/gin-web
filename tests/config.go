package tests

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configType = "yml"
	developmentConfig = "config.dev.yml"
	stagingConfig     = "config.stage.yml"
	productionConfig  = "config.prod.yml"
)

// 初始化配置文件
func Config() {
	if os.Getenv("TEST_CONF") == "" {
		panic("[单元测试]请检查环境变量TEST_CONF")
	}
	// 初始化配置盒子
	var box global.CustomConfBox
	ginWebConf := strings.ToLower(os.Getenv("TEST_CONF"))
	// 从环境变量中读取配置路径
	if ginWebConf != "" {
		if strings.HasPrefix(ginWebConf, "/") {
			// 指定的目录为绝对路径
			box.ConfEnv = ginWebConf
		} else {
			// 指定的目录为相对路径
			box.ConfEnv = utils.GetWorkDir() + "/" + ginWebConf
		}
	}
	// 获取viper实例(可创建多实例读取多个配置文件, 这里不做演示)
	box.ViperIns = viper.New()
	global.ConfBox = &box
	v := box.ViperIns

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
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v, 环境变量GIN_WEB_CONF: %s", err, global.ConfBox.ConfEnv))
	}

	// 表前缀去掉后缀_
	if strings.TrimSpace(global.Conf.Mysql.TablePrefix) != "" && strings.HasSuffix(global.Conf.Mysql.TablePrefix, "_") {
		global.Conf.Mysql.TablePrefix = strings.TrimSuffix(global.Conf.Mysql.TablePrefix, "_")
	}

	// 初始化OperationLogDisabledPaths
	global.Conf.System.OperationLogDisabledPathArr = make([]string, 0)
	if strings.TrimSpace(global.Conf.System.OperationLogDisabledPaths) != "" {
		global.Conf.System.OperationLogDisabledPathArr = strings.Split(global.Conf.System.OperationLogDisabledPaths, ",")
	}
	fmt.Println("[单元测试]初始化配置文件完成")
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config, err := global.ConfBox.Find(configFile)
	if err != nil {
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v, 环境变量TEST_CONF: %s", err, global.ConfBox.ConfEnv))
	}
	// 加载配置
	if err = v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("[单元测试]初始化配置文件失败: %v, 环境变量TEST_CONF: %s", err, global.ConfBox.ConfEnv))
	}
}
