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
	WeChat WeChatConfiguration `mapstructure:"we-chat" json:"weChat"`
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
	Uri           string `mapstructure:"uri" json:"uri"`
	BinlogPos     string `mapstructure:"binlog-pos" json:"binlogPos"`
	Enable        bool   `mapstructure:"enable" json:"enable"`
	EnableService bool   `mapstructure:"enable-service" json:"enableService"`
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
	Minio                        UploadOssMinioConfiguration `mapstructure:"oss-minio" json:"ossMinio"`
	SaveDir                      string                      `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize                uint                        `mapstructure:"single-max-size" json:"singleMaxSize"`
	MergeConcurrentCount         uint                        `mapstructure:"merge-concurrent-count" json:"mergeConcurrentCount"`
	CompressImageCronTask        string                      `mapstructure:"compress-image-cron-task" json:"compressImageCronTask"`
	CompressImageRootDir         string                      `mapstructure:"compress-image-root-dir" json:"compressImageRootDir"`
	CompressImageOriginalSaveDir string                      `mapstructure:"compress-image-original-save-dir" json:"compressImageOriginalSaveDir"`
}

type UploadOssMinioConfiguration struct {
	Enable   bool   `mapstructure:"enable" json:"enable"`
	Bucket   string `mapstructure:"bucket" json:"bucket"`
	Endpoint string `mapstructure:"endpoint" json:"endpoint"`
	AccessId string `mapstructure:"access-id" json:"accessId"`
	Secret   string `mapstructure:"secret" json:"secret"`
	UseHttps bool   `mapstructure:"use-https" json:"useHttps"`
}

type WeChatConfiguration struct {
	Official WeChatOfficialConfiguration `mapstructure:"official" json:"official"`
}

type WeChatOfficialConfiguration struct {
	AppId              string                                        `mapstructure:"app-id" json:"appId"`
	AppSecret          string                                        `mapstructure:"app-secret" json:"appSecret"`
	Encoding           string                                        `mapstructure:"encoding" json:"encoding"`
	TplMessageCronTask WeChatOfficialTplMessageCronTaskConfiguration `mapstructure:"tpl-message-cron-task" json:"tplMessageCronTask"`
}

type WeChatOfficialTplMessageCronTaskConfiguration struct {
	Expr                string `mapstructure:"expr" json:"expr"`
	Users               string `mapstructure:"users" json:"users"`
	TemplateId          string `mapstructure:"template-id" json:"templateId"`
	MiniProgramAppId    string `mapstructure:"mini-program-app-id" json:"miniProgramAppId"`
	MiniProgramPagePath string `mapstructure:"mini-program-page-path" json:"miniProgramPagePath"`
}
