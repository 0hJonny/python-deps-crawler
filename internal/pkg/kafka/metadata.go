package kafka

import (
	"context"
)

type MetadataProducer struct {
	base      Producer
	extractor MetadataExtractor
}

func NewMetadataProducer(base Producer, extractor MetadataExtractor) *MetadataProducer {
	return &MetadataProducer{
		base:      base,
		extractor: extractor,
	}
}

func (p *MetadataProducer) SendData(ctx context.Context, topic string, data any, value []byte) error {
	key := p.extractor.ExtractKey(data)
	headers := p.extractor.ExtractHeaders(data)

	return p.base.SendMessage(ctx, topic, key, value, headers)
}

func (p *MetadataProducer) SendMessage(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error {
	return p.base.SendMessage(ctx, topic, key, value, headers)
}

func (p *MetadataProducer) Close() error {
	return p.base.Close()
}
