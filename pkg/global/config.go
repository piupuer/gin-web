package global

import (
	"github.com/go-sql-driver/mysql"
	"github.com/piupuer/go-helper/pkg/log"
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
	Base                 string `mapstructure:"-" json:"-"`
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
	Category                 string                   `mapstructure:"category" json:"category"`
	Level                    log.Level                `mapstructure:"level" json:"level"`
	Json                     bool                     `mapstructure:"json" json:"json"`
	LineNum                  LogsLineNumConfiguration `mapstructure:"line-num" json:"lineNum"`
	OperationKey             string                   `mapstructure:"operation-key" json:"operationKey"`
	OperationAllowedToDelete bool                     `mapstructure:"operation-allowed-to-delete" json:"operationAllowedToDelete"`
}

type LogsLineNumConfiguration struct {
	Disable bool `mapstructure:"disable" json:"disable"`
	Level   int  `mapstructure:"level" json:"level"`
	Version bool `mapstructure:"version" json:"version"`
	Source  bool `mapstructure:"source" json:"source"`
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
	Minio                UploadOssMinioConfiguration `mapstructure:"oss-minio" json:"ossMinio"`
	SaveDir              string                      `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize        int64                       `mapstructure:"single-max-size" json:"singleMaxSize"`
	MergeConcurrentCount int                         `mapstructure:"merge-concurrent-count" json:"mergeConcurrentCount"`
}
