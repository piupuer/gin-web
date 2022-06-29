package tests

import (
	"bytes"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
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

var ctx = tracing.NewId(nil)

func Config() {
	if os.Getenv("TEST_CONF") == "" {
		panic("[unit test]check environment TEST_CONF")
	}
	confDir := os.Getenv("TEST_CONF")
	var box ms.ConfBox
	box.Dir = confDir
	global.ConfBox = box
	v := viper.New()
	// read development config as default config
	readConfig(box, v, developmentConfig)
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	env := strings.ToLower(os.Getenv(fmt.Sprintf("%s_MODE", global.ProProdName)))
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
		readConfig(box, v, configName)
	}
	// unmarshal to global.Conf
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(errors.Wrapf(err, "[unit test]initialize config failed, config env: TEST_CONF: %s", box.Dir))
	}

	// read env to global.Conf: config.yml system.port => CFG_SYSTEM_PORT
	envPrefix := strings.ToUpper(os.Getenv(fmt.Sprintf("%s_ENV", global.ProEnvName)))
	if envPrefix == "" {
		envPrefix = "CFG"
	}
	utils.EnvToInterface(
		utils.WithEnvObj(&global.Conf),
		utils.WithEnvPrefix(envPrefix),
	)
	// remove suffix "_"
	if strings.TrimSpace(global.Conf.Mysql.TablePrefix) != "" && strings.HasSuffix(global.Conf.Mysql.TablePrefix, "_") {
		global.Conf.Mysql.TablePrefix = strings.TrimSuffix(global.Conf.Mysql.TablePrefix, "_")
	}

	// change default logger
	log.DefaultWrapper = log.NewWrapper(log.New(
		log.WithCategory(global.Conf.Logs.Category),
		log.WithLevel(global.Conf.Logs.Level),
		log.WithJson(global.Conf.Logs.Json),
		log.WithLineNumPrefix(global.RuntimeRoot),
		log.WithLineNum(!global.Conf.Logs.LineNum.Disable),
		log.WithLineNumLevel(global.Conf.Logs.LineNum.Level),
		log.WithLineNumVersion(global.Conf.Logs.LineNum.Version),
		log.WithLineNumSource(global.Conf.Logs.LineNum.Source),
	))

	log.Info("[unit test]initialize config success, config env: `TEST_CONF: %s`", box.Dir)
}

func readConfig(box ms.ConfBox, v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := box.Get(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("[unit test]initialize config failed, config env: `TEST_CONF: %s`", box.Dir))
	}
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(errors.Wrapf(err, "[unit test]initialize config failed, config env: `TEST_CONF: %s`", box.Dir))
	}
}
