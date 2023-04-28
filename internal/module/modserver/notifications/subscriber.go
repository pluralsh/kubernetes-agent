package notifications

import (
	"context"

	"google.golang.org/protobuf/proto"
)

const (
	GitPushEventsChannel = "git_push_events"
)

type Callback func(ctx context.Context, message proto.Message)

// Subscriber provides a `Subscribe` interface to receive messages that have
// been published using `Publisher` from the `notifications` module.
type Subscriber interface {
	Subscribe(ctx context.Context, channel string, callback Callback)
}
