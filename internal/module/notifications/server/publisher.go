package server

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Publisher provides a `Publish` interface to emit events.
// The `Subscriber` interface in the `modserver.notifications` module
// shall be used to subscribe to such events.
type Publisher interface {
	Publish(ctx context.Context, channel string, message proto.Message) error
}
