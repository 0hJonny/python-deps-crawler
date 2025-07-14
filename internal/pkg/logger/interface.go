package logger

import "go.uber.org/zap"

type LoggerInterface interface {
	WithRequestID(requestID string) LoggerInterface
	WithFields(fields map[string]any) LoggerInterface
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Sync() error
}
