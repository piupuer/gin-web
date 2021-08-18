package global

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"reflect"
	"strings"
	"time"
)

// zap日志自定义本地时间格式
func ZapLogLocalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(MsecLocalTimeFormat))
}

// zap logger for gorm
type GormZapLogger struct {
	log *zap.Logger
	logger.Config
	normalStr, traceStr, traceErrStr, traceWarnStr string
}

// New logger like gorm2
func NewGormZapLogger(zapLogger *zap.Logger, config logger.Config) *GormZapLogger {
	var (
		normalStr    = "%v%s "
		traceStr     = "%v%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%v%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%v%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		normalStr = logger.Cyan + "%v" + logger.Blue + "%s " + logger.Reset
		traceStr = logger.Cyan + "%v" + logger.Blue + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Cyan + "%v" + logger.Blue + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.Cyan + "%v" + logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	l := &GormZapLogger{
		log:          zapLogger,
		Config:       config,
		normalStr:    normalStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	return l
}

// LogMode gorm log mode
// LogMode log mode
func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Debug print info
func (l GormZapLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.DebugLevel) {
		requestId := getRequestId(ctx)
		l.log.Sugar().Debugf(l.normalStr+format, append([]interface{}{requestId, removePrefix(utils.FileWithLineNum())}, args...)...)
	}
}

// Info print info
func (l GormZapLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.InfoLevel) {
		requestId := getRequestId(ctx)
		l.log.Sugar().Infof(l.normalStr+format, append([]interface{}{requestId, removePrefix(utils.FileWithLineNum())}, args...)...)
	}
}

// Warn print warn messages
func (l GormZapLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.WarnLevel) {
		requestId := getRequestId(ctx)
		l.log.Sugar().Warnf(l.normalStr+format, append([]interface{}{requestId, removePrefix(utils.FileWithLineNum())}, args...)...)
	}
}

// Error print error messages
func (l GormZapLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.ErrorLevel) {
		requestId := getRequestId(ctx)
		l.log.Sugar().Errorf(l.normalStr+format, append([]interface{}{requestId, removePrefix(utils.FileWithLineNum())}, args...)...)
	}
}

// Trace print sql message
func (l GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.log.Core().Enabled(zapcore.DPanicLevel) || l.LogLevel <= logger.Silent {
		return
	}
	lineNum := removePrefix(utils.FileWithLineNum())
	elapsed := time.Since(begin)
	elapsedF := float64(elapsed.Nanoseconds()) / 1e6
	sql, rows := fc()
	row := "-"
	if rows > -1 {
		row = fmt.Sprintf("%d", rows)
	}
	requestId := getRequestId(ctx)
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

func getRequestId(ctx context.Context) string {
	var v interface{}
	vi := reflect.ValueOf(ctx)
	if vi.Kind() == reflect.Ptr {
		if !vi.IsNil() {
			v = ctx.Value(RequestIdContextKey)
		}
	}
	requestId := ""
	if v != nil {
		requestId = fmt.Sprintf("%v ", v)
	}
	return requestId
}

func removePrefix(s string) string {
	s = strings.TrimPrefix(s, RuntimeRoot)
	return s
}
