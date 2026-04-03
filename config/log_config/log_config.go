package log_config

import (
	"os"
	"path/filepath"

	conf "github.com/MagicPig9898/familychefassistant_server/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var sugar *zap.SugaredLogger

// NewLogConfig 根据配置初始化日志
func NewLogConfig() error {
	c := conf.Cfg.Log

	// 确保日志目录存在
	logDir := filepath.Dir(c.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 日志切割：基于文件大小 + 保留天数
	lumberjackLogger := &lumberjack.Logger{
		Filename:   c.FilePath,
		MaxSize:    c.MaxSize,    // MB
		MaxAge:     c.MaxAge,     // 天
		MaxBackups: c.MaxBackups, // 保留旧文件数
		Compress:   c.Compress,
		LocalTime:  true,
	}

	// 日志级别
	level := parseLevel(c.Level)

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 文件输出：JSON 格式
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		level,
	)

	// 控制台输出：彩色可读格式
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 同时输出到文件和控制台
	core := zapcore.NewTee(fileCore, consoleCore)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = Logger.WithOptions(zap.AddCallerSkip(1)).Sugar()

	return nil
}

// Close 刷新日志缓冲
func Close() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ========== 便捷日志方法（printf 风格） ==========

func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}
