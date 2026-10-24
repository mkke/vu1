package conn

import "context"

// Conn is a line-based connection to a VU-1 hub.
type Conn interface {
	// Send writes the request line req and does not wait for a reply.
	Send(ctx context.Context, req string) error
	// SendAndRecv writes the request line req and returns the lines read
	// up to and including the first line that starts with '<'.
	SendAndRecv(ctx context.Context, req string) (resp []string, err error)
	// Close closes the underlying port.
	Close() error
}
