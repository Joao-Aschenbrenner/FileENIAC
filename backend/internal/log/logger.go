// SPDX-License-Identifier: MIT
package log

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

type contextKey string

const correlationIDKey contextKey = "correlation_id"

// NewID returns a short random identifier suitable for request correlation.
func NewID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "xxxxxxxx"
	}
	return hex.EncodeToString(b)
}

// WithCorrelationID returns a new context carrying the provided correlation ID.
// If id is empty, a new ID is generated automatically.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	if id == "" {
		id = NewID()
	}
	return context.WithValue(ctx, correlationIDKey, id)
}

// WithContext returns a logger enriched with the correlation ID from ctx, if any.
func WithContext(ctx context.Context) *zap.Logger {
	logger := L()
	if id, ok := ctx.Value(correlationIDKey).(string); ok && id != "" {
		return logger.With(zap.String("correlation_id", id))
	}
	return logger
}
