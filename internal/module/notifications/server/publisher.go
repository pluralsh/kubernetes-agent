package server

import "context"

type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}
