package mocks

import (
	"context"

	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/kafka"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

type MockKafkaProducer struct {
	mock.Mock
}

var _ kafka.Producer = (*MockKafkaProducer)(nil) // Compile-time check

func NewMockKafkaProducer() *MockKafkaProducer {
	return &MockKafkaProducer{}
}

func (m *MockKafkaProducer) PublishEvent(ctx context.Context, message proto.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}
