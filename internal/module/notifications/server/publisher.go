package server

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Publisher interface {
	Publish(ctx context.Context, channel string, message proto.Message) error
}
