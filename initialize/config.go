package initialize

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"os"
	"strconv"
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

var ctx *gin.Context // 生成启动时request id

// 初始化配置文件
func Config(c *gin.Context) {
	ctx = c
	// 初始化配置盒子
	var box global.CustomConfBox
	confDir := strings.ToLower(os.Getenv(fmt.Sprintf("%s_CONF", global.ProEnvName)))
	// 从环境变量中读取配置路径
	if confDir != "" {
		if strings.HasPrefix(confDir, "/") {
			// 指定的目录为绝对路径
			box.ConfEnv = confDir
		} else {
			// 指定的目录为相对路径
			box.ConfEnv = utils.GetWorkDir() + "/" + confDir
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
	env := strings.ToLower(os.Getenv(fmt.Sprintf("%s_MODE", global.ProEnvName)))
	configName := ""
	if env == global.Stage {
		configName = stagingConfig
	} else if env == global.Prod {
		configName = productionConfig
	} else {
		env = global.Dev
	}
	global.Mode = env
	if configName != "" {
		// 读取不同环境中的差异部分
		readConfig(v, configName)
	}
	// 转换为结构体
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 环境变量%s_CONF: %s", err, global.ProEnvName, global.ConfBox.ConfEnv))
	}

	// 从环境变量中加载配置: 如config.yml中system.port, 对应的环境变量为CFG_SYSTEM_PORT
	readConfigFromEnv(&global.Conf)

	if global.Conf.System.ConnectTimeout < 1 {
		global.Conf.System.ConnectTimeout = defaultConnectTimeout
	}

	if strings.TrimSpace(global.Conf.System.UrlPathPrefix) == "" {
		global.Conf.System.UrlPathPrefix = "api"
	}

	if strings.TrimSpace(global.Conf.System.ApiVersion) == "" {
		global.Conf.System.UrlPathPrefix = "v1"
	}

	global.Conf.Redis.BinlogPos = fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, global.Conf.Redis.BinlogPos)

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
	publicBytes := global.ConfBox.Find(global.Conf.System.RSAPublicKey)
	if len(publicBytes) == 0 {
		fmt.Println("RSA公钥未能加载, 请检查路径: ", global.Conf.System.RSAPublicKey)
	} else {
		global.Conf.System.RSAPublicBytes = publicBytes
	}
	privateBytes := global.ConfBox.Find(global.Conf.System.RSAPrivateKey)
	if len(privateBytes) == 0 {
		fmt.Println("RSA私钥未能加载, 请检查路径: ", global.Conf.System.RSAPrivateKey)
	} else {
		global.Conf.System.RSAPrivateBytes = privateBytes
	}

	fmt.Printf("初始化配置文件完成, 环境变量%s_CONF: %s\n", global.ProEnvName, global.ConfBox.ConfEnv)
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := global.ConfBox.Find(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("初始化配置文件失败, 环境变量%s_CONF: %s", global.ProEnvName, global.ConfBox.ConfEnv))
	}
	// 加载配置
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 环境变量%s_CONF: %s", err, global.ProEnvName, global.ConfBox.ConfEnv))
	}
}

// 从环境变量中加载配置(适用于docker镜像中不方便临时修改配置, 直接修改环境变量重启即可)
func readConfigFromEnv(defaultConfig *global.Configuration) {
	cfgMap := make(map[string]interface{}, 0)
	utils.Struct2StructByJson(defaultConfig, &cfgMap)
	newMap := parseCfgMap("", cfgMap)
	utils.Struct2StructByJson(newMap, &defaultConfig)
}

func parseCfgMap(parentKey string, m map[string]interface{}) map[string]interface{} {
	if parentKey == "" {
		parentKey = "CFG"
	}
	newMap := make(map[string]interface{}, 0)
	// json的几种基础类型(string/bool/float64)
	for key, item := range m {
		newKey := strings.ToUpper(fmt.Sprintf("%s_%s", utils.SnakeCase(parentKey), utils.SnakeCase(key)))
		switch item.(type) {
		case map[string]interface{}:
			// 仍然是map, 继续向下解析
			itemM, _ := item.(map[string]interface{})
			newMap[key] = parseCfgMap(newKey, itemM)
			continue
		case string:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				newMap[key] = env
				fmt.Println(fmt.Sprintf("[从环境变量中加载配置]读取到%s: %v", newKey, newMap[key]))
				continue
			}
		case bool:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				itemB, ok := item.(bool)
				b, err := strconv.ParseBool(env)
				if ok && err == nil {
					if itemB && !b {
						// 原值为true, 现为false
						newMap[key] = false
						fmt.Println(fmt.Sprintf("[从环境变量中加载配置]读取到%s: %v", newKey, newMap[key]))
						continue
					} else if !itemB && b {
						// 原值为false, 现为true
						newMap[key] = true
						fmt.Println(fmt.Sprintf("[从环境变量中加载配置]读取到%s: %v", newKey, newMap[key]))
						continue
					}
				}
			}
		case float64:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				v, err := strconv.ParseFloat(env, 64)
				if err == nil {
					newMap[key] = v
					fmt.Println(fmt.Sprintf("[从环境变量中加载配置]读取到%s: %v", newKey, newMap[key]))
					continue
				}
			}
		}
		// 值没有发生变化
		newMap[key] = item
	}
	return newMap
}
