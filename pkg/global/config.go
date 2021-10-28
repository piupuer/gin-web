package global

import (
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap/zapcore"
)

// config from conf/config.dev.yml
type Configuration struct {
	System SystemConfiguration `mapstructure:"system" json:"system"`
	Logs   LogsConfiguration   `mapstructure:"logs" json:"logs"`
	Mysql  MysqlConfiguration  `mapstructure:"mysql" json:"mysql"`
	Redis  RedisConfiguration  `mapstructure:"redis" json:"redis"`
	Jwt    JwtConfiguration    `mapstructure:"jwt" json:"jwt"`
	Upload UploadConfiguration `mapstructure:"upload" json:"upload"`
}

type SystemConfiguration struct {
	MachineId            uint32 `mapstructure:"machine-id" json:"machineId"`
	UrlPrefix            string `mapstructure:"url-prefix" json:"urlPrefix"`
	ApiVersion           string `mapstructure:"api-version" json:"apiVersion"`
	Port                 int    `mapstructure:"port" json:"port"`
	PprofPort            int    `mapstructure:"pprof-port" json:"pprofPort"`
	ConnectTimeout       int    `mapstructure:"connect-timeout" json:"connectTimeout"`
	IdempotenceTokenName string `mapstructure:"idempotence-token-name" json:"idempotenceTokenName"`
	CasbinModelPath      string `mapstructure:"casbin-model-path" json:"casbinModelPath"`
	RateLimitMax         int64  `mapstructure:"rate-limit-max" json:"rateLimitMax"`
	AmapKey              string `mapstructure:"amap-key" json:"amapKey"`
}

type LogsConfiguration struct {
	Level                    zapcore.Level `mapstructure:"level" json:"level"`
	Path                     string        `mapstructure:"path" json:"path"`
	MaxSize                  int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups               int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge                   int           `mapstructure:"max-age" json:"maxAge"`
	Compress                 bool          `mapstructure:"compress" json:"compress"`
	OperationKey             string        `mapstructure:"operation-key" json:"operationKey"`
	OperationDisabledPaths   string        `mapstructure:"operation-disabled-paths" json:"operationDisabledPaths"`
	OperationDisabledPathArr []string      `mapstructure:"-" json:"-"`
	OperationAllowedToDelete bool          `mapstructure:"operation-allowed-to-delete" json:"operationAllowedToDelete"`
}

type MysqlConfiguration struct {
	Uri         string       `mapstructure:"uri" json:"uri"`
	TablePrefix string       `mapstructure:"table-prefix" json:"tablePrefix"`
	NoSql       bool         `mapstructure:"no-sql" json:"noSql"`
	Transaction bool         `mapstructure:"transaction" json:"transaction"`
	InitData    bool         `mapstructure:"init-data" json:"initData"`
	DSN         mysql.Config `json:"-"`
}

type RedisConfiguration struct {
	Uri          string `mapstructure:"uri" json:"uri"`
	Enable       bool   `mapstructure:"enable" json:"enable"`
}

type JwtConfiguration struct {
	Realm           string `mapstructure:"realm" json:"realm"`
	Key             string `mapstructure:"key" json:"key"`
	Timeout         int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh      int    `mapstructure:"max-refresh" json:"maxRefresh"`
	RSAPublicKey    string `mapstructure:"rsa-public-key" json:"rsaPublicKey"`
	RSAPrivateKey   string `mapstructure:"rsa-private-key" json:"rsaPrivateKey"`
	RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
}

type UploadConfiguration struct {
	SaveDir                      string                      `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize                int64                       `mapstructure:"single-max-size" json:"singleMaxSize"`
	MergeConcurrentCount         int                         `mapstructure:"merge-concurrent-count" json:"mergeConcurrentCount"`
	CompressImageCronTask        string                      `mapstructure:"compress-image-cron-task" json:"compressImageCronTask"`
	CompressImageRootDir         string                      `mapstructure:"compress-image-root-dir" json:"compressImageRootDir"`
	CompressImageOriginalSaveDir string                      `mapstructure:"compress-image-original-save-dir" json:"compressImageOriginalSaveDir"`
}
