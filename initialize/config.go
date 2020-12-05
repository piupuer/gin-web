package initialize

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configBoxName         = "gin-conf-box"
	configType            = "yml"
	configPath            = "../conf" // 配置文件目录, packr.Box基于当前包目录, 文件名需要写完整, 即使viper可以自动获取
	developmentConfig     = "config.dev.yml"
	stagingConfig         = "config.stage.yml"
	productionConfig      = "config.prod.yml"
	defaultConnectTimeout = 5
)

// 初始化配置文件
func Config() {
	// 初始化配置盒子
	var box global.CustomConfBox
	ginWebConf := strings.ToLower(os.Getenv("GIN_WEB_CONF"))
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
	// 环境变量不存在, 需要打包到二进制文件中
	if box.ConfEnv == "" {
		// 使用packr将配置文件打包到二进制文件中, 如果以docker镜像方式运行将会非常舒服
		box.PackrBox = packr.New(configBoxName, configPath)
	}
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
		panic(fmt.Sprintf("初始化配置文件失败: %v, 环境变量GIN_WEB_CONF: %s", err, global.ConfBox.ConfEnv))
	}

	if global.Conf.System.ConnectTimeout < 1 {
		global.Conf.System.ConnectTimeout = defaultConnectTimeout
	}

	if strings.TrimSpace(global.Conf.System.UrlPathPrefix) == "" {
		global.Conf.System.UrlPathPrefix = "api"
	}

	if strings.TrimSpace(global.Conf.System.ApiVersion) == "" {
		global.Conf.System.UrlPathPrefix = "v1"
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

	// 加载rsa公私钥(优先从configBox中读取)
	publicBytes, err := global.ConfBox.Find(global.Conf.System.RSAPublicKey)
	if err != nil || len(publicBytes) == 0 {
		publicBytes = utils.RSAReadKeyFromFile(global.Conf.System.RSAPublicKey)
	}
	if len(publicBytes) == 0 {
		fmt.Println("RSA公钥未能加载, 请检查路径: ", global.Conf.System.RSAPublicKey)
	} else {
		global.Conf.System.RSAPublicBytes = publicBytes
	}
	privateBytes, err := global.ConfBox.Find(global.Conf.System.RSAPrivateKey)
	if err != nil || len(privateBytes) == 0 {
		privateBytes = utils.RSAReadKeyFromFile(global.Conf.System.RSAPrivateKey)
	}
	if len(privateBytes) == 0 {
		fmt.Println("RSA私钥未能加载, 请检查路径: ", global.Conf.System.RSAPrivateKey)
	} else {
		global.Conf.System.RSAPrivateBytes = privateBytes
	}

	fmt.Println("初始化配置文件完成, 环境变量GIN_WEB_CONF: ", global.ConfBox.ConfEnv)
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config, err := global.ConfBox.Find(configFile)
	if err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 环境变量GIN_WEB_CONF: %s", err, global.ConfBox.ConfEnv))
	}
	// 加载配置
	if err = v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 环境变量GIN_WEB_CONF: %s", err, global.ConfBox.ConfEnv))
	}
}
