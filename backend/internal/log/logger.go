package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func init() {
	Init("info", false)
}

func Init(level string, dev bool) error {
	cfg := zap.NewProductionConfig()
	if dev {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	defaultLogger = logger
	return nil
}

func L() *zap.Logger {
	if defaultLogger == nil {
		Init("info", false)
	}
	return defaultLogger
}

func Sync() {
	if defaultLogger != nil {
		defaultLogger.Sync()
	}
}

func SetOutput(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := zapcore.AddSync(f)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writer,
		defaultLogger.Level(),
	)
	defaultLogger = zap.New(core)
	return nil
}
