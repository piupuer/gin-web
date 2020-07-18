package global

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

/**
 * 初始化日志
 * filename 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 * 由于zap不具备日志切割功能, 这里使用lumberjack配合
 */
func InitLogger() {
	now := time.Now()
	filename := fmt.Sprintf("%s/%04d-%02d-%02d.log", Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	hook := &lumberjack.Logger{
		Filename:   filename,             // 日志文件路径
		MaxSize:    Conf.Logs.MaxSize,    // 最大尺寸, M
		MaxBackups: Conf.Logs.MaxBackups, // 备份数
		MaxAge:     Conf.Logs.MaxAge,     // 存放天数
		Compress:   Conf.Logs.Compress,   // 是否压缩
	}
	defer hook.Close()
	// zap 的 Config 非常的繁琐也非常强大，可以控制打印 log 的所有细节，因此对于我们开发者是友好的，有利于二次封装。
	// 但是对于初学者则是噩梦。因此 zap 提供了一整套的易用配置，大部分的姿势都可以通过一句代码生成需要的配置。
	enConfig := zap.NewProductionEncoderConfig() // 生成配置

	// 时间格式
	enConfig.EncodeTime = ZapLogLocalTimeEncoder
	// level字母大写
	enConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enConfig),                                            // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		Conf.Logs.Level,                                                                // 日志等级
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()
	Log.Debug("初始化日志完成")
}

// zap日志自定义本地时间格式
func ZapLogLocalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(MsecLocalTimeFormat))
}
