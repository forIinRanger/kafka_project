package sender

import "context"

type MappedMessage struct {
	Value string
}

type Sender interface {
	Send(ctx context.Context, v MappedMessage) error
}
