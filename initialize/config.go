package initialize

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configType            = "yml"
	configDir             = "conf"
	developmentConfig     = "config.dev.yml"
	stagingConfig         = "config.stage.yml"
	productionConfig      = "config.prod.yml"
	defaultConnectTimeout = 5
)

var ctx context.Context

func Config(c context.Context, conf embed.FS) {
	ctx = c
	confDir := os.Getenv(fmt.Sprintf("%s_CONF", global.ProEnvName))
	var box ms.ConfBox
	box.Ctx = ctx
	box.Fs = conf
	if confDir == "" {
		confDir = configDir
	}
	box.Dir = confDir
	global.ConfBox = box
	v := viper.New()
	// read development config as default config
	readConfig(box, v, developmentConfig)
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	// project mode
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
		// read diff config
		readConfig(box, v, configName)
	}
	// unmarshal to global.Conf
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(errors.Wrapf(err, "initialize config failed, config env: %s_CONF: %s", global.ProEnvName, box.Dir))
	}

	// read env to global.Conf: config.yml system.port => CFG_SYSTEM_PORT
	envPrefix := strings.ToUpper(os.Getenv(fmt.Sprintf("%s_ENV", global.ProEnvName)))
	if envPrefix == "" {
		envPrefix = "CFG"
	}
	utils.EnvToInterface(
		utils.WithEnvObj(&global.Conf),
		utils.WithEnvPrefix(envPrefix),
		utils.WithEnvFormat(func(key string, val interface{}) string {
			if utils.Contains([]string{
				// hidden val
				"CFG_MYSQL_URI",
				"CFG_REDIS_URI",
				"CFG_JWT_REALM",
				"CFG_JWT_KEY",
				"CFG_UPLOAD_OSS_MINIO_SECRET",
			}, key) {
				val = "******"
			}
			return fmt.Sprintf("%s: %v", key, val)
		}),
	)

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
	publicBytes := box.Get(global.Conf.Jwt.RSAPublicKey)
	if len(publicBytes) == 0 {
		log.WithContext(ctx).Warn("read rsa public file failed, please check path: %s", global.Conf.Jwt.RSAPublicKey)
	} else {
		global.Conf.Jwt.RSAPublicBytes = publicBytes
	}
	privateBytes := box.Get(global.Conf.Jwt.RSAPrivateKey)
	if len(privateBytes) == 0 {
		log.WithContext(ctx).Warn("read rsa private file failed, please check path: %s", global.Conf.Jwt.RSAPrivateKey)
	} else {
		global.Conf.Jwt.RSAPrivateBytes = privateBytes
	}

	log.WithContext(ctx).Info("initialize config success, config env: `%s_CONF: %s`", global.ProEnvName, box.Dir)
}

func readConfig(box ms.ConfBox, v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := box.Get(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("initialize config failed, config env: `%s_CONF: %s`", global.ProEnvName, box.Dir))
	}
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(errors.Wrapf(err, "initialize config failed, config env: `%s_CONF: %s`", global.ProEnvName, box.Dir))
	}
}
