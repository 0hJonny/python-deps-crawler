package kafka

import (
	"context"
)

type Producer interface {
	SendMessage(
		ctx context.Context,
		topic string,
		key string,
		value []byte,
		headers map[string]string,
	) error
	Close() error
}

type MessageHandler func(ctx context.Context, message *Message) error

type Message struct {
	Topic     string
	Key       string
	Value     []byte
	Headers   map[string]string
	Partition int32
	Offset    int64
}

type Consumer interface {
	Subscribe(
		ctx context.Context,
		topics []string,
		handler MessageHandler,
	) error
	Close() error
}

type MetadataExtractor interface {
	ExtractKey(data any) string
	ExtractHeaders(data any) map[string]string
}
