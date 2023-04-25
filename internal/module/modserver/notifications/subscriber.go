package notifications

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Callback func(ctx context.Context, message proto.Message) error

type Subscriber interface {
	Subscribe(ctx context.Context, channel string, callback Callback) error
}
