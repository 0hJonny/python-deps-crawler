package kafka

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Producer interface {
	PublishEvent(ctx context.Context, message proto.Message) error
	Close() error
}
