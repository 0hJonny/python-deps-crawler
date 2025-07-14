package mocks

import (
	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger мок для Logger, реализует LoggerInterface
type MockLogger struct {
	mock.Mock
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) WithRequestID(requestID string) logger.LoggerInterface {
	args := m.Called(requestID)
	return args.Get(0).(logger.LoggerInterface)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) logger.LoggerInterface {
	args := m.Called(fields)
	return args.Get(0).(logger.LoggerInterface)
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}
