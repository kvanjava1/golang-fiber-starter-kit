package mail

import "context"

type Driver interface {
	Name() string
	Send(ctx context.Context, msg *Message) error
}
