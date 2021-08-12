package global

import (
	"context"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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

	l := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	Log = NewGormZapLogger(l, logger.Config{}).log.Sugar()
	Log.Debug("初始化日志完成")
}

// zap日志自定义本地时间格式
func ZapLogLocalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(MsecLocalTimeFormat))
}

// zap logger for gorm
type GormZapLogger struct {
	log *zap.Logger
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// New logger like gorm2
func NewGormZapLogger(zapLogger *zap.Logger, config logger.Config) *GormZapLogger {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%v%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%v%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%v%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = "%v" + logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = "%v" + logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = "%v" + logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	l := &GormZapLogger{
		log: zapLogger,
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Warn,
		},
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	return l
}

// LogMode gorm log mode
func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info print info
func (l GormZapLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.DebugLevel) {
		l.log.Sugar().Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Info print info
func (l GormZapLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.InfoLevel) {
		l.log.Sugar().Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Warn print warn messages
func (l GormZapLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.WarnLevel) {
		l.log.Sugar().Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Error print error messages
func (l GormZapLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.ErrorLevel) {
		l.log.Sugar().Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Trace print sql message
func (l GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.log.Core().Enabled(zapcore.DPanicLevel) || l.LogLevel <= logger.Silent {
		return
	}
	lineNum := utils.FileWithLineNum()
	elapsed := time.Since(begin)
	elapsedF := float64(elapsed.Nanoseconds()) / 1e6
	sql, rows := fc()
	row := "-"
	if rows > -1 {
		row = fmt.Sprintf("%d", rows)
	}
	v := ctx.Value(RequestIdContextKey)
	requestId := ""
	if v != nil {
		requestId = fmt.Sprintf("%v ", v)
	}
	switch {
	case l.log.Core().Enabled(zapcore.ErrorLevel) && err != nil && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.log.Error(fmt.Sprintf(l.traceErrStr, requestId, lineNum, err, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.WarnLevel) && elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.log.Warn(fmt.Sprintf(l.traceWarnStr, requestId, lineNum, slowLog, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.DebugLevel):
		l.log.Debug(fmt.Sprintf(l.traceStr, requestId, lineNum, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.InfoLevel):
		l.log.Info(fmt.Sprintf(l.traceStr, requestId, lineNum, elapsedF, row, sql))
	}
}
