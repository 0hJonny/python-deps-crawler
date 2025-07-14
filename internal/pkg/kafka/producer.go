package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type BaseProducer struct {
	producer sarama.SyncProducer
	config   *ProducerConfig
}

type ProducerConfig struct {
	Brokers           []string
	RequiredAcks      sarama.RequiredAcks
	RetryMax          int
	CompressionType   sarama.CompressionCodec
	EnableIdempotence bool
}

func NewBaseProducer(config *ProducerConfig) (*BaseProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = config.RequiredAcks
	saramaConfig.Producer.Retry.Max = config.RetryMax
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Compression = config.CompressionType
	saramaConfig.Producer.Idempotent = config.EnableIdempotence
	saramaConfig.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &BaseProducer{
		producer: producer,
		config:   config,
	}, nil
}

func (
	p *BaseProducer,
) SendMessage(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
	headers map[string]string,
) error {
	saramaHeaders := make([]sarama.RecordHeader, 0, len(headers))

	for k, v := range headers {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: saramaHeaders,
	}

	partition, offset, err := p.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("✅ Message sent: topic=%s, key=%s, partition=%d, offset=%d",
		topic, key, partition, offset)

	return nil
}

func (p *BaseProducer) Close() error {
	return p.producer.Close()
}
