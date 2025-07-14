package logger

import (
	"fmt"

	"github.com/0hJonny/python-deps-crawler/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(cfg *config.LoggerConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Development,
		DisableCaller:     false,
		DisableStacktrace: !cfg.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         cfg.Encoding,
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig:    getEncoderConfig(cfg.Encoding),
	}

	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return logger, nil
}

func getEncoderConfig(encoding string) zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()

	if encoding == "console" {
		config = zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	config.CallerKey = "caller"
	config.EncodeCaller = zapcore.ShortCallerEncoder

	return config
}

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	zapLogger, err := NewZapLogger(cfg)
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

func (l *Logger) WithRequestID(requestID string) LoggerInterface {
	return &Logger{
		Logger: l.With(zap.String("request_id", requestID)),
	}
}

func (l *Logger) WithFields(fields map[string]any) LoggerInterface {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return &Logger{
		Logger: l.With(zapFields...),
	}
}
