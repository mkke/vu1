package conn

import "context"

type Conn interface {
	Send(ctx context.Context, req string) error
	SendAndRecv(ctx context.Context, req string) (resp []string, err error)
	Close() error
}
