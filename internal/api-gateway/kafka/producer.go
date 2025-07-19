package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/0hJonny/python-deps-crawler/internal/pkg/kafka"
	eventspb "github.com/0hJonny/python-deps-crawler/pkg/proto/api_gateway_kafka_events"
	"github.com/IBM/sarama"
)

type APIGatewayProducer struct {
	producer *kafka.MetadataProducer
	topic    string
}

// interface check
var _ Producer = (*APIGatewayProducer)(nil)

func NewAPIGatewayProducer(brokers []string, topic string) (*APIGatewayProducer, error) {
	baseProducer, err := kafka.NewBaseProducer(&kafka.ProducerConfig{
		Brokers:           brokers,
		RequiredAcks:      sarama.WaitForAll,
		RetryMax:          3,
		CompressionType:   4, // LZ4
		EnableIdempotence: true,
	})
	if err != nil {
		return nil, err
	}

	// Оборачиваем в retry декоратор
	retryProducer := kafka.NewRetryProducer(baseProducer, 3, 1*time.Second)

	// Оборачиваем в metadata декоратор
	metadataProducer := kafka.NewMetadataProducer(
		retryProducer,
		&protobufMetadataExtractor{},
	)

	return &APIGatewayProducer{
		producer: metadataProducer,
		topic:    topic,
	}, nil
}

// PublishEvent отправляет protobuf событие
func (p *APIGatewayProducer) PublishEvent(ctx context.Context, event proto.Message) error {
	// Сериализуем protobuf
	data, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	// Используем декоратор для автоматического извлечения метаданных
	return p.producer.SendData(ctx, p.topic, event, data)
}

func (p *APIGatewayProducer) Close() error {
	return p.producer.Close()
}

type protobufMetadataExtractor struct{}

func (e *protobufMetadataExtractor) ExtractKey(data any) string {
	switch event := data.(type) {
	case *eventspb.AnalysisStartedEvent:
		return event.RequestId
	case *eventspb.AnalysisStatusEvent:
		return event.RequestId
	default:
		log.Printf("⚠️  Unknown event type: %T", data)
		return "unknown"
	}
}

func (e *protobufMetadataExtractor) ExtractHeaders(data any) map[string]string {
	switch event := data.(type) {
	case *eventspb.AnalysisStartedEvent:
		return map[string]string{
			"content-type": "application/x-protobuf",
			"event-type":   "AnalysisStartedEvent",
			"producer":     "api-gateway",
			"user-id":      event.UserId,
		}
	case *eventspb.AnalysisStatusEvent:
		return map[string]string{
			"content-type": "application/x-protobuf",
			"event-type":   "AnalysisStatusEvent",
			"producer":     "api-gateway",
			"service":      event.ServiceName,
		}
	default:
		return map[string]string{
			"content-type": "application/x-protobuf",
			"event-type":   "UnknownEvent",
			"producer":     "api-gateway",
		}
	}
}
