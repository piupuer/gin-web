package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

// 初始化日志
func Logger() {
	now := time.Now()
	filename := fmt.Sprintf("%s/%04d-%02d-%02d.log", global.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	hook := &lumberjack.Logger{
		Filename:   filename,                    // 日志文件路径
		MaxSize:    global.Conf.Logs.MaxSize,    // 最大尺寸, M
		MaxBackups: global.Conf.Logs.MaxBackups, // 备份数
		MaxAge:     global.Conf.Logs.MaxAge,     // 存放天数
		Compress:   global.Conf.Logs.Compress,   // 是否压缩
	}
	defer hook.Close()
	// zap 的 Config 非常的繁琐也非常强大，可以控制打印 log 的所有细节，因此对于我们开发者是友好的，有利于二次封装。
	// 但是对于初学者则是噩梦。因此 zap 提供了一整套的易用配置，大部分的姿势都可以通过一句代码生成需要的配置。
	enConfig := zap.NewProductionEncoderConfig() // 生成配置

	// 时间格式
	enConfig.EncodeTime = global.ZapLogLocalTimeEncoder
	// level字母大写
	enConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enConfig),                                            // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		global.Conf.Logs.Level,                                                         // 日志等级
	)

	l := zap.New(core)
	colorful := false
	if global.Conf.Logs.Level <= zapcore.DebugLevel {
		colorful = true
	}
	global.Log = global.NewGormZapLogger(l, logger.Config{
		Colorful: colorful,
	})
	global.Log.Debug(ctx, "初始化日志完成")
}
