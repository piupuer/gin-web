package tests

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/query"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configType        = "yml"
	developmentConfig = "config.dev.yml"
	stagingConfig     = "config.stage.yml"
	productionConfig  = "config.prod.yml"
)

var ctx = query.NewRequestId(nil, constant.MiddlewareRequestIdCtxKey)

func Config() {
	if os.Getenv("TEST_CONF") == "" {
		panic("[unit test]check environment TEST_CONF")
	}
	var box global.CustomConfBox
	confDir := strings.ToLower(os.Getenv("TEST_CONF"))
	if confDir != "" {
		if strings.HasPrefix(confDir, "/") {
			box.ConfEnv = confDir
		} else {
			box.ConfEnv = utils.GetWorkDir() + "/" + confDir
		}
	}
	box.ViperIns = viper.New()
	global.ConfBox = &box
	v := box.ViperIns

	readConfig(v, developmentConfig)
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	env := strings.ToLower(os.Getenv(fmt.Sprintf("%s_MODE", global.ProEnvName)))
	configName := ""
	if env == constant.Stage {
		configName = stagingConfig
	} else if env == constant.Prod {
		configName = productionConfig
	} else {
		env = constant.Dev
	}
	global.Mode = env
	if configName != "" {
		readConfig(v, configName)
	}
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("[unit test]initialize config failed: %v, config env: TEST_CONF: %s", err, global.ConfBox.ConfEnv))
	}

	// read env to global.Conf: config.yml system.port => CFG_SYSTEM_PORT
	utils.EnvToInterface(&global.Conf, "CFG")

	// remove suffix "_"
	if strings.TrimSpace(global.Conf.Mysql.TablePrefix) != "" && strings.HasSuffix(global.Conf.Mysql.TablePrefix, "_") {
		global.Conf.Mysql.TablePrefix = strings.TrimSuffix(global.Conf.Mysql.TablePrefix, "_")
	}

	fmt.Printf("[unit test]initialize config success, config env: TEST_CONF: %s\n", global.ConfBox.ConfEnv)
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := global.ConfBox.Find(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("[unit test]initialize config failed, config env: TEST_CONF: %s", global.ConfBox.ConfEnv))
	}
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("[unit test]initialize config failed: %v, config env: TEST_CONF: %s", err, global.ConfBox.ConfEnv))
	}
}
