package initialize

import (
	"bytes"
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/gobuffalo/packr/v2"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configBoxName         = "gin-conf-box"
	configType            = "yml"
	configPath            = "../conf" // relative path: initialize/config.go
	developmentConfig     = "config.dev.yml"
	stagingConfig         = "config.stage.yml"
	productionConfig      = "config.prod.yml"
	defaultConnectTimeout = 5
)

var ctx context.Context

func Config(c context.Context) {
	ctx = c
	var box global.CustomConfBox
	// read config dir form env
	confDir := strings.ToLower(os.Getenv(fmt.Sprintf("%s_CONF", global.ProEnvName)))
	if confDir != "" {
		if strings.HasPrefix(confDir, "/") {
			// absolute path
			box.ConfEnv = confDir
		} else {
			// relative path: add work dir prefix
			box.ConfEnv = utils.GetWorkDir() + "/" + confDir
		}
	}
	box.ViperIns = viper.New()
	if box.ConfEnv == "" {
		// packr config files to binary executable file
		box.PackrBox = packr.New(configBoxName, configPath)
	}
	global.ConfBox = &box
	v := box.ViperIns

	// read development config as default config
	readConfig(v, developmentConfig)
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	// project mode
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
		// read diff config
		readConfig(v, configName)
	}
	// unmarshal to global.Conf
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(errors.Wrapf(err, "initialize config failed, config env: %s_CONF: %s", global.ProEnvName, global.ConfBox.ConfEnv))
	}

	// read env to global.Conf: config.yml system.port => CFG_SYSTEM_PORT
	utils.EnvToInterface(&global.Conf, "CFG")

	if global.Conf.System.ConnectTimeout < 1 {
		global.Conf.System.ConnectTimeout = defaultConnectTimeout
	}

	if strings.TrimSpace(global.Conf.System.UrlPrefix) == "" {
		global.Conf.System.UrlPrefix = "api"
	}

	if strings.TrimSpace(global.Conf.System.ApiVersion) == "" {
		global.Conf.System.ApiVersion = "v1"
	}

	global.Conf.System.Base = fmt.Sprintf("/%s/%s", global.Conf.System.UrlPrefix, global.Conf.System.ApiVersion)

	// remove suffix "_"
	if strings.TrimSpace(global.Conf.Mysql.TablePrefix) != "" && strings.HasSuffix(global.Conf.Mysql.TablePrefix, "_") {
		global.Conf.Mysql.TablePrefix = strings.TrimSuffix(global.Conf.Mysql.TablePrefix, "_")
	}

	if !global.Conf.Redis.Enable {
		global.Conf.Redis.EnableBinlog = false
	}

	// read rsa files
	publicBytes := global.ConfBox.Find(global.Conf.Jwt.RSAPublicKey)
	if len(publicBytes) == 0 {
		fmt.Println("read rsa public file failed, please check path: ", global.Conf.Jwt.RSAPublicKey)
	} else {
		global.Conf.Jwt.RSAPublicBytes = publicBytes
	}
	privateBytes := global.ConfBox.Find(global.Conf.Jwt.RSAPrivateKey)
	if len(privateBytes) == 0 {
		fmt.Println("read rsa private file failed, please check path: ", global.Conf.Jwt.RSAPrivateKey)
	} else {
		global.Conf.Jwt.RSAPrivateBytes = privateBytes
	}

	// change default logger
	log.DefaultWrapper = log.NewWrapper(log.New(
		log.WithCategory(constant.LogCategoryLogrus),
		log.WithLevel(global.Conf.Logs.Level),
		log.WithJson(false),
		log.WithLineNumPrefix(global.RuntimeRoot),
		log.WithLineNumLevel(1),
		log.WithKeepVersion(false),
		log.WithKeepSourceDir(false),
	))

	log.Info("initialize config success, config env: %s_CONF: %s", global.ProEnvName, global.ConfBox.ConfEnv)
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := global.ConfBox.Find(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("initialize config failed, config env: %s_CONF: %s", global.ProEnvName, global.ConfBox.ConfEnv))
	}
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(errors.Wrapf(err, "initialize config failed, config env: %s_CONF: %s", global.ProEnvName, global.ConfBox.ConfEnv))
	}
}
