package kafka

import (
	"context"
	"fmt"
	"log"
	"time"
)

type RetryProducer struct {
	base       Producer
	maxRetries int
	retryDelay time.Duration
}

func NewRetryProducer(base Producer, maxRetries int, retryDelay time.Duration) *RetryProducer {
	return &RetryProducer{
		base:       base,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func (p *RetryProducer) SendMessage(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error {
	var lastErr error

	for i := 0; i <= p.maxRetries; i++ {
		if i > 0 {
			log.Printf("⏳ Retrying send message, attempt %d/%d", i, p.maxRetries)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.retryDelay):
			}
		}

		if err := p.base.SendMessage(ctx, topic, key, value, headers); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to send message after %d retries: %w", p.maxRetries, lastErr)
}

func (p *RetryProducer) Close() error {
	return p.base.Close()
}
