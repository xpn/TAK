package dialer

import "context"

type Dialer interface {
	Dial(ctx context.Context, target string) (Connection, error)
	Close() error
}

type Connection interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
}
