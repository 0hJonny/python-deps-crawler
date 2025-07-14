package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type BaseConsumer struct {
	consumer sarama.ConsumerGroup
	config   *ConsumerConfig
}

type ConsumerConfig struct {
	Brokers       []string
	GroupID       string
	AutoCommit    bool
	InitialOffset int64
}

func NewBaseConsumer(config *ConsumerConfig) (*BaseConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = config.InitialOffset
	saramaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, saramaConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &BaseConsumer{
		consumer: consumer,
		config:   config,
	}, nil
}

type consumerGroupHandler struct {
	handler MessageHandler
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// Setup implements sarama.ConsumerGroupHandler.
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (
	h *consumerGroupHandler,
) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for message := range claim.Messages() {
		msg := &Message{
			Topic:     message.Topic,
			Key:       string(message.Key),
			Value:     message.Value,
			Headers:   h.convertHeaders(message.Headers),
			Partition: message.Partition,
			Offset:    message.Offset,
		}

		if err := h.handler(context.Background(), msg); err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func (h *consumerGroupHandler) convertHeaders(headers []*sarama.RecordHeader) map[string]string {
	result := make(map[string]string, len(headers))
	for _, header := range headers {
		result[string(header.Key)] = string(header.Value)
	}
	return result
}

func (
	c *BaseConsumer,
) Subscribe(
	ctx context.Context,
	topics []string,
	handler MessageHandler,
) error {
	consumerHandler := &consumerGroupHandler{handler: handler}

	for {
		if err := c.consumer.Consume(ctx, topics, consumerHandler); err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}
	}
}

func (c *BaseConsumer) Close() error {
	return c.consumer.Close()
}
